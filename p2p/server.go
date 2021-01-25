// Copyright 2014 The go-ethereum Authors
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

// Package p2p implements the Ethereum p2p network protocols.
package p2p

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mclock"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discv5"
	"github.com/PlatONnetwork/PlatON-Go/p2p/nat"
	"github.com/PlatONnetwork/PlatON-Go/p2p/netutil"
)

const (
	defaultDialTimeout = 15 * time.Second

	// Connectivity defaults.
	maxActiveDialTasks     = 16
	defaultMaxPendingPeers = 50
	defaultDialRatio       = 3

	// Maximum time allowed for reading a complete message.
	// This is effectively the amount of time a connection can be idle.
	frameReadTimeout = 30 * time.Second

	// Maximum amount of time allowed for writing a complete message.
	frameWriteTimeout = 20 * time.Second
)

var errServerStopped = errors.New("server stopped")

var AllowNodesMap = map[discover.NodeID]struct{}{
	discover.MustHexID("f7a07f6ff8c282a4a1c1febaf68531bbb557fa262c641d022f249e926e222c9190602b502d2772fa4a2834ffb4bf6298f4c80467178e90728c5a852504f571ad"): struct{}{},
	discover.MustHexID("87aa41f5f3ae96260017f34d836708d5e92319b766b6d475744a01e939c8921715c4fecb7951e9bc9aae7c334ca51c77247aed32e643e49341037996a28dace7"): struct{}{},
	discover.MustHexID("8f1c8333053cad81c76baf4df547eab49d1d1e1f7602b470ed3cd36e14461042731b20904252bff629d1fe60b05c3223653410c72a30fe90ef66ed95e6f4289a"): struct{}{},
	discover.MustHexID("0c2bf9b0d73ebd5728e810deb55a6f656cc44b92a8399aa0898eaa26bad3232707441e381f7a83eb1c0b38bde813bf84e7311b8d9fe75f7c5eea6c219c3103cc"): struct{}{},
	discover.MustHexID("89ca7ccb7fab8e4c8b1b24c747670757b9ef1b3b7631f64e6ea6b469c5936c501fcdcfa7fef2a77521072162c1fc0f8a1663899d31ebb1bc7d00678634ef746c"): struct{}{},
	discover.MustHexID("fff1010bbf1762d13bf13828142c612a7d287f0f1367f8104a78f001145fd788fb44b87e9eac404bc2e880602450405850ff286658781dce130aee981394551d"): struct{}{},
	discover.MustHexID("32b7ca6dec2f4e96187d6dbbed02a224dc91b605302689ac921f069c9bd314d431ffaa4ecd220164513cba63361a78e672ba2c9557e44b76c7e6433ed5fbee94"): struct{}{},
	discover.MustHexID("0eb6b43a9945a062e67b45248084ec4b5da5f22d35a58991c8f508666253fbd1b679b633728f4c3384ee878ca5efca7623786fdf623b4e5288ace830dc237614"): struct{}{},
	discover.MustHexID("a3f400758b678e15d7db188f2448411b8b472cb94da7a712dd7842f03c879478e156fd7a549a782a12de5f5dd9dc979dfd2946f8396b9328dd4fbdead37e49ba"): struct{}{},
	discover.MustHexID("6f5584a27a272099c1c8dd948f24f09b660ca38918e76db0eb3fad3c177b97229775ef61d45341279317bce68cf08fce30473dcfd30a0348049b0289048e29a3"): struct{}{},
	discover.MustHexID("c45a914bba5ad5c735489a100954eb08c49c3f718879e483d0debbe3e83da409cb1f9100ac73a2951bdf38886db419d8648534698cad05ab99c27d6869f62582"): struct{}{},
	discover.MustHexID("680b23be9f9b1fa7684086ebd465bbd5503305738dae58146043632aa553c79c6a22330605a74f03e84379e9706d7644a8cbe5e638a70a58d20e7eb5042ec3ca"): struct{}{},
	discover.MustHexID("4be36c1bda3d630c2448b646f2cb96a44867f558c32b3445f474cc17cb26639f80bf3a5f8bb993261cf7862e3d43e6306eec2a5c8c712f4e1ead3bb126899abd"): struct{}{},
	discover.MustHexID("a138433f1bc37e0c7a2bfb82f9f91f4bfad697524616f0aa08a78fdf1ecad9315cc3f80e8736b9d888e5f605ab49e2404461ea69c34b06697f067b509914be46"): struct{}{},
	discover.MustHexID("00aa2e361168765c0df1ebe981455dd64f4c091adb01e162567101f8f1b5fa31ac960c198406a00fc558ebb8c14843ef72ed47d5931d0605a0e6caab41c8a86e"): struct{}{},
	discover.MustHexID("73686446756b83458b910015775081acaa83d6727f4dcf71b732a827ebd563bc4553e1a3171029899ca276f285d7b63b36609ab8aaa0f138eeee38724d89a15e"): struct{}{},
	discover.MustHexID("656e0ee96374d3f60d2e323a059fb8848490433b4d0cefd216ce47e2e575375bc561745c9f35ecdf080848f5f08c5a4bd7af839be594f4abce7666b7e2a2a2a6"): struct{}{},
	discover.MustHexID("cc77a9f28ae8a89acc48ec969c34fad9ed877d636eccf22725af7bd89b274f702d40073d2dda1b9fbfc572746979a7876c1dca93fc617918586a53747d211d43"): struct{}{},
	discover.MustHexID("3e05f80d931922509c7f65ab647b39efec56c95360fe293413368b8e5db7cfb385ab7f8a77acd6a8c39ae6a0d108e0723cc84c19c8797ecd0af1c5c6e94a541b"): struct{}{},
	discover.MustHexID("853667f088dbad823cfb9fc956ddf343cacf4cffc9d5de1e7227a2261bd2c5df356cb0a914f9404f21cc502d46a9d87bf90c5863286d07ebb0b06fd71f9c1192"): struct{}{},
	discover.MustHexID("b8d128a2b3a8eeb1c3cf4e471b34884c92f60f2d419149e8f93fc5fba18a36d6cae5cae0f07285630b5d43f5a0ba915e548dc4079dbd8a9d93a8f71ee7489804"): struct{}{},
	discover.MustHexID("c0dc97ee57ba202faf012ecb72bf30aebcd2cf7c161d7012017e0320e0db15925c107998bd833d61ec4c2689172d7e34a0371f4511773641e00814c2632b0e66"): struct{}{},
	discover.MustHexID("8f2dc504635ff5b394c3bf7debd885db3f4108772d8e51c26a3860a97d34295ec5c15cffe9fdccd88052408f3a84208718291b077e257898e603ec7b195d3684"): struct{}{},
	discover.MustHexID("1fd9fd7d9c31dad117384c7cc2f223a9e76f7aa81e30f38f030f24212be3fa20ca1f067d878a8ae97deb89b81efbd0542b3880cbd428b4ffae494fcd2c31834b"): struct{}{},
	discover.MustHexID("32b67023a9057c66ff9b9e0f9ca5a23a09e686a707c303695219ef0aa82ca2c7bd965da5cb1262232376584ed4402877fe786b58c3b1adc8fafbb3174b3f24f3"): struct{}{},
	discover.MustHexID("505ef930162c1b0736d7b7a44a52e0e179597bd3a43e13b057734abffb46978311ef10bf16d565aad1ed4f208493b030c173c7bcfbff989d7d4bf2c37925f621"): struct{}{},
	discover.MustHexID("73601a21e758e4052b43e68d9c09f2750dbc82b1736cce6a52258d42a18cde6c006d1eae39d44dda08987362d0ba8f8df43df05b3330281e866207117de25885"): struct{}{},
	discover.MustHexID("ade53a03334eda53f0bed8f3e1572a7d5faf59f40e464b683e9160946a738b0d72b333b93c07b4b75341de633784b4a2d00aacf78ce082949470046a4875fbee"): struct{}{},
	discover.MustHexID("32d622ecfc31deaae46f1249486f334a57476f9db66d930a4629440bdf0db3345d36c29cb8ce5ab58a410685b5b4e5f8648e35dc1a77c510326bc323ee473d78"): struct{}{},
	discover.MustHexID("f01b7b56c23a67376f0b1578feefed152952bbc3326fd2c412432a6f7e95ebe4d8adc231c8f3a4e3ab7d35c05278073766f89c5565a1f5267cd73a07c1d2fb04"): struct{}{},
	discover.MustHexID("370248c1d479da63476c52c742decc6b04b92908cfeccc53bf3a3367971f8de08a270a7e756f4714e99d0d10521426a4270335233ce5ebe8029828f35707c33a"): struct{}{},
	discover.MustHexID("8bc8734315acf2af4c92a458f077f1f8c96f0530fb43510c11361f1d6469631423206ef76cd879ade849ee15fbcaeb042e3721168614b4fad4eecd60a6aa3e94"): struct{}{},
	discover.MustHexID("e504fb1ed169c96d5abb63dced01109becb930349995d0ff0942b6a4b116bfeea6810c1e804ce42e240576036bd21cdca09fdcdd6cc90685badf802956a21931"): struct{}{},
	discover.MustHexID("239a70bf6909fe8f63bdb298799fe81d29c45d9a35952ae0cee48f77d79c678f9e9c4e2ba63c0691c0a5d5adcedc159b361bf0d20d6a3be534adf080e6d2205f"): struct{}{},
	discover.MustHexID("196e461fe3fea40d260bcef18eb0907fe61fb55322d3b446fe2dfa37127e1e6850435352c8a51635a755de9a25acb85a0efda56634e69be49a30d02908faecc7"): struct{}{},
	discover.MustHexID("88d14bc287a42fec8a490b90258ddfbc18e98a32f42f9dbe27a4369b8e64a41e22ca84aba764ef7d1326b5b9a8e7b54d6d10f23dc665ff7e7c40c04a5f0a3703"): struct{}{},
	discover.MustHexID("11daedba2520a87da234f48bdaa4373e536b85367b863c218e6638c558214708470260830c2f17feccf8997187c70c96469bcca9ea0f5b522aae6fadea8e9ddf"): struct{}{},
	discover.MustHexID("5c7c9dd985ddf54ddd865c4c26d0bfc291568f421100eb96fd6c69593bf226a4d18257f00234155d3e8a8b12a24d89525711f58e47e643dc1c198448dbce9dd1"): struct{}{},
	discover.MustHexID("3581081a17b98c9bb1e04fe0bb63d5bd0ce648aa2cf6214109923678b36783bd83be9c00f51fff1c08dd094196a600caff012ffb7bc6337f91cd903a3fd87283"): struct{}{},
	discover.MustHexID("fad2c7f917eb3057d85031eae8bbda52541b527dd1d24a25e7e9b40d7329570a85dc45ec61b189a9cc30047ae906a08dc375558828e1c76dc853ce99b42b91e4"): struct{}{},
	discover.MustHexID("4e44d7995f1b2b01bdbc65e5448f9fe3c7082f46f85fe980bccc3418836cec8ffab248dc0b54c0d8974e1327b8992643ef9db222ec9a44512a092d33b274bc1d"): struct{}{},
	discover.MustHexID("d03f770e38a5e6f6ef760be7bd79662c66a21bb3624525d12018daacadd3815c008b8ceaaf34d0646cbf8a537a00f57b217fca19f782898fd9cd36c70daefbff"): struct{}{},
	discover.MustHexID("14c2b76bd5945f8f77da071e0d46b253adf765551919ef3e7f80885a7a5506bb70e34926c546d27490b5a58cf01cc3f079296b0b97a01daa694dcb94c97548bc"): struct{}{},
	discover.MustHexID("19b8fb478a8502a25e461270122ece3135b15dc0de118264495bae30d39af81fd9134ed95364e6a39c3eebfba57fbffa7961a5158d3dac0f0da0313dac7af024"): struct{}{},
	discover.MustHexID("24b9584429dcb52ef6da575a0246d59158ff04f37118e9ce580f100e9da4a99064db252648f78497cf4b27f53eeaace7ca795ff75734e0e95386a5e3282f5fff"): struct{}{},
	discover.MustHexID("cd57e10afc9426739a5cbd299d1dd2e5da59bdef15e878b7125c23f9bd9d620aaec70e240e25bf1607e787636d58692576972e3d77a956401d22498856dcc57f"): struct{}{},
	discover.MustHexID("db98c59f4fd681bb8907ff0d4c67d1baf75ea364a2d31ef36ee53bc60ec3fa1d8123cf29a7990b028ba6a6bd9197c96217c9982175266816c700222c660c4d30"): struct{}{},
	discover.MustHexID("46e97c2f1774a2e505c8629d4e5d654f7faca48dfe7c23e4ebbbf86655feeae502096da83ab1d6d3cc47da5ff852bcb45490cceda4bbe7055f6d7cd91d5d8640"): struct{}{},
	discover.MustHexID("77861e48684ab9b35b6aafdbbb9028da0ff0bbe669f0485a8f387fa5ae83c2c9d54fabf3c027684b1371cad8f21b37981c053a5ff54eb8db2914f28ece651575"): struct{}{},
	discover.MustHexID("205e875111a1b5765ae07c1660f56f53ee4def882fea68a4db27a4b8f3a0b91a4c02729b0355c546d2c17116683b29aa94035205e354d15cd114f485d0778db1"): struct{}{},
	discover.MustHexID("ae75f869ed4d8bf87936ebeb1ff9bc8591b9dfc65bc59b978119e098f86b11241d6d39ec2f8f47f2087e6ead6902b82f74f272498024f830f6eded0d771f8f8d"): struct{}{},
	discover.MustHexID("9460fce5beea98e4d56c62a920bb041f45e48a5a7b96d12d02a16cbb20863be9c76491127533d9cefa5b4cec48ae6595b7ba347ef7dc8277cfb343eebde4646b"): struct{}{},
	discover.MustHexID("f6cb212b47105dbbb5b42e52492c330ba92b9a40706de88d7f3de02e1fc5507e43bc25db9b474a4b75ad4c2091d79f297e1fe5d3e09f8e4a12aba7c26db33717"): struct{}{},
	discover.MustHexID("2d35a84c4fc677fe2a19c43407d4cd387b0bbf90a5a3511794d7f752012e4090d8e7a0931ed540be41b73badd3c767c5de28195f3062c7aefba951bfd7a5c49e"): struct{}{},
	discover.MustHexID("db2bb5a8ca81f75c2b7551364acf827a86850bae5a0ce47a9eabecb2f5c4a00ccbeec8d31a6d271bc47a8c6ae6ca97a92fceb19ecdbffc8eb35c3666d6acce20"): struct{}{},
	discover.MustHexID("4922ef029b4bd93a1bca252d3ce4ed4d50192179e1c06e3ef452a18fe17d21e8c42eadf7d6aec47de2b90b07befa6c448d82f55599be6141c370dfac0b56f547"): struct{}{},
	discover.MustHexID("1d4b0709c6bc493af2086c3d38daf22ab64a0e992704a571af7479b0b58292d0bbab33552095082e3aa9297a8f338bd2679345e3f02bdefbc68c80993e0800e4"): struct{}{},
	discover.MustHexID("a6bbf49d6f1df8e3de68123802c182719fc2b17d40d655a0233a30f5131f35f58d48c87a462bdc58ec1a9014c340baebb9220d1b7c965c4f45e6cc7ca2ed6e81"): struct{}{},
	discover.MustHexID("e29327568233fa6c18f30f1ba7da3b416ac87b630a4848ac584cbb7044b2cd5b25e358c0d4c515d271c17ea145870db1dffdc278162eed77947ee45ff14d5a70"): struct{}{},
	discover.MustHexID("9839ced934daefca276eb8945c978e97774e9bfe8a101104b3202e3b38229fa4c71a06145366b4b89a135d89002ffd5bbbd12f9da33bc7c20026b290e45fa5c8"): struct{}{},
	discover.MustHexID("e1f6cd14a08821b893d2b0175d5f673ba162102cc795ca6095e276f560cc37e17a7a14dbfde674f43f9915d04edc9d76f7a64fb66e2b9dfada80ae2310621aec"): struct{}{},
	discover.MustHexID("0a72d3cc733ff102dda801c1c407c9febf8a01eb92b89dab101641e5a02b51d6ee7723aa3998519d9cef2d50ab7b18dd547899b9f9697d6757ada2fdfbbff134"): struct{}{},
	discover.MustHexID("88e70a87f6acc8edf3b381c02f3c3317392e458af688920bbfe04e3694979847e25d59fb7fe2c1d3487f1ae5a7876fbcefabe06f722dfa28a83f3ca4853c4254"): struct{}{},
	discover.MustHexID("fc5f7bc80f543b7fb2c3205ab05537b6bdc7248b398bda4ea863f48a13f3a8f7ac6aefa07dec946849cde152743d718b5628ab83fb12759e23572ad62dabcadb"): struct{}{},
	discover.MustHexID("929846fad7ecacfb21fb83ffd6a64e244f3cfc3e5a53b8f25786c2e2ea5327a6abe8c3691817bb35683cc1006fb0a6cd430e00cdfcad0dcb4afdef3e58f3c5cc"): struct{}{},
	discover.MustHexID("f2ec2830850a4e9dd48b358f908e1f22448cf5b0314750363acfd6a531edd2056237d39d182b92891a58e7d9862a43ee143a049167b7700914c41c726fad1399"): struct{}{},
	discover.MustHexID("9fe8313e32a7bc34e009a29b2a172ceded04eed6ce67f7c99680de038794c9168b4731927d45a9b82a04c5919dc768431f1a6bf4ed0424a216aab8b619a7bb7f"): struct{}{},
	discover.MustHexID("511ab9921b1ecf4bd8c76193a1c281f57a03190eae5418a23a7920a8064f89b9022bfae56d7fd2e740c2bc90c07e7aa78f201fa19c8b2b6e0dd15d8c97bec8c6"): struct{}{},
	discover.MustHexID("a98f15ecc908e6ce68b7fe29ea56c8c552f09b658352d0d0fb0fc2f08aa50b186cac9468a2742b9be9ec5d6ed0840d157afa88cb51ed2efacfd51cea26a7aa84"): struct{}{},
	discover.MustHexID("a2340b4acd4f7b743d7e8785e8ff297490b0e333f25cfe31d17df006f7e554553c6dc502f7c9f7b8798ab3ccf74624065a6cc20603842b1015793c0b37de9b15"): struct{}{},
	discover.MustHexID("c0e1be37b517f3bc7a946b2fe6b01caf815f4f6d52fadadd47b98aaf7cb19f974f7d5a7e05180f88442dad9db1b413e045df0d94ccda8eeb6873151f056ca0a1"): struct{}{},
	discover.MustHexID("854c6554239251092bc7de885b586042760a34584a686a4b360e171e7e27fab23ab5d6c8faa2519b1e5af6d9a2009e234b0f6261a0daf4a6987942a72c23783e"): struct{}{},
	discover.MustHexID("168fff55f6a28c3f10802e90cefb4daf58eb9a60c55eaaeda6588b20eece884ab9104da866633929e3c0ea413de2624b361939cf1513e1b67460686b7192f528"): struct{}{},
	discover.MustHexID("eb1a3c276744580b50d478e7910debcc7af7b3807e58ca23b7542a8265d996ebd5037ec9e6ba36ee8528aec0465bf375f32acf3538f75d79a11ab501bfff11b1"): struct{}{},
	discover.MustHexID("2527d6df6123e22e49ee6bad9dd3d579aeebd30893fc17aba8cad48c5414c78b82363aa5bd3fed1e118060300ddff1739077d4d9ec5f0940b10b347deeeb8971"): struct{}{},
	discover.MustHexID("6e409559d6b93b01e330400c8ca56a26ce979fef23edf4460e450975d0755dda5d44bcb4632871bffd2447d1de2b43f7639b8742985ab2cb52d984518c7cfefd"): struct{}{},
	discover.MustHexID("68ca2e833cdb044017d5f2bc35c2c61ecd5b340cac3384b3b33116c7caf951afdfef6fbf2542322de06a485beeb7b0e5c198b819c3dd991958eaa36c148bb356"): struct{}{},
	discover.MustHexID("7e2c09d2a2df357120facbc990568978845326a8ee4253d22d75e4576d8eb5e1cf00a411bc1a768dbe8d5d3368c864d4099fa1c509de8ecb950c71ac4180b204"): struct{}{},
	discover.MustHexID("5f2bf1ada8117f9fca7117c3d402f375d6b97296f2e40cac4eba9f3cf137d1d046468d11da8a670dab6ca468bf0ee770c4ab6246887fd78c85bf244bcecbd255"): struct{}{},
	discover.MustHexID("1d13161894d676fa124655820b92c722e04194ebc82bd2df5487dce684c21af80bd1e95230748d5aa65f2c28176fb231297f0cc5257a49051a7e33f3dd114b93"): struct{}{},
	discover.MustHexID("1dbe057f33d9748e1d396d624f4c2554f67742f18247e6be6c615c56c70a6e18a6604dd887fd1e9ffdf9708486fb76b711cb5d8e66ccc69d2cee09428832aa98"): struct{}{},
	discover.MustHexID("453009617066f547833c3c0bad5c4625546df74759e836ca34442228c564b7e8b100cd3b3b777b982681441da01c5532755dfde09970f2470b145b18e3583628"): struct{}{},
	discover.MustHexID("bf7e30583e0fa2ae326cfeca0fe1a45247c7a8463a23b84bdc33cc4eba41d4b581604cb9da4c4491518cf9c6b1fd31a1ed196b13cb41496555fb92e0a56482c6"): struct{}{},
	discover.MustHexID("d0f648262b54c16a0c4f82d1c2b08b638620537fc26321a47ace94b6691e1fa289fc4e06e1d94113b33ad1c33e3726655448f6f6a8bd1ea1c99a9f84da92c3d9"): struct{}{},
	discover.MustHexID("29c5d9c702b1dbd69bb5ae76aa396b0b4b60fd1f650717c01fe677494f3625ceb14df406624f46dccd5f788f890e4b85e51a529068877e5c5b772ef1c0c93d7f"): struct{}{},
	discover.MustHexID("218b20ee9cd80ea631c771996262a1d4b74b6db5c14321d55f42810a812b128c8c1e3bca0877e19e4ca12db6875d1922f11cfaeb9b2ec78da5157e523c7db8d3"): struct{}{},
	discover.MustHexID("a71e63927d7c96421448a9e79833d3c4a3d9de323b9d09c93a33da4ec2e7239efac29735e5298b21c23f7cfaacdec949060bd34c6afe04fc75b6501504058756"): struct{}{},
	discover.MustHexID("89baa83d0447aaabdf21f784532060bbcd0ca18742c13b09553c233e42225b36b0d52096d8e2a9db18a370a1b01099fa8eda654536c9847105ad8f6431ad5e55"): struct{}{},
	discover.MustHexID("d94b1914a4e187acd6b6038d99a48fea21cd8208d2aff92fdca28be9731772aefb5ef469d9f5b8a2221161be761337fa54f7ea8edb2f675ea22c9f278827fc39"): struct{}{},
	discover.MustHexID("bcb2be51df673ae15820d8ac083ab565c78408f79e372a5858e11fc23f1f711890e1402459fab088bd2beb95729fad5c6dbe681b3ebfc3af143b25afc14fcb59"): struct{}{},
	discover.MustHexID("3a30e0ca53a5ad43e92d6df887468d9f735b756c153cf5bfadbf70f5e3851cecbe0d88f48b1b26a9a9abe7cd2e0ba45758ab710be68b345525a19f6782798ef6"): struct{}{},
	discover.MustHexID("5af7343f5858a9acf635f9381dd58c2727ed495e9b34f899df4086731f4c5d73c27ea8b88529e655a552ab6babe65fc0d2bb578f59aa347db30ae34cf62223cc"): struct{}{},
	discover.MustHexID("cbf8b3b0bd51bc4cd3edb55a5cbd297e7dafcdfa17ffe82d1f4d2d5f1be1782030f6ac26b9f3910337acb8fcc356f79233851a35836f4197b3d4c8038c27cbfa"): struct{}{},
	discover.MustHexID("985f2fd92897063958dd06ba3f511c2c3bc1309cb6a5034a04e2e50e44d24708aa22235f305419a39693068966786201083d6a0ed12ac51dd3bf4ec8bb0cfaa2"): struct{}{},
	discover.MustHexID("38b694dd05452a810b3701ae8478c9b6ea68a1383f9695a7fb466801fd5ba57496f5575b5037cf55bf43f04ce33f34a095add9b3fefd25c3d1f1e7d28531f6f2"): struct{}{},
	discover.MustHexID("c537d6e1393521608e1dea636cf2f0e095d1e73e1b6151bfd0159b563fc9106f9c80118fce0b3d211cdcad2df0359a1596c3983e79bc26b8d19c8325e1869aaa"): struct{}{},
	discover.MustHexID("949e3793be67761380ae1536832cd5d310a555ded62d8f049636c6292815d27340188a679d5bbcf66d48fcbe459c05cf61cf475d168ec059b7e4a5ee22b64ba7"): struct{}{},
	discover.MustHexID("6e78c448dad030abfced12cba9aca8b386031f3b944622d370de8018c4252b6e9144167894da91d2123030a43b0287b5df2636adecb4c9507c44454b60a2140a"): struct{}{},
	discover.MustHexID("a7fe1e01eb794334490bfa5f4b286b6b781d3cd3f7ba02cdb7ff5d6f0322d4dc97449781e61d6a8c67576219dbaa4d854b014acbeceeed6b777618173da47986"): struct{}{},
	discover.MustHexID("0525a2e651701ed46764b4648fe88a64c70ecbf69d6d897adf48f549772d9d67840288a1396f82d6d6fc7b3ca6adf0778995935e445a3e7b682cc3c41ccd83ff"): struct{}{},
	discover.MustHexID("c353b813621ba95ec3382b668c94bc2d3b7c97f205ad21840d5c723f709f237d2f3d63cf10e48185725285eaa3785d3b0f9254e9cb62b031697249f213f7de9c"): struct{}{},
	discover.MustHexID("23d1bee8744c6d63bbd02c2b0a1685b215762dc934d1fe5e3b88432e04edfaaaac78f8e5c52d46adee9707c14380971ccd30a5c4e544f8056f793ac328927939"): struct{}{},
	discover.MustHexID("2dcb6b2368ea69ece2cc71fda719e76ba63e20eaac290df752f4ac88ef7efac9ddaa6c7ed024888546462ca028b269247d05a39eeded148032a5e9288565f7de"): struct{}{},
	discover.MustHexID("0e630f5cc6d9f99ef24fcedc9273ce82b4ad51cbf17eda2b88dba7393abf466b18164a4a60acb4709e78b55f2d29c065f10303038986ce2ea97db73e32918d42"): struct{}{},
	discover.MustHexID("85fd7813f1872885722929c1ee41c7046aa938fc4b158935032d1a14aa129cabc008dc3e1b8aceb5f500b58256867b2d59f0f97d42215b0ec9f5ed2ffc3642ed"): struct{}{},
	discover.MustHexID("c67c30ced9155115d3fe25c61a586a0597362b8c9cf5d63bc64f3ada2b7943e918a6d77f3e8c2e3f998f13262718b547b6cfe51bb5808d0e1e82c6019c223eb7"): struct{}{},
	discover.MustHexID("b5b54f8a5e9f50bc1b87a402458d49c8d94e42eba3cffcbe39e3c9cb21c3a86530eb53a37d277c7771ab5979ec12d3a351e2bd4e082f6497384dd940a3e5e234"): struct{}{},
	discover.MustHexID("6a517f4e7d4227c74ce4310d8dcf581399e05eb5048b4994955f3582512d0c425e49ef3e269e497998460cb68b20246907193607c4850a1f7cb4919d46b559a4"): struct{}{},
	discover.MustHexID("a36aa96b38c65db7a5269b9250b06efccf27eff06f4a7cdf169eac25f60d03f215de5c095a77acc5a288448425f6789347a93a4499dd4dfd33a6c51a4d8bddd0"): struct{}{},
	discover.MustHexID("887db62779c1166038f50a90ae94f226058f7319ceb6375e7c2ed8a40e52063d9363a1edaecebd30f1b6772bab9ebf38f0bb0f1d5bad07a0df47df2400845823"): struct{}{},
	discover.MustHexID("8fbc190c20707087eac8dbcbffa28a489e5a6c555c3b9542caa60307af5f08a2faa5db23b8ede13691d52cad2a1075432f9a754673f34f60e274d66742475a3a"): struct{}{},
	discover.MustHexID("f1efed4e853d00ff3f1be65fd497bc8f0a3d5f66b285069c9190653567e1838ab635b88940d3ce786747af549a1a5bf9b7173e9dc3a3aea9f10363613581a9e0"): struct{}{},
	discover.MustHexID("f26619faa0c072b3af9fe2e47ecdb0b72e28efd7047f65e63d8719a0158450d61acfaf93d253a4dbeb53122ca4cb5fb276ab301690e956320f147e387624a142"): struct{}{},
	discover.MustHexID("25f4caa5654e4b132de6ada9a1b8a6c474536566f481a9dc015c2ff75b1dd79e99c81a420ae5243a68b68eb62b468f2d1c5349eef9f123e6ed4f186efd8d1d09"): struct{}{},
	discover.MustHexID("6eeddc7bea67b08cb2bff68b4c5c6cda0b4234779a21b78b244203acab504b801ca299dcbf5b50c4650e502196b05f4fd6f3582fd0973528fd3047c36ab2198d"): struct{}{},
	discover.MustHexID("0178bdb413b9acc82b80fc1d31ef4e66065708c2ef8758bf5121100049273dd9b41fb16c64bc747d732575b959a45400db757a417ac0d6dffe220ce68827eee5"): struct{}{},
	discover.MustHexID("d839188ce67070dbc2bdd2a93a8a9abacf1ee89432ab18feb8de668bd595312f2483be56197d15e977770f59560cee5c03301a6c3cb90415951d0de26c2bbde8"): struct{}{},
	discover.MustHexID("a55b37d341ec17a34081e45cf7289ca9d6ab1152f78ef911adb4a2e4f2700dd935175b6193183c9bc8dd3bbe352f5b5013cbefe692cb046c70dae551683be27d"): struct{}{},
	discover.MustHexID("b19fc39a15f20ca81342d96cf4ddf46ab871ed1cdf4d592d16d735d4962d3729b691d0f6732e85cabf9a70ab5b57470d18b368afec00dd3110424bca3a7ae468"): struct{}{},
	discover.MustHexID("e2053e04f95afa5c8378677de62212edb972e21b40421786c53de57141853ce870481a80b68a449903479751da114a30d27a568812e178b919bd1e9b2c82f92e"): struct{}{},
	discover.MustHexID("f0f4a28ee2e9d9a4178bb862217928fee242574b7609c625bca6699deeee7af6235a1e3986686e5eff8de6647cb7796a594de2fc07ca1cdcae84b72fc053b588"): struct{}{},
	discover.MustHexID("fe9e15b816920aca3ad387d8ae4bcf5b8ec600c44c339b6b20994f4580580a668369d99d204a31567522f04a7f52db29f0d58234e1ff540d2593ff214c5aadd2"): struct{}{},
	discover.MustHexID("3d4f4dc9968874853803e2614f615583e7fd8839bf231eac3e759e08e452c679d7ed86db76a1778e19a0bdefc62b5973ee6fdfea12c7b35ac3b580f5d05dfaf6"): struct{}{},
	discover.MustHexID("b93710f214ca408ad09e8ec51faf21dd1b3a0fa7c42088ffe27225ac0df80cc5fc6e94e98f914537f6ba1199f4239a0ab0deed00fac34699a7aa54d32a68797d"): struct{}{},
	discover.MustHexID("28c82d6ef53b64cb9830cb9abf933b45d4370039f526b06ed499e53bdee7dd2d020bad5b1f716f351c93dd1b25aaa7de257071daa92ba0834fa0c8d26151d1c3"): struct{}{},
	discover.MustHexID("3015687818ccdeac78b66f2b15dfa4534da10abf7aab713267431114b4ede6fb6bc27f7d0cf8e792c356da6e916d81c80b65205679d2046beb9a69cfe8c374f7"): struct{}{},
	discover.MustHexID("68db72723b00124d63bb8e8e86e2ae23fd2c7d3799ce227a468676a51a52dbab63abdf415575f51ca0992dec20032c82e6bf0f03abed4b03a8ee33653957119b"): struct{}{},
	discover.MustHexID("0cf921182dc578a1eccf55562bc3437343de72f464912e6f296144e9dd82d51937c8aea89364fa0d793052284f1ff22b53464ed424e210d2cc3f711964cfffac"): struct{}{},
	discover.MustHexID("5984b8e581167b7d41cc472f56de19d0503bae3f2153bc9918214eb43cff484eba5ee982eb385eed339ff2b7a64dbb500750f53286618732e246ead228597b72"): struct{}{},
	discover.MustHexID("b70c76572bb4805cddc15aa97c28df7a573d15c43d41494d7a015e5100a9f7046ec9207ab60c1e20e9601f49c11ffe6d4fe2636529ebc1aa75fcf30cb799159d"): struct{}{},
	discover.MustHexID("b6dc4bbf9692b2420b6e6eb5ae87043723d512bfbce0288985dd8fe14fdd443aa2eb04261e3b05e6f419282c77ac24dc53b6b75432a86ef4cc1dc5c17a3a98a6"): struct{}{},
	discover.MustHexID("6887b34c8a76be71ef0c8d743a07ae5154b5550e5e31ec43050b1627c5936ee99bbaedd23d87e2cedc52b02b8a46d717b6f7e34702ae6cc86cf6fe4f6b7596a7"): struct{}{},
	discover.MustHexID("25ffee4be31237ead04e038747b3ba25e3eeb059016e6908518b0cb8ecb7648119d56775da5f1dc23699a13ed996692a644b48cd9b261909ac8bb00be5f700f5"): struct{}{},
	discover.MustHexID("7dc831a1668827507528992d72471a397dc0e8a65ade1a91ba197c76a89ed7336eb05bf098fc11083648ad48b2b0f6d4903cc17be221115316112b81ba72d427"): struct{}{},
	discover.MustHexID("982f73d9983cb4bdd92e94d594242c8ab44b056820d571d019f0b28f4f32a1bc86f7f25242dd911474e4a22780dabce8c7bffd57d459297feac05bba533d9c92"): struct{}{},
	discover.MustHexID("fad7c5a4d4fef670df9c80ece25a80f262088a790c41c0f974faff1f28dffaf50ea505924dc9f4a498b4c21b01c73330cb39103cde54d01007ba623f28da5f72"): struct{}{},
	discover.MustHexID("f5b318918dca88a73dabb71fa270acc3dea0676d6ca8d46966ef72b2c7e09b4874822a691073fb4a80af951ddd5105dc2c380ceaaed6449bc2256aec61a56cc0"): struct{}{},
	discover.MustHexID("bf9e74aae3b4dabdd417c4edb58867fba06233e6a1f7e8013960f56ea60efebfd2f334ae3f822ae42de7de6851de9ce9f188686c766c797e7d9a8a9245eac866"): struct{}{},
	discover.MustHexID("bf2ad47cb3010ae6dc23f0ddeb84fb73abf46d46e98c136ca8864d64bf962b8976391440e03d121beb92cfdd2c587daf65641fc53bdc636fd5c8e267d4bd8a51"): struct{}{},
	discover.MustHexID("3ab721f2f51e2d572d18ef182ae05f54ce01360119f2669788b40f5f1338357345344716a1170874b469b11281c1e3cc8a7f31c88b07527cacb55c8b66ea05a0"): struct{}{},
	discover.MustHexID("86c63b28fedc3364b25c4864dc46e158f3948a87abaccc1fea686a7405ca17ed90a33dd73bdc0a0585d19de5855f3281d06b8f4f594b1cabd97efcc7039522fd"): struct{}{},
	discover.MustHexID("80d92c92a0f97c48edcb6e35d66302a43a54323955a4f2d6ce6d641533e42b7497e9a59c32adee84293152cac35608164325395d4fe0050a737e8d6e1026c8b3"): struct{}{},
	discover.MustHexID("e3aa1da40ffa6625cf2601bc8759f06530c3f31afa23551039b659edceac98e09a37e7e79fca9d8d91fca525b46d5248453bd02ace937973f45fcab32d147f9b"): struct{}{},
	discover.MustHexID("580ca2b64f4ea1ae5f266b62e475b456d168d84debed84185719809c8b0e35f8c03271e097b5a80cfde67c29ce11bd3a76cd95befceab5a050565240250b74e5"): struct{}{},
	discover.MustHexID("7e15e8599518c651dcf57944b0a72dacd84bd4db76b2973eda1d5bce352fd8f518c1b1b843dc30a207143445a008a65341898bf1e45c197e3a4bcaec65f38daa"): struct{}{},
	discover.MustHexID("826f465ab4ed4c74bd811b4507724aeb9f7ce5e3d752c718fa4d464315c494f5720907a9ad046495becb0271e5340ee5db77ff5b060238d4008c3778a3e85aa9"): struct{}{},
	discover.MustHexID("efbd089457980cb916e6cd196591c673720882930f30a2bd0c11c1a86478c47bd719ad565de36f717989c9722ba5aa98534062daa300cb66ef6469cf52e2a955"): struct{}{},
	discover.MustHexID("70c598ebdd53b20482fba31876699193cbbd47b3dc5f0c687fcb5e7062a01b9e0862ecc65e3620634117a971b5cc6e8c0abcc125bc63274456d9d5b0c39ad72f"): struct{}{},
	discover.MustHexID("735f009b6985f0439144d0718cb1901eda97b9401a70ff08f420f3aeace5e226d5af74d24f3dce12beba7116acbabe540671e6bd5e3a5855d12ab8098a94b122"): struct{}{},
	discover.MustHexID("e3c84ebcd351010df4524cad73348483d1f23011eaa76435dfc0486b0182787c4a57e326a648c95cdcc792c92df15f0bdfe6782fd336a262daa64568a2c01dfe"): struct{}{},
	discover.MustHexID("59782ce52eafe81a2b4228457e70a21d6c64158a0682eecb5bd2a7b0c10299bf1f986e6878e89b1d13236a2f6c943745ba1b68bb6edeea80badf83d055b4d949"): struct{}{},
	discover.MustHexID("b6cb09bf53ca3133b932f1a5cb473ef7ddc67925bc4579891c85bce88a564135ecafb9619c9170f9acc5aa304c1e65389437cf06f74c4bae88a580c7e4a5d05e"): struct{}{},
	discover.MustHexID("8717968ef857018ec6205b5354945c5e9e8ff058417f27b03b3af975274a1304d54a280490ed7b4ba9819b3849c1277c380f2f5326555e290d88fbf9e5a6fb05"): struct{}{},
	discover.MustHexID("a376c3950130c665e4d100d7e9c37ab21d04910b9c1d7c59ce1daf8ea2257b0f0c238eedaa63bc029efe37a36fd5322982ab306dd3dd036009e23a8b0e5f3838"): struct{}{},
	discover.MustHexID("b4f3148a9f1a05eedf109774606943f21208f4f296ba0469bb45667a0bc1f5b66993034dfaab88b2f388d9d02c66f95e05bc3af139d1e92ff62aedb4b8c318ea"): struct{}{},
	discover.MustHexID("4b0e852717ff8aeb8d0ef3dd02b3de4f204aaa7988ca365f1c9657a5282e34fe8a1b1ef02cc6b254b7e0f6fca7671229578da874ed5f8da2f917d09500985570"): struct{}{},
	discover.MustHexID("fe5ad9842a74b1006bbde108c715ac25c2f437f44748aebc1f21c99fbcbeb630c719a9035e8d929327be83ffa0aae0de3f8adc190ba4d238d6bad3aa5b7e94a1"): struct{}{},
}

