BINARY = goromdb
WORKDIR = /go/src/github.com/yowcow/$(BINARY)
OS = 1.8-jessie 1.9-alpine 1.9-stretch

all: $(addprefix docker-pull/,$(OS)) $(addprefix docker-build/,$(OS))

docker-build/%:
	docker build -t $(BINARY):$* -f Dockerfile.$* .

docker-pull/%:
	docker pull golang:$*

build: $(addprefix go-build/,$(OS))

go-build/%:
	docker run --rm \
		-v `pwd`/../:$(WORKDIR) \
		-w $(WORKDIR) \
		$(BINARY):$* \
		sh -c 'mkdir -p docker/build && go build -o docker/build/$(BINARY).$*'

clean:
	rm -rf build

.PHONY: all build docker-build/% docker-pull/% go-build/% clean
