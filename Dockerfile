FROM golang:alpine as golang
WORKDIR /go/src/app
COPY . .
ARG VERSION=docker
RUN CGO_ENABLED=0 go install -ldflags "-s -w -X main.Version=${VERSION}"

FROM scratch
MAINTAINER Pavel Korotkiy <paul.korotkiy@gmail.com>

COPY LICENSE /
COPY README.md /
COPY CHANGELOG.md /
COPY rcon.yaml /
COPY --from=golang /go/bin/rcon-cli /rcon

CMD ["/rcon"]