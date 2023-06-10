all: build

build: bin
	go build \
	    -ldflags="-s -w" \
	    -o bin/gosmc \
	    cmd/main.go

bin:
	mkdir -p bin
