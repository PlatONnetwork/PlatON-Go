// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package node

import (
	"crypto/ecdsa"
	"encoding/hex"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

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

func (chandler *CryptoHandler) Sign(data interface{}) ([]byte, error) {
	if chandler == nil || chandler.privateKey == nil {
		return nil, common.InternalError.Wrap("PrivateKey missed")
	}
	return crypto.Sign(RlpHash(data).Bytes(), chandler.privateKey)
}

func (chandler *CryptoHandler) MustSign(data interface{}) []byte {
	if chandler == nil || chandler.privateKey == nil {
		panic("Private key is missed")
	}
	sig, err := crypto.Sign(RlpHash(data).Bytes(), chandler.privateKey)
	if err != nil {
		panic(err)
	}
	return sig
}

func (chandler *CryptoHandler) IsSignedByNodeID(data interface{}, sig []byte, nodeID discover.NodeID) bool {
	pubKey, err := crypto.SigToPub(RlpHash(data).Bytes(), sig)
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

func RlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
