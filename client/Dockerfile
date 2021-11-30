# Starts a client in a docker container
FROM golang:1.17-alpine AS base

COPY go.mod go.sum /src/
WORKDIR /src/
RUN go mod download

COPY . /src/

RUN go build -o "bin/out" github.com/toffernator/miniproject3/client
EXPOSE 50000
ENTRYPOINT ["bin/out"]
