VERSION ?= v0.0.x
DIST ?= development

.PHONY: clean
.PHONY: copy
.PHONY: build-all
.PHONY: release

all: test      \
     benchmark \
     coverage

clean:
	rm -rf dist/*

update:

build-all:
	cd webgpu && make release

release: 
	@echo "... releasing wav2png_$(VERSION)"
	rm -rf dist/$(DIST)
	mkdir -p dist/$(DIST)
	cp -r  ./webgpu/html   dist/$(DIST)/cloudflare
	cp -r  ./cloudflare/*  dist/$(DIST)/cloudflare
	cd dist/$(DIST)/cloudflare; zip --recurse-paths ../../cloudflare.zip .

