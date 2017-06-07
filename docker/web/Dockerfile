FROM golang:1.8.3

ENV GOPATH /go

RUN go get github.com/vetheslav/Social-sournament-service-API
COPY . /$GOPATH/
RUN go test $GOPATH/*.go
RUN go build -o http_api $GOPATH/*.go


ENTRYPOINT /go/bin/http_api
EXPOSE 1323