all: test      \
	 benchmark \
     coverage

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -o ./bin ./... 

test: build
	go clean -testcache
	go test -v ./...

benchmark: build
	go test -bench ./...

coverage: build
	go clean -testcache
	go test -cover ./...

vet: test
	go vet ./...

lint: vet
	golint ./...

clean:
	go clean
	rm -rf bin/*

debug: build
	./bin/wav2png --height 390 --width 640 --padding 0 -out ./runtime ./runtime/noise.wav
	open ./runtime/noise.png

run: build
	./bin/wav2png --height 256 --width 1024 --padding 4 -out ./runtime ./runtime/entangled.wav
	open ./runtime/entangled.png

noise: build
	./bin/wav2png --height 390 --width 640 --padding 0 -out ./runtime ./runtime/noise.wav
	open ./runtime/noise.png

entangled: build
	./bin/wav2png --height 256 --width 1024 --padding 4 -out ./runtime ./runtime/entangled.wav
	open ./runtime/entangled.png

