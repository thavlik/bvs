ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"

#FROM inputoutput/cardano-db-sync:latest as dbsync
#FROM nixos/nix as nix
#RUN apk add --no-cache git
#WORKDIR /
#RUN git clone -b release/11.0.x https://github.com/input-output-hk/cardano-db-sync \
#    && cd cardano-db-sync \
#    && nix-build -A cardano-db-sync -o db-sync-node

FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG} AS builder
WORKDIR /go/src/github.com/thavlik/bvs
COPY components/node components/node
RUN cd components/node/cmd/node \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM ${KINDEST_REPOSITORY}thavlik/cardano-base:${KINDEST_TAG}
#RUN apt-get update \
#    && apt-get install -y postgresql-13 \
#    && apt-get clean \
#    && rm -rf /var/lib/apt/lists/*
#USER postgres
#RUN /etc/init.d/postgresql start \
#    && psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" \
#    && createdb -O docker docker \
#    && echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/13/main/pg_hba.conf \
#    && echo "listen_addresses='*'" >> /etc/postgresql/13/main/postgresql.conf
#USER root
#COPY --from=dbsync /bin/cardano-db-sync /usr/local/bin
#COPY --from=dbsync /lib/libecpg.so.6.11 /lib
#COPY --from=dbsync /lib/libecpg_compat.so.3.11 /lib
#COPY --from=dbsync /lib/libpgcommon.a /lib
#COPY --from=dbsync /lib/libpgport.a /lib
#COPY --from=dbsync /lib/libpqxx-6.4.so /lib
#COPY --from=dbsync /lib/libpqxx.la /lib
#WORKDIR /lib
#RUN ln -s libecpg_compat.so.3.11 libecpg_compat.so.3 \
#    && ln -s libecpg_compat.so.3.11 libecpg_compat.so \
#    && ln -s libecpg.so.6.11 libecpg.so \
#    && ln -s libecpg.so.6.11 libecpg.so.6 \
#    && ln -s libpqxx-6.4.so libpqxx.so
#ENV POSTGRES_BIN_PATH="/usr/lib/postgresql/13/bin/postgres"
#ENV POSTGRES_DB_PATH="/mnt/postgres"
#ENV PATH="/usr/lib/postgresql/13/bin:$PATH"

COPY --from=builder /go/src/github.com/thavlik/bvs/components/node/cmd/node/node /usr/local/bin
COPY components/node/configs /configs
ENV CARDANO_CONFIG_PATH="/configs"
ENV CARDANO_DB_PATH="/mnt/db"
ENV CARDANO_NODE_SOCKET_PATH="/shared/node.socket"
ENV TESTNET_MAGIC="1097911063"

CMD ["node", "server"]

