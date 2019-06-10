// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

import (
	"compress/gzip"
	"encoding/xml"
	"io"
	"log"
	"os"
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
	Num         string     `xml:"num"`
	Description string     `xml:"description"`
	DatePosted  string     `xml:"date-posted>date"`
	Splits      []xmlSplit `xml:"splits>split"`
}

type xmlSplit struct {
	ReconciledState string `xml:"reconciled-state"`
	ReconcileDate   string `xml:"reconcile-date>date"`
	Value           string `xml:"value"`
	Account         string `xml:"account"`
	Memo            string `xml:"memo"`
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
					log.Printf("Account '%s' has no ParentID", xmlact.Name) // Template Root & ses fils ?
					continue
				}
				//act.Parent = parent
				act := Account{ID: xmlact.ID, Name: xmlact.Name, Type: xmlact.Type, Parent: parent}
				parent.Children = append(parent.Children, &act)

				index[xmlact.ID] = &act
				continue
			}

			// if se.Name.Local == "transaction" {
			// 	var xtrn XMLTransaction
			// 	decoder.DecodeElement(&xtrn, &se)
			// 	trnsRead++
			// 	for _, split := range xtrn.Splits {
			// 		fmt.Println(split)
			// 		fmt.Printf("%s;%s;%s;%s;%s;%s;%s;%s\n",
			// 			xtrn.Num,
			// 			xtrn.Description,
			// 			xtrn.DatePosted,
			// 			split.ReconciledState,
			// 			split.ReconcileDate,
			// 			split.Value,
			// 			split.Account,
			// 			split.Memo)
			// 	}
			// }

		}
	}

	if actsRead != actsExpected {
		log.Printf("Read %d accounts when %d were expected\n", actsRead, actsExpected)
	}
	if trnsRead != trnsExpected {
		log.Printf("Read %d transactions when %d were expected\n", trnsRead, trnsExpected)
	}

	t2 := time.Now()
	duration := t2.Sub(t1)
	log.Printf("Gnucash data loaded in %s", duration)

	return data, index, nil
}
