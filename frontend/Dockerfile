# Starts a frontend in a docker container
FROM golang:1.17-alpine AS base

COPY . /src
WORKDIR /src

RUN go mod download
RUN go build -o "bin/out" main.go
EXPOSE 50000
ENTRYPOINT ["bin/out"]
