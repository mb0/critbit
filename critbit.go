// Copyright 2013 Martin Schnabel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package critbit provides an example implementation of a critbit tree.
// It is based on Adam Langley's well documented critbit C implementation at https://github.com/agl/critbit
package critbit

// Tree represents a string set.
// This is meant as a reference for custom implementations.
type Tree struct {
	root   *ref
	length int
}

// ref holds either a string key or a node pointer
type ref struct {
	string
	*node
}

type node struct {
	child [2]ref
	// off is the offset of the differing byte
	off int
	// bit contains the single crit bit in the differing byte
	bit byte
}

// dir calculates the direction for the given key
func (n *node) dir(key string) byte {
	if n.off < len(key) && key[n.off]&n.bit != 0 {
		return 1
	}
	return 0
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
	p := *t.root
	for p.node != nil {
		// try next node
		p = p.node.child[p.node.dir(key)]
	}
	// check for membership
	return key == p.string
}

// Insert returns whether the key was inserted into the tree.
// Otherwise the tree already contained the key.
func (t *Tree) Insert(key string) bool {
	// test for empty tree
	if t.root == nil {
		t.root = &ref{key, nil}
		t.length++
		return true
	}
	// walk for best member
	p := *t.root
	for p.node != nil {
		// try next node
		p = p.node.child[p.node.dir(key)]
	}
	// find critical bit
	var off int
	var ch, bit byte
	// find differing byte
	for off = 0; off < len(key); off++ {
		if ch = 0; off < len(p.string) {
			ch = p.string[off]
		}
		if keych := key[off]; ch != keych {
			bit = ch ^ keych
			goto ByteFound
		}
	}
	if off < len(p.string) {
		ch = p.string[off]
		bit = ch
		goto ByteFound
	}
	return false
ByteFound:
	// find differing bit
	bit |= bit >> 1
	bit |= bit >> 2
	bit |= bit >> 4
	bit = bit &^ (bit >> 1)
	var ndir byte
	if ch&bit != 0 {
		ndir++
	}
	// insert new node
	nn := &node{off: off, bit: bit}
	nn.child[1-ndir].string = key
	// walk for best insertion node
	wp := t.root
	for wp.node != nil {
		p = *wp
		if p.node.off > off || p.node.off == off && p.node.bit < bit {
			break
		}
		// try next node
		wp = &p.node.child[p.node.dir(key)]
	}
	nn.child[ndir] = *wp
	wp.node = nn
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
	var dir byte
	var wp *ref
	p := t.root
	for p.node != nil {
		wp = p
		// try next node
		dir = p.node.dir(key)
		p = &p.node.child[dir]
	}
	// check for membership
	if key != p.string {
		return false
	}
	// delete from tree
	t.length--
	if wp == nil {
		t.root = nil
		return true
	}
	*wp = wp.node.child[1-dir]
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
		return iterate(*t.root, handler)
	}
	// walk for best member
	p, top := *t.root, *t.root
	for p.node != nil {
		newtop := p.node.off < len(prefix)
		// try next node
		p = p.node.child[p.node.dir(prefix)]
		if newtop {
			top = p
		}
	}
	if len(p.string) < len(prefix) {
		return true
	}
	for i := 0; i < len(prefix); i++ {
		if p.string[i] != prefix[i] {
			return true
		}
	}
	return iterate(top, handler)
}

// iterate calls the key handler or traverses both node children unless aborted.
func iterate(p ref, h func(string) bool) bool {
	if p.node != nil {
		return iterate(p.node.child[0], h) && iterate(p.node.child[1], h)
	}
	return h(p.string)
}

// Keys returns all keys, as a slice of strings, in sorted order.
func (t *Tree) Keys() []string {
	keys := make([]string, 0, t.length)

	// empty tree?
	if t.root == nil {
		return keys
	}

	// Walk the tree without function recursion
	to_visit := make([]*ref, 1)
	p := t.root
	to_visit[0] = p

	for len(to_visit) > 0 {
		// shift the list to get the first item
		p, to_visit = to_visit[0], to_visit[1:]

		// leaf?
		if p.node == nil {
			keys = append(keys, p.string)
		} else {
			// unshift the children and continue
			to_visit = append([]*ref{&p.node.child[0], &p.node.child[1]},
				to_visit...)
		}
	}
	return keys
}

// Copy makes a new deep copy of the tree
func (t *Tree) Copy() *Tree {

	// Copy a ref, but don't deep-copy into orig.node.child
	shallow_copy := func(orig *ref) *ref {
		var nn *node
		if orig.node != nil {
			nn = &node{orig.child, orig.off, orig.bit}
		}
                // NOTE: if your implementation changes the key to some
                // allocated value, you may want to deep-copy that key
                // here and pass its value here to &ref{}.
		return &ref{orig.string, nn}
	}

	new_tree := &Tree{}
	new_tree.root = shallow_copy(t.root)
	new_tree.length = t.length

	// Walk the tree without function recursion
	to_visit := make([]*ref, 1)
	to_visit[0] = new_tree.root
	var p *ref

	for len(to_visit) > 0 {
		// shift the list to get the first item
		p, to_visit = to_visit[0], to_visit[1:]

		// We only need to worry about internal nodes,
		// not leaf nodes, because we shallow_copy the children
		// of each ref we visit.
		if p.node != nil {
			// Make shallow copies of the children
			left_copy := shallow_copy(&p.node.child[0])
			right_copy := shallow_copy(&p.node.child[1])

			// Fix the parent to point to the copies of the children
			p.node.child[0] = *left_copy
			p.node.child[1] = *right_copy

			// unshift onto stack, and continue
			to_visit = append([]*ref{left_copy, right_copy},
				to_visit...)
		}
	}
	return new_tree
}

// Dump is useful for debugging. It println()'s the entire tree
func (t *Tree) Dump() {
	println("*Tree length=", t.length)
        if t.root.node == nil {
            println(" Root: string=", t.root.string)
        } else {
            t.root.dump("Root:", "")
        }
}

func (r *ref) dump(title string, indent string) {
    if r.node == nil {
        println(indent, title, "String:", r.string)
    } else {
        println(indent, title, "String:", r.string, "Offset:", r.off, "Bit:", r.bit)
    }
    indent += "  "
    if r.node != nil {
        r.child[0].dump("Left", indent)
        r.child[1].dump("Right", indent)
    }
}
