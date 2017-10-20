FROM golang:1.9.1-alpine

RUN set -eux; \
    apk add --no-cache \
        db-dev \
        g++ \
        gcc \
        git \
        make \
        putty
RUN go get github.com/golang/dep/cmd/dep

COPY ./ /go/src/github.com/yowcow/go-romdb
WORKDIR /go/src/github.com/yowcow/go-romdb

RUN make clean && make

CMD ["./romdb"]
