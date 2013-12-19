// Copyright 2013 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package critbit

import (
	"os"
	"testing"
	"text/scanner"
)

var words, tests []string

func initdata(b *testing.B) {
	if words != nil {
		return
	}
	var err error
	words, err = scan("critbit.go")
	if err != nil {
		b.Fatal(err)
	}
	tests, err = scan("critbit_test.go")
	if err != nil {
		b.Fatal(err)
	}
}

func scan(name string) (w []string, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s scanner.Scanner
	s.Init(f)
	for t := s.Scan(); t != scanner.EOF; t = s.Scan() {
		w = append(w, s.TokenText())
	}
	return
}

func BenchmarkMap(b *testing.B) {
	initdata(b)
	var count int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]struct{})
		for _, w := range words {
			m[w] = struct{}{}
		}
		count = 0
		for _, w := range tests {
			if _, ok := m[w]; ok {
				count++
			}
		}
	}
}

func BenchmarkTree(b *testing.B) {
	initdata(b)
	var count int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := new(Tree)
		for _, w := range words {
			t.Insert(w)
		}
		count = 0
		for _, w := range tests {
			if t.Contains(w) {
				count++
			}
		}
	}
}
