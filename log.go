// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package main

import (
	"io"
	"log"
	"os"
	"time"
)

func getLogFilePath() string {
	return os.Getenv("LOG_FILE_PATH")
}

type writer struct {
	io.Writer
}

func (w writer) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format("2006-01-02T15:04:05-07:00 ")), b...))
}

func setupLog() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(writer{os.Stdout})
	if logFilePath := getLogFilePath(); logFilePath != "" {
		f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(writer{f})
	}
}
