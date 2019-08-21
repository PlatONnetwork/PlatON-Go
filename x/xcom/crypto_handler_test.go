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
	priKey = crypto.HexMustToECDSA("8e1477549bea04b97ea15911e2e9b3041b7a9921f80bd6ddbe4c2b080473de22")
	nodeID = discover.MustHexID("3e7864716b671c4de0dc2d7fd86215e0dcb8419e66430a770294eb2f37b714a07b6a3493055bb2d733dee9bfcc995e1c8e7885f338a69bf6c28930f3cf341819")
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
	sig, err := chandler.Sign(uint32(1792))
	if err != nil {
		log.Error("sign err", "err", err)
	}
	log.Error("Decode hex String", "sig", hex.EncodeToString(sig), "src", "b195655dd28594ead36cd9291a9b09b29630e188a3dff96ed7af145167ae86e648698f8b2dcc89a6f73018b7a69549b97bda6c61d10fbbd46d8a70b867b3be2b00")
	version, err := hex.DecodeString("b195655dd28594ead36cd9291a9b09b29630e188a3dff96ed7af145167ae86e648698f8b2dcc89a6f73018b7a69549b97bda6c61d10fbbd46d8a70b867b3be2b00")
	if err != nil {
		log.Error("Decode hex String", "err", err)
		return
	}
	if !chandler.IsSignedByNodeID(uint32(1792), version, nodeID) {
		t.Error("verify sign error")
	} else {
		t.Error("verify sign OK")
	}

}
