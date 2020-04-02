package state

import "github.com/PlatONnetwork/PlatON-Go/common"

// BAD CONTRACT ADDRESS
var badContracts = map[common.Address]struct{}{
	//common.HexToAddress(""): {},

	common.HexToAddress("0x0f2a1795f7605a96910c6c782ea5d3d291fd77fc"): {},
	common.HexToAddress("0x18253041ce1d238f42f685ff1714153ea9c97699"): {},
	common.HexToAddress("0x187b3a7d5f790a30f59703338eba42ee11e584fa"): {},
	common.HexToAddress("0x1e4419a4a0c96bab21004c258586c9e172c98ed6"): {},
	common.HexToAddress("0x3800f5390e5921059a1ae817e12c224813cdd33a"): {},
	common.HexToAddress("0x471b4f5e00bf612766b1ead6df6e658244e6c179"): {},
	common.HexToAddress("0xab782161cf50b8282afd717d27b5a99bc80f909b"): {},
	common.HexToAddress("0xb3dce06660b0d1b6c89e217a6ec61a074029e5d3"): {},
	common.HexToAddress("0xc1f4f1fa20461564e031c6e189018adc11270156"): {},
	common.HexToAddress("0xdd26d10d54e3860a62d268bb877952292a6fcc48"): {},
	common.HexToAddress("0xfbe13f9f86a7bb272f0d6479beb2b0f4b4114ed8"): {},
	common.HexToAddress("0xfd783cbf2603b1ddf0740512aac47dc954a13897"): {},
}

func IsBadContract(addr common.Address) bool {
	if _, ok := badContracts[addr]; ok {
		return true
	}
	return false
}
