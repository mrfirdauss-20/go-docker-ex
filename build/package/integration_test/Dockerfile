FROM golang:1.15.15-alpine3.14 as build
WORKDIR /go/src/github.com/ghazlabs/hex-mathrush

COPY go.mod go.sum ./
RUN go mod download

COPY ./internal ./internal
