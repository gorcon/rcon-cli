#!/usr/bin/env bash

VERSION="$1"
if [ -z "${VERSION}" ]; then echo "VERSION is not set. Use ./compile.sh 0.0.0" >&2; exit 1; fi

RELEASE_DIR=".tmp/release"

rm -r "${RELEASE_DIR}"
mkdir "${RELEASE_DIR}"
touch "${RELEASE_DIR}/checksum.txt"

make_release() {
    local arch="$1"
    local os="$2"
    local release_name="$3"
    if [ -z "${arch}" ] || [ -z "${os}" ] || [ -z "${release_name}" ]; then echo "args are not set" >&2; return 1; fi

    local ext="$4"

    local dir="${RELEASE_DIR}/${release_name}"

    mkdir -p "${dir}"
    env GOARCH="${arch}" GOOS="${os}" CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=${VERSION}" -o "${dir}/rcon${ext}" ./cmd/gorcon/main.go

    cp LICENSE "${dir}"
    cp README.md "${dir}"
    cp CHANGELOG.md "${dir}"
    cp rcon.yaml "${dir}"

    cd "${RELEASE_DIR}/"
    case "${os}" in
        linux | darwin)
            tar -zcvf "${release_name}.tar.gz" "${release_name}"
            md5sum "${release_name}.tar.gz" >> checksum.txt
            ;;
        windows)
            zip -r "${release_name}.zip" "${release_name}"
            md5sum "${release_name}.zip" >> checksum.txt
            ;;
    esac
    rm -r "${release_name}"
    cd ../../
}

make_release 386 linux "rcon-${VERSION}-i386_linux"
make_release amd64 linux "rcon-${VERSION}-amd64_linux"
make_release 386 windows "rcon-${VERSION}-win32" .exe
make_release amd64 windows "rcon-${VERSION}-win64" .exe
make_release amd64 darwin "rcon-${VERSION}-amd64_darwin"
make_release arm64 darwin "rcon-${VERSION}-arm64_darwin"

env GOARCH="amd64" GOOS="linux" CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=${VERSION}" -o gorcon ./cmd/gorcon/main.go