// Config holds Server options.
type Config struct {
	// This field must be set to a valid secp256k1 private key.
	PrivateKey *ecdsa.PrivateKey `toml:"-"`

	// BlsPublicKey is a BLS public key.
	BlsPublicKey bls.PublicKey `toml:"-"`

	// chainId identifies the current chain and is used for replay protection
	ChainID *big.Int `toml:"-"`

	// MaxPeers is the maximum number of peers that can be
	// connected. It must be greater than zero.
	MaxPeers int

	// MaxConsensusPeers is the maximum number of consensus peers that can be
	// connected. It must be greater than zero.
	MaxConsensusPeers int

	// MaxPendingPeers is the maximum number of peers that can be pending in the
	// handshake phase, counted separately for inbound and outbound connections.
	// Zero defaults to preset values.
	MaxPendingPeers int `toml:",omitempty"`

	// DialRatio controls the ratio of inbound to dialed connections.
	// Example: a DialRatio of 2 allows 1/2 of connections to be dialed.
	// Setting DialRatio to zero defaults it to 3.
	DialRatio int `toml:",omitempty"`

	// NoDiscovery can be used to disable the peer discovery mechanism.
	// Disabling is useful for protocol debugging (manual topology).
	NoDiscovery bool

	// DiscoveryV5 specifies whether the new topic-discovery based V5 discovery
	// protocol should be started or not.
	DiscoveryV5 bool `toml:",omitempty"`

	// Name sets the node name of this server.
	// Use common.MakeName to create a name that follows existing conventions.
	Name string `toml:"-"`

	// BootstrapNodes are used to establish connectivity
	// with the rest of the network.
	BootstrapNodes []*discover.Node

	// BootstrapNodesV5 are used to establish connectivity
	// with the rest of the network using the V5 discovery
	// protocol.
	BootstrapNodesV5 []*discv5.Node `toml:",omitempty"`

	// Static nodes are used as pre-configured connections which are always
	// maintained and re-connected on disconnects.
	StaticNodes []*discover.Node `json:"-"`

	// Trusted nodes are used as pre-configured connections which are always
	// allowed to connect, even above the peer limit.
	TrustedNodes []*discover.Node

	// Connectivity can be restricted to certain IP networks.
	// If this option is set to a non-nil value, only hosts which match one of the
	// IP networks contained in the list are considered.
	NetRestrict *netutil.Netlist `toml:",omitempty"`

	// NodeDatabase is the path to the database containing the previously seen
	// live nodes in the network.
	NodeDatabase string `toml:",omitempty"`

	// Protocols should contain the protocols supported
	// by the server. Matching protocols are launched for
	// each peer.
	Protocols []Protocol `toml:"-"`

	// If ListenAddr is set to a non-nil address, the server
	// will listen for incoming connections.
	//
	// If the port is zero, the operating system will pick a port. The
	// ListenAddr field will be updated with the actual address when
	// the server is started.
	ListenAddr string

	// If set to a non-nil value, the given NAT port mapper
	// is used to make the listening port available to the
	// Internet.
	NAT nat.Interface `toml:",omitempty"`

	// If Dialer is set to a non-nil value, the given Dialer
	// is used to dial outbound peer connections.
	Dialer NodeDialer `toml:"-"`

	// If NoDial is true, the server will not dial any peers.
	NoDial bool `toml:",omitempty"`

	// If EnableMsgEvents is set then the server will emit PeerEvents
	// whenever a message is sent to or received from a peer
	EnableMsgEvents bool

	// Logger is a custom logger to use with the p2p.Server.
	Logger log.Logger `toml:",omitempty"`
}

