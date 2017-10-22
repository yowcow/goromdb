BINARY = romdb
CIDFILE = .romdb-cid

ifeq ($(shell uname -s),Darwin)
MD5 = md5 -r
else
MD5 = md5sum
endif

DB_FILES = sample-data.json sample-bdb.db sample-memcachedb-bdb.db
DB_DIR = data/store
DB_PATHS = $(addprefix $(DB_DIR)/,$(DB_FILES))
MD5_PATHS = $(foreach path,$(DB_PATHS),$(path).md5)

all: dep $(DB_DIR) $(DB_PATHS) $(MD5_PATHS) $(BINARY)

dep:
	dep ensure

test:
	go test ./...

$(DB_DIR):
	mkdir -p $@

$(DB_DIR)/%.md5: $(DB_DIR)/%
	$(MD5) $< > $@

$(DB_DIR)/sample-data.json: data/sample-data.json
	cp $< $@

$(DB_DIR)/sample-bdb.db: data/sample-data.json
	go run ./cmd/sample-data/bdb/bdb.go -input-from $< -output-to $@

$(DB_DIR)/sample-memcachedb-bdb.db: data/sample-data.json
	go run ./cmd/sample-data/memcachedb-bdb/memcachedb-bdb.go -input-from $< -output-to $@

bench:
	go test -bench .

$(BINARY):
	go build -o $@ ./cmd/server

clean:
	rm -rf $(BINARY) $(DB_PATHS) $(MD5_PATHS)

realclean: clean
	rm -rf vendor

docker/build:
	docker build -t $(BINARY) .

docker/run:
	-docker run \
		--rm \
		-v `pwd`:/go/src/github.com/yowcow/go-romdb \
		--cidfile=$(CIDFILE) \
		-it $(BINARY) sh
	rm -f $(CIDFILE)

docker/exec:
	test -f $(CIDFILE) && docker exec -it `cat $(CIDFILE)` sh

docker/rmi:
	docker rmi $(BINARY)

.PHONY: dep test bench clean realclean docker/build docker/run docker/exec docker/rmi
