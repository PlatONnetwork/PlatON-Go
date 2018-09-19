#!/bin/sh

set -e

echo "into env.sh..."

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

echo "into env.sh - 02..."

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"

#for test
echo "into env.sh - 03..."
echo "$workspace"

root="$PWD"
atpdir="$workspace/src"

echo "$atpdir ||| $root"

if [ ! -L "$atpdir/Platon-go" ]; then
    mkdir -p "$atpdir"
    cd "$atpdir"
    # /cygdrive/c/sunzone/MyDocument/liteide/src/Platon-go/build/_workspace/src
    ln -s ../../../. Platon-go
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# fot test
echo "into env.sh - 04..."
echo "GOPATH: $GOPATH"
echo "$@"

# Run the command inside the workspace.
cd "$atpdir/Platon-go"
PWD="$atpdir/Platon-go"

go version

# Launch the arguments with the configured environment.
exec "$@"
