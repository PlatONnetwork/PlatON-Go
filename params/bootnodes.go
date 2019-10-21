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
// the main Ethereum network.
var MainnetBootnodes = []string{
	"enode://4b11f1b1269b0d2ad0294a39a6b146574f74f97aa576dfec41f45591872a4ac1017f14e3188fe77984115af714cf9014e1d58f2582268fd03db895e83e089891@40.115.117.118:16789",
	"enode://a7c93028b7075d6cb8b345939511e44949c43bd2a891553dc51edd4bfcc9e9f0e4af1c1ae72c06040f94d2bd81202aac0cac7e1b96025e05fc26a3c9b6469ef3@52.175.21.166:16789",
	"enode://a06721ca35534e06b68bdf25a66cf113feaa170ba07faad915b32de2b1c402b6764b80220136cd92cccdd633f09ff8acbda3f11c54754e9dc5fd4fa02afce915@13.72.228.149:16789",
	"enode://1ce60ab861c13a9d97a8619e850be00728aa78ce59550819e8257b4b5b779ac8a742fc1240c6fc8430e936375f3cbe012a425b21f812f675476a36dee707ad7d@52.233.38.169:16789",
	"enode://ee9655cab5e6a6be56457fc2ab06ed38b05fef8a0f754730e064513aa821ac0bc82ba9c986ef2ceeec69b98651f15a8769a0a085d73e9df1c885f1e253661295@13.235.119.163:16789",
	"enode://3d716e226d2c92f90acad316f641691e76c5562158c8b1e7ba9deb6b26f13bfc0d4d1ce4fb6d196851c6362a2dc8c673ec8c9aef12b82bd9ad2826eb5fc3cb45@18.130.156.184:16789",
	"enode://18bf33a6ef0cdb125f356dbe18489921c288f03245b6b43ff07ca918dc6a57f86fb5958ac79654a774452079b9e889c88d856baaefc8ba0c81620dee149064ed@18.197.168.156:16789",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// test network.
var TestnetBootnodes = []string{
	"enode://6fbdf43ad1b95fd4d5ecb88946a4efcb5cfb526c1db4683ff7f9898a9a235d97b08f180a076545e622205b9129f3ca7a7a12c701a84d0044f682ae3b4f17f965@40.113.90.184:16789",
	"enode://6886e23fe0602fcdede8e2e5d00f7672255094a47fb79c044a16c41d8ede18ba5de5015cbcc7284658bbd67d94e7da46c7dd3a4979244268e03bda23bcf16475@52.236.129.182:16789",
	"enode://818e750342d00c690f34a32e420985965f32eeb41d8aad20d0cb46caf3b09446b661c4c810088f618048004e00938ab8501a0757fee0febef35db59b57c4e6d4@52.228.24.46:16789",
	"enode://dc2b2a6250cc4b7df75894e488cc83cca2f36fedbf8eea159a76265248e14930f95e2be52564f079e3f508422fdd420d69ea974bbb17b7f5607e55af955d7fb3@51.105.52.146:16789",
	"enode://d5e7124a8ad5ee5d086fe273b2d1988eca7fc8dca5cf8b0f3f78b12446d37cde30934490058344cf3c9f79af6d500ea05819d764a1a56706aab4c5311f860e49@18.196.12.156:16789",
	"enode://c38a7f2f7f5bd559a55530496379d60428717453ae779467e8343239ec0907895ae403284c8d901c3b40499d86fff19054ddd97f813ac8b31d2847fdda909a6d@13.235.239.114:16789",
	"enode://060cbd0ded75d37bea9b7aaa927f09dd353998488c734d83fdcfb5ee0f80ef06b809b2d52ff8124310daafae78ce08e2477d917cbf06406f103c4a827c12cad0@18.139.233.254:16789",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
