#!/usr/bin/env bash

set -e

root=`pwd | awk '{split($0, path, "Platon-go"); print path[1]}'`

root=$root/Platon-go/life/resolver

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
