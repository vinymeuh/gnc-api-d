// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package api

import (
	"log"
	"net/http"
)

func httpBadRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s 400 Bad Request", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad Request"))
}

func httpInternalServerError(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s 500 Internal Server Error", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func httpMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Allow", "GET")
	log.Printf("%s %s 405 Method Not Allowed", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 Method Not Allowed"))
}

func httpNotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s 404 Not Found", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
