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
	"enode://3fec5e5982a0b32a25168dae575c4705ab8509f266947cb8b16b62ac9eafb78d3e7efce2c31bac447edce3446a12b71383a41dcbdbe80fa856d8739b0214ff35@ts1.7757.platon.network:16789",
	"enode://da6e4c410df2128a79e13179e744b5d28a704c2ea8d729bf63aa3c6f75afed24ef01385380560a5d607f63bda23a9a1e6bf2848e7793f3ea9b59e12f35e662a1@ts2.c72f.platon.network:16789",
	"enode://bf4c161ec0d0854bc4161a4831f15312f9115bdeb0e262022d24c3f17ca06f693f0e36c09412d165155232b11c2d77b5e5c58a805fb4d4533f648b1babd51040@ts3.0464.platon.network:16789",
	"enode://9f3ebc1ddf763e47626d7a13d4f501fcd0d83db8cc4fa26b7a2658ad41e76d29ffa8bb6354ff71d5e22d079a570a07333a3312a59350a75b625850e5819fcbd0@ts4.31e0.platon.network:16789",
	"enode://18820d5df0684e7871643351573f095d99c5c6ecb26021627c60ff55c6d146f5b9984d43cc55f9efaaf043b40628a1c2a770e7e891878e1dc168a3acb7e05dc6@ts5.5eb0.platon.network:16789",
	"enode://991f577221be7a18800ab8ae8eb3214ef3e1fa0be4a161c9cd95100b56e95df677aa12dc624aaedc073c5b48367df0351a4955e91290f19fba8c37413cf5b5e0@ts6.8847.platon.network:16789",
	"enode://86fbd54d03786d924ee33d1fe15542c4ceb528a058c955e42327f0c528db7b2494aecc814b361b4f501950f46aca4c2926e255e6542db991546e0540d24a9d35@ts7.72f5.platon.network:16789",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
