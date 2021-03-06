FROM golang:1.15.8-alpine3.13 AS builder

ENV GO111MODULE=on

WORKDIR /go/src/proxy
COPY . .
RUN apk update && apk add git && rm -rf /var/cache/apk/* \
    go get github.com/aliyun/aliyun-oss-go-sdk/oss; \
    go get github.com/spf13/viper; \
    CGO_ENABLED=0 GOOS=linux go build -a -o proxy .

FROM alpine:3.13.1
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /opt
COPY --from=builder /go/src/proxy/proxy .
COPY --from=builder /go/src/proxy/config.yaml .
CMD ["./proxy"]