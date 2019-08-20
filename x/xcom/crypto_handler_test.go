package xcom

import (
	"encoding/hex"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	priKey = crypto.HexMustToECDSA("76711a880d0b2fc40167428005aa80bdeb66ada7a82d3e9c78d93201022161e2")
	nodeID = discover.MustHexID("4b6083b5d2fa4638690e54e3ea37771f42776c044c76fd021931c476dc480492264ffaacaf59259438c16e404366ace3ce2fbbf19d230a8417a04ddc2f8be3c6")
)

func initChandlerHandler() {
	chandler = GetCryptoHandler()
	chandler.SetPrivateKey(priKey)
}

func TestCryptoHandler_IsSignedByNodeID(t *testing.T) {
	initChandlerHandler()
	version := uint32(1<<16 | 1<<8 | 0)
	sig := chandler.MustSign(version)

	versionSign := common.VersionSign{}
	versionSign.SetBytes(sig)

	t.Log("...", "version", version, "sig", hex.EncodeToString(sig))

	if !chandler.IsSignedByNodeID(version, versionSign.Bytes(), nodeID) {
		t.Fatal("verify sign error")
	}
}

func Test_Decode(t *testing.T) {
	initChandlerHandler()

	/*	sig := chandler.MustSign(1792)
		t.Fatal("sign", hex.EncodeToString(sig))
	*/
	bytes, _ := hex.DecodeString("155fe2050c65ff4633499f9d81acf3a0f185a6110e6ed1459f7c5ac95925fb284fdc3c2f299ea71ab176f41003383caf345e2ded72e4e5dc568f5d01e982d1cd01")

	if !chandler.IsSignedByNodeID(common.Uint32ToBytes(1792), bytes, nodeID) {
		t.Fatal("verify sign error")
	} else {
		t.Fatal("verify sign OK")
	}

}
