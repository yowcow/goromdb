all: dep

dep:
	dep ensure -update

test:
	go test ./...

.PHONY: dep test
