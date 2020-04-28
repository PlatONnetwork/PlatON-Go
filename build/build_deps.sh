#!/usr/bin/env bash

if [ ! -f "build/build_deps.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi
root=`pwd`
bls_build=$root/build/


$bls_build/build_bls.sh


