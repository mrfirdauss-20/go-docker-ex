FROM golang:1.17.7-alpine3.15 as build
WORKDIR /go/src/github.com/ghazlabs/hex-mathrush

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

WORKDIR /go/src/github.com/ghazlabs/hex-mathrush/cmd/redis_server
RUN go build -o server

FROM alpine:3.15
RUN apk add ca-certificates tzdata
COPY --from=build /go/src/github.com/ghazlabs/hex-mathrush/cmd/redis_server/server .
COPY ./build/package/shared/web ./web

EXPOSE 9190

ENTRYPOINT [ "./server" ]