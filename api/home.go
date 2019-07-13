// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("/accounts\n"))
	w.Write([]byte("/accounts/{id}\n"))
	w.Write([]byte("/accountypes\n"))
	w.Write([]byte("/balance/{id}\n"))
}
