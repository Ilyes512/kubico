FROM node:18.18.1-bullseye AS node

WORKDIR /src/assets

COPY ./assets/package*.json .

RUN npm ci

COPY ./assets .
COPY ./templates ../templates

RUN npm run prod

FROM golang:1.21.2-bookworm AS builder

WORKDIR /src

RUN apt-get update \
    && apt-get install --assume-yes --no-install-recommends \
        ca-certificates

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

COPY --from=node /src/dist dist

ARG GOARCH=amd64

RUN go generate \
    && CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -a -tags netgo -ldflags '-extldflags "-static" -s -w' -o ./kubico

FROM debian:12.1-slim AS debian

COPY --from=builder /src/kubico /usr/local/bin

CMD ["kubico"]

FROM scratch AS scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /src/kubico /
COPY --from=builder /etc/passwd /etc/passwd

CMD ["/kubico"]
