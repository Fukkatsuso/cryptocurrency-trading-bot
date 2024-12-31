# ref: https://github.com/mattn/go-sqlite3/blob/master/_example/simple/Dockerfile
FROM golang:alpine3.20 AS base

ENV CGO_ENABLED=1

ENV PORT=8080
EXPOSE 8080

RUN apk add --no-cache tzdata && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  apk del tzdata

RUN apk add --no-cache \
  gcc \
  musl-dev

WORKDIR /go/src/github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard

COPY dashboard/go.* ./


FROM base AS dev

RUN mkdir -p /var/sqlite
RUN apk add --no-cache sqlite

RUN go mod download

RUN go install github.com/air-verse/air@latest
RUN air init

CMD [ "air" ]


FROM base AS builder

COPY dashboard/ ./

RUN GOOS=linux go build -ldflags='-s -w -extldflags "-static"' -mod=readonly -x -o /go/bin/server


FROM debian:buster-slim AS release

RUN mkdir -p /var/sqlite
COPY ./sqlite/trading-sqlite3.db /var/sqlite/trading-sqlite3.db

COPY --from=builder /go/bin/server /go/bin/server
COPY dashboard/view/ /go/bin/view/

WORKDIR /go/bin
CMD [ "./server" ]
