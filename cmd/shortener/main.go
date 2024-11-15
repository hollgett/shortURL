package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var routes = map[string]http.HandlerFunc{}

func randomID() string {
	//create random URL path
	shortLink := make([]byte, 8)
	for i := range shortLink {
		if rand.Intn(2) == 0 {
			shortLink[i] = byte(rand.Intn(26) + 65)
		} else {
			shortLink[i] = byte(rand.Intn(26) + 97)
		}
	}
	return string(shortLink)
}

func shortURLPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.Header.Set("Content-Type", "text/plain")
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//create short route
	link := randomID()
	//handler body short URL
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
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, `http://localhost:8080/%s`, link)
}

func shortURLGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if handler, ok := routes[r.URL.Path[1:]]; ok {
		// starting handler if exists
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	rtr := http.NewServeMux()
	rtr.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			shortURLPost(w, r)
		} else {
			shortURLGet(w,r)
		}
	})
	if err := http.ListenAndServe(`:8080`, rtr); err != nil {
		panic(err)
	}
}
