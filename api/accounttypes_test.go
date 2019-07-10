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

func TestAccountTypesHandler(t *testing.T) {

	acts := models.Account{
		ID:   "0",
		Type: "ROOT",
		Children: []*models.Account{
			{
				ID:   "1",
				Type: "BANK",
			},
			{
				ID:   "2",
				Type: "BANK",
			},
		},
	}
	h := AccountTypesHandler{Data: &acts}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, nil)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is wrong.")

	var body map[string]interface{}
	json.NewDecoder(res.Body).Decode(&body)

	assert.Equal(t, 2, len(body), "number of account types does not match")
	assert.Equal(t, float64(1), body["ROOT"], "number of ROOT accounts does not match")
	assert.Equal(t, float64(2), body["BANK"], "number of BANK accounts does not match")
}
