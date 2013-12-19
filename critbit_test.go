// Copyright 2013 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package critbit

import "testing"

func keys(tr *Tree) (s []string) {
	tr.Iterate("", func(key string) bool {
		s = append(s, key)
		return true
	})
	return
}

func TestEmpty(t *testing.T) {
	tr := &Tree{}
	if keys(tr) != nil {
		t.Error("must be empty")
	}
	if tr.Contains("a") {
		t.Error("cannot contain a key")
	}
	if tr.Delete("a") {
		t.Error("cannot delete a key")
	}
}

func TestKeyOrder(t *testing.T) {
	tests := []struct {
		ins []string
		res []string
	}{
		{
			[]string{"x", "y", "z", "c", "c", "b", "b", "a", "a"},
			[]string{"a", "b", "c", "x", "y", "z"},
		},
		{
			[]string{"aaa", "aa", "a"},
			[]string{"a", "aa", "aaa"},
		},
		{
			[]string{"b", "a", "aa"},
			[]string{"a", "aa", "b"},
		},
		{
			[]string{"aa", "aaa", "aab", "ab", "ba", "bb", "bba", "bbb"},
			[]string{"aa", "aaa", "aab", "ab", "ba", "bb", "bba", "bbb"},
		},
	}
	for i, test := range tests {
		tr := &Tree{}
		for _, s := range test.ins {
			tr.Insert(s)
			if tr.Contains(s) {
				continue
			}
			t.Errorf("test %d did not contain %q after insert", i, s)
		}
		res := keys(tr)
		if len(res) != len(test.res) || tr.Len() != len(test.res) {
			t.Errorf("test %d unexpected length %d", i, len(res))
			continue
		}
		for j, s := range test.res {
			if res[j] == s {
				continue
			}
			t.Errorf("test %d unexpected element %q at %d", i, res[j], j)
		}
		for j, s := range test.res {
			if tr.Delete(s) {
				continue
			}
			t.Errorf("test %d could not delete %q at %d", i, res[j], j)
		}
	}
}

func TestDeleteUnknownKey(t *testing.T) {
	tr := &Tree{}
	if !tr.Insert("aa") {
		t.Error("failed to insert into empty tree")
	}
	if tr.Delete("ab") {
		t.Error("deleted unknown key")
	}
}

func TestIterate(t *testing.T) {
	tr := &Tree{}
	keys := []string{"aa", "aaa", "aab", "ab", "ba", "bb", "bba", "bbb"}

	for _, s := range keys {
		tr.Insert(s)
	}
	tests := []struct {
		prefix string
		keys   []string
	}{
		{"", keys},
		{"a", []string{"aa", "aaa", "aab", "ab"}},
		{"aa", []string{"aa", "aaa", "aab"}},
		{"aaa", []string{"aaa"}},
		{"aaaa", nil},
		{"c", nil},
	}
	for i, test := range tests {
		s := test.keys
		tr.Iterate(test.prefix, func(key string) bool {
			if len(s) < 1 {
				t.Errorf("test %d superfluous key %q", i, key)
				return true
			}
			if s[0] != key {
				t.Errorf("test %d got key %q expected %q", i, key, s[0])
			}
			s = s[1:]
			return true
		})
	}
}
