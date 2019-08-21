package xcom

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/log"

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
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	version, err := hex.DecodeString("0x05af3bbd099562e520ddb824199182dcd8249bc91274afbcce4be24bd0fbf8c259cd403738923722163fa7493ea8ef8725d52f3d1bb3fd2713592ac135d0f85200")
	if err != nil {
		log.Error("Decode hex String", "err", err)
		return
	}

	if !chandler.IsSignedByNodeID(66048, version, nodeID) {
		t.Error("verify sign error")
	} else {
		t.Error("verify sign OK")
	}

}
