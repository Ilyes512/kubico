FROM golang:1.13.1-alpine3.10 AS builder

WORKDIR /src

RUN apk add --no-cache --upgrade \
        git \
        musl-dev \
        gcc \
        libc-dev \
        upx \
    && go get -u github.com/gobuffalo/packr/v2/packr2

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go generate \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-extldflags "-static" -s -w' -o ./build/kubico

ARG COMPRESS=false
ENV COMPRESS=$COMPRESS

RUN if [ $COMPRESS = true ]; then upx --brute --no-progress ./build/kubico; fi
