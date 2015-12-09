package main

import (
	"log"
	"net/http"
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

func main() {
	finalHander := http.HandlerFunc(final)

	http.Handle("/", middlewareOne(middlewareTwo(finalHander)))
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
