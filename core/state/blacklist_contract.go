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
// "0x18f6338136a6dc661af751b9ee5d194cba2e6a60"
// "0x07b65aacc07446a547640ff9672b70065bfa3337"
// "0xaf601b5acaa9d4e1afea9e692a59930a2e936831"


// JUNGE
//
// "0x976f3f063681dd29bbe50eafd861f17b576085f7"
// "0xaae9b8d53050db35ead5cd3508c6e836f3af715b"
// "0xf3c7c0ac0f893048b068bbfcc0260cef5d4d37bf"
// "0x01add8ca3ad86c7336736318796fd6ef3f05df70"
// "0x8077bd265da7d18c552c8eac7320e8172c133196"
// "0x5d349a62345387dec164de900fcedff416bdb394"
// "0xdbcefc73194253b28d99410997938bf6c00233b2"
// "0xa0be6bd120cd094a70232973313be36c07ccfc23"
// "0xaa67e189495968b2b5fea7225710f7bef19c12e9"
// "0xc7226f6faa4d531c7cb63842f8f5b0c66322ae61"
//

var badContracts  = map[common.Address]struct{} {
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
	common.HexToAddress("0x18f6338136a6dc661af751b9ee5d194cba2e6a60"): {},
	common.HexToAddress("0x07b65aacc07446a547640ff9672b70065bfa3337"): {},
	common.HexToAddress("0xaf601b5acaa9d4e1afea9e692a59930a2e936831"): {},
	common.HexToAddress("0x976f3f063681dd29bbe50eafd861f17b576085f7"): {},
	common.HexToAddress("0xaae9b8d53050db35ead5cd3508c6e836f3af715b"): {},
	common.HexToAddress("0xf3c7c0ac0f893048b068bbfcc0260cef5d4d37bf"): {},
	common.HexToAddress("0x01add8ca3ad86c7336736318796fd6ef3f05df70"): {},
	common.HexToAddress("0x8077bd265da7d18c552c8eac7320e8172c133196"): {},
	common.HexToAddress("0x5d349a62345387dec164de900fcedff416bdb394"): {},
	common.HexToAddress("0xdbcefc73194253b28d99410997938bf6c00233b2"): {},
	common.HexToAddress("0xa0be6bd120cd094a70232973313be36c07ccfc23"): {},
	common.HexToAddress("0xaa67e189495968b2b5fea7225710f7bef19c12e9"): {},
	common.HexToAddress("0xc7226f6faa4d531c7cb63842f8f5b0c66322ae61"): {},

}


func IsBadContract (addr common.Address) bool {
	if _, ok := badContracts[addr]; ok {
		return true
	}
	return false
}

