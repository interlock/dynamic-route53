FROM golang:1.13 AS build

WORKDIR /go/src/app
COPY *.go go.mod go.sum /go/src/app/
RUN pwd
RUN go get -d -v ./...
RUN go build -o dynamic-router53 .


FROM ubuntu:20.04

COPY --from=build /go/src/app/dynamic-router53 .
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates && rm -rf /var/lib/apt/lists/*

CMD ./dynamic-router53
