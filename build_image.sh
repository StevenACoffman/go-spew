#!/usr/bin/env bash
set -x

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o main main.go

# brew install upx

# upx is slow to run, so disable during dev cycle
upx --brute main

# Building GO binaries
# + you get static executables unless there is a need to make them dynamic;
# + you force dynamic executables with ``-ldflags -linkmode=external`;
# + `CGO_ENABLED=0` will disable cgo-support, making a static binary more likely;
# + `-tags netgo` will disable netcgo-support, making a static binary more likely.
# + `-a` force rebuilding of packages that are already up-to-date.
# + `-installsuffix cgo` That causes the packages to be built in ${GOROOT}/pkg/<arch>_cgo
# By default it is empty, but custom builds that need to keep their outputs separate can set InstallSuffix to do so.

# Mac build:
# GOOS=darwin -e GOARCH=amd64

# Linux build:
# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Alternatively, for absolutley minimal sizes:
# CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-w -s' -o main .

# You will get the smallest binaries if you compile with -ldflags '-w -s'.
#
# The -w turns off DWARF debugging information: you will not be able to use gdb
#  on the binary to look at specific functions or set breakpoints or
#  get stack traces, because all the metadata gdb needs will not be included.
# You will also not be able to use other tools that depend on the information,
# like pprof profiling.
#
# The -s turns off generation of the Go symbol table:
# you will not be able to use 'go tool nm' to list the symbols in the binary.
# Strip -s is like passing -s to -ldflags but it doesn't strip quite as much.
# 'Go tool nm' might still work after 'strip -s'. I am not completely sure.

# Building with musl on alpine:
# CC=/usr/local/musl/bin/musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' main.go

IMAGE_NAME=go-spew
REPOSITORY_NAMESPACE=${1:-stevenacoffman}
REPOSITORY="${REPOSITORY_NAMESPACE}/${IMAGE_NAME}"
# A docker tag name must be valid ASCII and may contain lowercase and uppercase letters,
# digits, underscores, periods and dashes.
# A docker tag name may not start with a period or a dash and may contain a maximum of 128 characters.
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD | sed 's/[^\w.-]+//g')
GIT_REVISION=$(git rev-parse HEAD)
BUILD_TIME=$(date +'%s')

docker build \
    -t "${REPOSITORY}:SHA1-${GIT_REVISION}" \
    -t "${REPOSITORY}:latest" \
    -t "${REPOSITORY}:GIT_BRANCH_${GIT_BRANCH}" \
    --build-arg "GIT_BRANCH=${GIT_BRANCH}" \
    --build-arg "BUILD_TIME=${BUILD_TIME}" \
    --build-arg "GIT_COMMIT=${GIT_REVISION}" \
    -f Dockerfile .

docker push "${REPOSITORY}"
