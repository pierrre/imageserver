# Image Server
An image server toolkit in Go (Golang)

[![GoDoc](https://godoc.org/github.com/pierrre/imageserver?status.svg)](https://godoc.org/github.com/pierrre/imageserver)
[![Build Status](https://travis-ci.org/pierrre/imageserver.svg)](https://travis-ci.org/pierrre/imageserver)
[![codecov.io](https://codecov.io/github/pierrre/imageserver/coverage.svg)](https://codecov.io/github/pierrre/imageserver)

## Features
- Http server
- Resize / convert ([nfnt resize](https://github.com/nfnt/resize), [Graphicsmagick](http://www.graphicsmagick.org/))
- Cache ([groupcache](https://github.com/golang/groupcache), [Redis](https://github.com/garyburd/redigo), [Memcache](https://github.com/bradfitz/gomemcache), in memory)
- Gamma correction
- Fully modular

## Examples
- [Simple](https://github.com/pierrre/imageserver/blob/master/examples/simple/simple.go)
- [Advanced](https://github.com/pierrre/imageserver/blob/master/examples/advanced/advanced.go)

## Demo

### Normal
![Normal](http://fuckingfrogs.fr:8080/?source=small.jpg)
```
http://fuckingfrogs.fr:8080/?source=small.jpg
```

### Resize (width=400)
![Resize](http://fuckingfrogs.fr:8080/?source=large.jpg&width=400)
```
http://fuckingfrogs.fr:8080/?source=large.jpg&width=400
```

### Thumbnail (100x100)
![Thumbnail 1](http://fuckingfrogs.fr:8080/?source=small.jpg&width=100&height=100&mode=thumbnail)
![Thumbnail 2](http://fuckingfrogs.fr:8080/?source=medium.jpg&width=100&height=100&mode=thumbnail)
![Thumbnail 3](http://fuckingfrogs.fr:8080/?source=large.jpg&width=100&height=100&mode=thumbnail)
![Thumbnail 4](http://fuckingfrogs.fr:8080/?source=huge.jpg&width=100&height=100&mode=thumbnail)
```
http://fuckingfrogs.fr:8080/?source=small.jpg&width=100&height=100&mode=thumbnail
http://fuckingfrogs.fr:8080/?source=medium.jpg&width=100&height=100&mode=thumbnail
http://fuckingfrogs.fr:8080/?source=large.jpg&width=100&height=100&mode=thumbnail
http://fuckingfrogs.fr:8080/?source=huge.jpg&width=100&height=100&mode=thumbnail
```

### Gamma correction
Original | Disabled | Enabled
----- | ----- | -----
[Dalai Gamma](http://fuckingfrogs.fr:8080/?source=dalai_gamma.jpg) | ![Disabled](http://fuckingfrogs.fr:8080/?source=dalai_gamma.jpg&width=200&gamma_correction=false) | ![Disabled](http://fuckingfrogs.fr:8080/?source=dalai_gamma.jpg&width=200&gamma_correction=true)
[Gray squares](http://fuckingfrogs.fr:8080/?source=gray_squares.jpg) | ![Disabled](http://fuckingfrogs.fr:8080/?source=gray_squares.jpg&width=200&gamma_correction=false) | ![Disabled](http://fuckingfrogs.fr:8080/?source=gray_squares.jpg&width=200&gamma_correction=true)
[Random](http://fuckingfrogs.fr:8080/?source=random.png) | ![Disabled](http://fuckingfrogs.fr:8080/?source=random.png&width=200&gamma_correction=false) | ![Disabled](http://fuckingfrogs.fr:8080/?source=random.png&width=200&gamma_correction=true)

### Quality (JPEG 5%)
![Convert](http://fuckingfrogs.fr:8080/?source=large.jpg&width=400&quality=5)
```
http://fuckingfrogs.fr:8080/?source=large.jpg&width=400&quality=5
```

### Convert (JPEG to GIF)
![Convert](http://fuckingfrogs.fr:8080/?source=large.jpg&width=400&format=gif)
```
http://fuckingfrogs.fr:8080/?source=large.jpg&width=400&format=gif
```
