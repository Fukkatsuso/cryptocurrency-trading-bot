# https://hub.docker.com/layers/golang/library/golang/1.14-buster/images/sha256-93dd6ad8909cb192a2a0f113521d4ec08e9dd35d02d33e4ae1eca2db3c269eff
FROM golang:1.14-buster AS base

RUN apt-get update && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

ENV PORT 8080
EXPOSE 8080

WORKDIR /go/src/github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard


FROM base AS dev

RUN go get -v github.com/pilu/fresh

CMD [ "fresh" ]


# ビルドの参考
# https://cloud.google.com/run/docs/quickstarts/build-and-deploy?hl=ja#containerizing
