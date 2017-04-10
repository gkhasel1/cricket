FROM golang

MAINTAINER Don Omead <donomead@gmail.com>

RUN go get github.com/tools/godep

ADD . /go/src/github.com/donomead/cricket

WORKDIR /go/src/github.com/donomead/cricket

RUN godep get

ENTRYPOINT /go/bin/cricket

EXPOSE 8080
