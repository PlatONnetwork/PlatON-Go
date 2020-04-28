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
	"enode://a0230e771ffed3b4174d23a51c719a2502d5d56d021da542823e48457d255fca759e1249775da34559afd249467cd651cd861679cf78070c801d0f575dcaea6a@13.237.72.220:16789",
	"enode://2f786a7e07cfab42e0960a04adf430be9a70a8fdfaf7729c0edc98d4a745625adcbe17a5213441184bfa4810b1094505db379b872678c0498dc121f84109ab9b@13.232.13.105:16789",
	"enode://5ec9e5be4cf6eca39262db9a32f90e9c9eca2f36ebfabcfba3b527b4e869ba0694a704e00ff87e4124e75ae6fc6906014ff9c0bf0035a07b75a502290dc353f3@15.222.206.58:16789",
	"enode://9bce02d3e306842f52e75ad1a212946b3f190a4f566d40c25fb5040382d9236f7a1482f5d1b6327632f8ae7012394c235968d1f4c3d2fb17e4517988fafa5cbc@3.11.110.94:16789",
	"enode://0131dd963346b1a669e8c674d6ad9d76e77b59896feda09afd39c2dc7a8c70986fa5b08b3be69eafdf4f84e4618303b76fa1e34cbee68ce5ac8c75903d091d57@3.122.78.237:16789",
}

// DemonetBootnodes are the enode URLs of the P2P bootstrap nodes running on the demo network.
var DemonetBootnodes = []string{}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