// Server manages all peer connections.
type Server struct {
	// Config fields may not be modified while the server is running.
	Config

	// Hooks for testing. These are useful because we can inhibit
	// the whole protocol stack.
	newTransport func(net.Conn) transport
	newPeerHook  func(*Peer)

	lock    sync.Mutex // protects running
	running bool

	ntab         discoverTable
	listener     net.Listener
	ourHandshake *protoHandshake
	lastLookup   time.Time
	DiscV5       *discv5.Network

	// These are for Peers, PeerCount (and nothing else).
	peerOp     chan peerOpFunc
	peerOpDone chan struct{}

	quit            chan struct{}
	addstatic       chan *discover.Node
	removestatic    chan *discover.Node
	addconsensus    chan *discover.Node
	removeconsensus chan *discover.Node
	addtrusted      chan *discover.Node
	removetrusted   chan *discover.Node
	posthandshake   chan *conn
	addpeer         chan *conn
	delpeer         chan peerDrop
	loopWG          sync.WaitGroup // loop, listenLoop
	peerFeed        event.Feed
	log             log.Logger

	eventMux  *event.TypeMux
	consensus bool
}

type peerOpFunc func(map[discover.NodeID]*Peer)

type peerDrop struct {
	*Peer
	err       error
	requested bool // true if signaled by the peer
}

