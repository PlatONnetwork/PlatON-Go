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
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	sig := chandler.MustSign(uint32(1792))

	log.Error("Decode hex String", "sig", hex.EncodeToString(sig), "src", "c6c027a49c04afb7daecdaaa03590a374746d65e000b8dce4542afdc985106fe4ef9477cb0a697340097c6a786b59a4c090075a592a51a337b9b8c299cc8c6d401")
	version, err := hex.DecodeString("c6c027a49c04afb7daecdaaa03590a374746d65e000b8dce4542afdc985106fe4ef9477cb0a697340097c6a786b59a4c090075a592a51a337b9b8c299cc8c6d401")
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
