BINARY = goromdb
WORKDIR = /go/src/github.com/yowcow/$(BINARY)

OS = alpine debian-jessie debian-stretch

all: $(addprefix docker-build/,$(OS))

docker-build/%:
	docker build -t $(BINARY):$* -f Dockerfile.$* .

build: $(addprefix go-build/,$(OS))

go-build/%:
	docker run --rm \
		-v `pwd`/../:$(WORKDIR) \
		-w $(WORKDIR) \
		$(BINARY):$* \
		sh -c 'mkdir -p docker/build && go build -o docker/build/$(BINARY).$*'

.PHONY: all build docker-build/% go-build/%
