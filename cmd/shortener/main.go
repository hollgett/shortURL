package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var routes = map[string]http.HandlerFunc{}

// func shortURLPost(w http.ResponseWriter, r *http.Request) {
// 	if handler, ok := routes[r.URL.Path]; ok {
// 		fmt.Fprintln(w, r.URL.Path)
// 		handler(w, r)
// 	} else {
// 		if r.Method != http.MethodPost {
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 			return
// 		}
// 		w.WriteHeader(http.StatusCreated)
// 		r.Header.Set("Content-Type", "text/plain")
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		shortLink := make([]byte, 8)
// 		for i := range shortLink {
// 			if rand.Intn(2) == 0 {
// 				shortLink[i] = byte(rand.Intn(26) + 65)
// 			} else {
// 				shortLink[i] = byte(rand.Intn(26) + 97)
// 			}
// 		}
// 		link := `/` + string(shortLink)
// 		routes[link] = func(w http.ResponseWriter, r *http.Request) {
// 			stockURL := string(body)
// 			if r.Method != http.MethodGet {
// 				w.WriteHeader(http.StatusMethodNotAllowed)
// 				fmt.Fprintln(w, "get")
// 				return
// 			}
// 			w.Header().Set("Location", stockURL)
// 			w.WriteHeader(http.StatusTemporaryRedirect)
// 		}
// 		fmt.Fprintln(w, `http://localhost:8080`+link)
// 	}

// }

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
	if r.Method == http.MethodPost {
		r.Header.Set("Content-Type", "text/plain")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		shortLink := make([]byte, 8)
		for i := range shortLink {
			if rand.Intn(2) == 0 {
				shortLink[i] = byte(rand.Intn(26) + 65)
			} else {
				shortLink[i] = byte(rand.Intn(26) + 97)
			}
		}
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
		fmt.Fprintln(w, `http://localhost:8080`+link)
	}
}

func main() {
	rtr := http.NewServeMux()
	rtr.HandleFunc(`/`, shortURLPost)
	if err := http.ListenAndServe(`:8080`, rtr); err != nil {
		panic(err)
	}
}
