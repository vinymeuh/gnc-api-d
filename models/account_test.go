// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	act := Account{
		ID:   "0",
		Name: "Dummy Account",
		Type: "ROOT",
		Transactions: []*Transaction{
			{Date: "2019-01-01", Value: 1000.0},
			{Date: "2019-01-02", Value: -9.5, Num: "X"},
			{Date: "2019-01-03", Value: -0.5, Num: "X"},
			{Date: "2019-02-01", Value: 1000.0},
			{Date: "2019-02-03", Value: -500.0},
		},
	}

	assert.Equal(t, 1490.0, act.Balance(BalanceOptions{}), "current Balance is incorrect")
	assert.Equal(t, 990.0, act.Balance(BalanceOptions{To: "2019-01-03"}), "Balance at a defined Date is incorrect")
	assert.Equal(t, 499.5, act.Balance(BalanceOptions{From: "2019-01-03", To: "2019-02-03"}), "Balance between 2 Dates is incorrect")
	assert.Equal(t, -10.0, act.Balance(BalanceOptions{Type: "X"}), "current Balance for a defined Type is incorrect")
}
