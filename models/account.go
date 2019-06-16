// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package models

type Account struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Parent       *Account       `json:"-"`
	Children     []*Account     `json:"-"`
	Transactions []*Transaction `json:"-"`
}

type Transaction struct {
	Num   string  `json:"-"`
	Date  string  `json:"-"` // YYYY-MM-DD
	Value float64 `json:"-"`
}

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

func (a *Account) Descendants() []*Account {
	return a.WalkBFS(func(act *Account) bool { return true })
}

func (a *Account) FindByName(name string) []*Account {
	return a.WalkBFS(func(act *Account) bool { return act.Name == name })
}

func (a *Account) FindByType(atype string) []*Account {
	return a.WalkBFS(func(act *Account) bool { return act.Type == atype })
}
