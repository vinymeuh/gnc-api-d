// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinymeuh/gnc-api-d/models"
)

var accountsByIDTests = []struct {
	method string
	path   string
	status int
}{
	{"GET", "/accounts/0", http.StatusOK},
	{"GET", "/accounts/666", http.StatusNotFound},
	{"GET", "/accounts/", http.StatusBadRequest},
}

func TestAccountsHandlerByID(t *testing.T) {
	acts := models.Account{
		ID:   "0",
		Type: "ROOT",
	}
	h := AccountsHandler{Data: &acts}

	for _, tt := range accountsByIDTests {
		req, err := http.NewRequest(tt.method, tt.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, tt.status, res.StatusCode, "Status code is wrong.")
	}
}

var accountsByNameOrTypeTests = []struct {
	method string
	path   string
	status int
	count  int
}{
	{"GET", "/accounts?name=Dummy", http.StatusOK, 1},
	{"GET", "/accounts?name=NotExisting", http.StatusOK, 0},
	{"GET", "/accounts?type=ROOT", http.StatusOK, 1},
	{"GET", "/accounts?type=FAKE", http.StatusOK, 0},
	{"GET", "/accounts?type=ROOT&Name=Dummy", http.StatusBadRequest, 0},
}

func TestAccountsHandlerByNameOrType(t *testing.T) {
	acts := models.Account{
		ID:   "0",
		Type: "ROOT",
		Name: "Dummy",
	}
	h := AccountsHandler{Data: &acts}

	for _, tt := range accountsByNameOrTypeTests {
		req, err := http.NewRequest(tt.method, tt.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, tt.status, res.StatusCode, "Status code is wrong.")

		var body []interface{}
		json.NewDecoder(res.Body).Decode(&body)
		assert.Equal(t, tt.count, len(body), "number of results does not match")
	}
}
