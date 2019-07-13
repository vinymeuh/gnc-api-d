// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	home(w, nil)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is wrong.")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	want := "/accounts\n/accounts/{id}\n/accountypes\n/balance/{id}\n"
	assert.Equal(t, want, string(body), "Response body is wrong.")
}
