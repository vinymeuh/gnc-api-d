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

func getListenAddress() string {
	addr := os.Getenv("LISTEN_ADDRESS")
	if addr == "" {
		return "localhost:8000"
	}
	return addr
}

func getGnuCashFile() string {
	file := os.Getenv("GNUCASH_FILE_PATH")
	if file == "" {
		log.Printf("variable GNUCASH_FILE_PATH is not defined")
		os.Exit(1)
	}
	return file
}

func main() {
	setupLog()

	// load Gnucash data
	gncfile := getGnuCashFile()
	log.Printf("Loading GnuCash file '%s'", gncfile)
	root, err := models.LoadFromFile(gncfile)
	if err != nil {
		log.Fatal(err)
	}

	// start HTTP server
	r := api.NewRouter(root)
	addr := getListenAddress()
	log.Printf("Starting HTTP server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
