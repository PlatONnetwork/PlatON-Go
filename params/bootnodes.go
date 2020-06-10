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
	"enode://3fec5e5982a0b32a25168dae575c4705ab8509f266947cb8b16b62ac9eafb78d3e7efce2c31bac447edce3446a12b71383a41dcbdbe80fa856d8739b0214ff35@35.228.209.144:16789",
	"enode://da6e4c410df2128a79e13179e744b5d28a704c2ea8d729bf63aa3c6f75afed24ef01385380560a5d607f63bda23a9a1e6bf2848e7793f3ea9b59e12f35e662a1@34.76.112.182:16789",
	"enode://bf4c161ec0d0854bc4161a4831f15312f9115bdeb0e262022d24c3f17ca06f693f0e36c09412d165155232b11c2d77b5e5c58a805fb4d4533f648b1babd51040@52.64.155.63:16789",
	"enode://9f3ebc1ddf763e47626d7a13d4f501fcd0d83db8cc4fa26b7a2658ad41e76d29ffa8bb6354ff71d5e22d079a570a07333a3312a59350a75b625850e5819fcbd0@3.7.228.187:16789",
	"enode://18820d5df0684e7871643351573f095d99c5c6ecb26021627c60ff55c6d146f5b9984d43cc55f9efaaf043b40628a1c2a770e7e891878e1dc168a3acb7e05dc6@15.223.91.58:16789",
	"enode://991f577221be7a18800ab8ae8eb3214ef3e1fa0be4a161c9cd95100b56e95df677aa12dc624aaedc073c5b48367df0351a4955e91290f19fba8c37413cf5b5e0@18.132.106.72:16789",
	"enode://86fbd54d03786d924ee33d1fe15542c4ceb528a058c955e42327f0c528db7b2494aecc814b361b4f501950f46aca4c2926e255e6542db991546e0540d24a9d35@3.127.182.139:16789",
}

// DemonetBootnodes are the enode URLs of the P2P bootstrap nodes running on the demo network.
var DemonetBootnodes = []string{}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
