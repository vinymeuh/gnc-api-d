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

type BalanceHandler struct {
	Data *models.Account // Root !
}

func (bh *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := strings.Split(r.URL.Path, "/")[2]
	log.Printf("id = %s", id)
	act := bh.Data.FindByID(id)

	opts := models.BalanceOptions{Recursive: true}
	params := r.URL.Query()
	if _, ok := params["norecursive"]; ok {
		opts.Recursive = false
	}
	if from := params.Get("from"); from != "" {
		opts.From = from
	}
	if to := params.Get("to"); to != "" {
		opts.To = to
	}
	if aType := params.Get("type"); aType != "" {
		opts.Type = aType
	}
	log.Printf("%v", opts)

	value := act.Balance(opts)

	resp, err := json.Marshal(value)
	if err != nil {
		log.Printf("Unable to marshall balance: %s\n", err)
		httpInternalServerError(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", resp)
	log.Printf("%s", resp)
	return
}
