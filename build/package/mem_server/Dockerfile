FROM golang:1.15.15-alpine3.14 as build
WORKDIR /go/src/github.com/ghazlabs/hex-mathrush

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

WORKDIR /go/src/github.com/ghazlabs/hex-mathrush/cmd/mem_server
RUN go build -o server

FROM alpine:3.14
RUN apk add ca-certificates tzdata
COPY --from=build /go/src/github.com/ghazlabs/hex-mathrush/cmd/mem_server/server .
COPY ./build/package/mem_server/questions.json .
COPY ./build/package/shared/web ./web

EXPOSE 9190

ENTRYPOINT [ "./server" ]