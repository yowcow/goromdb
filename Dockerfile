FROM golang:1.9

RUN apt-get update \
    && apt-get -yqq install \
        libdb-dev \
        telnet \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /go/src/github.com/yowcow/go-romdb
RUN go get github.com/golang/dep/cmd/dep

COPY ./ /go/src/github.com/yowcow/go-romdb
WORKDIR /go/src/github.com/yowcow/go-romdb

RUN make clean && make

CMD ["./romdb"]
