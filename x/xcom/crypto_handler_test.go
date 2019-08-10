package xcom

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	priKey = crypto.HexMustToECDSA("d30b490011d2a08053d46506ae533ff96f2cf6a37f73be740f52ad24243c4958")
	nodeID = discover.MustHexID("a20aef0b2c6baeaa34be2848e7dfc04c899b5985adf6fa0e98b38f754f2bb0c47974506a8de13f2a2ae97c08bcb12b438b3dcbf237b7be58f6d6d8beb36dd235")
)

func initChandlerHandler() {
	chandler = GetCryptoHandler()
	chandler.SetPrivateKey(priKey)
}

func TestCryptoHandler_IsSignedByNodeID(t *testing.T) {
	initChandlerHandler()
	version := uint32(2<<16 | 0<<8 | 0)
	sig, err := chandler.Sign(version)

	if err != nil {
		t.Fatal("Sign error")
	} else {
		if !chandler.IsSignedByNodeID(version, sig, nodeID) {
			t.Fatal("verify sign error")
		}
	}
}
