// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

import (
	"compress/gzip"
	"encoding/xml"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type xmlCountData struct {
	Type  string `xml:"type,attr"`
	Value int    `xml:",chardata"`
}

type xmlAccount struct {
	Name      string `xml:"name"`
	ID        string `xml:"id"`
	Type      string `xml:"type"`
	ParentID  string `xml:"parent"`
	Commodity string `xml:"commodity>space"`
	Parent    *xmlAccount
	Children  []*xmlAccount
}

type xmlTransaction struct {
	Num        string     `xml:"num"`
	DatePosted string     `xml:"date-posted>date"`
	Splits     []xmlSplit `xml:"splits>split"`
}

type xmlSplit struct {
	Value   string `xml:"value"`
	Account string `xml:"account"`
}

func LoadFromFile(path string) (*Account, map[string]*Account, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	zr, err := gzip.NewReader(f)
	switch err {
	case nil:
		defer zr.Close()
		return LoadFrom(zr)
	case gzip.ErrHeader: // uncompressed file
		f.Seek(0, 0)
		return LoadFrom(f)
	default:
		return nil, nil, err
	}
}

func LoadFrom(r io.Reader) (*Account, map[string]*Account, error) {
	var data *Account
	var index map[string]*Account

	var actsExpected int
	var actsRead int
	var trnsExpected int
	var trnsRead int

	rejected := make(map[string]int) // track rejected accounts and reject their transactions accordingly

	t1 := time.Now()

	decoder := xml.NewDecoder(r)
	for {
		token, err := decoder.Token()
		if err != nil && err != io.EOF {
			return nil, nil, err
		}
		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Space != "http://www.gnucash.org/XML/gnc" {
				continue
			}

			if se.Name.Local == "count-data" {
				var cd xmlCountData
				decoder.DecodeElement(&cd, &se)

				switch cd.Type {
				case "account":
					actsExpected = cd.Value
					// initialize accounts Index
					index = make(map[string]*Account, actsExpected)
				case "transaction":
					trnsExpected = cd.Value
				}
				continue
			}

			if se.Name.Local == "account" {
				var xmlact xmlAccount
				decoder.DecodeElement(&xmlact, &se)

				// Skip Account templates used in schedule action
				// See "<cmdty:space>template</cmdty:space>"
				if xmlact.Commodity == "template" {
					rejected[xmlact.ID] = 1
					continue
				}
				actsRead++

				// I hope Root Account is always the first account encountered
				if data == nil && xmlact.Type == "ROOT" && xmlact.Name == "Root Account" {
					data = &Account{ID: xmlact.ID, Name: xmlact.Name, Type: xmlact.Type}
					index[xmlact.ID] = data
					continue
				}

				// Attach this node to the accounts tree
				parent := index[xmlact.ParentID]
				if parent == nil {
					log.Printf("ParentID not found in index for Account '%s'", xmlact.Name)
					continue
				}
				act := Account{ID: xmlact.ID, Name: xmlact.Name, Type: xmlact.Type, Parent: parent}
				parent.Children = append(parent.Children, &act)

				index[xmlact.ID] = &act
				continue
			}

			if se.Name.Local == "transaction" {
				var xtrn xmlTransaction
				decoder.DecodeElement(&xtrn, &se)
				trnsRead++
				for _, split := range xtrn.Splits {
					act := index[split.Account]
					if act == nil {
						if rejected[split.Account] == 0 {
							log.Printf("Account '%s' not found in index for transaction", split.Account)
						}
						continue
					}
					trn := Transaction{
						Num:   xtrn.Num,
						Date:  strings.TrimSpace(strings.Split(xtrn.DatePosted, " ")[0]), // '2014-07-30 00:00:00 +0200', we keep only '2014-07-30'
						Value: stringToFloat(split.Value),
					}
					act.Transactions = append(act.Transactions, &trn)
				}
			}

		}
	}

	if actsRead != actsExpected {
		log.Printf("Read %d accounts when %d were expected", actsRead, actsExpected)
	}
	if trnsRead != trnsExpected {
		log.Printf("Read %d transactions when %d were expected", trnsRead, trnsExpected)
	}

	t2 := time.Now()
	duration := t2.Sub(t1)
	log.Printf("Gnucash data loaded in %s (%d accounts, %d transactions)", duration, actsRead, trnsRead)

	return data, index, nil
}

func stringToFloat(v string) float64 {
	i := strings.Split(v, "/")
	n, _ := strconv.ParseFloat(i[0], 10)
	d, _ := strconv.ParseFloat(i[1], 10)
	return n / d
}
