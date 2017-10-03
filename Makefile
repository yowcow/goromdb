all: dep

dep:
	dep ensure

test:
	go test ./...

build:
	go build -o romdb ./cmd/server

docker/%:
	make -C docker $(notdir $@)

.PHONY: dep test docker/%
