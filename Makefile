all: test lint

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

.PHONY: test lint clean
