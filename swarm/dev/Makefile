.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    Platon-go/cmd/swarm
	go build -o bin/geth     Platon-go/cmd/geth
	go build -o bin/bootnode Platon-go/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v Platon-go/swarm/...
