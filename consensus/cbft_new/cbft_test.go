package cbft_new

import (
	"Platon-go/common/hexutil"
	"Platon-go/crypto"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {
	keys, err := hexutil.Decode("0x8b54398b67e656dcab213c1b5886845963a9ab0671786eefaf6e241ee9c8074f")
	if err != nil {
		log.Info("error", "err", err)
	}

	privateKey, err := crypto.ToECDSA(keys)

	if err != nil {
		log.Info("error", "err", err)
	}
	publicKey := privateKey.PublicKey
	nodeID := discover.PubkeyID(&publicKey)

	fmt.Println(nodeID.String())

	sign, _ := hexutil.Decode("0x005866959f620eb627a17d8ad34b6e6db27334dd67774336d90d3a36f44e5c623df7b73452c85b505aa0bcbd9c693ded5466c6c9a1ada83b5654502ae1fb97c500")
	hash, _ := hexutil.Decode("0x1267fc0161987579e9b5a84fb5f1c783ac53aceee77e1bcff75d167aab76881c")

	pubKey, _ := crypto.Ecrecover(hash, sign)
	fmt.Println(hexutil.Encode(pubKey[1:]))
}

func TestSync(t *testing.T) {
	m := make(map[int]int)
	go func() {
		for {
			_ = m[1]
		}
	}()
	go func() {
		for {
			m[2] = 2
		}
	}()
	select {}
}
