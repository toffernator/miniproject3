# Starts a client in a docker container
FROM golang:1.17-alpine AS base

COPY . /src
WORKDIR /src

RUN go mod download
RUN go build -o "bin/out" main.go
ENTRYPOINT ["bin/out"]
