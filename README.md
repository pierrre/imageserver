# Image Server
An image server written in Go (Golang)

## Features
- Http server
- Resize / convert / process (Graphicsmagick)
- Cache (Redis, Memcache, in memory)
- Fully modular

## Demo
```
Normal
http://fuckingfrogs.fr:8080/?source=https://raw.github.com/pierrre/imageserver/master/testdata/small.jpg
```
![Normal](http://fuckingfrogs.fr:8080/?source=https://raw.github.com/pierrre/imageserver/master/testdata/small.jpg)

```
Resize animated gif
http://fuckingfrogs.fr:8080/?source=https://raw.github.com/pierrre/imageserver/master/testdata/animated.gif&width=300&height=300
```
![Resize animated gif](http://fuckingfrogs.fr:8080/?source=https://raw.github.com/pierrre/imageserver/master/testdata/animated.gif&width=300&height=300)

```
Resize and crop
http://fuckingfrogs.fr:8080/?source=https://raw.github.com/pierrre/imageserver/master/testdata/medium.jpg&width=200&height=200&extent=1&fill=1
```
![Resize and crop](http://fuckingfrogs.fr:8080/?source=https://raw.github.com/pierrre/imageserver/master/testdata/medium.jpg&width=200&height=200&extent=1&fill=1)

```
Resize jpeg low quality
http://fuckingfrogs.fr:8080/?source=https://raw.github.com/pierrre/imageserver/master/testdata/large.jpg&width=400&format=jpeg&quality=50
```
![Resize jpeg low quality](http://fuckingfrogs.fr:8080/?source=https://raw.github.com/pierrre/imageserver/master/testdata/large.jpg&width=400&format=jpeg&quality=50)

## Resize benchmark
```
go test -bench=. -benchtime=10s ./processor/graphicsmagick
testing: warning: no tests to run
PASS
BenchmarkResizeSmallWorker1-8	    5000	   8188875 ns/op
BenchmarkResizeMediumWorker1-8	    1000	  23143921 ns/op
BenchmarkResizeLargeWorker1-8	     500	  52679041 ns/op
BenchmarkResizeHugeWorker1-8	      50	 389060935 ns/op
BenchmarkResizeAnimatedWorker1-8	      50	 311136430 ns/op
BenchmarkResizeSmallWorker2-8	    5000	   5786665 ns/op
BenchmarkResizeMediumWorker2-8	    1000	  16159015 ns/op
BenchmarkResizeLargeWorker2-8	     500	  36070204 ns/op
BenchmarkResizeHugeWorker2-8	      50	 247750150 ns/op
BenchmarkResizeAnimatedWorker2-8	      50	 221887473 ns/op
BenchmarkResizeSmallWorker4-8	    5000	   5674571 ns/op
BenchmarkResizeMediumWorker4-8	    2000	  13427690 ns/op
BenchmarkResizeLargeWorker4-8	    1000	  29846192 ns/op
BenchmarkResizeHugeWorker4-8	      50	 217280287 ns/op
BenchmarkResizeAnimatedWorker4-8	      50	 243820932 ns/op
BenchmarkResizeSmallWorker8-8	    5000	   4770737 ns/op
BenchmarkResizeMediumWorker8-8	    2000	  12349114 ns/op
BenchmarkResizeLargeWorker8-8	    1000	  27195685 ns/op
BenchmarkResizeHugeWorker8-8	      50	 214468274 ns/op
BenchmarkResizeAnimatedWorker8-8	      50	 200378466 ns/op
ok  	github.com/pierrre/imageserver/processor/graphicsmagick	443.701s
```

## Status
[![Build Status](https://travis-ci.org/pierrre/imageserver.png?branch=master)](https://travis-ci.org/pierrre/imageserver)

## Usage / Build
You have to compile/configure your own image server.

See examples: 
- [Simple](https://github.com/pierrre/imageserver/blob/master/examples/simple/simple.go)
- [Advanced](https://github.com/pierrre/imageserver/blob/master/examples/advanced/advanced.go)

## Documentation
http://godoc.org/github.com/pierrre/imageserver

## Help
- Twitter: @pierredurand87
- Github issue

## TODO
- more tests
- source provider
    - dispatch (uri scheme)
    - limit concurrent
    - timeout
- processor:
    - chain
	- native / imagemagick
- regroup requests?
- thread count problem with system calls http://misfra.me/post/52148362774/callback-magic-with-go
