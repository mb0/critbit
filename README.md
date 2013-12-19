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
Benchmark against map[string]struct{} on an i5-2400.

The first benchmark checks simple key insertion and checks.

	BenchmarkMap		   10000	    103719 ns/op	    5065 B/op	      10 allocs/op
	BenchmarkTree		   10000	    174412 ns/op	    5193 B/op	      81 allocs/op

The second one assummes you want distinct and sorted keys. The tree is already sorted. The map uses the sort package.

	BenchmarkMapSort	   20000	     94981 ns/op	    6401 B/op	      12 allocs/op
	BenchmarkTreeSort	   10000	    112375 ns/op	    6494 B/op	      82 allocs/op

In both cases the map approach is faster, that demonstrates how awesome the go standard library is.
The tree might still be faster than a map if you need a constantly sorted set of keys.