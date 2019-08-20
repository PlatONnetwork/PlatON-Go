package xcom

import (
	"encoding/hex"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	priKey = crypto.HexMustToECDSA("8c56e4a0d8bb1f82b94231d535c499fdcbf6e06221acf669d5a964f5bb974903")
	nodeID = discover.MustHexID("0x3a06953a2d5d45b29167bef58208f1287225bdd2591260af29ae1300aeed362e9b548369dfc1659abbef403c9b3b07a8a194040e966acd6e5b6d55aa2df7c1d8")
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

	/*sig := chandler.MustSign(1792)
	t.Fatal("sign", hex.EncodeToString(sig))
	*/
	bytes, _ := hex.DecodeString("7c03cb6bc1c6103b979c7842cced4cdec5abfe8506b03ac9a714d1c31344f29e1185af720beae23a023f8ed02a9eb6eefaec0a381dda63c425d2af7a0c7ca7e501")

	if !chandler.IsSignedByNodeID(common.Uint32ToBytes(2048), bytes, nodeID) {
		t.Fatal("verify sign error")
	} else {
		t.Fatal("verify sign OK")
	}

}
