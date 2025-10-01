FROM golang:1.25-alpine AS base

LABEL org.opencontainers.image.source=https://github.com/komadiina/spelltext
LABEL org.opencontainers.image.description="spelltext storeserver dockerized image (golang service)"

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build /server/store/cmd/storeserver/main.go

FROM alpine:3.22.1
WORKDIR /

COPY --from=base ./main ./storeserver

ENTRYPOINT [ "./storeserver" ]