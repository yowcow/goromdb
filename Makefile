BINARY = romdb
CIDFILE = .romdb-cid
DB_FILES = data/sample-bdb.db data/sample-memcachedb-bdb.db

all: dep $(DB_FILES) $(BINARY)

dep:
	dep ensure

test: $(DB_FILES)
	go test ./...

data/sample-bdb.db: data/sample-data.json
	go run ./cmd/sample-data/bdb/bdb.go -input-from $< -output-to $@

data/sample-memcachedb-bdb.db: data/sample-data.json
	go run ./cmd/sample-data/memcachedb-bdb/memcachedb-bdb.go -input-from $< -output-to $@

bench:
	go test -bench .

$(BINARY):
	go build -o $@ ./cmd/server

clean:
	rm -rf $(BINARY) $(DB_FILES)

realclean: clean
	rm -rf vendor

docker/build:
	docker build -t $(BINARY) .

docker/run:
	-docker run \
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
