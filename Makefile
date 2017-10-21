build: \
	build/example-simple \
	build/example-cache \
	build/example-httpsource \
	build/example-groupcache \
	build/example-advanced

build/example-simple:
	go build -v -i -o build/example-simple ./examples/simple

build/example-cache:
	go build -v -i -o build/example-cache ./examples/cache

build/example-httpsource:
	go build -v -i -o build/example-httpsource ./examples/httpsource

build/example-groupcache:
	go build -v -i -o build/example-groupcache ./examples/groupcache

build/example-advanced:
	go build -v -i -o build/example-advanced ./examples/advanced

test:
	go test -v ./...

lint:
	go get -v -u github.com/alecthomas/gometalinter
	gometalinter --install
	GOGC=800 gometalinter --enable-all -D dupl -D lll -D gas -D goconst -D gotype -D interfacer -D safesql -D test -D testify -D vetshadow\
	 --tests --warn-unmatched-nolint --deadline=10m --concurrency=4 --enable-gc ./...

clean:
	rm -rf build

.PHONY: build test lint clean
