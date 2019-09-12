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
// Alpha test network.
var TestnetBootnodes = []string{
	"enode://6fbdf43ad1b95fd4d5ecb88946a4efcb5cfb526c1db4683ff7f9898a9a235d97b08f180a076545e622205b9129f3ca7a7a12c701a84d0044f682ae3b4f17f965@40.113.90.184:16789",
	"enode://6886e23fe0602fcdede8e2e5d00f7672255094a47fb79c044a16c41d8ede18ba5de5015cbcc7284658bbd67d94e7da46c7dd3a4979244268e03bda23bcf16475@52.236.129.182:16789",
	"enode://818e750342d00c690f34a32e420985965f32eeb41d8aad20d0cb46caf3b09446b661c4c810088f618048004e00938ab8501a0757fee0febef35db59b57c4e6d4@52.228.24.46:16789",
	"enode://dc2b2a6250cc4b7df75894e488cc83cca2f36fedbf8eea159a76265248e14930f95e2be52564f079e3f508422fdd420d69ea974bbb17b7f5607e55af955d7fb3@51.105.52.146:16789",
	"enode://d5e7124a8ad5ee5d086fe273b2d1988eca7fc8dca5cf8b0f3f78b12446d37cde30934490058344cf3c9f79af6d500ea05819d764a1a56706aab4c5311f860e49@18.196.12.156:16789",
	"enode://c38a7f2f7f5bd559a55530496379d60428717453ae779467e8343239ec0907895ae403284c8d901c3b40499d86fff19054ddd97f813ac8b31d2847fdda909a6d@13.235.239.114:16789",
	"enode://060cbd0ded75d37bea9b7aaa927f09dd353998488c734d83fdcfb5ee0f80ef06b809b2d52ff8124310daafae78ce08e2477d917cbf06406f103c4a827c12cad0@18.139.233.254:16789",
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
