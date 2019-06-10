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

type AccountTypesHandler struct {
	Index map[string]*models.Account
}

func (ath *AccountTypesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	types := make(map[string]int)
	for _, act := range ath.Index {
		types[act.Type]++
	}

	resp, err := json.Marshal(types)
	if err != nil {
		log.Printf("Unable to marshall all accounts to JSON: %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", resp)
}
