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

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=9632c7e59d791cd231f4e2a79310d139&q=" + city)
	if err != nil {
		return weatherData{}, nil
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, nil
	}
	return d, nil
}

func main() {
	finalHander := http.HandlerFunc(final)

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset-utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.Handle("/", middlewareOne(middlewareTwo(finalHander)))
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
