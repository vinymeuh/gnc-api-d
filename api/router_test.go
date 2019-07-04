// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validRoutes = []string{
	"/",
	//"/accounts",
	//"/accounttypes",
	"/balance",
}

func TestValidRoutes(t *testing.T) {
	r := NewRouter(nil)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range validRoutes {
		res, err := http.Get(ts.URL + tt)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, res.StatusCode, "Status code for %s is wrong.", tt)
	}
}

func TestInvalidRoutes(t *testing.T) {
	r := NewRouter(nil)
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusNotFound, res.StatusCode, "Status code for /not-exists is wrong.")
}