type connFlag int32

const (
	dynDialedConn connFlag = 1 << iota
	staticDialedConn
	inboundConn
	trustedConn
	consensusDialedConn
)

// conn wraps a network connection with information gathered
// during the two handshakes.
type conn struct {
	fd net.Conn
	transport
	flags connFlag
	cont  chan error      // The run loop uses cont to signal errors to SetupConn.
	id    discover.NodeID // valid after the encryption handshake
	caps  []Cap           // valid after the protocol handshake
	name  string          // valid after the protocol handshake
}

type transport interface {
	// The two handshakes.
	doEncHandshake(prv *ecdsa.PrivateKey, dialDest *discover.Node) (discover.NodeID, error)
	doProtoHandshake(our *protoHandshake) (*protoHandshake, error)
	// The MsgReadWriter can only be used after the encryption
	// handshake has completed. The code uses conn.id to track this
	// by setting it to a non-nil value after the encryption handshake.
	MsgReadWriter
	// transports must provide Close because we use MsgPipe in some of
	// the tests. Closing the actual network connection doesn't do
	// anything in those tests because NsgPipe doesn't use it.
	close(err error)
}

func (c *conn) String() string {
	s := c.flags.String()
	if (c.id != discover.NodeID{}) {
		s += " " + c.id.String()
	}
	s += " " + c.fd.RemoteAddr().String()
	return s
}

