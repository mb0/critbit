// Copyright 2013 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package critbit provides an example implementation of a critbit tree.
// It is based on Adam Langley's well documented critbit C implementation at https://github.com/agl/critbit
package critbit

// Tree represents a string set.
// This is meant as a reference for custom implementations.
type Tree struct {
	// root is either a string key or an internal node
	root   interface{}
	length int
}

type node struct {
	// children are either string keys or internal nodes
	child [2]interface{}
	// at contains both the byte offset and the inverted critical bit: (nbyte<<8) | ^byte(critbit)
	// "aa" and "ab" would have the first critical bit in the second byte in the second bit.
	// with omitting the first six zero bytes and adding spaces for readability: 0000 0010 1111 1101
	at uint64
}

// Len returns the number of keys in the tree.
func (t *Tree) Len() int {
	return t.length
}

// Contains returns whether the tree contains the key.
func (t *Tree) Contains(key string) bool {
	// test for empty tree
	if t.root == nil {
		return false
	}
	// walk for best member
	p := t.root
	n, ok := p.(*node)
	for ok {
		// calculate direction
		var c byte
		if b := n.at >> 8; b < uint64(len(key)) {
			c = key[b]
		}
		dir := (1 + uint32(c|byte(n.at))) >> 8
		// try next node
		p = n.child[dir]
		n, ok = p.(*node)
	}
	// check for membership
	return key == p.(string)
}

// Insert returns whether the key was inserted into the tree.
// Otherwise the tree already contained the key.
func (t *Tree) Insert(key string) bool {
	// test for empty tree
	if t.root == nil {
		t.root = key
		t.length++
		return true
	}
	// walk for best member
	p := t.root
	n, ok := p.(*node)
	for ok {
		// calculate direction
		var c byte
		if b := n.at >> 8; b < uint64(len(key)) {
			c = key[b]
		}
		dir := (1 + uint32(c|byte(n.at))) >> 8
		// try next node
		p = n.child[dir]
		n, ok = p.(*node)
	}
	s := p.(string)
	// find critical bit
	var nbyte int
	var bits byte
	var c byte
	// find differing byte
	for nbyte = 0; nbyte < len(key); nbyte++ {
		if c = 0; nbyte < len(s) {
			c = s[nbyte]
		}
		if kbits := key[nbyte]; c != kbits {
			bits = c ^ kbits
			goto ByteFound
		}
	}
	if nbyte < len(s) {
		c = s[nbyte]
		bits = c
		goto ByteFound
	}
	return false
ByteFound:
	// find differing bit
	bits |= bits >> 1
	bits |= bits >> 2
	bits |= bits >> 4
	bits = ^(bits &^ (bits >> 1))
	ndir := (1 + uint32(c|bits)) >> 8
	// insert new node
	nn := &node{at: (uint64(nbyte) << 8) | uint64(bits)}
	nn.child[1-ndir] = key
	// walk for best insertion node
	wp := &t.root
	for {
		p = *wp
		n, ok := p.(*node)
		if !ok {
			break
		}
		b := n.at >> 8
		if b > uint64(nbyte) {
			break
		}
		if b == uint64(nbyte) && byte(n.at) > bits {
			break
		}
		// calculate direction
		var c byte
		if b < uint64(len(key)) {
			c = key[b]
		}
		dir := (1 + uint32(c|byte(n.at))) >> 8
		// try next node
		wp = &n.child[dir]
	}
	nn.child[ndir] = *wp
	*wp = nn
	t.length++
	return true
}

// Delete returns whether the key was deleted from the tree.
// Otherwise the tree does not contain the key.
func (t *Tree) Delete(key string) bool {
	// test for empty tree
	if t.root == nil {
		return false
	}
	// walk for best member
	var dir uint32
	var wp *interface{}
	p := &t.root
	n, ok := t.root.(*node)
	for ok {
		wp = p
		// calculate direction
		var c byte
		if b := n.at >> 8; b < uint64(len(key)) {
			c = key[b]
		}
		dir = (1 + uint32(c|byte(n.at))) >> 8
		// try next node
		p = &n.child[dir]
		n, ok = (*p).(*node)
	}
	// check for membership
	if key != (*p).(string) {
		return false
	}
	// delete from tree
	t.length--
	if wp == nil {
		t.root = nil
		return true
	}
	n = (*wp).(*node)
	*wp = n.child[1-dir]
	return true
}

// Iterate calls the handler for all keys with the given prefix.
// It returns whether all prefixed keys were iterated.
// The handler can continue the process by returning true or abort with false.
func (t *Tree) Iterate(prefix string, handler func(key string) bool) bool {
	// test empty tree
	if t.root == nil {
		return true
	}
	// shortcut for empty prefix
	if prefix == "" {
		return iterate(t.root, handler)
	}
	// walk for best member
	p, top := t.root, t.root
	n, ok := p.(*node)
	for ok {
		// calculate direction
		var c byte
		b := n.at >> 8
		if b < uint64(len(prefix)) {
			c = prefix[b]
		}
		dir := (1 + uint32(c|byte(n.at))) >> 8
		// try next node
		p = n.child[dir]
		if b < uint64(len(prefix)) {
			top = p
		}
		n, ok = p.(*node)
	}
	s := p.(string)
	if len(s) < len(prefix) {
		return true
	}
	for i := 0; i < len(prefix) && i < len(s); i++ {
		if s[i] != prefix[i] {
			return true
		}
	}
	return iterate(top, handler)
}

// iterate calls the key handler or traverses both node children unless aborted.
func iterate(p interface{}, h func(string) bool) bool {
	if n, ok := p.(*node); ok {
		return iterate(n.child[0], h) && iterate(n.child[1], h)
	}
	return h(p.(string))
}
