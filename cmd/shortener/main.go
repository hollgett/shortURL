package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var routes = map[string]http.HandlerFunc{}

func shortURLPost(w http.ResponseWriter, r *http.Request) {
	// checking for method get
	if r.Method == http.MethodGet {
		if handler, ok := routes[r.URL.Path]; ok {
			// starting handler if exists
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//checking for method post 
	if r.Method == http.MethodPost {
		r.Header.Set("Content-Type", "text/plain")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//create random URL path
		shortLink := make([]byte, 8)
		for i := range shortLink {
			if rand.Intn(2) == 0 {
				shortLink[i] = byte(rand.Intn(26) + 65)
			} else {
				shortLink[i] = byte(rand.Intn(26) + 97)
			}
		}
		//create short route   
		link := `/` + string(shortLink)
		routes[link] = func(w http.ResponseWriter, r *http.Request) {
			stockURL := string(body)
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Location", stockURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `http://localhost:8080`+link)
	}
}

func main() {
	rtr := http.NewServeMux()
	rtr.HandleFunc(`/`, shortURLPost)
	if err := http.ListenAndServe(`:8080`, rtr); err != nil {
		panic(err)
	}
}
