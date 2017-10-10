BINARY = romdb

all: dep $(BINARY)

dep:
	dep ensure

test:
	go test ./...

$(BINARY):
	go build -o $@ ./cmd/server

clean:
	rm -rf $(BINARY)

realclean: clean
	rm -rf vendor

docker:
	docker build -t $(BINARY) .

.PHONY: dep test clean realclean docker
