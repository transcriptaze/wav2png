DIST ?= development

.DEFAULT_GOAL := buil

.PHONY: sass

all: test      \
	 benchmark \
     coverage

clean:

format: 
	find html/javascript -name "*.js" -exec npx eslint --fix {} +

sass: 
	find sass -name "*.scss" | entr sass --no-source-map sass/themes:html/css

build: format
	sass --no-source-map sass/themes:html/css

test: build

benchmark: build

coverage: build

vet: 

lint:

build-all: test vet lint

release: build-all
	rm -rf dist/html
	mkdir -p dist/html/css
	mkdir -p dist/html/fonts
	mkdir -p dist/html/images
	mkdir -p dist/html/javascript
	cp    html/index.html  dist/html
	cp    html/favicon.ico dist/html
	cp -r html/css         dist/html
	cp -r html/fonts       dist/html
	cp -r html/images      dist/html

debug: build
	python3 httpd.py --directory ./html

run: build
	python3 httpd.py

