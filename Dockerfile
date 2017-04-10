FROM golang

MAINTAINER Don Omead <donomead@gmail.com>

RUN go get github.com/tools/godep

ADD . /go/src/github.com/donomead/deadmoon

WORKDIR /go/src/github.com/donomead/deadmoon

RUN godep get

ENTRYPOINT /go/bin/deadmoon

EXPOSE 8080
