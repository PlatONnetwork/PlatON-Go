#!/bin/bash

#!/usr/bin/env bash

if [ ! -f "build/cbft_test.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi
root=`pwd`

hash=$(go test -v  -tags=test github.com/PlatONnetwork/PlatON-Go/consensus/cbft   -run TestCbft_CreateGenesis | sed -n '2p')
echo "replace root $hash"

tmp='Root:        common.BytesToHash(hexutil.MustDecode("HASH")),'
relpac=${tmp/HASH/$hash}

sed -i "/Root/c$relpac"  $root/consensus/cbft/cbft_common_util.go
go fmt $root/consensus/cbft/cbft_common_util.go