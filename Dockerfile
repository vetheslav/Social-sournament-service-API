FROM golang:1.8

ENV GOPATH /go

RUN go get github.com/vetheslav/Social-sournament-service-API
COPY . $GOPATH/src/github.com/vetheslav/Social-sournament-service-API
RUN ls $GOPATH/src/github.com/vetheslav/Social-sournament-service-API
RUN ls
RUN go build -o http_api $GOPATH/*.go


ENTRYPOINT /go/bin/http_api
EXPOSE 80