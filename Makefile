BINARY = romdb
CIDFILE = .romdb-cid

all: dep $(BINARY)

dep:
	dep ensure

test:
	go test ./...

bench:
	go test -bench .

$(BINARY):
	go build -o $@ ./cmd/server

clean:
	rm -rf $(BINARY)

realclean: clean
	rm -rf vendor

docker/build:
	docker build -t $(BINARY) .

docker/run:
	docker run \
		--rm \
		-v `pwd`:/go/src/github.com/yowcow/go-romdb \
		--cidfile=$(CIDFILE) \
		-it $(BINARY) bash
	rm -f $(CIDFILE)

docker/exec:
	test -f $(CIDFILE) && docker exec -it `cat $(CIDFILE)` bash

docker/rmi:
	docker rmi $(BINARY)

.PHONY: dep test bench clean realclean docker/build docker/run docker/exec docker/rmi
