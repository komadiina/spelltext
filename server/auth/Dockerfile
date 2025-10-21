FROM golang:1.25-alpine AS base

LABEL org.opencontainers.image.source=https://github.com/komadiina/spelltext
LABEL org.opencontainers.image.description="spelltext authserver dockerized image"

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build /server/auth/cmd/main.go

FROM alpine:3.22.1
WORKDIR /

ENV CONFIG_FILE=./config.yml
ENV VERSION=0.4.2

COPY --from=base ./main ./authserver
COPY --from=base ./config.yml ./config.yml

ENTRYPOINT [ "./authserver" ]