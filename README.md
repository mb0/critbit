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

	BenchmarkMap		   20000	     96068 ns/op	    5023 B/op	       9 allocs/op
	BenchmarkTree		   10000	    118832 ns/op	    4863 B/op	      76 allocs/op

The second one assummes you want distinct and sorted keys. The tree is already sorted. The map uses the sort package.

	BenchmarkMapSort	   20000	     85485 ns/op	    6360 B/op	      12 allocs/op
	BenchmarkTreeSort	   20000	     75553 ns/op	    6167 B/op	      77 allocs/op


In the second case the tree approach is now faster, that demonstrates the strength of the critbit tree.
This shows the applicability of the critbit tree whenever you need a sorted set or map.
Not to mention the ability to do fast prefix iteration which would be even more difficult with the standard map.

This is a big improvement over the last commit. The important change was to use a condition for the critbit test
instead of an arithmetic solution. Although the stand-alone benchmark does favor the latter.

	BenchmarkBitCondInv	  500000	      3066 ns/op
	// if dir := 0; byteToCheck&^invertedBit != 0 { dir++ }
	BenchmarkBitCond	  500000	      2923 ns/op
	// if dir := 0; byteToCheck&bit != 0 { dir++ }
	BenchmarkBitArith	 1000000	      2814 ns/op
	// dir := (1 + uint32(byteToCheck|invertedBit)) >> 8

Another minor improvement was that whenever the key is shorter than the node to test,
the `byteToCheck` must be 0 and thus we already know the direction.
