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
	Name     string `xml:"name"`
	ID       string `xml:"id"`
	Type     string `xml:"type"`
	ParentID string `xml:"parent"`
	Parent   *xmlAccount
	Children []*xmlAccount
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
func LoadFromFile(path string) (*Account, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
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
		return nil, err
	}
}

// Load loads GnuCash account hierarchy from a XML document
// Returns a pointer to the root account of the hierarchy
func Load(r io.Reader) (*Account, error) {
	var root *Account
	var actsIndex map[string]*Account

	type countData struct {
		acts int
		trns int
	}
	var expected countData
	var read countData

	t1 := time.Now()

	decoder := xml.NewDecoder(r)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
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
					expected.acts = cd.Value
					// initialize accounts Index
					actsIndex = make(map[string]*Account, expected.acts)
				case "transaction":
					expected.trns = cd.Value
				}
				continue
			}

			if se.Name.Local == "account" {
				var xmlact xmlAccount
				decoder.DecodeElement(&xmlact, &se)
				read.acts++

				// I hope Root Account is always the first account encountered
				if root == nil {
					if xmlact.Type == "ROOT" && xmlact.Name == "Root Account" {
						root = &Account{ID: xmlact.ID, Name: xmlact.Name, Type: xmlact.Type}
						actsIndex[xmlact.ID] = root
						continue
					}
					return root, errors.New("Unable to initialize accounts hierarchy with Root Account")
				}

				// Attach this node to the accounts tree
				parent := actsIndex[xmlact.ParentID]
				if parent == nil {
					log.Printf("ParentID not found in index for Account '%s'", xmlact.Name)
					continue
				}
				act := Account{ID: xmlact.ID, Name: xmlact.Name, Type: xmlact.Type, Parent: parent}
				parent.Children = append(parent.Children, &act)

				actsIndex[xmlact.ID] = &act
				continue
			}

			if se.Name.Local == "transaction" {
				var xtrn xmlTransaction
				decoder.DecodeElement(&xtrn, &se)
				read.trns++
				for _, split := range xtrn.Splits {
					act := actsIndex[split.Account]
					if act == nil {
						log.Printf("Account '%s' not found in index for transaction", split.Account)
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

			// Skip all accounts and transactions templates used in schedule action
			if se.Name.Local == "template-transactions" {
				decoder.Skip()
			}

		}
	}

	if read.acts != expected.acts {
		log.Printf("Read %d accounts when %d were expected", read.acts, expected.acts)
	}
	if read.trns != expected.trns {
		log.Printf("Read %d transactions when %d were expected", read.trns, expected.trns)
	}

	t2 := time.Now()
	duration := t2.Sub(t1)
	log.Printf("Gnucash data loaded in %s (%d accounts, %d transactions)", duration, read.acts, read.trns)

	if root == nil {
		return root, errors.New("Unable to parse XML file")
	}
	return root, nil
}

func stringToFloat(v string) float64 {
	i := strings.Split(v, "/")
	n, _ := strconv.ParseFloat(i[0], 10)
	d, _ := strconv.ParseFloat(i[1], 10)
	return n / d
}