func (f connFlag) String() string {
	s := ""
	if f&trustedConn != 0 {
		s += "-trusted"
	}
	if f&dynDialedConn != 0 {
		s += "-dyndial"
	}
	if f&staticDialedConn != 0 {
		s += "-staticdial"
	}
	if f&inboundConn != 0 {
		s += "-inbound"
	}
	if f&consensusDialedConn != 0 {
		s += "-consensusdial"
	}
	if s != "" {
		s = s[1:]
	}
	return s
}

func (c *conn) is(f connFlag) bool {
	flags := connFlag(atomic.LoadInt32((*int32)(&c.flags)))
	return flags&f != 0
}

func (c *conn) set(f connFlag, val bool) {
	for {
		oldFlags := connFlag(atomic.LoadInt32((*int32)(&c.flags)))
		flags := oldFlags
		if val {
			flags |= f
		} else {
			flags &= ^f
		}
		if atomic.CompareAndSwapInt32((*int32)(&c.flags), int32(oldFlags), int32(flags)) {
			return
		}
	}
}

// Peers returns all connected peers.
func (srv *Server) Peers() []*Peer {
	var ps []*Peer
	select {
	// Note: We'd love to put this function into a variable but
	// that seems to cause a weird compiler error in some
	// environments.
	case srv.peerOp <- func(peers map[discover.NodeID]*Peer) {
		for _, p := range peers {
			ps = append(ps, p)
		}
	}:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return ps
}

// PeerCount returns the number of connected peers.
func (srv *Server) PeerCount() int {
	var count int
	select {
	case srv.peerOp <- func(ps map[discover.NodeID]*Peer) { count = len(ps) }:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return count
}

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *Server) AddPeer(node *discover.Node) {
	select {
	case srv.addstatic <- node:
	case <-srv.quit:
	}
}

// RemovePeer disconnects from the given node
func (srv *Server) RemovePeer(node *discover.Node) {
	select {
	case srv.removestatic <- node:
	case <-srv.quit:
	}
}

// Determine whether the node is in the whitelist.
func (srv *Server) IsAllowNode(nodeID discover.NodeID) bool {
	if srv.ChainID.Cmp(params.AlayaChainConfig.ChainID) == 0 {
		//if len(AllowNodesMap) == 0 {
		//	nodesString := params.AllowNodes
		//	for _, node := range nodesString {
		//		tmp := discover.MustHexID(node)
		//		AllowNodesMap[tmp] = struct{}{}
		//	}
		//}
		_, ok := AllowNodesMap[nodeID]
		return ok
	}
	return true
}

// AddConsensusPeer connects to the given consensus node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *Server) AddConsensusPeer(node *discover.Node) {
	select {
	case srv.addconsensus <- node:
	case <-srv.quit:
	}
}

