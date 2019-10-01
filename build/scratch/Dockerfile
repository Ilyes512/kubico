ARG KUBICO_VERSION=0.1.0

FROM ilyes512/kubico-builder:$KUBICO_VERSION-builder as src

FROM scratch

COPY --from=src /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=src /src/build/kubico /
COPY --from=src /etc/passwd /etc/passwd

CMD ["/kubico"]
