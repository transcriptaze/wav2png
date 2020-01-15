all: test      \
	 benchmark \
     coverage

format: 
	go fmt ./...

build: format
	go build -o bin ./... 

test: build
	go clean -testcache
	go test ./...

benchmark: build
	go test -bench ./...

coverage: build
	go clean -testcache
	go test -cover ./...

clean:
	go clean
	rm -rf bin/*

run: build
	./bin/wav2png -out runtime runtime/entangled.wav

