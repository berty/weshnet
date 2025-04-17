# builder
FROM golang:1.19-alpine3.16 as builder
MAINTAINER gfanton <8671905+gfanton@users.noreply.github.com>


ARG GIT_REPOS=https://github.com/libp2p/go-libp2p-relay-daemon.git
ARG GIT_TAG=v0.3.0

RUN apk add --no-cache git
RUN git clone --depth 1 --branch "${GIT_TAG}" "${GIT_REPOS}" /app

WORKDIR /app

RUN go build -o /go/bin/daemon -v -ldflags="-s -w" -v ./cmd/libp2p-relay-daemon

# runner
FROM alpine:3.16

COPY --from=builder /go/bin/daemon /usr/local/bin

ENTRYPOINT ["/usr/local/bin/daemon"]
