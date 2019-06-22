// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vinymeuh/gnc-api-d/models"
)

type AccountHandler struct {
	Data *models.Account
	ID   string
}

func (ah *AccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	act := ah.Data.FindByID(ah.ID)

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
