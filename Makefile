all: build test lint

build: \
	build-example-simple \
	build-example-cache \
	build-example-httpsource \
	build-example-groupcache \
	build-example-advanced

build-example-simple:
	go build -v -i -o build/example-simple ./examples/simple

build-example-cache:
	go build -v -i -o build/example-cache ./examples/cache

build-example-httpsource:
	go build -v -i -o build/example-httpsource ./examples/httpsource

build-example-groupcache:
	go build -v -i -o build/example-groupcache ./examples/groupcache

build-example-advanced:
	go build -v -i -o build/example-advanced ./examples/advanced

test:
	mkdir -p build
	echo 'mode: set' > build/coverage.txt
	go list ./... | xargs -n1 -I{} sh -c\
	 'rm -f build/coverage.tmp && touch build/coverage.tmp &&\
	 go test -v -covermode=set -coverprofile=build/coverage.tmp {} &&\
	 tail -n +2 build/coverage.tmp >> build/coverage.txt'
	rm build/coverage.tmp
	go tool cover -html=build/coverage.txt -o=build/coverage.html

lint:
	go get -v github.com/alecthomas/gometalinter
	gometalinter --install
	gometalinter -E gofmt -D gotype -D vetshadow -D dupl -D goconst -D interfacer -D gas -D gocyclo\
	 --tests --deadline=10m --concurrency=2 --enable-gc ./...

clean:
	rm -rf build

.PHONY: \
	build \
	build-example-simple \
	build-example-cache \
	build-example-httpsource \
	build-example-groupcache \
	build-example-advanced \
	test \
	lint \
	clean
