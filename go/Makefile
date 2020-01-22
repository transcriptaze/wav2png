all: test      \
	 benchmark \
     coverage

format: 
	go fmt ./...

build: format
	go build -o ../bin ./... 

test: build
	go clean -testcache
	go test -v ./...

benchmark: build
	go test -bench ./...

coverage: build
	go clean -testcache
	go test -cover ./...

clean:
	go clean
	rm -rf bin/*

run: build
	../bin/wav2png --height 256 --width 1024 --padding 4 -out ../runtime ../runtime/entangled.wav

