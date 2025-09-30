FROM golang:1.25-alpine AS base

LABEL org.opencontainers.image.source="https://github.com/komadiina/spelltext"
LABEL org.opencontainers.image.description="spelltext chatserver dockerized image (golang service)"

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" /server/chat/cmd/chatserver/main.go

FROM alpine:3.22.1
WORKDIR /

COPY --from=base ./main ./chatserver

ENTRYPOINT [ "./chatserver" ]