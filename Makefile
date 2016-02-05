GB=gb

all: build

lint: gofmt

gofmt:
	gofmt -w src

build: lint
	$(GB) build all

deploy:
	GOOS=linux $(GB) build all
	scp -C words users bin/weaselbot-linux-amd64 nohya.net:~/weaselbot/
