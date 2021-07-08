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
	mkdir -p dist/windows/wav2png
	mkdir -p dist/darwin/wav2png
	mkdir -p dist/linux/wav2png
	env GOOS=linux   GOARCH=amd64  go build -o dist/linux/wav2png   ./...
	env GOOS=darwin  GOARCH=amd64  go build -o dist/darwin/wav2png  ./...
	env GOOS=windows GOARCH=amd64  go build -o dist/windows/wav2png ./...

release: build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist/linux  --exclude=".DS_Store" -cvzf dist/wav2png_$(DIST)-linux.tar.gz  wav2png
	tar --directory=dist/darwin --exclude=".DS_Store" -cvzf dist/wav2png_$(DIST)-darwin.tar.gz wav2png
	cd dist/windows; zip --recurse-paths ../wav2png_$(DIST)-windows.zip wav2png

debug: build
	./bin/wav2png version

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

# wav2mp4
wav2mp4: build
	rm -f ./runtime/frames/*
	rm -f ./runtime/chirp.mp4
	./bin/wav2mp4 --debug --mix L+R --width 640 --padding 8 --out ./runtime/chirp.mp4 --window 1s --fps 30 --cursor red:right ./samples/chirp.wav
	cd ./runtime/frames; \
	ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y out.mp4; \
	ffmpeg -i out.mp4 -i ../../samples/chirp.wav -c:v copy -c:a aac -y ../chirp.mp4
	open ./runtime/chirp.mp4
# 	open ./runtime/frames/frame-00001.png
