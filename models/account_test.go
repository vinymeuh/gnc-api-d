// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var accountBalanceTests = []struct {
	options		BalanceOptions
	expected 	float64
	errmsg  	string
}{
	{BalanceOptions{}, 1490.0, "Non recursive current Balance is incorrect"},
	{BalanceOptions{Recursive: true}, 1374.5, "Recursive current Balance is incorrect"},
	{BalanceOptions{To: "2019-01-03"}, 990.0, "Balance at a defined Date is incorrect"},
	{BalanceOptions{From: "2019-01-03", To: "2019-02-03"}, 499.5, "Balance between 2 Dates is incorrect"},
	{BalanceOptions{Type: "X"}, -10.0, "Current Balance for a defined Type is incorrect"},
}

func TestAccount(t *testing.T) {
	acts := Account{
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
		Children: []*Account{
			{
				ID:   "1",
				Name: "Account 1",
				Type: "BANK",
				Transactions: []*Transaction{
					{Date: "2019-01-05", Value: -100.0},
				},
			},
			{
				ID:   "2",
				Name: "Account 2",
				Type: "BANK",
				Transactions: []*Transaction{
					{Date: "2019-02-20", Value: -15.5},
				},
			},
		},
	}

	root := acts.FindByType("ROOT")[0]
	assert.Equal(t, &acts, root, "Problem while retrieve root account by Type")

	assert.Nil(t, acts.FindByID("not existing"), "Problem while retrieve not existing account by ID")

	act1 := acts.FindByID("1")
	assert.Equal(t, root.Children[0], act1, "Problem while retrieve account 1 by ID")

	act2 := acts.FindByName("Account 2")[0]
	assert.Equal(t, root.Children[1], act2, "Problem while retrieve account 2 by Name")

	children := root.Descendants()
	assert.Equal(t, 2, len(children), "Problem with root account descendants")

	for _, tt := range accountBalanceTests {
		assert.Equal(t, tt.expected, root.Balance(tt.options).Value, tt.errmsg)
	}
}
