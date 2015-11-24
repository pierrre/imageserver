# Image Server
An image server toolkit in Go (Golang)

[![GoDoc](https://godoc.org/github.com/pierrre/imageserver?status.svg)](https://godoc.org/github.com/pierrre/imageserver)
[![Build Status](https://travis-ci.org/pierrre/imageserver.svg)](https://travis-ci.org/pierrre/imageserver)
[![codecov.io](https://codecov.io/github/pierrre/imageserver/coverage.svg)](https://codecov.io/github/pierrre/imageserver)

## Features
- HTTP server
- Resize / convert ([GIFT](https://github.com/disintegration/gift), [nfnt resize](https://github.com/nfnt/resize), [Graphicsmagick](http://www.graphicsmagick.org/))
- Cache ([groupcache](https://github.com/golang/groupcache), [Redis](https://github.com/garyburd/redigo), [Memcache](https://github.com/bradfitz/gomemcache), in memory)
- Gamma correction
- Fully modular

## Examples
- [Simple](https://github.com/pierrre/imageserver/blob/master/examples/simple/simple.go)
- [Advanced](https://github.com/pierrre/imageserver/blob/master/examples/advanced/advanced.go)

## Demo
*Click the images to see the URL parameters.*

### Normal
<a href="http://fuckingfrogs.fr:8080/small.jpg"><img src="http://fuckingfrogs.fr:8080/small.jpg" /></a>

### Resize
<table>
    <thead>
        <tr>
            <th>Options</th>
            <th>Result</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><code>width=200</code><br />(preserve aspect ratio)</td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200" /></a></td>
        </tr>
        <tr>
            <td><code>height=200</code><br />(preserve aspect ratio)</td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?height=200"><img src="http://fuckingfrogs.fr:8080/large.jpg?height=200" /></a></td>
        </tr>
        <tr>
            <td><code>width=200&height=200</code><br />(change aspect ratio)</td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&height=200"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&height=200" /></a></td>
        </tr>
        <tr>
            <td><code>width=200&height=200&mode=fit</code><br />(fit in 200x200)</td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&height=200&mode=fit"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&height=200&mode=fit" /></a></td>
        </tr>
        <tr>
            <td><code>width=200&height=200&mode=fill</code><br />(fill 200x200 and crop)</td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&height=200&mode=fill"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&height=200&mode=fill" /></a></td>
        </tr>
    </tbody>
</table>

### Gamma correction
<table>
    <thead>
        <tr>
            <th>Original</th>
            <th>Disabled</th>
            <th>Enabled</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><a href="http://fuckingfrogs.fr:8080/dalai_gamma.jpg">Dalai Gamma</a></td>
            <td><a href="http://fuckingfrogs.fr:8080/dalai_gamma.jpg?width=200&gamma_correction=false"><img src="http://fuckingfrogs.fr:8080/dalai_gamma.jpg?width=200&gamma_correction=false" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/dalai_gamma.jpg?width=200&gamma_correction=true"><img src="http://fuckingfrogs.fr:8080/dalai_gamma.jpg?width=200&gamma_correction=true" /></a></td>
        </tr>
        <tr>
            <td><a href="http://fuckingfrogs.fr:8080/gray_squares.jpg">Gray squares</a></td>
            <td><a href="http://fuckingfrogs.fr:8080/gray_squares.jpg?width=200&gamma_correction=false"><img src="http://fuckingfrogs.fr:8080/gray_squares.jpg?width=200&gamma_correction=false" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/gray_squares.jpg?width=200&gamma_correction=true"><img src="http://fuckingfrogs.fr:8080/gray_squares.jpg?width=200&gamma_correction=true" /></a></td>
        </tr>
        <tr>
            <td><a href="http://fuckingfrogs.fr:8080/random.png">Random</a></td>
            <td><a href="http://fuckingfrogs.fr:8080/random.png?width=200&gamma_correction=false"><img src="http://fuckingfrogs.fr:8080/random.png?width=200&gamma_correction=false" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/random.png?width=200&gamma_correction=true"><img src="http://fuckingfrogs.fr:8080/random.png?width=200&gamma_correction=true" /></a></td>
        </tr>
        <tr>
            <td><a href="http://fuckingfrogs.fr:8080/random.png">Rings</a></td>
            <td><a href="http://fuckingfrogs.fr:8080/rings.png?width=200&gamma_correction=false"><img src="http://fuckingfrogs.fr:8080/rings.png?width=200&gamma_correction=false" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/rings.png?width=200&gamma_correction=true"><img src="http://fuckingfrogs.fr:8080/rings.png?width=200&gamma_correction=true" /></a></td>
        </tr>
    </tbody>
</table>

### Resampling
<table>
    <thead>
        <tr>
            <th>Resampling</th>
            <th><a href="http://fuckingfrogs.fr:8080/rings.png">Rings</a></th>
            <th><a href="http://fuckingfrogs.fr:8080/large.jpg">Large</a></th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>Nearest neighbor</td>
            <td><a href="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=nearest_neighbor&format=png"><img src="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=nearest_neighbor&format=png" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=nearest_neighbor&format=png"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=nearest_neighbor&format=png" width="400" /></a></td>
        </tr>
        <tr>
            <td>Box</td>
            <td><a href="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=box&format=png"><img src="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=box&format=png" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=box&format=png"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=box&format=png" width="400" /></a></td>
        </tr>
        <tr>
            <td>Linear</td>
            <td><a href="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=linear&format=png"><img src="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=linear&format=png" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=linear&format=png"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=linear&format=png" width="400" /></a></td>
        </tr>
        <tr>
            <td>Cubic</td>
            <td><a href="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=cubic&format=png"><img src="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=cubic&format=png" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=cubic&format=png"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=cubic&format=png" width="400" /></a></td>
        </tr>
        <tr>
            <td>Lanczos</td>
            <td><a href="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=lanczos&format=png"><img src="http://fuckingfrogs.fr:8080/rings.png?width=200&resampling=lanczos&format=png" /></a></td>
            <td><a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=lanczos&format=png"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&resampling=lanczos&format=png" width="400" /></a></td>
        </tr>
    </tbody>
</table>

### Quality
<table>
    <thead>
        <tr>
            <th>JPEG quality</th>
            <th>Result</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>5%</td>
            <td><a href="http://fuckingfrogs.fr:8080/medium.jpg?width=200&quality=5"><img src="http://fuckingfrogs.fr:8080/medium.jpg?width=200&quality=5" width="400" /></a></td>
        </tr>
        <tr>
            <td>10%</td>
            <td><a href="http://fuckingfrogs.fr:8080/medium.jpg?width=200&quality=10"><img src="http://fuckingfrogs.fr:8080/medium.jpg?width=200&quality=10" width="400" /></a></td>
        </tr>
        <tr>
            <td>50%</td>
            <td><a href="http://fuckingfrogs.fr:8080/medium.jpg?width=200&quality=50"><img src="http://fuckingfrogs.fr:8080/medium.jpg?width=200&quality=50" width="400" /></a></td>
        </tr>
    </tbody>
</table>

### Convert (JPEG to GIF)
<a href="http://fuckingfrogs.fr:8080/large.jpg?width=200&format=gif"><img src="http://fuckingfrogs.fr:8080/large.jpg?width=200&format=gif" width="600" /></a>
