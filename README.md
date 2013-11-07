# Image Server
An image server written in Go (Golang)

## Features
- Http server
- Resize / convert / process (Graphicsmagick)
- Cache (Redis, Memcache, in memory)
- Fully modular

## Demo
```
http://fuckingfrogs.fr:8080/?source=https://www.google.com/images/srpr/logo6w.png
```
![Normal](http://fuckingfrogs.fr:8080/?source=https://www.google.com/images/srpr/logo6w.png)

```
http://fuckingfrogs.fr:8080/?source=https://www.google.com/images/srpr/logo6w.png&width=400&format=jpeg&quality=50
```
![Resized jpeg low quality](http://fuckingfrogs.fr:8080/?source=https://www.google.com/images/srpr/logo6w.png&width=400&format=jpeg&quality=50)

```
http://fuckingfrogs.fr:8080/?source=https://www.google.com/images/srpr/logo6w.png&width=200&height=200&fill=1&extent=1
```
![Crop](http://fuckingfrogs.fr:8080/?source=https://www.google.com/images/srpr/logo6w.png&width=200&height=200&fill=1&extent=1)

```
http://fuckingfrogs.fr:8080/?source=https://www.google.com/images/srpr/logo6w.png&width=200&height=200&extent=1&background=000000
```
![Extent background](http://fuckingfrogs.fr:8080/?source=https://www.google.com/images/srpr/logo6w.png&width=200&height=200&extent=1&background=000000)

## Status
[![Build Status](https://travis-ci.org/pierrre/imageserver.png?branch=master)](https://travis-ci.org/pierrre/imageserver)

## Usage / Build
You have to compile/configure your own image server.

See example: https://github.com/pierrre/imageserver/blob/master/example/main.go

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
- Go 1.2:
	- Use crypto/sha256.Sum256() new function
