FROM flyway/flyway:11 AS builder

LABEL org.opencontainers.image.source=https://github.com/komadiina/spelltext
LABEL org.opencontainers.image.description="spelltext chatserver flyway migrations image (baked-in sql scripts)"

WORKDIR /flyway
COPY ./sql/ /flyway/sql/

ENTRYPOINT [ "flyway" ]
CMD ["-baselineOnMigrate=true", "migrate"]