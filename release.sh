#!/usr/bin/env bash

VER="0.1.0"

rm -fr release
mkdir release
print > release/rcon-$VER-md5.txt

function make_release() {
    local arch="$1"
    local os="$2"
    local release_name="$3"
    if [ -z "${arch}" ] || [ -z "${os}" ] || [ -z "${release_name}" ]; then
       echo "args are not set" >&2
       return 1
    fi
    local ext="$4"

    local dir=release/$release_name

    mkdir -p $dir
    env GOARCH=$arch GOOS=$os go build -ldflags "-s -w" -o $dir/rcon$ext
    upx-ucl --best $dir/rcon$ext -o $dir/rcon-upx$ext

    cp LICENSE $dir
    cp README.md $dir

    cd release/
    case $os in
        linux)
            tar -zcvf $release_name.tar.gz $release_name
            md5sum $release_name.tar.gz >> rcon-$VER-md5.txt
            ;;
        windows)
            zip -r $release_name.zip $release_name
            md5sum $release_name.zip >> rcon-$VER-md5.txt
            ;;
    esac
    rm -r $release_name
    cd ../
}

function source_code() {
    local release_name=rcon-${VER}-src
    local dir=release/$release_name

    mkdir -p $dir

    cp -r vendor $dir
    cp .gitignore $dir
    cp Gopkg.lock $dir
    cp Gopkg.toml $dir
    cp LICENSE $dir
    cp main.go $dir
    cp README.md $dir
    cp release.sh $dir

    cd release/
    tar -zcvf $release_name.tar.gz $release_name
    zip -r $release_name.zip $release_name
    rm -r $release_name
    cd ../
}

make_release 386 linux rcon-$VER-i386_linux
make_release amd64 linux rcon-$VER-amd64_linux
make_release 386 windows rcon-$VER-win32 .exe
make_release amd64 windows rcon-$VER-win64 .exe

source_code
