// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main PlatON network.
var MainnetBootnodes = []string{
	"enode://f1409fe5a87442808ee99bd244447c5df362c8da3e2a9a136f1a5ebb6ddf41cf65c53342aa054345ffcc44a4bec081c03ccca65ae470d3d2ef7cfc9a4f594830@40.115.117.118:16789",
	"enode://c1011f4956790caf8f40ecb719a123eca1ff90ce4e54e6799495c05382d7f40e00e654215c2740de38737d2c1142a81d1e5cfc858769eda53d0d389abde6caab@52.175.21.166:16789",
	"enode://1f47c61b520f9c4809acc89a0a9e8e924537e884b46dde5123cfa9b2d8a2c8b0ef65e9418b1557b4691f9ba961110e88b7e1578cadb3c3eef32e6648f7a0d71e@13.72.228.149:16789",
	"enode://1c728f1444f42373b20305692c5eeed76a21c925ea5ac78a3adaf2696dc33d3fbc7e06bcb5f825a82189c2dbab1383f3d6f4ec5f72f317ea5e5c38bcc4fb8537@52.233.38.169:16789",
	"enode://3f569dced6c677c035e08bf01205c5810900019b0987464741723dfc1a51cd3d3d4be1976e54e72a9e1923d3c57d1889651aefb22a8abce043d9f2807511237e@13.235.119.163:16789",
	"enode://17a476dbd2016efab80682d587b8e3de6ed94a021d42c8a357e9ac0f534b7564cd6b7c0d8e798ad59121a899075a2b264f56c698e7c5b489f92d5cf58dfd717a@18.130.156.184:16789",
	"enode://6a656ad24cf282b17df422c0820cd60bd02cf578ea3d19831ef56767d7a452299a3b43ae5c88fdd8da6a0393d71932f2490e41652952fbd27e7f2344fc96ea21@18.197.168.156:16789",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the test network.
var TestnetBootnodes = []string{
	"enode://0dcb00d1fb3d9dfa319922c4e996e3183243ef537753cc95db77891cc14d48e2b9cbe6cbff7ceddfdff04c98d5cfbe47e48357ca4ccba146d1406cd89931e986@35.228.209.144:16789",
	"enode://b48ccda0d0134126ee68dae5d8b63ad118dece802cc7f36d4453c5a3a268dfd0319c22b3c0225516cac2395e7a70dac4c93f13b730936f07d7abd0694bf1b12e@34.76.112.182:16789",
	"enode://31f9dc63aa7502bae76d39d95868aa76b5239894f0fef4e3bf2dfca0d26c0f47bac4a79067b05d859305a46ad6ab9af6b639e86dd066ea4df1bff189cd4d251b@13.237.72.220:16789",
	"enode://8b84524cc17bf715cec3134abd22ee45c91b80443499316f3bebb7281d24c2d878ad406d8d2474c7c0740350742c55e8211c6b1b7213b90f2adde5c061c0ec0d@13.232.13.105:16789",
	"enode://a39d2345a3fbe7b45ef568d96a1ea6d8b9e110ec7be031d214dbe618b41a94f55c6ee6aaae39efc1b2e72a290a5682c6dbbc91eb1b5bc58138c7434e73244970@15.222.206.58:16789",
	"enode://abefe91fed4b0e79dc0f157e8b51ef8e17eceb22fffc74f2b9bd979812195ecc05cff199db4a52add80bedad4bab07ebc72a87f222942246fdd790487a011122@3.11.110.94:16789",
	"enode://d8db40aadeef04749d055372c4041a5e45e6142868acbdf8b2abfa8fb1cda2f626bd937910722cb36e2ad226bdf8e87189f6f457820dc04622363ac6dcf4a5a3@3.122.78.237:16789",
}

// DemonetBootnodes are the enode URLs of the P2P bootstrap nodes running on the demo network.
var DemonetBootnodes = []string{}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
