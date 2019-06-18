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
			&Transaction{Date: "2019-01-01", Value: 1000.0},
			&Transaction{Date: "2019-01-02", Value: -9.5},
			&Transaction{Date: "2019-01-03", Value: -0.5},
			&Transaction{Date: "2019-02-01", Value: 1000.0},
			&Transaction{Date: "2019-02-03", Value: -500.0},
		},
	}

	assert.Equal(t, 1490.0, act.Balance(""), "current Balance is incorrect")
	assert.Equal(t, 990.0, act.Balance("2019-01-03"), "Balance at 2019-01-03 is incorrect")
}
