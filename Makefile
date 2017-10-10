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

docker/build:
	docker build -t $(BINARY) .

docker/run:
	docker run --rm $(BINARY)

docker/rmi:
	docker rmi $(BINARY)

.PHONY: dep test clean realclean docker/build docker/run docker/rmi
