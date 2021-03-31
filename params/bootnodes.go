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
	"enode://81dd24640878badc06de82a82fdb0fe55f24d27877144261d81b1f39faa6686ac8f5d2489dbf97cd44d583b4b00976a8f92845378084d25c7a8bae671a543983@ms1.bfa6.platon.network:16789",
	"enode://32d628cfd32d3f464666792f4fa0bf097c723045f8fe415a8015f1c3cbd0a1bba23e7c76defac277967101dc73e3bd8fc255febb7a52d77c1018ed0cbf8d3ad4@ms2.6cc3.platon.network:16789",
	"enode://3b2ca03c94a2b8f36b88983d8666947dd08e15347980f95c395b36a4e69218c902894e9e2e92c5a2e0fe8b5c137732d2df40a118766245fdac88c480eb120c18@ms3.cd41.platon.network:16789",
	"enode://7b5323a73e9cbffd1e6d9178f0b1d55e92649aa71ebe55a0a9c577d374a9ae21ee4980aef2a3214b6e16aa9928ee48df65a382bd2d7ec19f7b87e6d993654d17@ms4.1fda.platon.network:16789",
	"enode://e6123b585a8e030b42d873d7d09b68847d1f3bba86fab84490fc29acf332a94682a8f8e1518ca857fc75391d62eaf2117703dfeed386b4e0926bf017b5cae445@ms5.ee7a.platon.network:16789",
	"enode://ab2f7cdf347d4ca26f4fdf5657d7b669464c5712cddc42609ad2060691226187815f0ce87f4dca2cac3ee618d4beeeba9618dbd31c54f97af21d16b7cbf0dccd@ms6.63a8.platon.network:16789",
	"enode://4c5a092156c43d5aa3dc71f9dc11d304d7631d393725b09e574577c583759e58ddc245e38f993cc0f32fe873cc782bfc1c62fbd49097eec3278b240de785800b@ms7.66dc.platon.network:16789",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the test network.
var TestnetBootnodes = []string{
	"enode://3fec5e5982a0b32a25168dae575c4705ab8509f266947cb8b16b62ac9eafb78d3e7efce2c31bac447edce3446a12b71383a41dcbdbe80fa856d8739b0214ff35@127.0.0.1:16789",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
