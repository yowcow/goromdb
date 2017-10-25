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

COPY ./ /go/src/github.com/yowcow/goromdb
WORKDIR /go/src/github.com/yowcow/goromdb

RUN make clean && make

CMD ["./romdb"]
