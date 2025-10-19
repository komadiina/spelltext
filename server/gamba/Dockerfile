FROM golang:1.25-alpine AS base

LABEL org.opencontainers.image.source=https://github.com/komadiina/spelltext
LABEL org.opencontainers.image.description="spelltext gambaserver dockerized image"

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build /server/gamba/cmd/main.go

FROM alpine:3.22.1
WORKDIR /

ENV CONFIG_FILE=./config.yml
ENV VERSION=0.3.0

COPY --from=base ./main ./gambaserver
COPY --from=base ./config.yml ./config.yml

ENTRYPOINT [ "./gambaserver" ]