// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

import (
	"time"
)

// Account is a node of the accounts hierarchy.
// Each account has its own  list of transactions
type Account struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Parent       *Account       `json:"-"`
	Children     []*Account     `json:"-"`
	Transactions []*Transaction `json:"-"`
}

// Transaction keeps data for a transaction
type Transaction struct {
	Num   string  `json:"-"`
	Date  string  `json:"-"` // YYYY-MM-DD
	Value float64 `json:"-"`
}

// WalkAccountFunc is the type of the function called for each account visited by WalkBFS
type WalkAccountFunc func(act *Account) bool

// WalkBFS traverses the tree of accounts using Breadth-first search algorithm.
// Starting at node a, returns the list of accounts for which walkFunc is true.
func (a *Account) WalkBFS(walkFunc WalkAccountFunc) []*Account {
	acts := make([]*Account, 0)

	queue := make([]*Account, 0)
	queue = append(queue, a)
	for len(queue) > 0 {
		act := queue[0]
		if walkFunc(act) == true {
			acts = append(acts, act)
		}
		queue = queue[1:]
		if len(act.Children) > 0 {
			for _, child := range act.Children {
				queue = append(queue, child)
			}
		}
	}

	return acts
}

// Descendants return the list of sub-accounts for an account
func (a *Account) Descendants() []*Account {
	return a.WalkBFS(func(act *Account) bool { return true })[1:]
}

// FindByID returns an accounts matching ID
func (a *Account) FindByID(ID string) *Account {
	acts := a.WalkBFS(func(act *Account) bool { return act.ID == ID })
	if len(acts) == 0 {
		return nil
	}
	return acts[0]
}

// FindByName returns a list of accounts matching name
func (a *Account) FindByName(name string) []*Account {
	return a.WalkBFS(func(act *Account) bool { return act.Name == name })
}

// FindByType returns a list of accounts matching type
func (a *Account) FindByType(atype string) []*Account {
	return a.WalkBFS(func(act *Account) bool { return act.Type == atype })
}

// BalanceOptions is the type used as input parameters for the Balance function
type BalanceOptions struct {
	From      string
	To        string
	Type      string
	Recursive bool
}

// Balance is the type used to return result for the Balance function
type Balance struct {
	Date  string
	Value float64
}

// Balance returns the amount of the account
func (a *Account) Balance(opts BalanceOptions) Balance {

	if opts.To == "" {
		opts.To = time.Now().Format("2006-01-02")
	}

	var b float64
	// transactions directly attached to the account
	for _, t := range a.Transactions {
		if opts.Type != "" && t.Num != opts.Type {
			continue
		}
		if t.Date >= opts.From && t.Date <= opts.To {
			b = b + t.Value
		}
	}
	// transactions on sub-accounts
	if opts.Recursive {
		for _, sa := range a.Children {
			b = b + sa.Balance(opts).Value
		}
	}

	return Balance{Date: opts.To, Value: b}
}
