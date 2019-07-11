// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/vinymeuh/gnc-api-d/models"
)

// Router will send incoming requests to dedicated handler
type Router struct {
	root *models.Account
}

// NewRouter returns a new Router instance
func NewRouter(root *models.Account) *Router {
	return &Router{root: root}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.RequestURI())

	if r.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}

	// returns api documentation
	if r.URL.Path == "/" {
		home(w, r)
		return
	}

	// main routing
	path := strings.Split(r.URL.Path, "/")
	switch path[1] {
	case "accounts":
		switch len(path) {
		case 2, 3: // /accounts or /accounts/{:id}
			h := AccountsHandler{Data: router.root}
			h.ServeHTTP(w, r)
			return
		default:
			log.Printf("%s %s 400 Bad Request", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request"))
			return
		}
	case "accounttypes":
		switch len(path) {
		case 2: // /accounttypes
			h := AccountTypesHandler{Data: router.root}
			h.ServeHTTP(w, r)
			return
		default:
			log.Printf("%s %s 400 Bad Request", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request"))
			return
		}
	case "balance":
		if len(path) == 2 { // /balance
			log.Printf("not yet implemented")
			return
		}
	}

	log.Printf("%s %s 404 Not Found", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
