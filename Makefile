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
	go test ./...

benchmark: build
	go test -bench . ./...

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

run: build
	./bin/wav2png --height 256 --width 1024 --padding 4 -out ./runtime ./samples/entangled.wav
	open ./runtime/entangled.png

entangled: build
	./bin/wav2png --height 256 --width 1024 --padding 4 -out ./runtime ./samples/entangled.wav
	open ./runtime/entangled.png

# 16-bit signed integer PCM WAV file
noise: build
	./bin/wav2png --debug --height 390 --width 641 --padding 0 -out ./runtime ./samples/noise.wav
	open ./runtime/noise.png

# float32 WAV format file
noise-float32: build
	./bin/wav2png --debug --height 390 --width 641 --padding 0 -out ./runtime ./samples/noise-float32.wav
	open ./runtime/noise-float32.png

