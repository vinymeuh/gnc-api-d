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
		httpMethodNotAllowed(w, r)
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
		default:
			httpBadRequest(w, r)
		}
		return
	case "accounttypes":
		switch len(path) {
		case 2: // /accounttypes
			h := AccountTypesHandler{Data: router.root}
			h.ServeHTTP(w, r)
		default:
			httpBadRequest(w, r)
		}
		return
	case "balance":
		switch len(path) {
		case 3: // /balance/{:id}
			h := BalanceHandler{Data: router.root}
			h.ServeHTTP(w, r)
		default:
			httpBadRequest(w, r)
		}
		return
	}

	httpNotFound(w, r)
}
