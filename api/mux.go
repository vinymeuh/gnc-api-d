// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/vinymeuh/gnc-api-d/models"
)

type Router struct {
	Data *models.Account
}

func NewRouter(data *models.Account) *Router {
	return &Router{Data: data}
}

func (mux *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.RequestURI())

	if r.Method != "GET" {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}

	if r.URL.Path == "/" {
		w.Header().Set("content-type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
		return
	}

	path := strings.Split(r.URL.Path, "/")
	switch path[1] {
	case "accounts":
		switch len(path) {
		case 2: // /accounts
			h := AccountsHandler{Data: mux.Data}
			h.ServeHTTP(w, r)
			return
		case 3:
			if path[2] != "" { // /accounts/{:id}
				h := AccountHandler{Data: mux.Data, ID: path[2]}
				h.ServeHTTP(w, r)
				return
			}
		}
	case "accounttypes":
		if len(path) == 2 { // /accounttypes
			h := AccountTypesHandler{Data: mux.Data}
			h.ServeHTTP(w, r)
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
