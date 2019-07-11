// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/vinymeuh/gnc-api-d/models"
)

type AccountsHandler struct {
	Data *models.Account
}

func (ah *AccountsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	switch len(path) { // +1 for leading /
	case 2:
		ah.serveAccountsByNameOrType(w, r)
	case 3:
		id := path[2]
		if id == "" {
			log.Printf("%s %s 400 Bad Request", r.Method, r.URL.Path)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		ah.serveAccountByID(w, r, id)
	default:
		log.Printf("%s %s 400 Bad Request", r.Method, r.URL.Path)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func (ah *AccountsHandler) serveAccountsByNameOrType(w http.ResponseWriter, r *http.Request) {
	var acts []*models.Account

	params := r.URL.Query()
	switch len(params) {
	case 0:
		acts = ah.Data.Descendants()
	case 1:
		if name := params.Get("name"); name != "" {
			acts = ah.Data.FindByName(name)
		}
		if atype := params.Get("type"); atype != "" {
			acts = ah.Data.FindByType(atype)
		}
	}

	if acts == nil {
		log.Printf("%s %s 400 Bad Request", r.Method, r.URL.Path)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(acts)
	if err != nil {
		log.Printf("Unable to marshall all accounts to JSON: %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", resp)
}

func (ah *AccountsHandler) serveAccountByID(w http.ResponseWriter, r *http.Request, id string) {
	act := ah.Data.FindByID(id)
	if act != nil {
		resp, err := json.Marshal(act)
		if err != nil {
			log.Printf("Unable to marshall all accounts to JSON: %s\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", resp)
		return
	}
	log.Printf("%s %s 404 Not Found", r.Method, r.URL.Path)
	http.Error(w, "Not Found", http.StatusNotFound)
}
