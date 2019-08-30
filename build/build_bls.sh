#!/usr/bin/env bash

if [ ! -f "build/build_bls.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

PLATON_ROOT=`pwd`
BLS_ROOT=$PLATON_ROOT/crypto/bls

if [ "`ls $BLS_ROOT/bls_win`" = "" ]; then
    # pull bls
    git submodule update --init
fi

if [ `expr substr $(uname -s) 1 5` == "MINGW" ]; then
    echo "not support system $(uname -s)"
    exit 0
fi


# sudo apt install libgmp-dev
# sudo apt install libssl-dev
# the above are prerequisites

cd $BLS_ROOT
mkdir -p bls_linux_darwin
cd bls_linux_darwin
mkdir -p include
mkdir -p lib
mkdir -p src
#cd src
#git clone https://github.com/herumi/mcl.git
#git clone https://github.com/herumi/bls.git
# below is only for Windows
# git clone https://github.com/herumi/cybozulib_ext.git

set -e

# Build and test bls lib
MAKE="make"
cd $BLS_ROOT/bls_linux_darwin/src/bls
$MAKE -j 4

# copy bls  header and lib files to destination directory
cp -r $BLS_ROOT/bls_linux_darwin/src/bls/include/bls $BLS_ROOT/bls_linux_darwin/include/
rm -rf $BLS_ROOT/bls_linux_darwin/src/bls/ffi
cp $BLS_ROOT/bls_linux_darwin/src/bls/lib/*.a $BLS_ROOT/bls_linux_darwin/lib/

# copy mcl header and lib files to destination directory
cd $BLS_ROOT/bls_linux_darwin/src/mcl
rm -rf ffi
cp -r $BLS_ROOT/bls_linux_darwin/src/mcl/include/mcl $BLS_ROOT/bls_linux_darwin/include/
cp -r /$BLS_ROOT/bls_linux_darwin/src/mcl/include/cybozu $BLS_ROOT/bls_linux_darwin/include/
cp $BLS_ROOT/bls_linux_darwin/src/mcl/lib/*.a $BLS_ROOT/bls_linux_darwin/lib/
