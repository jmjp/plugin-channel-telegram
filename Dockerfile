# Imagem final: binário pre-buildado (Dokploy/VPS faz go build antes).
# Estrutura esperada:
#   ./plugin             (binario compilado localmente)
#   ./plugin.yaml        (manifesto)
#   ./Dockerfile         (este arquivo)
FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 10000 fluxa && adduser -D -G fluxa -u 10000 fluxa

COPY plugin /usr/local/bin/plugin
COPY plugin.yaml /home/fluxa/plugin.yaml

# Diretorio de sockets Unix compartilhado com o core (volume mount).
RUN mkdir -p /plugins/workdir && chown -R fluxa:fluxa /plugins /home/fluxa

USER fluxa
WORKDIR /home/fluxa

# Plugin recebe FLUXA_PLUGIN_SOCKET por env (injetado pelo core supervisor).
ENTRYPOINT ["/usr/local/bin/plugin"]