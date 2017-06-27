# Image Server
An image server toolkit in Go (Golang)

[![GoDoc](https://godoc.org/github.com/pierrre/imageserver?status.svg)](https://godoc.org/github.com/pierrre/imageserver)
[![Build Status](https://travis-ci.org/pierrre/imageserver.svg?branch=master)](https://travis-ci.org/pierrre/imageserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/pierrre/imageserver)](https://goreportcard.com/report/github.com/pierrre/imageserver)

## Features
- HTTP server
- Resize ([GIFT](https://github.com/disintegration/gift), [nfnt resize](https://github.com/nfnt/resize), [Graphicsmagick](http://www.graphicsmagick.org/))
- Rotate
- Crop
- Convert (JPEG, GIF (animated), PNG , BMP, TIFF, ...)
- Cache ([groupcache](https://github.com/golang/groupcache), [Redis](https://github.com/garyburd/redigo), [Memcache](https://github.com/bradfitz/gomemcache), in memory)
- Gamma correction
- Fully modular

## Examples
- [Simple](https://github.com/pierrre/imageserver/blob/master/examples/simple/simple.go)
- [Advanced](https://github.com/pierrre/imageserver/blob/master/examples/advanced/advanced.go)
- [Cache](https://github.com/pierrre/imageserver/blob/master/examples/cache/cache.go)
- [Groupcache](https://github.com/pierrre/imageserver/blob/master/examples/groupcache/groupcache.go)
- [HTTP Source](https://github.com/pierrre/imageserver/blob/master/examples/httpsource/httpsource.go)
- [Mandelbrot](https://github.com/pierrre/mandelbrot/blob/master/examples/httpserver/httpserver.go) ([see interactive demo](https://mandelbrot.pierredurand.fr)) <img src="https://mandelbrot.pierredurand.fr/i?x=0&y=0&z=0" width="32" height="32" />

## Demos
These demos use the "advanced" example.

*Click the images to see the URL parameters.*

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
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200" /></a></td>
		</tr>
		<tr>
			<td><code>height=200</code><br />(preserve aspect ratio)</td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?height=200"><img src="https://imageserver.pierredurand.fr/large.jpg?height=200" /></a></td>
		</tr>
		<tr>
			<td><code>width=200&height=200</code><br />(change aspect ratio)</td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200&height=200"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&height=200" /></a></td>
		</tr>
		<tr>
			<td><code>width=200&height=200&mode=fit</code><br />(fit in 200x200)</td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200&height=200&mode=fit"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&height=200&mode=fit" /></a></td>
		</tr>
		<tr>
			<td><code>width=200&height=200&mode=fill</code><br />(fill 200x200 and crop)</td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200&height=200&mode=fill"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&height=200&mode=fill" /></a></td>
		</tr>
	</tbody>
</table>

### Rotate
<table>
	<thead>
		<tr>
			<th>Options</th>
			<th>Result</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td><code>rotation=90</code><br />(counterclockwise)</td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?rotation=90&width=200"><img src="https://imageserver.pierredurand.fr/large.jpg?rotation=90&width=200" /></a></td>
		</tr>
		<tr>
			<td><code>rotation=45&background=ffaa88</code><br />(background)</td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?rotation=45&background=ffaa88&width=200"><img src="https://imageserver.pierredurand.fr/large.jpg?rotation=45&background=ffaa88&width=200" /></a></td>
		</tr>
	</tbody>
</table>

### Crop
Format: `min_x,min_y|max_x,max_y`

<table>
	<thead>
		<tr>
			<th>Options</th>
			<th>Result</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td><code>crop=556,111|2156,1711</code></td>
			<td><a href="https://imageserver.pierredurand.fr/huge.jpg?crop=556,111|2156,1711&width=200"><img src="https://imageserver.pierredurand.fr/huge.jpg?crop=556,111|2156,1711&width=200" /></a></td>
		</tr>
		<tr>
			<td><code>crop=956,511|1756,1311</code></td>
			<td><a href="https://imageserver.pierredurand.fr/huge.jpg?crop=956,511|1756,1311&width=200"><img src="https://imageserver.pierredurand.fr/huge.jpg?crop=956,511|1756,1311&width=200" /></a></td>
		</tr>
		<tr>
			<td><code>crop=1252,799|1460,1022</code></td>
			<td><a href="https://imageserver.pierredurand.fr/huge.jpg?crop=1252,799|1460,1022"><img src="https://imageserver.pierredurand.fr/huge.jpg?crop=1252,799|1460,1022" /></a></td>
		</tr>
	</tbody>
</table>


### Animated GIF
<table>
	<thead>
		<tr>
			<th>Original</th>
			<th>Resized</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<th><a href="https://imageserver.pierredurand.fr/animated.gif">Animated</a></th>
			<th><a href="https://imageserver.pierredurand.fr/animated.gif?width=300&height=300&mode=fill"><img src="https://imageserver.pierredurand.fr/animated.gif?width=300&height=300&mode=fill" /></a></th>
		</tr>
		<tr>
			<th><a href="https://imageserver.pierredurand.fr/spaceship.gif">Spaceship</a></th>
			<th><a href="https://imageserver.pierredurand.fr/spaceship.gif?width=300"><img src="https://imageserver.pierredurand.fr/spaceship.gif?width=300" /></a></th>
		</tr>
	</tbody>
</table>

### Gamma correction ([more info](http://www.ericbrasseur.org/gamma.html))
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
			<td><a href="https://imageserver.pierredurand.fr/dalai_gamma.jpg">Dalai Gamma</a></td>
			<td><a href="https://imageserver.pierredurand.fr/dalai_gamma.jpg?width=200&gamma_correction=false"><img src="https://imageserver.pierredurand.fr/dalai_gamma.jpg?width=200&gamma_correction=false" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/dalai_gamma.jpg?width=200&gamma_correction=true"><img src="https://imageserver.pierredurand.fr/dalai_gamma.jpg?width=200&gamma_correction=true" /></a></td>
		</tr>
		<tr>
			<td><a href="https://imageserver.pierredurand.fr/gray_squares.jpg">Gray squares</a></td>
			<td><a href="https://imageserver.pierredurand.fr/gray_squares.jpg?width=200&gamma_correction=false"><img src="https://imageserver.pierredurand.fr/gray_squares.jpg?width=200&gamma_correction=false" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/gray_squares.jpg?width=200&gamma_correction=true"><img src="https://imageserver.pierredurand.fr/gray_squares.jpg?width=200&gamma_correction=true" /></a></td>
		</tr>
		<tr>
			<td><a href="https://imageserver.pierredurand.fr/random.png">Random</a></td>
			<td><a href="https://imageserver.pierredurand.fr/random.png?width=200&gamma_correction=false"><img src="https://imageserver.pierredurand.fr/random.png?width=200&gamma_correction=false" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/random.png?width=200&gamma_correction=true"><img src="https://imageserver.pierredurand.fr/random.png?width=200&gamma_correction=true" /></a></td>
		</tr>
		<tr>
			<td><a href="https://imageserver.pierredurand.fr/rings.png">Rings</a></td>
			<td><a href="https://imageserver.pierredurand.fr/rings.png?width=200&gamma_correction=false"><img src="https://imageserver.pierredurand.fr/rings.png?width=200&gamma_correction=false" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/rings.png?width=200&gamma_correction=true"><img src="https://imageserver.pierredurand.fr/rings.png?width=200&gamma_correction=true" /></a></td>
		</tr>
		<tr>
			<td><a href="https://imageserver.pierredurand.fr/rules_sucks.png">Rules / sucks</a></td>
			<td><a href="https://imageserver.pierredurand.fr/rules_sucks.png?width=200&gamma_correction=false"><img src="https://imageserver.pierredurand.fr/rules_sucks.png?width=200&gamma_correction=false" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/rules_sucks.png?width=200&gamma_correction=true"><img src="https://imageserver.pierredurand.fr/rules_sucks.png?width=200&gamma_correction=true" /></a></td>
		</tr>
	</tbody>
</table>

### Resampling
<table>
	<thead>
		<tr>
			<th>Resampling</th>
			<th><a href="https://imageserver.pierredurand.fr/rings.png">Rings</a></th>
			<th><a href="https://imageserver.pierredurand.fr/large.jpg">Large</a></th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>Nearest neighbor</td>
			<td><a href="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=nearest_neighbor&format=png"><img src="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=nearest_neighbor&format=png" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=nearest_neighbor&format=png"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=nearest_neighbor&format=png" width="400" /></a></td>
		</tr>
		<tr>
			<td>Box</td>
			<td><a href="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=box&format=png"><img src="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=box&format=png" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=box&format=png"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=box&format=png" width="400" /></a></td>
		</tr>
		<tr>
			<td>Linear</td>
			<td><a href="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=linear&format=png"><img src="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=linear&format=png" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=linear&format=png"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=linear&format=png" width="400" /></a></td>
		</tr>
		<tr>
			<td>Cubic</td>
			<td><a href="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=cubic&format=png"><img src="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=cubic&format=png" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=cubic&format=png"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=cubic&format=png" width="400" /></a></td>
		</tr>
		<tr>
			<td>Lanczos</td>
			<td><a href="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=lanczos&format=png"><img src="https://imageserver.pierredurand.fr/rings.png?width=200&resampling=lanczos&format=png" /></a></td>
			<td><a href="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=lanczos&format=png"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&resampling=lanczos&format=png" width="400" /></a></td>
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
			<td><a href="https://imageserver.pierredurand.fr/medium.jpg?width=200&quality=5"><img src="https://imageserver.pierredurand.fr/medium.jpg?width=200&quality=5" width="400" /></a></td>
		</tr>
		<tr>
			<td>10%</td>
			<td><a href="https://imageserver.pierredurand.fr/medium.jpg?width=200&quality=10"><img src="https://imageserver.pierredurand.fr/medium.jpg?width=200&quality=10" width="400" /></a></td>
		</tr>
		<tr>
			<td>50%</td>
			<td><a href="https://imageserver.pierredurand.fr/medium.jpg?width=200&quality=50"><img src="https://imageserver.pierredurand.fr/medium.jpg?width=200&quality=50" width="400" /></a></td>
		</tr>
	</tbody>
</table>

### Convert (JPEG to GIF)
<a href="https://imageserver.pierredurand.fr/large.jpg?width=200&format=gif"><img src="https://imageserver.pierredurand.fr/large.jpg?width=200&format=gif" width="600" /></a>

## Backward compatibility
There is no backward compatibility promises.
If you want to use it, vendor it.
It's always OK to change things to make things better.
The API is not 100% correct in the first commit.
