// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/vinymeuh/gnc-api-d/models"
)

type BalanceHandler struct {
	Data *models.Account
}

func (bh *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
