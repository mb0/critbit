// Copyright 2013 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package critbit

import (
	"os"
	"sort"
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

func BenchmarkMapSort(b *testing.B) {
	initdata(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]struct{})
		for _, w := range words {
			m[w] = struct{}{}
		}
		s := make([]string, 0, len(m))
		for w := range m {
			s = append(s, w)
		}
		sort.Strings(s)
	}
}

func BenchmarkTreeSort(b *testing.B) {
	initdata(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := new(Tree)
		for _, w := range words {
			t.Insert(w)
		}
		s := make([]string, 0, t.Len())
		t.Iterate("", func(key string) bool {
			s = append(s, key)
			return true
		})
	}
}

var ib = []byte{^byte(1 << 0), ^byte(1 << 1), ^byte(1 << 2), ^byte(1 << 3), ^byte(1 << 4), ^byte(1 << 5), ^byte(1 << 6), ^byte(1 << 7)}

func BenchmarkBitCond(b *testing.B) {
	var count uint32
	for i := 0; i < b.N; i++ {
		count = 0
		for j := byte(0); j < 0xff; j++ {
			for k := 0; k < 8; k++ {
				if j&^ib[k] != 0 {
					count++
				}
			}
		}
	}
}

func BenchmarkBitArith(b *testing.B) {
	var count uint32
	for i := 0; i < b.N; i++ {
		count = 0
		for j := byte(0); j < 0xff; j++ {
			for k := 0; k < 8; k++ {
				count += (1 + uint32(j|ib[k])) >> 8
			}
		}
	}
}
