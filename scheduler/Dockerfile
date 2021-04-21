FROM golang:1.14-buster

RUN apt-get update && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

ENV PORT 8080
EXPOSE 8080

WORKDIR /go/src/github.com/Fukkatsuso/cryptocurrency-trading-bot/scheduler
