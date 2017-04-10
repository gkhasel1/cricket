FROM golang

MAINTAINER Gabriel Khaselev <gkhasel1@gmail.com>

RUN go get github.com/tools/godep

ADD . /go/src/github.com/gkhasel1/deadmoon

WORKDIR /go/src/github.com/gkhasel1/deadmoon

RUN godep get

ENTRYPOINT /go/bin/deadmoon

EXPOSE 8080
