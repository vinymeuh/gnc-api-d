// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadGnuCashFile(t *testing.T) {
	data, index, err := LoadFromFile("testdata/empty.gnucash")
	if assert.NoError(t, err) {
		root := data.FindByType("ROOT")
		assert.Equal(t, 1, len(root), "Problem while retrieve the account of type ROOT")
		assert.Equal(t, "Root Account", root[0].Name, "Problem while retrieve the Root Account")

		alls := root[0].All()
		assert.Equal(t, len(index), len(alls), "accounts index is incoherent")
	}
}
