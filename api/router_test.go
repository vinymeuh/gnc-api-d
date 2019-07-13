// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinymeuh/gnc-api-d/models"
)

var routesToTest = []struct {
	method string
	path   string
	status int
}{
	// OK
	{"GET", "/", http.StatusOK},
	{"GET", "/accounts", http.StatusOK},
	{"GET", "/accounts/0", http.StatusOK},
	{"GET", "/accounttypes", http.StatusOK},
	{"GET", "/balance/0", http.StatusOK},
	// Not Found
	{"GET", "/not-exists", http.StatusNotFound},
	// Not Allowed
	{"POST", "/", http.StatusMethodNotAllowed},
	// Bad Request
	{"GET", "/accounts/0/1", http.StatusBadRequest},
	{"GET", "/accounttypes/0", http.StatusBadRequest},
	{"GET", "/balance", http.StatusBadRequest},
}

func TestRoutes(t *testing.T) {

	var root = models.Account{
		ID:   "0",
		Name: "Dummy",
		Type: "ROOT",
	}
	r := NewRouter(&root)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var res *http.Response
	var err error
	for _, tt := range routesToTest {
		switch tt.method {
		case "GET":
			res, err = http.Get(ts.URL + tt.path)
		case "POST":
			res, err = http.Post(ts.URL+tt.path, "text/plain", nil)
		default:
			t.Fatalf("Testing HTTP method %s is not supported", tt.method)
		}
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tt.status, res.StatusCode, "Status code for %s is wrong.", tt)
	}
}
