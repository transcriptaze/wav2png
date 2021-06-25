DIST ?= development

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin/*

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

build-all: build test
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	env GOOS=linux   GOARCH=amd64  go build -o dist/$(DIST)/linux   ./...
	env GOOS=darwin  GOARCH=amd64  go build -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64  go build -o dist/$(DIST)/windows ./...

release: build-all

debug: build
	./bin/wav2png --debug --height 390 --width 641 --padding 0 -out ./runtime ./samples/noise.wav
	open ./runtime/noise.png

run: build test
	./bin/wav2png --out runtime ./samples/noise.wav
	./bin/wav2png --out runtime ./samples/noise-float32.wav
	./bin/wav2png --out runtime ./samples/acoustic.wav
	./bin/wav2png --out runtime ./samples/chirp.wav
	./bin/wav2png --debug \
	              --settings './runtime/settings.json' \
	              --height   390                       \
	              --width    641                       \
	              --padding  2                         \
	              --palette './runtime/amber.png'      \
	              --fill 'solid:#0000ffff'             \
	              --grid 'rectangular:#800000ff:~32x128:overlay' \
	              --antialias 'soft' \
	              --scale 0.9        \
	              --start 2.5s       \
	              --end   7.5s       \
	              --mix   L+R        \
	              --out   ./runtime  \
	              ./samples/chirp.wav
	open ./runtime/chirp.png

acoustic: build
	./bin/wav2png --height 256 --width 1024 --padding 4 -out ./runtime ./samples/acoustic.wav
	open ./runtime/acoustic.png

# 16-bit signed integer PCM WAV file
noise: build
	./bin/wav2png --debug --height 390 --width 641 --padding 0 -out ./runtime ./samples/noise.wav
	open ./runtime/noise.png

# float32 WAV format file
noise-float32: build
	./bin/wav2png --debug --height 390 --width 641 --padding 0 -out ./runtime ./samples/noise-float32.wav
	open ./runtime/noise-float32.png

# stereo  WAV file
chirp: build
	./bin/wav2png --debug --mix L+R --out ./runtime ./samples/chirp.wav
	open ./runtime/chirp.png
