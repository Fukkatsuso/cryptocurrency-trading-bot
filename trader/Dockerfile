FROM golang:1.17-buster AS base

RUN apt-get update && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

ENV PORT 8080
EXPOSE 8080

WORKDIR /go/src/github.com/Fukkatsuso/cryptocurrency-trading-bot/trader

COPY go.* ./
RUN go mod download


FROM base AS dev

RUN go install github.com/pilu/fresh@latest

CMD [ "fresh" ]


# 参考: https://cloud.google.com/run/docs/quickstarts/build-and-deploy?hl=ja#containerizing
FROM base AS builder

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -x -o /go/bin/server


FROM debian:buster-slim AS release

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/server /go/bin/server

CMD [ "/go/bin/server" ]
