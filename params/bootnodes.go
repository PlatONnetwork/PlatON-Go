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
	"enode://582306a348dd89fe8b38acf414e854ffe7f7b2b92ba91bcf65239876475e521ff3a6c3f20095548eb0169fc63594a45dabb48b9b89207f421a45771d8b4afc04@ms1.bfa6.platon.network:16789",
	"enode://a79e50bfed6c000cb03c69e915eee07e804524e438920623c4472145aaae6825a4a04eb689c0615d979feb400f4c6062940e936eb7b3f90833d58e39a890d629@ms2.6cc3.platon.network:16789",
	"enode://a15d782889b4fc1af0120018a1025d4397fb61c47e0d45c7f3e06c81b6d5f47b629988b7a369a08d674bcadfd0817e00f1ab5ffc0f56037ac392934d2f44a94c@ms3.cd41.platon.network:16789",
	"enode://e4f1b343a571b234fc9df7080f51317e1eda311d0cabbd794028e4082c82524bf4ae9f9a1c0beb6245602a9770d71e153760c31468935973bd03bdecb0dd2c23@ms4.1fda.platon.network:16789",
	"enode://6dca11ad5a82680b938766b1458056dee4622d9fb174627e731901600c35940962b4bbee1cac5273b9685f6156f81955491fd1b570e940f865d882134c0ce9e6@ms5.ee7a.platon.network:16789",
	"enode://b174e899d9baa278fb10a705c130bedf78cdf057bbacee6f1a7fbe96429c4f39ab38043c069560155e836d052bd4a53b7c3bb83ffb3dff4f717b9f3cedfbdd9f@ms6.63a8.platon.network:16789",
	"enode://236a28f0d80ed8dff82bd4523aac056b8a24cb6b97b2d8a75bc2efa7dfdda844c82d0c15d84c819e4e2c2eeb2d4b2a5a88dae2b299afe7a65c11833f61be2d50@ms7.66dc.platon.network:16789",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the test network.
var TestnetBootnodes = []string{
	"enode://3fec5e5982a0b32a25168dae575c4705ab8509f266947cb8b16b62ac9eafb78d3e7efce2c31bac447edce3446a12b71383a41dcbdbe80fa856d8739b0214ff35@127.0.0.1:16789",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
