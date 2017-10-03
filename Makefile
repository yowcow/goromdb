all: dep

dep:
	dep ensure

test:
	go test ./...

.PHONY: dep test
