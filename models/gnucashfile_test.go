// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGnuCashFile(t *testing.T) {
	data, err := LoadFromFile("testdata/empty.gnucash")
	if assert.NoError(t, err) {
		root := data.FindByType("ROOT")
		assert.Equal(t, 1, len(root), "Problem while retrieve the account of type ROOT")
		assert.Equal(t, "Root Account", root[0].Name, "Problem while retrieve the Root Account")

		books := data.FindByName("Books")
		assert.Equal(t, 1, len(books), "Problem while retrieve the account 'Books'")

		actBooks := books[0]
		assert.Equal(t, "97c2d5b268164b479944e221ae0267f1", actBooks.ID, "Problem with 'Books' account ID")
		assert.Equal(t, 1, len(actBooks.Transactions), "Problem with 'Books' account transactions")

		trnBooks := actBooks.Transactions[0]
		assert.Equal(t, "2019-06-10", trnBooks.Date, "Problem with 'Books' transaction date")
		assert.Equal(t, 30.05, trnBooks.Value, "Problem with 'Books' transaction value")
	}
}

func TestLoadInvalidGnuCashFile(t *testing.T) {
	_, err := LoadFromFile("gnucashfile_test.go")
	assert.Error(t, err)
}

func TestLoadNotExistingFile(t *testing.T) {
	_, err := LoadFromFile("i_do_not_exist")
	assert.Error(t, err)
}
