RUN go get
RUN go build -o http_api *.go

ENTRYPOINT /go/bin/http_api
EXPOSE 80