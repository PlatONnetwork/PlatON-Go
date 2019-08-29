#!/usr/bin/env bash

if [ ! -f "build/build_deps.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi
root=`pwd`
bls_build=$root/build/
root=$root/life/resolver


if [ "`ls $root/softfloat`" = "" ]; then
    # pull softfloat
    git submodule update --init
fi

$bls_build/build_bls.sh

# Build softfloat
SF_BUILD=$root/softfloat/build
CMAKE_GEN="Unix Makefiles"
MAKE="make"
if [ "$(uname)" = "Darwin" ]; then
    SF_BUILD=$SF_BUILD/Linux-x86_64-GCC
elif [ `expr substr $(uname -s) 1 5` = "Linux" ]; then
    SF_BUILD=$SF_BUILD/Linux-x86_64-GCC
elif [ `expr substr $(uname -s) 1 10` = "MINGW64_NT" ] || [ `expr substr $(uname -s) 1 7` = "MSYS_NT" ]; then
    SF_BUILD=$SF_BUILD/Win64-MinGW-w64
    CMAKE_GEN="MinGW Makefiles"
    MAKE="mingw32-make.exe"

    x86_64-w64-mingw32-ar V
    if [ $? -ne 0 ]; then
        x86_64-w64-mingw32-gcc-ar V
        if [ $? -ne 0 ]; then
            echo 'not found x86_64-w64-mingw32-ar'
            exit 1
        fi
        sed -i "s/x86_64-w64-mingw32-ar/x86_64-w64-mingw32-gcc-ar/g" $SF_BUILD/Makefile
    fi
else
    echo "not support system $(uname -s)"
    exit 1
fi

set -e

cd $SF_BUILD
#$MAKE clean
$MAKE
cp ./softfloat.a ../libsoftfloat.a

# Build builtins
cd $root/builtins
mkdir -p build
cd build
#rm -rf *
cmake .. -G "$CMAKE_GEN" -DCMAKE_SH="CMAKE_SH-NOTFOUND" -Wno-dev
$MAKE
