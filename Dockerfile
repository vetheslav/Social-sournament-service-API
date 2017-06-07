FROM golang:1.8

ENV GOPATH /go

RUN go get github.com/vetheslav/Social-sournament-service-API
RUN ls
RUN ls src/github.com/vetheslav/Social-sournament-service-API
COPY src/github.com/vetheslav/Social-sournament-service-API /$GOPATH/
RUN ls
RUN go build -o http_api $GOPATH/*.go


ENTRYPOINT /go/bin/http_api
EXPOSE 80