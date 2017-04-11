FROM golang

MAINTAINER Gabriel Khaselev <gkhasel1@gmail.com>

RUN go get github.com/tools/godep

ADD . /go/src/github.com/gkhasel1/cricket

WORKDIR /go/src/github.com/gkhasel1/cricket

RUN godep get

ENTRYPOINT /go/bin/cricket

EXPOSE 8080
