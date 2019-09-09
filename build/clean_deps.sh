#!/usr/bin/env bash

set -e

if [ ! -f "build/clean_deps.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

root=`pwd`
BLS_BUILD=$root/crypto/bls/bls_linux_darwin
root=$root/life/resolver

# Build softfloat
SF_BUILD=$root/softfloat/build
CMAKE_GEN="Unix Makefiles"
MAKE="make"
if [ "$(uname)" = "Darwin" ]; then
    SF_BUILD=$SF_BUILD/Linux-x86_64-GCC
elif [ `expr substr $(uname -s) 1 5` = "Linux" ]; then
    SF_BUILD=$SF_BUILD/Linux-x86_64-GCC
elif [ `expr substr $(uname -s) 1 10` = "MINGW64_NT" ]; then
    SF_BUILD=$SF_BUILD/Win64-MinGW-w64
    CMAKE_GEN="MinGW Makefiles"
    MAKE="mingw32-make.exe"
else
    echo "not support system $(uname -s)"
    exit 0
fi

# Clean softfloat build files
cd $SF_BUILD
$MAKE clean
cd ..; rm -f libsoftfloatlib.a

cd $root/builtins/
rm -rf build



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
