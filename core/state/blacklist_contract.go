package state

import "github.com/PlatONnetwork/PlatON-Go/common"

// BAD CONTRACT ADDRESS
//
// "0x0f2a1795f7605a96910c6c782ea5d3d291fd77fc"
// "0x18253041ce1d238f42f685ff1714153ea9c97699"
// "0x187b3a7d5f790a30f59703338eba42ee11e584fa"
// "0x1e4419a4a0c96bab21004c258586c9e172c98ed6"
// "0x32bec384344c2dc1ea794a4e149c1b74dd8467ef"
// "0x3800f5390e5921059a1ae817e12c224813cdd33a"
// "0x471b4f5e00bf612766b1ead6df6e658244e6c179"
// "0xab782161cf50b8282afd717d27b5a99bc80f909b"
// "0xb3dce06660b0d1b6c89e217a6ec61a074029e5d3"
// "0xc1f4f1fa20461564e031c6e189018adc11270156"
// "0xdd26d10d54e3860a62d268bb877952292a6fcc48"
// "0xfbe13f9f86a7bb272f0d6479beb2b0f4b4114ed8"
// "0xfd783cbf2603b1ddf0740512aac47dc954a13897"

// QICHUAN
//
// 0xe0b2acfed581cffe318245dc5096870e73b2da15
// 0xb5ce170a4a728c6f62fa7e7227869a9d2b9f0f0e
// 0x533a923ffd5b7c83b48490a3ad939f75eb7fce62
//

// JUNGE
//
// 0x124071c61b83f71770798e7b757d14e497c558a5
// 0x7e11f9d863fb5982753eccd8c92207e172754075
// 0x13d8949cc9e4dc133f0d8070a54e9fca3d29f736
// 0xb04dd1a5dd2737af07d60854e0329b34d6f00505
// 0xbd891449a2403df312572e9a40f161547819dd71
// 0xde3185054e2ac47f036851ffab40bcf69e2075cf
// 0xa07c2a0eda390df0afbd8bf3b617b8c8b1192e0f
// 0xfd3af117c0fa61b92ee3aeaa04f754366cf46788
// 0x7437cc28c1a56450a1d7a57eeb096750f3835ebe
// 0x629183e4f3cd82bcd2f4278eb78da0844ad0d97c
// 0x6c25bbf05cb8952cba6f17d596f46b6fc286780e
// 0x1664a35d981b660f977f9d2f9f86f1c8ad1a5299
// 0x31faa0c561cbb81422cf05dabf166fb12d815948
//

var badContracts = map[common.Address]struct{}{
	//common.HexToAddress(""): {},

	common.HexToAddress("0x0f2a1795f7605a96910c6c782ea5d3d291fd77fc"): {},
	common.HexToAddress("0x18253041ce1d238f42f685ff1714153ea9c97699"): {},
	common.HexToAddress("0x187b3a7d5f790a30f59703338eba42ee11e584fa"): {},
	common.HexToAddress("0x1e4419a4a0c96bab21004c258586c9e172c98ed6"): {},
	common.HexToAddress("0x32bec384344c2dc1ea794a4e149c1b74dd8467ef"): {},
	common.HexToAddress("0x3800f5390e5921059a1ae817e12c224813cdd33a"): {},
	common.HexToAddress("0x471b4f5e00bf612766b1ead6df6e658244e6c179"): {},
	common.HexToAddress("0xab782161cf50b8282afd717d27b5a99bc80f909b"): {},
	common.HexToAddress("0xb3dce06660b0d1b6c89e217a6ec61a074029e5d3"): {},
	common.HexToAddress("0xc1f4f1fa20461564e031c6e189018adc11270156"): {},
	common.HexToAddress("0xdd26d10d54e3860a62d268bb877952292a6fcc48"): {},
	common.HexToAddress("0xfbe13f9f86a7bb272f0d6479beb2b0f4b4114ed8"): {},
	common.HexToAddress("0xfd783cbf2603b1ddf0740512aac47dc954a13897"): {},

	common.HexToAddress("0xe0b2acfed581cffe318245dc5096870e73b2da15"): {},
	common.HexToAddress("0xb5ce170a4a728c6f62fa7e7227869a9d2b9f0f0e"): {},
	common.HexToAddress("0x533a923ffd5b7c83b48490a3ad939f75eb7fce62"): {},
	common.HexToAddress("0x124071c61b83f71770798e7b757d14e497c558a5"): {},
	common.HexToAddress("0x7e11f9d863fb5982753eccd8c92207e172754075"): {},
	common.HexToAddress("0x13d8949cc9e4dc133f0d8070a54e9fca3d29f736"): {},
	common.HexToAddress("0xb04dd1a5dd2737af07d60854e0329b34d6f00505"): {},
	common.HexToAddress("0xbd891449a2403df312572e9a40f161547819dd71"): {},
	common.HexToAddress("0xde3185054e2ac47f036851ffab40bcf69e2075cf"): {},
	common.HexToAddress("0xa07c2a0eda390df0afbd8bf3b617b8c8b1192e0f"): {},
	common.HexToAddress("0xfd3af117c0fa61b92ee3aeaa04f754366cf46788"): {},
	common.HexToAddress("0x7437cc28c1a56450a1d7a57eeb096750f3835ebe"): {},
	common.HexToAddress("0x629183e4f3cd82bcd2f4278eb78da0844ad0d97c"): {},
	common.HexToAddress("0x6c25bbf05cb8952cba6f17d596f46b6fc286780e"): {},
	common.HexToAddress("0x1664a35d981b660f977f9d2f9f86f1c8ad1a5299"): {},
	common.HexToAddress("0x31faa0c561cbb81422cf05dabf166fb12d815948"): {},
}

func IsBadContract(addr common.Address) bool {
	if _, ok := badContracts[addr]; ok {
		return true
	}
	return false
}
