// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

import (
	"compress/gzip"
	"encoding/xml"
	"errors"
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

// LoadFromFile loads data from a GnuCash file compressed or not
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
		return Load(zr)
	case gzip.ErrHeader: // uncompressed file
		f.Seek(0, 0)
		return Load(f)
	default:
		return nil, nil, err
	}
}

// Load loads GnuCash account hierarchy from a XML document
func Load(r io.Reader) (*Account, map[string]*Account, error) {
	var root *Account
	var index map[string]*Account

	type countData struct {
		accounts     int
		transactions int
	}
	var expected countData
	var read countData

	actTemplates := make(map[string]int) // track template accounts and reject their transactions accordingly

	t1 := time.Now()

	decoder := xml.NewDecoder(r)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
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
					expected.accounts = cd.Value
					// initialize accounts Index
					index = make(map[string]*Account, expected.accounts)
				case "transaction":
					expected.transactions = cd.Value
				}
				continue
			}

			if se.Name.Local == "account" {
				var xmlact xmlAccount
				decoder.DecodeElement(&xmlact, &se)

				// Skip Account templates used in schedule action
				// See "<cmdty:space>template</cmdty:space>"
				if xmlact.Commodity == "template" {
					actTemplates[xmlact.ID] = 1
					continue
				}
				read.accounts++

				// I hope Root Account is always the first account encountered
				if root == nil && xmlact.Type == "ROOT" && xmlact.Name == "Root Account" {
					root = &Account{ID: xmlact.ID, Name: xmlact.Name, Type: xmlact.Type}
					index[xmlact.ID] = root
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
				read.transactions++
				for _, split := range xtrn.Splits {
					act := index[split.Account]
					if act == nil {
						if actTemplates[split.Account] == 0 {
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

	if read.accounts != expected.accounts {
		log.Printf("Read %d accounts when %d were expected", read.accounts, expected.accounts)
	}
	if read.transactions != expected.transactions {
		log.Printf("Read %d transactions when %d were expected", read.transactions, expected.transactions)
	}

	t2 := time.Now()
	duration := t2.Sub(t1)
	log.Printf("Gnucash data loaded in %s (%d accounts, %d transactions)", duration, read.accounts, read.transactions)

	if root == nil {
		return root, index, errors.New("Unable to parse XML file")
	}
	return root, index, nil
}

func stringToFloat(v string) float64 {
	i := strings.Split(v, "/")
	n, _ := strconv.ParseFloat(i[0], 10)
	d, _ := strconv.ParseFloat(i[1], 10)
	return n / d
}
