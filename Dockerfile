FROM scratch
MAINTAINER Pavel Korotkiy <paul.korotkiy@gmail.com>

COPY LICENSE /
COPY README.md /
COPY CHANGELOG.md /
COPY rcon.yaml /
COPY rcon-cli /rcon

CMD ["/rcon"]
