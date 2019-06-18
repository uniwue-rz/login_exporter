#!/bin/bash
#
# GoLang cross-compile snippet for Go 1.6+ based loosely on Dave Chaney's cross-compile script:
# http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go
#
# To use:
#
#   $ cd ~/path-to/my-awesome-project
#   $ go-build-all
#
# Features:
#
#   * Cross-compiles to multiple machine types and architectures.
#   * Uses the current directory name as the output name...
#     * ...unless you supply an source file: $ go-build-all main.go
#   * Windows binaries are named .exe.
#   * ARM v5, v6, v7 and v8 (arm64) support
#
# ARM Support:
#
# You must read https://github.com/golang/go/wiki/GoArm for the specifics of running
# Linux/BSD-style kernels and what kernel modules are needed for the target platform.
# While not needed for cross-compilation of this script, you're users will need to ensure
# the correct modules are included.
#
# Requirements:
#
#   * GoLang 1.6+ (for mips and ppc), 1.5 for non-mips/ppc.
#   * CD to directory of the binary you are compiling. $PWD is used here.
#
# For 1.4 and earlier, see http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go
#

# This PLATFORMS list is refreshed after every major Go release.
# Though more platforms may be supported (freebsd/386), they have been removed
# from the standard ports/downloads and therefore removed from this list.
#
PLATFORMS="darwin/amd64" # amd64 only as of go1.5
PLATFORMS="$PLATFORMS windows/amd64" # arm compilation not available for Windows
PLATFORMS="$PLATFORMS linux/amd64 linux/386"

# ARMBUILDS lists the platforms that are currently supported.  From this list
# we generate the following architectures:
#
#   ARM64 (aka ARMv8) <- only supported on linux and darwin builds (go1.6)
#   ARMv7
#   ARMv6
#   ARMv5
#
# Some words of caution from the master:
#
#   @dfc: you'll have to use gomobile to build for darwin/arm64 [and others]
#   @dfc: that target expects that you're bulding for a mobile phone
#   @dfc: iphone 5 and below, ARMv7, iphone 3 and below ARMv6, iphone 5s and above arm64
#
PLATFORMS_ARM="linux"

##############################################################
# Shouldn't really need to modify anything below this line.  #
##############################################################

type setopt >/dev/null 2>&1

SCRIPT_NAME=`basename "$0"`
FAILURES=""
SOURCE_FILE=`echo $@ | sed 's/\.go//'`
CURRENT_DIRECTORY=${PWD##*/}
OUTPUT=${SOURCE_FILE:-$CURRENT_DIRECTORY} # if no src file given, use current dir name

for PLATFORM in $PLATFORMS; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  BIN_FILENAME="${OUTPUT}_${GOOS}_${GOARCH}"
  if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
  if [[ "${GOOS}" == "linux" ]]; then BIN_FILENAME="${BIN_FILENAME}.run"; fi
  if [[ "${GOOS}" == "darwin" ]]; then BIN_FILENAME="${BIN_FILENAME}.run"; fi
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILENAME} $@"
  echo "${CMD}"
  eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
done

# ARM builds
#if [[ $PLATFORMS_ARM == *"linux"* ]]; then
#  CMD="GOOS=linux GOARCH=arm64 go build -o ${OUTPUT}-linux-arm64 $@"
#  echo "${CMD}"
#  eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
#fi
for GOOS in $PLATFORMS_ARM; do
  GOARCH="arm"
  # build for each ARM version
    BIN_FILENAME="${OUTPUT}_${GOOS}_${GOARCH}.run"
    CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILENAME} $@"
    echo "${CMD}"
    eval "${CMD}" || FAILURES="${FAILURES} ${GOOS}/${GOARCH}"
done

# eval errors
if [[ "${FAILURES}" != "" ]]; then
  echo ""
  echo "${SCRIPT_NAME} failed on: ${FAILURES}"
  exit 1
fi

sha256sum *.run *.exe > hashes_sha256.txt