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

package handler

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/vrf"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var (
	ErrInvalidVrfProve = errors.New("invalid vrf prove")
	ErrStorageNonce    = errors.New("storage previous nonce failed")

	NonceStorageKey = []byte("nonceStorageKey")

	once = sync.Once{}
)

type VrfHandler struct {
	db           snapshotdb.DB
	privateKey   *ecdsa.PrivateKey
	genesisNonce []byte
}

var vh *VrfHandler

func NewVrfHandler(genesisNonce []byte) *VrfHandler {
	once.Do(func() {
		vh = &VrfHandler{
			db:           snapshotdb.Instance(),
			genesisNonce: genesisNonce,
		}
	})
	return vh
}

func GetVrfHandlerInstance() *VrfHandler {
	return vh
}

func (vh *VrfHandler) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	vh.privateKey = privateKey
}

func (vh *VrfHandler) GenerateNonce(currentBlockNumber *big.Int, parentHash common.Hash) ([]byte, error) {
	parentNonce, err := vh.getParentNonce(currentBlockNumber, parentHash)
	if nil != err {
		return nil, err
	}
	log.Debug("Generate proof based on input", "currentBlockNumber", currentBlockNumber.Uint64(),
		"parentHash", hex.EncodeToString(parentHash.Bytes()), "parentNonce", hex.EncodeToString(parentNonce))
	if value, err := vrf.Prove(vh.privateKey, parentNonce); nil != err {
		log.Error("Generate proof failure", "currentBlockNumber", currentBlockNumber.Uint64(),
			"parentHash", hex.EncodeToString(parentHash.Bytes()), "parentNonce", hex.EncodeToString(parentNonce), "err", err)
		return nil, err
	} else {
		if len(value) > 0 {
			log.Info("Generate vrf proof Success", "blockNumber", currentBlockNumber.Uint64(),
				"parentHash", hex.EncodeToString(parentHash.Bytes()), "nonce", hex.EncodeToString(value),
				"nodeId", discover.PubkeyID(&vh.privateKey.PublicKey).String())
			return value, nil
		}
	}
	return nil, fmt.Errorf("generate proof failed, seed:%x", parentNonce)
}

func (vh *VrfHandler) VerifyVrf(pk *ecdsa.PublicKey, currentBlockNumber *big.Int, parentBlockHash common.Hash,
	blockHash common.Hash, proof []byte) error {
	// Verify VRF Proof
	log.Debug("Verification block vrf prove", "current blockNumber", currentBlockNumber.Uint64(),
		"current hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(parentBlockHash.Bytes()),
		"proof", hex.EncodeToString(proof))
	parentNonce, err := vh.getParentNonce(currentBlockNumber, parentBlockHash)
	if nil != err {
		return err
	}
	if value, err := vrf.Verify(pk, proof, parentNonce); nil != err {
		log.Error("Vrf proves verification failure", "current blockNumber", currentBlockNumber.Uint64(),
			"current hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(parentBlockHash.Bytes()),
			"proof", hex.EncodeToString(proof), "input", hex.EncodeToString(parentNonce), "err", err)
		return err
	} else if !value {
		log.Error("Vrf proves verification failure", "current blockNumber", currentBlockNumber.Uint64(),
			"current hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(parentBlockHash.Bytes()),
			"proof", hex.EncodeToString(proof), "input", hex.EncodeToString(parentNonce))
		return ErrInvalidVrfProve
	}
	log.Info("Vrf proves successful verification", "current blockNumber", currentBlockNumber.Uint64(),
		"current hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(parentBlockHash.Bytes()),
		"proof", hex.EncodeToString(proof), "input", hex.EncodeToString(parentNonce))
	return nil
}

