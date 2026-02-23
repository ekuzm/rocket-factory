FROM alpine:3.23

WORKDIR /migrate

RUN wget https://github.com/pressly/goose/releases/download/v3.27.0/goose_linux_x86_64 && \
    mv ./goose_linux_x86_64 ./goose && \
    chmod +x ./goose && \
    addgroup -S migrategroup && \
    adduser -S migrateuser -G migrategroup && \
    chown migrateuser:migrategroup ./goose

USER migrateuser

COPY ./order/migrations/*.sql ./migrations/

ENTRYPOINT [ "./goose", "up" ]
