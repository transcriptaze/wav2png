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

debug: build
	python3 httpd.py --directory ./html

run: build
	python3 httpd.py

