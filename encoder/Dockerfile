FROM golang:alpine AS builder
COPY dca.go /go/src/github.com/mitaka8/encode/dca.go
COPY go.mod /go/src/github.com/mitaka8/encode/go.mod
COPY go.sum /go/src/github.com/mitaka8/encode/go.sum
WORKDIR /go/src/github.com/mitaka8/encode
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /tmp/dca .

FROM alpine:latest
RUN apk add --no-cache ffmpeg minio-client rabbitmq-c-utils curl jq
COPY --from=builder /tmp/dca /usr/local/bin/dca
ADD encode.sh /usr/local/bin/encode.sh

CMD ["encode.sh"]
