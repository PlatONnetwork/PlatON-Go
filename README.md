## Go PlatON

Welcome to the PlatON source code repository! This software is Ethereum-based and it has changed some peculiarities according the PlatON's [whitepaper](https://www.platon.network/static/pdf/en/PlatON_A%20High-Efficiency%20Trustless%20Computing%20Network_Whitepaper_EN.pdf).

## Building the source

For prerequisites and detailed build instructions please read the
[Installation Instructions](https://github.com/PlatONnetwork/wiki/wiki)
on the wiki.

Building platon requires both a Go (version 1.7 or later) and a C compiler.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make platon

or, to build the full suite of utilities:

    make all

If you want to Building platon with MPC function, run

    make platon-with-mpc

or:

    make all-with-mpc

## Executables

The project comes with several executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **`platon`** | Our main PlatON CLI client. It is the entry point into the PlatON network |
| `ethkey`    | a key related tool. |

## Running a platon node

### Config the chain data

first, you need to get an account:

```
$ ./platon --datadir ./data account new
Your new account is locked with a password. Please give a password. Do not forget this password.
Passphrase:
Repeat passphrase:
Address: {566c274db7ac6d38da2b075b4ae41f4a5c481d21}
```

second, generate a private node's key pair and save the PrivateKey as a file named 'nodekey' into the ./data

```
$ ./ethkey genkeypair
Address   :  0xA9051ACCa5d9a7592056D07659f3F607923173ad
PrivateKey:  1abd1200759d4693f4510fbcf7d5caad743b11b5886dc229da6c0747061fca36
PublicKey :  8917c748513c23db46d23f531cc083d2f6001b4cc2396eb8412d73a3e4450ffc5f5235757abf9873de469498d8cf45f5bb42c215da79d59940e17fcb22dfc127
```

then, edit the following content and save it as json file, such as genesis.json:

```
{
    "config": {
    "chainId": 300,
    "homesteadBlock": 1,
    "eip150Block": 2,
    "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "eip155Block": 3,
    "eip158Block": 3,
    "byzantiumBlock": 4,
    "cbft": {
      "initialNodes": ["enode://8917c748513c23db46d23f531cc083d2f6001b4cc2396eb8412d73a3e4450ffc5f5235757abf9873de469498d8cf45f5bb42c215da79d59940e17fcb22dfc127@127.0.0.1:16789"]
      }
  },
  "nonce": "0x0",
  "timestamp": "0x5c074288",
  "extraData": "0x00",
  "gasLimit": "0x47b760",
  "difficulty": "0x40000",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc": {
    "0x566c274db7ac6d38da2b075b4ae41f4a5c481d21": {
      "balance": "999000000000000000000"
    }
  },
  "number": "0x0",
  "gasUsed": "0x0",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}
```

at last, init the chain as follow:

```
$ ./platon --datadir ./data init platon.json
```

and it will output msg as:

```
...
Successfully wrote genesis state
```

so we can launch the node: 

```
$ ./platon --identity "platon" --datadir ./data --nodekey ./data/platon/nodekey --rpcaddr 0.0.0.0 --rpcport 6789 --rpcapi "db,eth,net,web3,admin,personal" --rpc --nodiscover
```

### Send a transaction

```
> eth.sendTransaction({from:"0x566c274db7ac6d38da2b075b4ae41f4a5c481d21",to:"0x3dea985c48e82ce4023263dbb380fc5ce9de95fd",value:10,gas:88888,gasPrice:3333})
"0xa8a79933511158c2513ae3378ba780bf9bda9a12e455a7c55045469a6b856c1b"
```

Check the balance:

```
> eth.getBalance("0x3dea985c48e82ce4023263dbb380fc5ce9de95fd")
10
```
 
OK, it seems that the chain is running correctly

For more information, please visit our [wiki](https://github.com/PlatONnetwork/wiki/wiki)