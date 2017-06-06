FROM golang:1.8

ENV GOPATH /go

RUN go get github.com/vetheslav/Social-sournament-service-API
RUN go build -o http_api *.go

COPY . $GOPATH/src/github.com/spf13/hugo

ENTRYPOINT /go/bin/http_api
EXPOSE 80