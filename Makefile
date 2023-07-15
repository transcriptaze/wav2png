DIST ?= development

.DEFAULT_GOAL := debug

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

vet: 
	go vet ./...

lint:
	env GOOS=darwin  GOARCH=amd64 staticcheck ./...
	env GOOS=linux   GOARCH=amd64 staticcheck ./...
	env GOOS=windows GOARCH=amd64 staticcheck ./...

build-all: test vet lint
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
	./bin/wav2png --debug --style .styles/default.json --start 2.5s --end 7.5s --out runtime ./samples/noise.wav

version: build
	./bin/wav2png version
	./bin/wav2mp4 version

help: build
	./bin/wav2png help
	./bin/wav2mp4 help

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
	./bin/wav2png --height 480 --width 800 --padding 4  --palette horizon -out ./runtime ./samples/acoustic.wav
	open ./runtime/acoustic.png

# 16-bit signed integer PCM WAV file
noise: build
	./bin/wav2png --debug --height 390 --width 641 --padding 0 -out ./runtime ./samples/noise.wav
	open ./runtime/noise.png

# float32 WAV format file
noise-float32: build
	./bin/wav2png --debug --height 390 --width 641 --padding 0  -out ./runtime ./samples/noise-float32.wav
	open ./runtime/noise-float32.png

# stereo  WAV file
chirp: build
	./bin/wav2png --debug --mix L+R --palette aurora --out ./runtime ./samples/chirp.wav
	open ./runtime/chirp.png

# wav2mp4
wav2mp4: build
	rm -f ./runtime/frames/*
	rm -f ./runtime/chirp.mp4
	./bin/wav2mp4 --debug --mix L+R --width 640 --height 390 --padding 8 -palette ice --out ./runtime/chirp.mp4 \
	              --window 1s --fps 30 --cursor ./samples/cursors/gradient.png:erf \
	              ./samples/chirp.wav
	cd ./runtime/frames; \
	ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y out.mp4; \
	ffmpeg -i out.mp4 -i ../../samples/chirp.wav -c:v copy -c:a aac -y ../chirp.mp4
	open ./runtime/chirp.mp4

clocks:	build
	mkdir -p ./runtime/clocks
	rm -rf ./runtime/clocks/*
	rm -rf ./runtime/clocks/clocks.mp4
	./bin/wav2mp4 --debug --mix L --width 640 --height 390 --padding 8 -palette ice --out ./runtime/clocks/clocks.mp4 \
	              --window 1s --fps 30 --cursor ./samples/cursors/gradient.png:erf \
	              --start 0s --end 10s \
	              ./runtime/clocks.wav
	cd ./runtime/clocks/frames; \
	ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y out.mp4; \
	ffmpeg -i out.mp4 -i ../clocks.wav -c:v copy -c:a aac -y ../clocks.mp4
	open ./runtime/clocks.mp4

rondo:	build
	mkdir -p ./runtime/rondo
	rm -rf ./runtime/rondo/*
	rm -rf ./runtime/rondo/rondo.mp4
	./bin/wav2mp4 --debug --mix L --width 1920 --height 270 --padding 8 --palette ./runtime/palettes/fire.png --grid none \
	              --out ./runtime/rondo/rondo.mp4 \
	              --window 10s --fps 30 \
	              ./runtime/rondo.wav
	cd ./runtime/rondo/frames; \
	ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y out.mp4; \
	ffmpeg -i out.mp4 -i ../../rondo.wav -c:v copy -c:a aac -y ../rondo.mp4
	ffmpeg -i ./runtime/rondo/rondo.mp4 -vf "pad=1920:1080:0:810" -c copy -an rondo.mp4
	open rondo.mp4


inception:	build
	mkdir -p ./runtime/inception
	rm -rf ./runtime/inception/*

	./bin/wav2mp4 --debug --mix L --width 1920 --height 256 --padding 8 --palette ./runtime/palettes/yellow.png --grid none \
	              --window 30s --fps 30 \
	              --out ./runtime/inception/yellow/Y.mp4 ./runtime/time.wav
	cd ./runtime/inception/yellow/frames && ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y ../../yellow.mp4

	./bin/wav2mp4 --debug --mix L --width 1920 --height 256 --padding 8 --palette ./runtime/palettes/red.png --grid none \
	              --window 20s --fps 30 \
	              --out ./runtime/inception/red/R.mp4 ./runtime/time.wav
	cd ./runtime/inception/red/frames && ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y ../../red.mp4

	./bin/wav2mp4 --debug --mix L --width 1920 --height 256 --padding 8 --palette ./runtime/palettes/green.png --grid none \
	              --window 15s --fps 30 \
	              --out ./runtime/inception/green/G.mp4 ./runtime/time.wav
	cd ./runtime/inception/green/frames && ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y ../../green.mp4

	./bin/wav2mp4 --debug --mix L --width 1920 --height 256 --padding 8 --palette ./runtime/palettes/blue.png --grid none \
	              --window 20s --fps 30 \
	              --out ./runtime/inception/blue/B.mp4 ./runtime/time.wav
	cd ./runtime/inception/blue/frames && ffmpeg -framerate 30 -i frame-%05d.png -c:v libx264 -pix_fmt yuv420p -crf 23 -y ../../blue.mp4

	ffmpeg -i ./runtime/inception/yellow.mp4 -vf "pad=1920:1080:0:104" Y.mp4
	ffmpeg -i ./runtime/inception/red.mp4    -vf "pad=1920:1080:0:348" R.mp4
	ffmpeg -i ./runtime/inception/green.mp4  -vf "pad=1920:1080:0:592" G.mp4
	ffmpeg -i ./runtime/inception/blue.mp4   -vf "pad=1920:1080:0:836" B.mp4
	ffmpeg -i Y.mp4 -i R.mp4 -filter_complex "[1:v] colorkey=0x000000:0.1:0.2 [out]; [0:v][out] overlay=0:0" p.mp4
	ffmpeg -i p.mp4 -i G.mp4 -filter_complex "[1:v] colorkey=0x000000:0.1:0.2 [out]; [0:v][out] overlay=0:0" q.mp4
	ffmpeg -i q.mp4 -i B.mp4 -filter_complex "[1:v] colorkey=0x000000:0.1:0.2 [out]; [0:v][out] overlay=0:0" out.mp4
	ffmpeg -i out.mp4 -i inception.wav -c:v copy -c:a aac -y inception.mp4
	open inception.mp4



