package xcom

import (
	"crypto/ecdsa"
	"encoding/hex"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	cryptoHandlerOnce = sync.Once{}
)

var chandler *CryptoHandler

type CryptoHandler struct {
	privateKey *ecdsa.PrivateKey
}

func GetCryptoHandler() *CryptoHandler {
	cryptoHandlerOnce.Do(func() {
		log.Info("Init CryptoHandler ...")
		chandler = &CryptoHandler{}
	})
	return chandler
}

func (chandler *CryptoHandler) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	chandler.privateKey = privateKey
}

func (chandler *CryptoHandler) Sign(data []byte) ([]byte, error) {
	if chandler == nil || chandler.privateKey == nil {
		return nil, common.NewSysError("PrivateKey missed")
	}
	return crypto.Sign(data, chandler.privateKey)
}

func (chandler *CryptoHandler) MustSign(data []byte) []byte {
	if chandler == nil || chandler.privateKey == nil {
		panic("Private key is missed")
	}
	sig, err := crypto.Sign(data, chandler.privateKey)
	if err != nil {
		panic(err)
	}
	return sig
}

func (chandler *CryptoHandler) ValidateSign(data []byte, sig []byte, nodeID discover.NodeID) bool {
	pubKey, err := crypto.SigToPub(data, sig)
	if err != nil {
		log.Error("Check if the signature is signed by a node", "err", err)
		return false
	}
	id := discover.PubkeyID(pubKey)
	if id == nodeID {
		return true
	}
	log.Error("the signature is not signed by the node", "nodeID", hex.EncodeToString(nodeID.Bytes()[:8]))
	return false
}
