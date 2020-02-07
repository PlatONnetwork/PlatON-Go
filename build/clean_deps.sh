#!/usr/bin/env bash

set -e

if [ ! -f "build/clean_deps.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

root=`pwd`
BLS_BUILD=$root/crypto/bls/bls_linux_darwin


CMAKE_GEN="Unix Makefiles"
MAKE="make"

if [ `expr substr $(uname -s) 1 5` != "MINGW" ]; then
    if [ -d $BLS_BUILD/src/bls ]; then
        cd $BLS_BUILD/src/bls
        $MAKE clean
    fi
    if [ -d $BLS_BUILD/src/mcl ]; then
        cd $BLS_BUILD/src/mcl
        $MAKE clean
    fi
    if [ -d "$BLS_BUILD/include" ]; then
    rm -rf $BLS_BUILD/include/*
    fi
    if [ -d "$BLS_BUILD/lib" ]; then
    rm -rf $BLS_BUILD/lib/*
    fi
fi