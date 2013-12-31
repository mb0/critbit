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

	BenchmarkMap		   20000	     90154 ns/op	    5033 B/op	       9 allocs/op
	BenchmarkTree		   10000	    112353 ns/op	    4928 B/op	      77 allocs/op

The second one assummes you want distinct and sorted keys. The tree is already sorted. The map uses the sort package.

	BenchmarkMapSort	   20000	     78751 ns/op	    6367 B/op	      12 allocs/op
	BenchmarkTreeSort	   50000	     68961 ns/op	    6232 B/op	      78 allocs/op

In the second case the tree approach is now faster, that demonstrates the strength of the critbit tree.
This shows the applicability of the critbit tree whenever you need a sorted set or map.
Not to mention the ability to do fast prefix iteration which would be even more difficult with the standard map.
