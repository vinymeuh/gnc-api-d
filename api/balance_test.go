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

var balanceTests = []struct {
	path    string
	status  int
	balance float64
	errmsg  string
}{
	{"/balance/1", http.StatusOK, -100.0, "current balance for account 1 is wrong"},
	{"/balance/2", http.StatusOK, -15.5, "current balance for account 2 is wrong"},
	{"/balance/0", http.StatusOK, 1374.5, "recursive current balance for account 0 is wrong"},
	{"/balance/0?norecursive", http.StatusOK, 1490.0, "non recursive current balance for account 0 is wrong"},
	{"/balance/0?from=2019-02-01&to=2019-02-20", http.StatusOK, 484.5, "recursive account balance between two dates is wrong"},
	{"/balance/0?type=X", http.StatusOK, -10.0, "account balance for a given type is wrong"},
}

func TestBalance(t *testing.T) {
	acts := models.Account{
		ID:   "0",
		Name: "Dummy Account",
		Type: "ROOT",
		Transactions: []*models.Transaction{
			{Date: "2019-01-01", Value: 1000.0},
			{Date: "2019-01-02", Value: -9.5, Num: "X"},
			{Date: "2019-01-03", Value: -0.5, Num: "X"},
			{Date: "2019-02-01", Value: 1000.0},
			{Date: "2019-02-03", Value: -500.0},
		},
		Children: []*models.Account{
			{
				ID:   "1",
				Name: "Account 1",
				Type: "BANK",
				Transactions: []*models.Transaction{
					{Date: "2019-01-05", Value: -100.0},
				},
			},
			{
				ID:   "2",
				Name: "Account 2",
				Type: "BANK",
				Transactions: []*models.Transaction{
					{Date: "2019-02-20", Value: -15.5},
				},
			},
		},
	}
	h := BalanceHandler{Data: &acts}

	for _, tt := range balanceTests {
		req, err := http.NewRequest("GET", tt.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		res := w.Result()
		assert.Equal(t, tt.status, res.StatusCode, "Status code is wrong.")

		var balance models.Balance
		json.NewDecoder(res.Body).Decode(&balance)
		assert.Equal(t, tt.balance, balance.Value, tt.errmsg)
	}
}
