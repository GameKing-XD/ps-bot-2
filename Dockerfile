FROM golang:latest as builder

ADD . /usr/src/ps-bot-2
WORKDIR /usr/src/ps-bot-2
RUN --mount=type=cache,target=/go/pkg/mod \
      --mount=type=bind,source=go.mod,target=go.mod \
      --mount=type=bind,source=go.sum,target=go.sum \
      go mod download -x
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /usr/bin/ps-bot-2 .

FROM debian:stable-slim
WORKDIR /opt/discordbot
RUN apt-get update && apt-get install ca-certificates -y
COPY --from=builder /usr/bin/ps-bot-2 /usr/bin/ps-bot-2
COPY --from=builder /usr/src/ps-bot-2/web /opt/discordbot/web
ADD assets /opt/discordbot/assets

CMD ["/usr/bin/ps-bot-2", "discord"]
