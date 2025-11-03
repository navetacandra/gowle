#!/usr/bin/sh

mkdir -p build
tinygo=false
compress=false
print=false
clean=false
full=false

for arg in "$@"; do
    if [ "$arg" = "-tinygo" ]; then
        tinygo=true
    elif [ "$arg" = "-compress" ]; then
        compress=true
    elif [ "$arg" = "-print" ]; then
        print=true
    elif [ "$arg" = "-clean" ]; then
        clean=true
    elif [ "$arg" = "-full" ]; then
        full=true
    fi
done

if $tinygo; then
    echo "========= Build using TinyGo ========="
    build_tools="tinygo build -no-debug"
    filename="tinygowle"
else
    echo "========= Build using Go ========="
    build_tools="go build -ldflags=\"-s -w\" -trimpath -gcflags=\"-m\""
    build_tools="go1.17 build -ldflags=\"-s -w\" -trimpath -gcflags=\"-m\""
    filename="gowle"
fi

if $clean; then
    rm build/* -rf
fi

build() {
    local ext=""
    if [ "$1" = "windows" ]; then
        ext=".exe"
    fi
    local env="CGO_ENABLED=0 GOOS=$1 GOARCH=$2"
    local out=$(printf "build/%s_%s-%s%s" "$filename" "$1" "$2" "$ext")
    local cmd="$build_tools -o $out cmd/gowle.go"

    echo "Building $out"
    sh -c "$env $cmd"
}


if $full; then
    platforms="linux/amd64 linux/386 linux/arm64 linux/arm windows/amd64 windows/386 windows/arm64 darwin/amd64 darwin/arm64"
else
    platforms="linux/amd64"
fi
for platform in $platforms; do
    os=$(echo "$platform" | cut -d'/' -f1)
    arch=$(echo "$platform" | cut -d'/' -f2)
    build "$os" "$arch"
done


if $compress; then
    echo "Compressing."
    need_compress=$(find build ! -name "*darwin*" ! -name "*arm*" -name "$filename*")
    strip -v $need_compress
    upx --best --lzma $need_compress
fi

if $print; then
    echo ""
    find build ! -name "*darwin*" -name "$filename*" | xargs -I {} ls -lh "{}"
fi
