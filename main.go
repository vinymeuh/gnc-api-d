// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vinymeuh/gnc-api-d/api"
	"github.com/vinymeuh/gnc-api-d/models"
)

func getLogFilePath() string {
	return os.Getenv("LOG_FILE_PATH")
}

func getListenAddress() string {
	addr := os.Getenv("LISTEN_ADDRESS")
	if addr == "" {
		return "localhost:8000"
	}
	return addr
}

func getGnuCashFile() string {
	return os.Getenv("GNUCASH_FILE_PATH")
}

func main() {
	setupLogFile()

	// load Gnucash data
	gncfile := getGnuCashFile()
	log.Printf("Loading GnuCash file '%s'", gncfile)
	data, index, err := models.LoadFromFile(gncfile)
	if err != nil {
		log.Fatal(err)
	}

	// start HTTP server
	r := api.NewRouter(data, index)
	addr := getListenAddress()
	log.Printf("Starting HTTP server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