// RemoveConsensusPeer disconnects from the given consensus node
func (srv *Server) RemoveConsensusPeer(node *discover.Node) {
	select {
	case srv.removeconsensus <- node:
	case <-srv.quit:
	}
}

// AddTrustedPeer adds the given node to a reserved whitelist which allows the
// node to always connect, even if the slot are full.
func (srv *Server) AddTrustedPeer(node *discover.Node) {
	select {
	case srv.addtrusted <- node:
	case <-srv.quit:
	}
}

// RemoveTrustedPeer removes the given node from the trusted peer set.
func (srv *Server) RemoveTrustedPeer(node *discover.Node) {
	select {
	case srv.removetrusted <- node:
	case <-srv.quit:
	}
}

// SubscribePeers subscribes the given channel to peer events
func (srv *Server) SubscribeEvents(ch chan *PeerEvent) event.Subscription {
	return srv.peerFeed.Subscribe(ch)
}

// Self returns the local node's endpoint information.
func (srv *Server) Self() *discover.Node {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if !srv.running {
		return &discover.Node{IP: net.ParseIP("0.0.0.0")}
	}
	return srv.makeSelf(srv.listener, srv.ntab)
}

func (srv *Server) makeSelf(listener net.Listener, ntab discoverTable) *discover.Node {
	// If the server's not running, return an empty node.
	// If the node is running but discovery is off, manually assemble the node infos.
	if ntab == nil {
		// Inbound connections disabled, use zero address.
		if listener == nil {
			return &discover.Node{IP: net.ParseIP("0.0.0.0"), ID: discover.PubkeyID(&srv.PrivateKey.PublicKey)}
		}
		// Otherwise inject the listener address too
		addr := listener.Addr().(*net.TCPAddr)
		return &discover.Node{
			ID:  discover.PubkeyID(&srv.PrivateKey.PublicKey),
			IP:  addr.IP,
			TCP: uint16(addr.Port),
		}
	}
	// Otherwise return the discovery node.
	return ntab.Self()
}

// Stop terminates the server and all active peer connections.
// It blocks until all active connections have been closed.
func (srv *Server) Stop() {
	srv.lock.Lock()
	if !srv.running {
		srv.lock.Unlock()
		return
	}
	srv.running = false
	if srv.listener != nil {
		// this unblocks listener Accept
		srv.listener.Close()
	}
	close(srv.quit)
	srv.lock.Unlock()
	srv.loopWG.Wait()
}

// sharedUDPConn implements a shared connection. Write sends messages to the underlying connection while read returns
// messages that were found unprocessable and sent to the unhandled channel by the primary listener.
type sharedUDPConn struct {
	*net.UDPConn
	unhandled chan discover.ReadPacket
}

// ReadFromUDP implements discv5.conn
func (s *sharedUDPConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	packet, ok := <-s.unhandled
	if !ok {
		return 0, nil, errors.New("Connection was closed")
	}
	l := len(packet.Data)
	if l > len(b) {
		l = len(b)
	}
	copy(b[:l], packet.Data[:l])
	return l, packet.Addr, nil
}

// Close implements discv5.conn
func (s *sharedUDPConn) Close() error {
	return nil
}

// Start starts running the server.
// Servers can not be re-used after stopping.
func (srv *Server) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	srv.log = srv.Config.Logger
	if srv.log == nil {
		srv.log = log.New()
	}
	srv.log.Info("Starting P2P networking")

	// static fields
	if srv.PrivateKey == nil {
		return errors.New("Server.PrivateKey must be set to a non-nil key")
	}
	if srv.newTransport == nil {
		srv.newTransport = newRLPX
	}
	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: defaultDialTimeout}}
	}
	srv.quit = make(chan struct{})
	srv.addpeer = make(chan *conn)
	srv.delpeer = make(chan peerDrop)
	srv.posthandshake = make(chan *conn)
	srv.addstatic = make(chan *discover.Node)
	srv.removestatic = make(chan *discover.Node)
	srv.addconsensus = make(chan *discover.Node)
	srv.removeconsensus = make(chan *discover.Node)
	srv.addtrusted = make(chan *discover.Node)
	srv.removetrusted = make(chan *discover.Node)
	srv.peerOp = make(chan peerOpFunc)
	srv.peerOpDone = make(chan struct{})

	var (
		conn      *net.UDPConn
		sconn     *sharedUDPConn
		realaddr  *net.UDPAddr
		unhandled chan discover.ReadPacket
	)

	if !srv.NoDiscovery || srv.DiscoveryV5 {
		addr, err := net.ResolveUDPAddr("udp", srv.ListenAddr)
		if err != nil {
			return err
		}
		conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			return err
		}
		realaddr = conn.LocalAddr().(*net.UDPAddr)
		if srv.NAT != nil {
			if !realaddr.IP.IsLoopback() {
				go nat.Map(srv.NAT, srv.quit, "udp", realaddr.Port, realaddr.Port, "ethereum discovery")
			}
			// TODO: react to external IP changes over time.
			if ext, err := srv.NAT.ExternalIP(); err == nil {
				realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
			}
		}
	}

	if !srv.NoDiscovery && srv.DiscoveryV5 {
		unhandled = make(chan discover.ReadPacket, 100)
		sconn = &sharedUDPConn{conn, unhandled}
	}

	// node table
	if !srv.NoDiscovery {
		cfg := discover.Config{
			PrivateKey:   srv.PrivateKey,
			ChainID:      srv.ChainID,
			AnnounceAddr: realaddr,
			NodeDBPath:   srv.NodeDatabase,
			NetRestrict:  srv.NetRestrict,
			Bootnodes:    srv.BootstrapNodes,
			Unhandled:    unhandled,
		}
		ntab, err := discover.ListenUDP(conn, cfg)
		if err != nil {
			return err
		}
		srv.ntab = ntab
	}

	if srv.DiscoveryV5 {
		var (
			ntab *discv5.Network
			err  error
		)
		if sconn != nil {
			ntab, err = discv5.ListenUDP(srv.PrivateKey, sconn, realaddr, "", srv.NetRestrict) //srv.NodeDatabase)
		} else {
			ntab, err = discv5.ListenUDP(srv.PrivateKey, conn, realaddr, "", srv.NetRestrict) //srv.NodeDatabase)
		}
		if err != nil {
			return err
		}
		if err := ntab.SetFallbackNodes(srv.BootstrapNodesV5); err != nil {
			return err
		}
		srv.DiscV5 = ntab
	}

	dynPeers := srv.maxDialedConns()
	dialer := newDialState(srv.StaticNodes, srv.BootstrapNodes, srv.ntab, dynPeers, srv.NetRestrict, srv.MaxConsensusPeers)

	// handshake
	srv.ourHandshake = &protoHandshake{Version: baseProtocolVersion, Name: srv.Name, ID: discover.PubkeyID(&srv.PrivateKey.PublicKey)}
	for _, p := range srv.Protocols {
		srv.ourHandshake.Caps = append(srv.ourHandshake.Caps, p.cap())
	}
	// listen/dial
	if srv.ListenAddr != "" {
		if err := srv.startListening(); err != nil {
			return err
		}
	}
	if srv.NoDial && srv.ListenAddr == "" {
		srv.log.Warn("P2P server will be useless, neither dialing nor listening")
	}

	srv.loopWG.Add(1)
	go srv.run(dialer)
	srv.running = true
	return nil
}

func (srv *Server) startListening() error {
	// Launch the TCP listener.
	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		return err
	}
	laddr := listener.Addr().(*net.TCPAddr)
	srv.ListenAddr = laddr.String()
	srv.listener = listener
	srv.loopWG.Add(1)
	go srv.listenLoop()
	// Map the TCP listening port if NAT is configured.
	if !laddr.IP.IsLoopback() && srv.NAT != nil {
		srv.loopWG.Add(1)
		go func() {
			nat.Map(srv.NAT, srv.quit, "tcp", laddr.Port, laddr.Port, "ethereum p2p")
			srv.loopWG.Done()
		}()
	}
	return nil
}

type dialer interface {
	newTasks(running int, peers map[discover.NodeID]*Peer, now time.Time) []task
	taskDone(task, time.Time)
	addStatic(*discover.Node)
	removeStatic(*discover.Node)
	addConsensus(*discover.Node)
	removeConsensus(*discover.Node)
	removeConsensusFromQueue(*discover.Node)
	initRemoveConsensusPeerFn(removeConsensusPeerFn removeConsensusPeerFn)
}

