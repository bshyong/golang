package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// middleware
func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("executing mid one")
		next.ServeHTTP(w, r)
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("mid two")
		next.ServeHTTP(w, r)
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("final handler")
	w.Write([]byte("ok"))
}

func (w openWeatherMap) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=9632c7e59d791cd231f4e2a79310d139&q=" + city)
	if err != nil {
		return 0, nil
	}

	defer resp.Body.Close()

	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, nil
	}
	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}

func (w multiWeatherProvider) temperature(city string) (float64, error) {
	// channels for concurrency
	temps := make(chan float64, len(w))
	errs := make(chan error, len(w))

	// spawn goroutine with anonymous function for each provider
	for _, provider := range w {
		go func(p weatherProvider) {
			k, err := p.temperature(city)
			if err != nil {
				errs <- err
				return
			}
			temps <- k
		}(provider)
	}

	sum := 0.0

	for i := 0; i < len(w); i++ {
		select {
		case temp := <-temps:
			sum += temp
		case err := <-errs:
			return 0, err
		}
	}

	return sum / float64(len(w)), nil
}

func main() {
	mw := multiWeatherProvider{
		openWeatherMap{},
	}
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := mw.temperature(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset-utf-8")
		json.NewEncoder(w).Encode(data)
	})

	finalHander := http.HandlerFunc(final)
	http.Handle("/", middlewareOne(middlewareTwo(finalHander)))
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
