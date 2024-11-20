package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type Routers struct {
	routes map[string]string
}

func newRouters() *Routers {
	return &Routers{
		routes: make(map[string]string),
	}
}

func (r *Routers) randomID() string {
	//create random URL path
	shortLink := make([]byte, 8)
	for i := range shortLink {
		if rand.Intn(2) == 0 {
			shortLink[i] = byte(rand.Intn(26) + 65)
		} else {
			shortLink[i] = byte(rand.Intn(26) + 97)
		}
	}
	// search exist link
	if _, ok := r.routes[string(shortLink)]; ok {
		return r.randomID()
	} else {
		return string(shortLink)
	}
}

func (rtr *Routers) shortURLPost(w http.ResponseWriter, r *http.Request) {
	//read request body
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//create short route
	link := rtr.randomID()
	rtr.routes[link] = string(body)
	//return response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, `http://localhost:8080/%s`, link)
}

func (rtr *Routers) shortURLGet(w http.ResponseWriter, r *http.Request) {
	//search exist short url and return original URL
	path := r.URL.Path
	fmt.Println(path)
	if len(path) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if originalURL, ok := rtr.routes[path[1:]]; ok {
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func shortURLmiddleware(rtr *Routers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//checking http method and redirect
		switch r.Method {
		case http.MethodPost:
			rtr.shortURLPost(w, r)
		case http.MethodGet:
			rtr.shortURLGet(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}

func main() {
	//data base shortURL
	rtrServer := newRouters()
	//start serve
	rtr := http.NewServeMux()
	rtr.HandleFunc(`/`, shortURLmiddleware(rtrServer))
	if err := http.ListenAndServe(`:8080`, rtr); err != nil {
		panic(err)
	}
}
