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
	"enode://f7c33bd34b0e3c9a0317733ef3356409ff2eb009605cc357c213c367faf833ab48d557942731fd8dfdd39b92004e863c0fd8ebd01e79f69acd2b82f60ac63074@ms1.bfa6.platon.network:16789",
	"enode://7729f2313908670523d7babf227e69e93150aac916f3f372a36ee7f204ed55737cb667fc55c34c40c5499734ed08cc7c57800d1ac8131c0cb855768801b898e9@ms2.6cc3.platon.network:16789",
	"enode://4649f744f3e1d2400773fc48e057b96a8d4a10e00121f884f97b3182187ded0f89f5f4dbade55acaa4155e25c281f23a34587bad4fc3af2403eef9c130b57e5b@ms3.cd41.platon.network:16789",
	"enode://f77401d3dda6d0c58310744e9349c16c056f94179a4d7bdc3470b6461d7f64370fa21ebf380ab10e25834291215715550c6845442a0df88ab1c42d161d367626@ms4.1fda.platon.network:16789",
	"enode://8c71f4e1e795fc6e73144e4696a9fde3c3cdf6b99ab575357b77ab22542bc70c8f04e88f23eb6cdb225ad077aa67b71245da1ab1838bf8b362464ffd515ca3d6@ms5.ee7a.platon.network:16789",
	"enode://6c9f8a51ff27bb0e062952be1f1e3943847eaed1f14d54783238678767dd134ddedfc82c0f52696cee971ee1aabec00523cd634e41c079ee95407aa4ecb92c7b@ms6.63a8.platon.network:16789",
	"enode://0db310d4a6c429dcac973ff6433659ed710783872cf62bdcec09c76a8bb380d51f0c153401ceaa27e0b3de6f56fb939115c495862d15f6ccc8019702117be34d@ms7.66dc.platon.network:16789",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the test network.
var TestnetBootnodes = []string{
	"enode://3fec5e5982a0b32a25168dae575c4705ab8509f266947cb8b16b62ac9eafb78d3e7efce2c31bac447edce3446a12b71383a41dcbdbe80fa856d8739b0214ff35@127.0.0.1:16789",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