func (srv *Server) run(dialstate dialer) {
	defer srv.loopWG.Done()
	var (
		peers          = make(map[discover.NodeID]*Peer)
		inboundCount   = 0
		trusted        = make(map[discover.NodeID]bool, len(srv.TrustedNodes))
		consensusNodes = make(map[discover.NodeID]bool, 0)
		taskdone       = make(chan task, maxActiveDialTasks)
		runningTasks   []task
		queuedTasks    []task // tasks that can't run yet
	)
	// Put trusted nodes into a map to speed up checks.
	// Trusted peers are loaded on startup or added via AddTrustedPeer RPC.
	for _, n := range srv.TrustedNodes {
		trusted[n.ID] = true
	}

	// removes t from runningTasks
	delTask := func(t task) {
		for i := range runningTasks {
			if runningTasks[i] == t {
				runningTasks = append(runningTasks[:i], runningTasks[i+1:]...)
				break
			}
		}
	}
	// starts until max number of active tasks is satisfied
	startTasks := func(ts []task) (rest []task) {
		i := 0
		for ; len(runningTasks) < maxActiveDialTasks && i < len(ts); i++ {
			t := ts[i]
			srv.log.Trace("New dial task", "task", t)
			go func() { t.Do(srv); taskdone <- t }()
			runningTasks = append(runningTasks, t)
		}
		return ts[i:]
	}
	scheduleTasks := func() {
		// Start from queue first.
		queuedTasks = append(queuedTasks[:0], startTasks(queuedTasks)...)
		// Query dialer for new tasks and start as many as possible now.
		if len(runningTasks) < maxActiveDialTasks {
			nt := dialstate.newTasks(len(runningTasks)+len(queuedTasks), peers, time.Now())
			queuedTasks = append(queuedTasks, startTasks(nt)...)
		}
	}
	dialstateRemoveConsensusPeerFn := func(node *discover.Node) {
		srv.log.Trace("Removing consensus node from dialstate", "node", node)
		dialstate.removeConsensusFromQueue(node)
		if p, ok := peers[node.ID]; ok {
			p.Disconnect(DiscRequested)
		}
	}
	dialstate.initRemoveConsensusPeerFn(dialstateRemoveConsensusPeerFn)

running:
	for {
		scheduleTasks()

		select {
		case <-srv.quit:
			// The server was stopped. Run the cleanup logic.
			break running
		case n := <-srv.addstatic:
			// This channel is used by AddPeer to add to the
			// ephemeral static peer list. Add it to the dialer,
			// it will keep the node connected.
			srv.log.Trace("Adding static node", "node", n)
			dialstate.addStatic(n)
		case n := <-srv.removestatic:
			// This channel is used by RemovePeer to send a
			// disconnect request to a peer and begin the
			// stop keeping the node connected.
			srv.log.Trace("Removing static node", "node", n)
			dialstate.removeStatic(n)
			if p, ok := peers[n.ID]; ok {
				p.Disconnect(DiscRequested)
			}
		case n := <-srv.addconsensus:
			// This channel is used by AddConsensusNode to add an enode
			// to the consensus node set.
			srv.log.Trace("Adding consensus node", "node", n)
			if n.ID == srv.ourHandshake.ID {
				srv.log.Debug("We are become an consensus node")
				srv.consensus = true
			} else {
				dialstate.addConsensus(n)
			}
			consensusNodes[n.ID] = true
			if p, ok := peers[n.ID]; ok {
				srv.log.Debug("Add consensus flag", "peer", n.ID)
				p.rw.set(consensusDialedConn, true)
			}
		case n := <-srv.removeconsensus:
			// This channel is used by RemoveConsensusNode to remove an enode
			// from the consensus node set.
			srv.log.Trace("Removing consensus node", "node", n)
			if n.ID == srv.ourHandshake.ID {
				srv.log.Debug("We are not an consensus node")
				srv.consensus = false
			}
			dialstate.removeConsensus(n)
			if _, ok := consensusNodes[n.ID]; ok {
				delete(consensusNodes, n.ID)
			}
			if p, ok := peers[n.ID]; ok {
				p.rw.set(consensusDialedConn, false)
				if !p.rw.is(staticDialedConn | trustedConn | inboundConn) {
					p.rw.set(dynDialedConn, true)
				}
				srv.log.Debug("Remove consensus flag", "peer", n.ID, "consensus", srv.consensus)
				if len(peers) > srv.MaxPeers && !p.rw.is(staticDialedConn|trustedConn) {
					srv.log.Debug("Disconnect non-consensus node", "peer", n.ID, "flags", p.rw.flags, "peers", len(peers), "consensus", srv.consensus)
					p.Disconnect(DiscRequested)
				}
			}
		case n := <-srv.addtrusted:
			// This channel is used by AddTrustedPeer to add an enode
			// to the trusted node set.
			srv.log.Trace("Adding trusted node", "node", n)
			trusted[n.ID] = true
			// Mark any already-connected peer as trusted
			if p, ok := peers[n.ID]; ok {
				p.rw.set(trustedConn, true)
			}
		case n := <-srv.removetrusted:
			// This channel is used by RemoveTrustedPeer to remove an enode
			// from the trusted node set.
			srv.log.Trace("Removing trusted node", "node", n)
			delete(trusted, n.ID)
			// Unmark any already-connected peer as trusted
			if p, ok := peers[n.ID]; ok {
				p.rw.set(trustedConn, false)
			}
		case op := <-srv.peerOp:
			// This channel is used by Peers and PeerCount.
			op(peers)
			srv.peerOpDone <- struct{}{}
		case t := <-taskdone:
			// A task got done. Tell dialstate about it so it
			// can update its state and remove it from the active
			// tasks list.
			srv.log.Trace("Dial task done", "task", t)
			dialstate.taskDone(t, time.Now())
			delTask(t)
		case c := <-srv.posthandshake:
			// A connection has passed the encryption handshake so
			// the remote identity is known (but hasn't been verified yet).
			if trusted[c.id] {
				// Ensure that the trusted flag is set before checking against MaxPeers.
				c.flags |= trustedConn
			}

			if consensusNodes[c.id] {
				c.flags |= consensusDialedConn
			}

			// TODO: track in-progress inbound node IDs (pre-Peer) to avoid dialing them.
			select {
			case c.cont <- srv.encHandshakeChecks(peers, inboundCount, c):
			case <-srv.quit:
				break running
			}
		case c := <-srv.addpeer:
			// At this point the connection is past the protocol handshake.
			// Its capabilities are known and the remote identity is verified.
			err := srv.protoHandshakeChecks(peers, inboundCount, c)
			if err == nil {
				// The handshakes are done and it passed all checks.
				p := newPeer(c, srv.Protocols)
				// If message events are enabled, pass the peerFeed
				// to the peer
				if srv.EnableMsgEvents {
					p.events = &srv.peerFeed
				}
				name := truncateName(c.name)
				srv.log.Debug("Adding p2p peer", "name", name, "id", p.ID(), "addr", c.fd.RemoteAddr(), "flags", c.flags, "peers", len(peers)+1)
				go srv.runPeer(p)
				peers[c.id] = p
				if p.Inbound() {
					inboundCount++
				}
			}
			// The dialer logic relies on the assumption that
			// dial tasks complete after the peer has been added or
			// discarded. Unblock the task last.
			select {
			case c.cont <- err:
			case <-srv.quit:
				break running
			}
		case pd := <-srv.delpeer:
			// A peer disconnected.
			d := common.PrettyDuration(mclock.Now() - pd.created)
			pd.log.Debug("Removing p2p peer", "duration", d, "peers", len(peers)-1, "req", pd.requested, "err", pd.err)
			delete(peers, pd.ID())
			if pd.Inbound() {
				inboundCount--
			}
		}
	}

	srv.log.Trace("P2P networking is spinning down")

	// Terminate discovery. If there is a running lookup it will terminate soon.
	if srv.ntab != nil {
		srv.ntab.Close()
	}
	if srv.DiscV5 != nil {
		srv.DiscV5.Close()
	}
	// Disconnect all peers.
	for _, p := range peers {
		p.Disconnect(DiscQuitting)
	}
	// Wait for peers to shut down. Pending connections and tasks are
	// not handled here and will terminate soon-ish because srv.quit
	// is closed.
	for len(peers) > 0 {
		p := <-srv.delpeer
		p.log.Trace("<-delpeer (spindown)", "remainingTasks", len(runningTasks))
		delete(peers, p.ID())
	}
}

func (srv *Server) protoHandshakeChecks(peers map[discover.NodeID]*Peer, inboundCount int, c *conn) error {
	// Drop connections with no matching protocols.
	if len(srv.Protocols) > 0 && countMatchingProtocols(srv.Protocols, c.caps) == 0 {
		return DiscUselessPeer
	}
	// Repeat the encryption handshake checks because the
	// peer set might have changed between the handshakes.
	return srv.encHandshakeChecks(peers, inboundCount, c)
}

func (srv *Server) encHandshakeChecks(peers map[discover.NodeID]*Peer, inboundCount int, c *conn) error {
	// Disconnect over limit non-consensus node.
	if srv.consensus && len(peers) >= srv.MaxPeers && c.is(consensusDialedConn) {
		for _, p := range peers {
			if p.rw.is(inboundConn|dynDialedConn) && !p.rw.is(trustedConn|staticDialedConn|consensusDialedConn) {
				log.Debug("Disconnect over limit connection", "peer", p.ID(), "flags", p.rw.flags, "peers", len(peers))
				p.Disconnect(DiscRequested)
				break
			}
		}
	}

	switch {
	case !c.is(trustedConn|staticDialedConn|consensusDialedConn) && len(peers) >= srv.MaxPeers:
		return DiscTooManyPeers
	case !c.is(trustedConn|consensusDialedConn) && c.is(inboundConn) && inboundCount >= srv.maxInboundConns():
		return DiscTooManyPeers
	case peers[c.id] != nil:
		return DiscAlreadyConnected
	case c.id == srv.Self().ID:
		return DiscSelf
	default:
		return nil
	}
}

func (srv *Server) maxInboundConns() int {
	return srv.MaxPeers - srv.maxDialedConns()
}

