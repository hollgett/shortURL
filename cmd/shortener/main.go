package main

import (
	"fmt"
	"net/http"
)

func shortUrlPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, `http://localhost:8080/EwHXdJfB `)
}

func shortUrlGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Location", "https://practicum.yandex.ru/")
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func main() {
	rtr := http.NewServeMux()
	rtr.HandleFunc(`/`, shortUrlPost)
	rtr.HandleFunc(`/EwHXdJfB`, shortUrlGet)
	http.ListenAndServe(`:8080`, rtr)
}
