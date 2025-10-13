FROM golang:1.25-alpine AS base

LABEL org.opencontainers.image.source=https://github.com/komadiina/spelltext
LABEL org.opencontainers.image.description="spelltext chatserver dockerized image"

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

# TODO: break this down into smaller pieces
COPY . .
RUN go build /server/chat/cmd/chatserver/main.go

FROM alpine:3.22.1
WORKDIR /

ENV CONFIG_FILE=./config.yml
COPY --from=base ./main ./chatserver
COPY --from=base ./config.yml ./config.yml

ENTRYPOINT [ "./chatserver" ]