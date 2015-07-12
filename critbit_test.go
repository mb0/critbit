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
		for j := len(res) - 1; j >= 0; j-- {
			if tr.Delete(res[j]) {
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

func testStringsEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestKeys0(t *testing.T) {
	tr := &Tree{}
	expected := []string{}

	returned_keys := tr.Keys()
	if !testStringsEq(returned_keys, expected) {
		t.Errorf("Got: %q", returned_keys)
	}
}

func TestKeys1(t *testing.T) {
	tr := &Tree{}
	orig_keys := []string{"aa"}
	expected := []string{"aa"}

	for _, s := range orig_keys {
		tr.Insert(s)
	}
	returned_keys := tr.Keys()
	if !testStringsEq(returned_keys, expected) {
		t.Errorf("Got: %q", returned_keys)
	}
}

func TestKeysMany(t *testing.T) {
	tr := &Tree{}
	orig_keys := []string{"zz", "dd", "yy", "cc", "xx", "bb", "ww", "aa"}
	expected := []string{"aa", "bb", "cc", "dd", "ww", "xx", "yy", "zz"}

	for _, s := range orig_keys {
		tr.Insert(s)
	}
	returned_keys := tr.Keys()
	if !testStringsEq(returned_keys, expected) {
		t.Errorf("Got: %q", returned_keys)
	}
}

// Add all refs to the slice, in DFS order
func get_all_refs(t *Tree) []*ref {
	all_refs := make([]*ref, 0, t.length*2-1)
	recurse_refs(t.root, all_refs)
	return all_refs
}

func recurse_refs(p *ref, refs []*ref) {
	refs = append(refs, p)
	if p.node != nil {
		recurse_refs(&p.node.child[0], refs)
		recurse_refs(&p.node.child[1], refs)
	}
}

func TestCopy(t *testing.T) {
	orig_tree := &Tree{}
	orig_keys := []string{"bb", "dd", "aa", "cc"}

	expected := []string{"aa", "bb", "cc", "dd"}

	for _, s := range orig_keys {
		orig_tree.Insert(s)
	}
	new_tree := orig_tree.Copy()

	// Do they have the same "length" value?
	if orig_tree.Len() != new_tree.Len() {
		t.Errorf("New tree has %d keys", new_tree.Len())
	}

	// Are all the keys the same?
	new_keys := new_tree.Keys()
	if !testStringsEq(new_keys, expected) {
		t.Errorf("Got: %q", new_keys)
	}

	// Check the number of refs within the tree
	orig_refs := get_all_refs(orig_tree)
	new_refs := get_all_refs(new_tree)
	if len(orig_refs) != len(new_refs) {
		t.Errorf("Num orig refs: %d, Num new refs: %d", len(orig_refs), len(new_refs))
	}

	// Check each ref
	for i := 0; i < len(orig_refs); i++ {
		// Ensure that all refs were copied (different pointers)
		if orig_refs[i] == new_refs[i] {
			t.Errorf("Ref #%d is the same", i)
		}
		// Ensure that the old and new refs have the same key
		if orig_refs[i].string != new_refs[i].string {
			t.Errorf("Ref #%d orig string=%s new string=%s", i, orig_refs[i].string,
				new_refs[i].string)
		}
		// The refs should both either have nodes or not.
		if orig_refs[i].node != nil && new_refs[i].node == nil {
			t.Errorf("Ref %d orig has node, but new does not", i)
		} else if orig_refs[i].node == nil && new_refs[i].node != nil {
			t.Errorf("Ref %d orig has no node, but new does", i)
		}
		// If there is a node, ensure each node has the same data
		if orig_refs[i].node != nil {
			if orig_refs[i].node.off != new_refs[i].node.off {
				t.Errorf("Ref %d orig off=%d, new off=%d", i, orig_refs[i].node.off,
					new_refs[i].node.off)
			}
			if orig_refs[i].node.bit != new_refs[i].node.bit {
				t.Errorf("Ref %d orig bit=%d, new bit=%d", i, orig_refs[i].node.bit,
					new_refs[i].node.bit)
			}
		}
	}
}
