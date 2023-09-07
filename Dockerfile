FROM golang:1.20.8-bullseye AS builder

WORKDIR /src

RUN apt-get update \
    && apt-get install --assume-yes --no-install-recommends \
        ca-certificates \
        upx

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ARG GOARCH=amd64

RUN go generate \
    && CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -a -tags netgo -ldflags '-extldflags "-static" -s -w' -o ./kubico

ARG COMPRESS=false

RUN echo compress \
    && echo $COMPRESS

RUN if [ $COMPRESS = 1 ]; then upx --brute --no-progress ./kubico; fi

FROM debian:11.6-slim AS debian

COPY --from=builder /src/kubico /usr/local/bin

CMD ["kubico"]

FROM scratch AS scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /src/kubico /
COPY --from=builder /etc/passwd /etc/passwd

CMD ["/kubico"]
