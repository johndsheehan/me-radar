.DEFAULT_GOAL := build

TAG = $(shell date +"%Y%m%d-%H%M%S")

build: clean
	go get -u ;  go build

clean:
	rm -rf met-eireann-archive

docker:
	docker build  --force-rm  -t met-eireann-archive:$(TAG)  -f Dockerfile .

install: build
	cp ./met-eireann-archive  /usr/local/bin

uninstall:
	rm /usr/local/bin/met-eireann-archive
