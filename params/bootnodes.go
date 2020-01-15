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
	"enode://6fbdf43ad1b95fd4d5ecb88946a4efcb5cfb526c1db4683ff7f9898a9a235d97b08f180a076545e622205b9129f3ca7a7a12c701a84d0044f682ae3b4f17f965@40.113.90.184:16789",
	"enode://6886e23fe0602fcdede8e2e5d00f7672255094a47fb79c044a16c41d8ede18ba5de5015cbcc7284658bbd67d94e7da46c7dd3a4979244268e03bda23bcf16475@52.236.129.182:16789",
	"enode://818e750342d00c690f34a32e420985965f32eeb41d8aad20d0cb46caf3b09446b661c4c810088f618048004e00938ab8501a0757fee0febef35db59b57c4e6d4@52.228.24.46:16789",
	"enode://dc2b2a6250cc4b7df75894e488cc83cca2f36fedbf8eea159a76265248e14930f95e2be52564f079e3f508422fdd420d69ea974bbb17b7f5607e55af955d7fb3@51.105.52.146:16789",
	"enode://d5e7124a8ad5ee5d086fe273b2d1988eca7fc8dca5cf8b0f3f78b12446d37cde30934490058344cf3c9f79af6d500ea05819d764a1a56706aab4c5311f860e49@18.196.12.156:16789",
	"enode://c38a7f2f7f5bd559a55530496379d60428717453ae779467e8343239ec0907895ae403284c8d901c3b40499d86fff19054ddd97f813ac8b31d2847fdda909a6d@13.235.239.114:16789",
	"enode://060cbd0ded75d37bea9b7aaa927f09dd353998488c734d83fdcfb5ee0f80ef06b809b2d52ff8124310daafae78ce08e2477d917cbf06406f103c4a827c12cad0@18.139.233.254:16789",
}

// RallynetBootnodes are the enode URLs of the P2P bootstrap nodes running on the rally network.
var RallynetBootnodes = []string{
	"enode://36bd295049d060f1a772b0f67ac4f1e0d7741dc8a1e97c56c3d8f377a3838ee4a5f4afc90d97786f728f4cdd32383336a07d731059040bdf752197ce50237542@40.115.117.118:16789",
	"enode://212cfb0e4df433288f538049aa38fb0aa8dae566f92ee3a591e16dcf263628ceb5c568a105c5bda695594136c310bb5bf3a4ba89421ba29ba3232a8d9fcc9581@52.175.21.166:16789",
	"enode://7615eb03dcccbad592af88aee187be5fb92a9b11ca7e79570b2b413c8ff485b375f54240913fb368acc804714408a67783cb7cd09cd76811a7ae02cc5647d444@13.72.228.149:16789",
	"enode://fea5ee7462e1137545070af6d13f457b2de44aa9e6253aad140fcff84afcff1a12c83c888d3eb0211f87a866bf13e0a74fa6fba41d85c5b0553ab9aa4994f5f8@52.233.38.169:16789",
	"enode://2f9edd9644b9ba00516c5305a22b25f50be033432b7951881543b7594b748799da4a1e7af1596e6c744934edc083da7c0ef577844fc6fac08497c7e356412f59@13.235.119.163:16789",
	"enode://b1da85de28c06125846434d64417da7ff00bbeee95543ebc2b1545ea836c765598d90295f99cc980fe04b8fde390fa5d5db945bb6c58ebbf13b509713028c3bd@18.130.156.184:16789",
	"enode://ad2b1de0924443e157f95ccf4e402559914dfa730d9d0f60251fb35368de133b3d7e399695ff074351601e6507b85acb408587057840115291f8231cfa06cfb9@18.197.168.156:16789",
}

// UatnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the uat network.
var UatnetBootnodes = []string{
	"enode://1e1c6c4458ca9df01bc4c22cc6ee30757f2fa786bf7a3424267145ccb94ee326872ef9488bb28388e5889675ac32b06946bec3924b4a3696f48848f22d148538@52.164.228.203:16789",
	"enode://37bb0eaeda42bbc568a70cb20688aba97a8b7bbe59c8ad569087b7251cfd63eaca2d8bf27e080d02b29aeac8fff1ff2505ed556f2d0a35a8b6b5ba5f041b92ab@52.76.234.62:16789",
	"enode://d1dcf7484d9e3e8b4e3149b18a0d878b39114772f90003fc8ecd87246601ecaad40101f7ee476a72563076e89030d81753511e96746b74b337401e27c2aa932c@149.129.180.55:16789",
}

// DemonetBootnodes are the enode URLs of the P2P bootstrap nodes running on the demo network.
var DemonetBootnodes = []string{}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
