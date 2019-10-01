ARG KUBICO_VERSION=latest

FROM ilyes512/kubico-builder:$KUBICO_VERSION as src

FROM alpine:3.10.2

COPY --from=src /src/build/kubico /usr/local/bin/kubico

CMD ["/usr/local/bin/kubico"]
