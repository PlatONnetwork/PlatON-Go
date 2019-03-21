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
var MainnetBootnodes = []string{}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Alpha test network.
var TestnetBootnodes = []string{
	"enode://a6ef31a2006f55f5039e23ccccef343e735d56699bde947cfe253d441f5f291561640a8e2bbaf8a85a8a367b939efcef6f80ae28d2bd3d0b21bdac01c3aa6f2f@test-sea.platon.network:16789",       //TEST-SEA
	"enode://d124e660938dc3fd63d913ff753fafc262764b22294431e760b572b0b58d5e6b813b32ccbacc326c03171542ae0ff8ff6528625a2d612e0c49240f111eba3c22@test-sg.platon.network:16790",        //TEST-SG
	"enode://24b0c456ae5cad46c4fb9bc02c867b997e22f30696e6e330926f785ca2e7410baf1eb34ffd9b5b07b5ba6e02b693faf57afb33f7c66cfbcf4c9186b4bfac737d@test-na.platon.network:16789",        //TEST-NA
	"enode://c7fc34d6d8b3d894a35895aaf2f788ed445e03b7673f7ce820aa6fdc02908eeab6982b7eb97e983cc708bcec093b3bc512b0b1fbf668e6ab94cd91f2d642e591@test-us.platon.network:16790",        //TEST-US
	"enode://9871adb2f926dffa3ff6060e07ae85295ce4184d5881cc761e465ca59597a7c5fa46b589557b0be62b759344fec50313a69b5fbda8b420f058ede85dadcecc4a@test-sg-soga.platon.network:16789",   //TEST-SEA-SOGA
	"enode://73323061805daa21ad07aa31cf0cc8c2295b05cff47c9ecb25c7a215c1c720df6c8698e94632346654cc4d8c0e99688f367626f20db9be85f29e9f41c29ffb92@test-siga-soga.platon.network:16790", //TEST-SG-SOGA
	"enode://23aa343260d06e04107d1cd9a7d12c54cc238719a1523ffe42640210c913218b5940d41511c5adb716da38844a85cdab8b7db0600d242e24168d7df10aebd324@test-si-syde.platon.network:16789",   //TEST-SI-SYDE
}

// BetanetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Beta test network.
var BetanetBootnodes = []string{
	"enode://bcb7e49461cdd5f3227bb6cc6c36675cd936c11b69c3fd366c36997d514beabc423f8dfee6f91330a96273988bb68b1785161631181fd738d0f46d263b3ce8b3@54.176.216.82:16789",
	"enode://5449094bf985a688d378a90cf334d5a1abc55d694d6f2362899494d18048ef6b6bd724f4e51084bfe0563c732c481869c9da05d92e56f29f6880ad15ea851f13@54.176.216.82:16790",
	"enode://c0f7ae43af0605b80e35a5469adaa142059eaaf41d152613d74d42feffd6871f059f9ac4d596bd134bb1d6bbfbcea5391adff6f005ea9042c21797d51d0b7697@3.1.59.5:16789",
	"enode://b6883e86e833cec2405fb548405f7a1e693379f77ee8fc6bbf41b5c853d7ad654a2a3fb7ffbe57ae848509d1ed7e11acaf28666f8f81646eab575dafa8d51d0b@3.1.59.5:16790",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
