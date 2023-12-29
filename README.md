## Go PlatON

Welcome to the PlatON-Go source code repository! This is an Ethereum-based、high-performance and high-security implementation of the PlatON protocol.
Most of peculiarities according the PlatON's **whitepaper**([English](https://www.platon.network/pdf/en/PlatON_A_High-Efficiency_Trustless_Computing_Network_Whitepaper_EN.pdf)|[中文](https://www.platon.network/pdf/zh/PlatON_A_High-Efficiency_Trustless_Computing_Network_Whitepaper_ZH.pdf)) has been developed.

[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://pkg.go.dev/github.com/PlatONnetwork/PlatON-Go?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/PlatONnetwork/PlatON-Go)](https://goreportcard.com/report/github.com/PlatONnetwork/PlatON-Go)
[![Build Status](https://github.com/PlatONnetwork/PlatON-Go/workflows/unittest/badge.svg)](https://github.com/PlatONnetwork/PlatON-Go/actions)
[![codecov](https://codecov.io/gh/PlatONnetwork/PlatON-Go/branch/feature-mainnet-launch/graph/badge.svg)](https://codecov.io/gh/PlatONnetwork/PlatON-Go)
[![version](https://img.shields.io/github/v/tag/PlatONnetwork/PlatON-Go)](https://github.com/PlatONnetwork/PlatON-Go/releases/latest)
[![GitHub All Releases](https://img.shields.io/github/downloads/PlatONnetwork/PlatON-Go/total.svg)](https://github.com/PlatONnetwork/PlatON-Go)

## Building the source
The requirements to build PlatON-Go are:

- OS:Windows10/Ubuntu18.04
- [Golang](https://golang.org/doc/install) :version 1.17+
- [cmake](https://cmake.org/) :version 3.0+
- [g++&gcc](http://gcc.gnu.org/) :version 7.4.0+
> 'cmake' and 'gcc&g++' are usually built-in with Ubuntu

In addition, the following libraries needs to be installed manually

```
sudo apt install libgmp-dev libssl-dev
```
Then, clone the repository and download dependency

```
git clone https://github.com/PlatONnetwork/PlatON-Go.git --recursive

cd PlatON-Go && go mod download
```

Ubuntu:

```
make all
```

Windows:

```
go run build\ci.go install 
```

The resulting binary will be placed in '$PlatON-Go/build/bin' .

## Getting Started

The project comes with several executables found in the `build/bin` directory.

| Command    | Description |
|:----------:|-------------|
| **`platon`** | Our main PlatON CLI client. It is the entry point into the PlatON network |
| `platonkey`    | a key related tool. |

### Generate the keys

Each node requires two pairs of public&private keys, the one is called node's keypair, it's generated based on the secp256k1 curve for marking the node identity and signning the block, and the other is called node's blskeypair, it's based on the BLS_12_381 curve and is used for consensus verifing. These two pairs of public-private key need to be generated by the platonkey tool.

Switch to the directory where contains 'platonkey.exe'(Windows) or 'platonkey'(Ubuntu).
Node's keypair(Ubuntu for example):

```
platonkey genkeypair
PrivateKey:  1abd1200759d4693f4510fbcf7d5caad743b11b5886dc229da6c0747061fca36
PublicKey :  8917c748513c23db46d23f531cc083d2f6001b4cc2396eb8412d73a3e4450ffc5f5235757abf9873de469498d8cf45f5bb42c215da79d59940e17fcb22dfc127
```
Node's blskeypair:：

```
platonkey genblskeypair
PrivateKey:  7747ec6876bbf8ca0934f05e45917b4213afc5814639355868bbf06d0b3e0f19
PublicKey :  e5eb9915ed2b5fd52cf5ff760873a75a8562956e176968f3cbe5ea2b22e03a7b5efc07fdd5ad66d433b404cb880b560bed6295fa79f8fa649588be02231de2e70a782751dc28dbf516b7bb5d52053b5cdf985d8961a5baafa467e8dda55fe981
```

> Note: The PublicKey generated by the 'genkeypair' command is the ***NodeID*** we needed, the PrivateKey is the corresponding ***node private key***, and the PublicKey generated by the 'genblskeypair' command is the node ***BLS PublicKey***, used in the staking and consensus process, PrivateKey is the ***Node BLS PrivateKey***, these two keypairs are common in different operating systems, that is, the public and private keys generated in Windows above, can be used in Ubuntu.

store the two private keys in files:

```
mkdir -p ./data
touch ./data/nodekey 
echo "{your-nodekey}" > ./data/nodekey
touch ./data/blskey
echo "{your-blskey}" > ./data/blskey
```

### Generate a wallet

```
platon --datadir ./data account new
Your new account is locked with a password. Please give a password. Do not forget this password.
Passphrase:
Repeat passphrase:
Address: {lat1anp4tzmdggdrcf39qvshfq3glacjxcd5k60wg9}
```

> Do remember the password

### Connect to the PlatON network

| Options | description |
| :------------ | :------------ |
| --identity | Custom node name |
| --datadir  | Data directory for the databases and keystore |
| --rpcaddr  | HTTP-RPC server listening interface (default: "localhost") |
| --rpcport  | HTTP-RPC server listening port (default: 6789) |
| --rpcapi   | API's offered over the HTTP-RPC interface |
| --rpc      | Enable the HTTP-RPC server |
| --nodiscover | Disables the peer discovery mechanism (manual peer addition) |
| --nodekey | P2P node key file |
| --cbft.blskey | BLS key file |

Run the following command to launch a PlatON node connecting to the PlatON's mainnet:

```
platon --identity "platon" --datadir ./data --port {your-p2p-port} --rpcaddr 127.0.0.1 --rpcport {your-rpc-port} --rpcapi "platon,net,web3,admin,personal" --rpc --nodiscover --nodekey ./data/nodekey --cbft.blskey ./data/blskey
```

OK, it seems that the chain is running correctly, we can check it as follow:

```
platon attach http://127.0.0.1:6789
Welcome to the PlatON JavaScript console!

instance: PlatONnetwork/platon/v0.7.3-unstable/linux-amd64/go1.17
at block: 26 (Wed, 15 Dec 51802 20:22:44 CST)
 datadir: /home/develop/platon/data
 modules: admin:1.0 debug:1.0 miner:1.0 net:1.0 personal:1.0 platon:1.0 rpc:1.0 txgen:1.0 txpool:1.0 web3:1.0

> platon.blockNumber
29
```

For more information, please visit our [Docs](https://devdocs.platon.network/docs/en/).

## Contributing to PlatON-Go

All of codes for PlatON-Go are open source and contributing are very welcome! Before beginning, please take a look at our contributing [guidelines](https://github.com/PlatONnetwork/PlatON-Go/blob/develop/.github/CONTRIBUTING.md). You can also open an issue by clicking [here](https://github.com/PlatONnetwork/PlatON-Go/issues/new/choose).

## Support
If you have any questions or suggestions please contact us at support@platon.network.

## License
The PlatON-Go library (i.e. all code outside of the cmd directory) is licensed under the GNU Lesser General Public License v3.0, also included in our repository in the COPYING.LESSER file.

The PlatON-Go binaries (i.e. all code inside of the cmd directory) is licensed under the GNU General Public License v3.0, also included in our repository in the COPYING file.

