.DEFAULT_GOAL := build

TAG = $(shell date +"%Y%m%d-%H%M%S")

build: clean
	go get -u ;  go build

clean:
	rm -rf me-radar

docker:
	docker build  --network host  --force-rm  -t me-radar:$(TAG)  -f Dockerfile .

install: build
	cp ./me-radar  /usr/local/bin

uninstall:
	rm /usr/local/bin/me-radar
