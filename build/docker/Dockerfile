# Stage 1 - Build the code.
FROM golang:1.19.3-alpine as builder

RUN apk --no-cache add --update gcc musl-dev

ARG VERSION=docker

WORKDIR /build
COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=1 go build -ldflags "-s -w -X main.ServiceVersion=${VERSION}" -o gorcon ./cmd/gorcon/main.go

# Stage 2 - Create image.
FROM alpine as runner
MAINTAINER Pavel Korotkiy <paul.korotkiy@gmail.com>

COPY LICENSE /
COPY README.md /
COPY CHANGELOG.md /
COPY rcon.yaml /
COPY --from=builder /build/gorcon /rcon

CMD ["/rcon"]