func (vh *VrfHandler) Storage(blockNumber *big.Int, parentHash common.Hash, blockHash common.Hash, nonce []byte) error {
	log.Debug("Storage previous nonce", "blockNumber", blockNumber.Uint64(), "parentHash",
		hex.EncodeToString(parentHash.Bytes()), "current hash", hex.EncodeToString(blockHash.Bytes()), "nonce", hex.EncodeToString(nonce))
	nonces := make([][]byte, 0)

	maxValidatorsNum, err := gov.GovernMaxValidators(blockNumber.Uint64(), blockHash)
	if nil != err {
		log.Error("Failed to Storage VRF nonce", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "err", err)
		return err
	}

	if blockNumber.Cmp(common.Big1) > 0 {
		if value, err := vh.Load(parentHash); nil != err {
			return err
		} else {
			nonces = make([][]byte, len(value))
			copy(nonces, value)
			log.Debug("Storage previous nonce", "current blockNumber", blockNumber.Uint64(), "parentHash",
				hex.EncodeToString(parentHash.Bytes()), "current hash", hex.EncodeToString(blockHash.Bytes()), "valueLength",
				len(value), "MaxValidators", maxValidatorsNum)
			if uint64(len(nonces)) == maxValidatorsNum {
				nonces = nonces[1:]
			} else if uint64(len(nonces)) > maxValidatorsNum {
				nonces = nonces[uint64(len(nonces)+1)-maxValidatorsNum:]
			}
		}
	}
	nonces = append(nonces, vrf.ProofToHash(nonce))
	if enValue, err := rlp.EncodeToBytes(nonces); nil != err {
		log.Error("Storage previous nonce failed", "current blockNumber", blockNumber.Uint64(),
			"parentHash", hex.EncodeToString(parentHash.Bytes()), "current hash", hex.EncodeToString(blockHash.Bytes()),
			"key", string(NonceStorageKey), "valueLength", len(nonces), "nonce", hex.EncodeToString(nonce), "err", err)
		return err
	} else {
		if err := vh.db.Put(blockHash, NonceStorageKey, enValue); nil != err {
			log.Error("Storage previous nonce failed", "current blockNumber", blockNumber.Uint64(),
				"parentHash", hex.EncodeToString(parentHash.Bytes()), "current hash", hex.EncodeToString(blockHash.Bytes()),
				"key", string(NonceStorageKey), "valueLength", len(nonces), "nonce", hex.EncodeToString(nonce), "enValue", hex.EncodeToString(enValue), "err", err)
			return err
		}
		log.Info("Storage previous nonce Success", "current blockNumber", blockNumber.Uint64(),
			"parentHash", hex.EncodeToString(parentHash.Bytes()), "current hash", hex.EncodeToString(blockHash.Bytes()),
			"valueLength", len(nonces), "MaxValidators", maxValidatorsNum, "nonce", hex.EncodeToString(nonce),
			"firstNonce", hex.EncodeToString(nonces[0]), "lastNonce", hex.EncodeToString(nonces[len(nonces)-1]))
	}
	return nil
}

func (vh *VrfHandler) Load(hash common.Hash) ([][]byte, error) {
	if value, err := vh.db.Get(hash, NonceStorageKey); nil != err {
		log.Error("Loading previous nonce failed", "hash", hash.String(), "key", string(NonceStorageKey), "err", err)
		return nil, err
	} else {
		nonces := make([][]byte, 0)
		if err := rlp.DecodeBytes(value, &nonces); nil != err {
			log.Error("rlpDecode previous nonce failed", "hash", hash, "key", string(NonceStorageKey), "err", err)
			return nil, err
		}
		return nonces, nil
	}
}

func (vh *VrfHandler) getParentNonce(currentBlockNumber *big.Int, parentHash common.Hash) ([]byte, error) {
	// If it is the first block, take the random number from the Genesis block.
	log.Debug("Get the nonce of the previous block", "blockNumber", currentBlockNumber.Uint64(),
		"parentHash", hex.EncodeToString(parentHash.Bytes()))
	if currentBlockNumber.Cmp(common.Big1) == 0 && len(vh.genesisNonce) > 0 {
		log.Debug("Get the nonce of the genesis block", "nonce", hex.EncodeToString(vh.genesisNonce))
		return vrf.ProofToHash(vh.genesisNonce), nil
	} else {
		if value, err := vh.Load(parentHash); nil != err {
			return nil, err
		} else {
			if len(value) > 0 {
				return value[len(value)-1], nil
			}
		}
	}
	return nil, fmt.Errorf("nonce of the previous block could not be found, blockNumber：%v, parentHash：%x",
		currentBlockNumber.Uint64(), parentHash)
}
