GB=gb

all: build

lint: gofmt

gofmt:
	gofmt -w src

build: lint
	$(GB) build all

test: lint
	$(GB) test all
