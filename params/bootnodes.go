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
	"enode://cdf324ff821c314bb12949f99f7728c3600b1f4c4e384d8f4fa92a8df9fdc15ab36b0446ad6689a379daa3cd451f10b7030c7423305712d935a70f7690c7c5bb@40.113.90.184:16789",
	"enode://b00fbdcf21e9c9a983fa5fd3a3f1146c4c876f68bf074d6afdceef8df6d05b11cb914626d2b05dc5aebc3a5882e705315e1cb707b058065d906e71a06a711fa6@52.236.129.182:16789",
	"enode://50c64b8a047dc5dc76688bc45730d76b2d1fbdb886b964136188c0d8b4c45f44770efb109853714afc1111193fdffdcae039ad695afa6a886de636e8da77862d@52.228.24.46:16789",
	"enode://630ba358cef8c91e45718c19a42009ec65901b66787cd6410bbde4bef13f1279ae7d9dd4195b42c11a77039b49529b0825b820a20f64ccc3b296a0304c67dda9@51.105.52.146:16789",
	"enode://4a64b067c13e44ac582476322a7858498281e616fdc47424c89e1a13d6d8e703554255daeb3bf6e6799882d1e80d3c1d3e3df9007f6ccd3b5d637a901d97ed8c@18.196.12.156:16789",
	"enode://532382159638cd5fb87b03e6abaae2281f2b389439eb5b2e02cbd0f30d23cfb414b88833e504bd90defe3ca279cbb20949866a8c38b9bf002b809807f50293e6@13.235.239.114:16789",
	"enode://7028f8508eb6b65eba30b1a50b83411e96b05b00dcf26cedb1af30c24a4061189a0556c8e17386a50fbd643aa11fee66ec82a3b194c1aadd89a442a6bad3c64f@18.139.233.254:16789",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Alpha test network.
var TestnetBootnodes = []string{
	"enode://a6ef31a2006f55f5039e23ccccef343e735d56699bde947cfe253d441f5f291561640a8e2bbaf8a85a8a367b939efcef6f80ae28d2bd3d0b21bdac01c3aa6f2f@test-sea.platon.network:16791",       //TEST-SEA
	"enode://d124e660938dc3fd63d913ff753fafc262764b22294431e760b572b0b58d5e6b813b32ccbacc326c03171542ae0ff8ff6528625a2d612e0c49240f111eba3c22@test-sg.platon.network:16792",        //TEST-SG
	"enode://24b0c456ae5cad46c4fb9bc02c867b997e22f30696e6e330926f785ca2e7410baf1eb34ffd9b5b07b5ba6e02b693faf57afb33f7c66cfbcf4c9186b4bfac737d@test-na.platon.network:16791",        //TEST-NA
	"enode://c7fc34d6d8b3d894a35895aaf2f788ed445e03b7673f7ce820aa6fdc02908eeab6982b7eb97e983cc708bcec093b3bc512b0b1fbf668e6ab94cd91f2d642e591@test-us.platon.network:16792",        //TEST-US
	"enode://9871adb2f926dffa3ff6060e07ae85295ce4184d5881cc761e465ca59597a7c5fa46b589557b0be62b759344fec50313a69b5fbda8b420f058ede85dadcecc4a@test-sg-soga.platon.network:16791",   //TEST-SEA-SOGA
	"enode://73323061805daa21ad07aa31cf0cc8c2295b05cff47c9ecb25c7a215c1c720df6c8698e94632346654cc4d8c0e99688f367626f20db9be85f29e9f41c29ffb92@test-siga-soga.platon.network:16792", //TEST-SG-SOGA
	"enode://23aa343260d06e04107d1cd9a7d12c54cc238719a1523ffe42640210c913218b5940d41511c5adb716da38844a85cdab8b7db0600d242e24168d7df10aebd324@test-si-syde.platon.network:16790",   //TEST-SI-SYDE
}

// BetanetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Beta test network.
var BetanetBootnodes = []string{
	"enode://bcb7e49461cdd5f3227bb6cc6c36675cd936c11b69c3fd366c36997d514beabc423f8dfee6f91330a96273988bb68b1785161631181fd738d0f46d263b3ce8b3@54.176.216.82:16791",
	"enode://5449094bf985a688d378a90cf334d5a1abc55d694d6f2362899494d18048ef6b6bd724f4e51084bfe0563c732c481869c9da05d92e56f29f6880ad15ea851f13@54.176.216.82:16792",
	"enode://c0f7ae43af0605b80e35a5469adaa142059eaaf41d152613d74d42feffd6871f059f9ac4d596bd134bb1d6bbfbcea5391adff6f005ea9042c21797d51d0b7697@3.1.59.5:16791",
	"enode://b6883e86e833cec2405fb548405f7a1e693379f77ee8fc6bbf41b5c853d7ad654a2a3fb7ffbe57ae848509d1ed7e11acaf28666f8f81646eab575dafa8d51d0b@3.1.59.5:16792",
}

// InnerTestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Inner test network.
var InnerTestnetBootnodes = []string{
	"enode://97e424be5e58bfd4533303f8f515211599fd4ffe208646f7bfdf27885e50b6dd85d957587180988e76ae77b4b6563820a27b16885419e5ba6f575f19f6cb36b0@192.168.120.81:16789",
	"enode://3b53564afbc3aef1f6e0678171811f65a7caa27a927ddd036a46f817d075ef0a5198cd7f480829b53fe62bdb063bc6a17f800d2eebf7481b091225aabac2428d@192.168.120.82:16789",
	"enode://858d6f6ae871e291d3b7b2b91f7369f46deb6334e9dacb66fa8ba6746ee1f025bd4c090b17d17e0d9d5c19fdf81eb8bde3d40a383c9eecbe7ebda9ca95a3fb94@192.168.120.83:16789",
	"enode://e4556b211eb6712ab94d743990d995c0d3cd15e9d78ec0096bba24c48d34f9f79a52ca1f835cec589c5e7daff30620871ba37d6f5f722678af4b2554a24dd75c@192.168.120.84:16789",
	"enode://114e48f21d4d83ec9ac39a62062a804a0566742d80b191de5ba23a4dc25f7beda0e78dd169352a7ad3b11584d06a01a09ce047ad88de9bdcb63885e81de00a4d@192.168.120.85:16789",
	"enode://64ba18ce01172da6a95b0d5b0a93aee727d77e5b2f04255a532a9566edaee7808383812a860acf5e43efeca3d9321547bfcdefd89e9d0c605dcdb65ce0bbb617@192.168.120.86:16789",
	"enode://d31b3a7714610bd8e03b2c74aca4be16de7fcc319a1e577d50e5e8796680221b4b679bf1c37966d1a158902b8686f3ca2f41a89a7176e538141082540c4f6d66@192.168.120.87:16789",
	"enode://805b617b9d321a65d8936e758b5c60cd6e8c873b9f1e7c793ad5f887d26ce9667d0db2fe55a9aeb1cc81f9cf9a1e7c54473203473e3ebda89e63c03cbcfe5347@192.168.120.88:16789",
	"enode://fa147bc3625acc846a9f0e1e89172ca7470baa0f86516994f70860c6fb904ddbb1849e3cf2b40c58255e38401f40d2c3e4a3bd5c2f2849b98465a5bdb80ed6a0@192.168.120.89:16789",
	"enode://d8c4b58ae052ea9480577264bc1b2c09619757015849a4c92b71a4e4c8b5ede94f35d24107b1181d0711013ed7fdc068f21e6e6084b3e96750a571669715c0b1@192.168.120.90:16789",
}

// InnerDevnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Inner dev network.
var InnerDevnetBootnodes = []string{
	"enode://0abaf3219f454f3d07b6cbcf3c10b6b4ccf605202868e2043b6f5db12b745df0604ef01ef4cb523adc6d9e14b83a76dd09f862e3fe77205d8ac83df707969b47@192.168.9.76:16789",
	"enode://e0b6af6cc2e10b2b74540b87098083d48343805a3ff09c655eab0b20dba2b2851aea79ee75b6e150bde58ead0be03ee4a8619ea1dfaf529cbb8ff55ca23531ed@192.168.9.76:16790",
	"enode://15245d4dceeb7552b52d70e56c53fc86aa030eab6b7b325e430179902884fca3d684b0e896ea421864a160e9c18418e4561e9a72f911e2511c29204a857de71a@192.168.120.76:16789",
	"enode://fb886b3da4cf875f7d85e820a9b39df2170fd1966ffa0ddbcd738027f6f8e0256204e4873a2569ef299b324da3d0ed1afebb160d8ff401c2f09e20fb699e4005@192.168.120.76:16790",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}
