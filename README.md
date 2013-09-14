# Image Server
An image server written in Go (Golang)

## Features
- Resize / convert / process (Graphicsmagick)
- Cache (Redis, Memcache, in memory)
- Fully modular

## Usage / Build
You have to compile/configure your own image server.

See example: https://github.com/pierrre/imageserver/blob/master/example/main.go

## Help
- Twitter: @pierredurand87
- Github issue

## TODO
- documentation
- more tests
- source provider
    - dispatch (uri scheme)
    - limit concurrent
    - timeout
- processor:
    - chain
    - timeout
	- native / imagemagick
- regroup requests?
- thread count problem with system calls http://misfra.me/post/52148362774/callback-magic-with-go