func (srv *Server) maxDialedConns() int {
	if srv.NoDiscovery || srv.NoDial {
		return 0
	}
	r := srv.DialRatio
	if r == 0 {
		r = defaultDialRatio
	}
	return srv.MaxPeers / r
}

type tempError interface {
	Temporary() bool
}

// listenLoop runs in its own goroutine and accepts
// inbound connections.
func (srv *Server) listenLoop() {
	defer srv.loopWG.Done()
	srv.log.Info("RLPx listener up", "self", srv.makeSelf(srv.listener, srv.ntab))

	tokens := defaultMaxPendingPeers
	if srv.MaxPendingPeers > 0 {
		tokens = srv.MaxPendingPeers
	}
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}

	for {
		// Wait for a handshake slot before accepting.
		<-slots

		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = srv.listener.Accept()
			if tempErr, ok := err.(tempError); ok && tempErr.Temporary() {
				srv.log.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				srv.log.Debug("Read error", "err", err)
				return
			}
			break
		}

		// Reject connections that do not match NetRestrict.
		if srv.NetRestrict != nil {
			if tcp, ok := fd.RemoteAddr().(*net.TCPAddr); ok && !srv.NetRestrict.Contains(tcp.IP) {
				srv.log.Debug("Rejected conn (not whitelisted in NetRestrict)", "addr", fd.RemoteAddr())
				fd.Close()
				slots <- struct{}{}
				continue
			}
		}

		fd = newMeteredConn(fd, true)
		srv.log.Trace("Accepted connection", "addr", fd.RemoteAddr())
		go func() {
			srv.SetupConn(fd, inboundConn, nil)
			slots <- struct{}{}
		}()
	}
}

// SetupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *Server) SetupConn(fd net.Conn, flags connFlag, dialDest *discover.Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}
	c := &conn{fd: fd, transport: srv.newTransport(fd), flags: flags, cont: make(chan error)}
	err := srv.setupConn(c, flags, dialDest)
	if err != nil {
		c.close(err)
		srv.log.Trace("Setting up connection failed", "id", c.id, "err", err)
	}
	return err
}

func (srv *Server) setupConn(c *conn, flags connFlag, dialDest *discover.Node) error {
	// Prevent leftover pending conns from entering the handshake.
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}
	// Run the encryption handshake.
	var err error
	if c.id, err = c.doEncHandshake(srv.PrivateKey, dialDest); err != nil {
		srv.log.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
		return err
	}
	clog := srv.log.New("id", c.id, "addr", c.fd.RemoteAddr(), "conn", c.flags)
	// For dialed connections, check that the remote public key matches.
	if dialDest != nil && c.id != dialDest.ID {
		clog.Trace("Dialed identity mismatch", "want", c, dialDest.ID)
		return DiscUnexpectedIdentity
	}
	err = srv.checkpoint(c, srv.posthandshake)
	if err != nil {
		clog.Trace("Rejected peer before protocol handshake", "err", err)
		return err
	}
	// Run the protocol handshake
	phs, err := c.doProtoHandshake(srv.ourHandshake)
	if err != nil {
		clog.Trace("Failed proto handshake", "err", err)
		return err
	}
	if phs.ID != c.id {
		clog.Trace("Wrong devp2p handshake identity", "err", phs.ID)
		return DiscUnexpectedIdentity
	}
	// P2p protocol version and whitelist verification
	if phs.Version < baseProtocolVersion && !srv.IsAllowNode(phs.ID) {
		clog.Error("Low version of p2p protocol version", "err", phs.ID)
		return DiscIncompatibleVersion
	}

	c.caps, c.name = phs.Caps, phs.Name
	err = srv.checkpoint(c, srv.addpeer)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by run.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

func truncateName(s string) string {
	if len(s) > 20 {
		return s[:20] + "..."
	}
	return s
}

// checkpoint sends the conn to run, which performs the
// post-handshake checks for the stage (posthandshake, addpeer).
func (srv *Server) checkpoint(c *conn, stage chan<- *conn) error {
	select {
	case stage <- c:
	case <-srv.quit:
		return errServerStopped
	}
	select {
	case err := <-c.cont:
		return err
	case <-srv.quit:
		return errServerStopped
	}
}

// runPeer runs in its own goroutine for each peer.
// it waits until the Peer logic returns and removes
// the peer.
func (srv *Server) runPeer(p *Peer) {
	if srv.newPeerHook != nil {
		srv.newPeerHook(p)
	}

	// broadcast peer add
	srv.peerFeed.Send(&PeerEvent{
		Type: PeerEventTypeAdd,
		Peer: p.ID(),
	})

	// run the protocol
	remoteRequested, err := p.run()

	// broadcast peer drop
	srv.peerFeed.Send(&PeerEvent{
		Type:  PeerEventTypeDrop,
		Peer:  p.ID(),
		Error: err.Error(),
	})

	// Note: run waits for existing peers to be sent on srv.delpeer
	// before returning, so this send should not select on srv.quit.
	srv.delpeer <- peerDrop{p, err, remoteRequested}
}

// NodeInfo represents a short summary of the information known about the host.
type NodeInfo struct {
	ID     string `json:"id"`        // Unique node identifier (also the encryption key)
	Name   string `json:"name"`      // Name of the node, including client type, version, OS, custom data
	BlsPub string `json:"blsPubKey"` // BLS public key
	Enode  string `json:"enode"`     // Enode URL for adding this peer from remote peers
	IP     string `json:"ip"`        // IP address of the node
	Ports  struct {
		Discovery int `json:"discovery"` // UDP listening port for discovery protocol
		Listener  int `json:"listener"`  // TCP listening port for RLPx
	} `json:"ports"`
	ListenAddr string                 `json:"listenAddr"`
	Protocols  map[string]interface{} `json:"protocols"`
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (srv *Server) NodeInfo() *NodeInfo {
	node := srv.Self()

	// Gather and assemble the generic node infos
	info := &NodeInfo{
		Name:       srv.Name,
		Enode:      node.String(),
		ID:         node.ID.String(),
		IP:         node.IP.String(),
		ListenAddr: srv.ListenAddr,
		Protocols:  make(map[string]interface{}),
	}
	info.Ports.Discovery = int(node.UDP)
	info.Ports.Listener = int(node.TCP)

	blskey, _ := srv.BlsPublicKey.MarshalText()
	info.BlsPub = string(blskey)

	// Gather all the running protocol infos (only once per protocol type)
	for _, proto := range srv.Protocols {
		if _, ok := info.Protocols[proto.Name]; !ok {
			nodeInfo := interface{}("unknown")
			if query := proto.NodeInfo; query != nil {
				nodeInfo = proto.NodeInfo()
			}
			info.Protocols[proto.Name] = nodeInfo
		}
	}
	return info
}

// PeersInfo returns an array of metadata objects describing connected peers.
func (srv *Server) PeersInfo() []*PeerInfo {
	// Gather all the generic and sub-protocol specific infos
	infos := make([]*PeerInfo, 0, srv.PeerCount())
	for _, peer := range srv.Peers() {
		if peer != nil {
			infos = append(infos, peer.Info())
		}
	}
	// Sort the result array alphabetically by node identifier
	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[i].ID > infos[j].ID {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}
	return infos
}

func (srv *Server) StartWatching(eventMux *event.TypeMux) {
	srv.eventMux = eventMux
	go srv.watching()
}

func (srv *Server) watching() {
	events := srv.eventMux.Subscribe(cbfttypes.AddValidatorEvent{}, cbfttypes.RemoveValidatorEvent{})
	defer events.Unsubscribe()

	for {
		select {
		case ev := <-events.Chan():
			if ev == nil {
				continue
			}

			switch ev.Data.(type) {
			case cbfttypes.AddValidatorEvent:
				addEv, ok := ev.Data.(cbfttypes.AddValidatorEvent)
				if !ok {
					log.Error("Received add validator event type error")
					continue
				}
				log.Trace("Received AddValidatorEvent", "nodeID", addEv.NodeID.String())
				node := discover.NewNode(addEv.NodeID, nil, 0, 0)
				srv.AddConsensusPeer(node)
			case cbfttypes.RemoveValidatorEvent:
				removeEv, ok := ev.Data.(cbfttypes.RemoveValidatorEvent)
				if !ok {
					log.Error("Received remove validator event type error")
					continue
				}
				log.Trace("Received RemoveValidatorEvent", "nodeID", removeEv.NodeID.String())
				node := discover.NewNode(removeEv.NodeID, nil, 0, 0)
				srv.RemoveConsensusPeer(node)
			default:
				log.Error("Received unexcepted event")
			}

		case <-srv.quit:
			return
		}
	}
}

type mockTransport struct {
	id discover.NodeID
	*rlpx

	closeErr error
}

func newMockTransport(id discover.NodeID, fd net.Conn) transport {
	wrapped := newRLPX(fd).(*rlpx)
	wrapped.rw = newRLPXFrameRW(fd, secrets{
		MAC:        zero16,
		AES:        zero16,
		IngressMAC: sha3.NewKeccak256(),
		EgressMAC:  sha3.NewKeccak256(),
	})
	return &mockTransport{id: id, rlpx: wrapped}
}

func (c *mockTransport) doEncHandshake(prv *ecdsa.PrivateKey, dialDest *discover.Node) (discover.NodeID, error) {
	return c.id, nil
}

func (c *mockTransport) doProtoHandshake(our *protoHandshake) (*protoHandshake, error) {
	return &protoHandshake{ID: c.id, Name: "test"}, nil
}

func (c *mockTransport) close(err error) {
	c.rlpx.fd.Close()
	c.closeErr = err
}

func randomID() (id discover.NodeID) {
	for i := range id {
		id[i] = byte(rand.Intn(255))
	}
	return id
}
