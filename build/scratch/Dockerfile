ARG KUBICO_VERSION=latest

FROM ilyes512/kubico-builder:$KUBICO_VERSION as src

FROM scratch

COPY --from=src /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=src /src/build/kubico /
COPY --from=src /etc/passwd /etc/passwd

CMD ["/kubico"]
