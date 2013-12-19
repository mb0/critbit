critbit
=======

An example implementation of a critbit tree in Go

It is based on Adam Langley's well documented critbit C implementation at https://github.com/agl/critbit

This version is not intended to be used as is (unless you need a simple string set).
Instead just copy and modify it according to your needs.
You can change it to use other types of keys or use it as a map.

This package is BSD licensed, Copyright (c) 2013 Martin Schnabel

Benchmark
---------
Benchmark against map[string]struct{} on an i5-2400 with Go 1.2.

The first benchmark checks simple key insertion and checks.

	BenchmarkMap		   10000	    101667 ns/op	    5064 B/op	      10 allocs/op
	BenchmarkTree		   10000	    128732 ns/op	    5192 B/op	      81 allocs/op

The second one assummes you want distinct and sorted keys. The tree is already sorted. The map uses the sort package.

	BenchmarkMapSort	   20000	     92145 ns/op	    6403 B/op	      12 allocs/op
	BenchmarkTreeSort	   20000	     84976 ns/op	    6494 B/op	      82 allocs/op


In the second case the tree approach is now faster, that demonstrates the strength of the critbit tree.
This shows the applicability of the critbit tree whenever you need a sorted set or map.
Not to mention the ability to do fast prefix iteration which would be even more difficult with the standard map.

This is a big improvement over the last commit. The important change was to use a condition for the critbit test
instead of an arithmetic solution. Although the stand-alone benchmark does favor the latter.

Now the tree uses:

	//BenchmarkBitCond	  500000	      3181 ns/op
	var dir byte
	if byteToCheck&^invertedBit != 0 {
		dir++
	}

instead of (even if a benchmark indicates otherwise):

	//BenchmarkBitArith	 1000000	      2864 ns/op
	dir := (1 + uint32(byteToCheck|invertedBit)) >> 8

Another minor improvement was that whenever the key is shorter than the node to test,
the `byteToCheck` must be 0 and thus we already know the direction.
