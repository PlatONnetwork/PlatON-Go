package xcom

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/crypto/vrf"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"sync"
)

var (
	errInvalidVrfProve = errors.New("invalid vrf prove")
	errStorageNonce    = errors.New("storage previous nonce failed")

	nonceStorageKey = []byte("nonceStorageKey")

	once = sync.Once{}
)

type vrfHandler struct {
	db           snapshotdb.DB
	privateKey   *ecdsa.PrivateKey
	genesisNonce []byte
}

var vh *vrfHandler

func NewVrfHandler(genesisNonce []byte) *vrfHandler {
	once.Do(func() {
		vh = &vrfHandler{
			db:           snapshotdb.Instance(),
			genesisNonce: genesisNonce,
		}
	})
	return vh
}

func GetVrfHandlerInstance() *vrfHandler {
	return vh
}

func (vh *vrfHandler) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	vh.privateKey = privateKey
}

func (vh *vrfHandler) GenerateNonce(currentBlockNumber *big.Int, parentHash common.Hash) ([]byte, error) {
	parentNonce, err := vh.getParentNonce(currentBlockNumber, parentHash)
	if nil != err {
		return nil, err
	}
	log.Debug("Generate proof based on input", "currentBlockNumber", currentBlockNumber.Uint64(), "parentHash", hex.EncodeToString(parentHash.Bytes()), "parentNonce", hex.EncodeToString(parentNonce))
	if value, err := vrf.Prove(vh.privateKey, parentNonce); nil != err {
		log.Error("Generate proof failure", "currentBlockNumber", currentBlockNumber.Uint64(), "parentHash", hex.EncodeToString(parentHash.Bytes()), "parentNonce", hex.EncodeToString(parentNonce), "err", err)
		return nil, err
	} else {
		if len(value) > 0 {
			log.Info("Generate vrf proof Success", "blockNumber", currentBlockNumber.Uint64(), "parentHash", hex.EncodeToString(parentHash.Bytes()), "nonce", hex.EncodeToString(value))
			return value, nil
		}
	}
	return nil, fmt.Errorf("generate proof failed, seed:%x", parentNonce)
}

func (vh *vrfHandler) VerifyVrf(pk *ecdsa.PublicKey, currentBlockNumber *big.Int, parentBlockHash common.Hash, blockHash common.Hash, proof []byte) error {
	// Verify VRF Proof
	log.Debug("Verification block vrf prove", "blockNumber", currentBlockNumber.Uint64(), "hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(parentBlockHash.Bytes()), "proof", hex.EncodeToString(proof))
	parentNonce, err := vh.getParentNonce(currentBlockNumber, parentBlockHash)
	if nil != err {
		return err
	}
	if value, err := vrf.Verify(pk, proof, parentNonce); nil != err {
		log.Error("Vrf proves verification failure", "blockNumber", currentBlockNumber.Uint64(), "hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(parentBlockHash.Bytes()), "proof", hex.EncodeToString(proof), "input", hex.EncodeToString(parentNonce), "err", err)
		return err
	} else if !value {
		log.Error("Vrf proves verification failure", "blockNumber", currentBlockNumber.Uint64(), "hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(parentBlockHash.Bytes()), "proof", hex.EncodeToString(proof), "input", hex.EncodeToString(parentNonce))
		return errInvalidVrfProve
	}
	log.Info("Vrf proves successful verification", "blockNumber", currentBlockNumber.Uint64(), "hash", hex.EncodeToString(blockHash.Bytes()), "parentHash", hex.EncodeToString(parentBlockHash.Bytes()), "proof", hex.EncodeToString(proof), "input", hex.EncodeToString(parentNonce))
	return nil
}

func (vh *vrfHandler) Storage(currentBlockNumber *big.Int, parentHash common.Hash, hash common.Hash, nonce []byte) error {
	log.Debug("Storage previous nonce", "blockNumber", currentBlockNumber.Uint64(), "parentHash", hex.EncodeToString(parentHash.Bytes()), "hash", hex.EncodeToString(hash.Bytes()), "nonce", hex.EncodeToString(nonce))
	nonces := make([][]byte, 0)
	if currentBlockNumber.Cmp(common.Big1) > 0 {
		if value, err := vh.Load(parentHash); nil != err {
			return err
		} else {
			nonces = make([][]byte, len(value))
			copy(nonces, value)
			log.Debug("Storage previous nonce", "blockNumber", currentBlockNumber.Uint64(), "parentHash", hex.EncodeToString(parentHash.Bytes()), "hash", hex.EncodeToString(hash.Bytes()), "valueLength", len(value), "EpochValidatorNum", EpochValidatorNum)
			if uint64(len(nonces)) == EpochValidatorNum {
				nonces = nonces[1:]
			}
		}
	}
	nonces = append(nonces, vrf.ProofToHash(nonce))
	if enValue, err := rlp.EncodeToBytes(nonces); nil != err {
		log.Error("Storage previous nonce failed", "blockNumber", currentBlockNumber.Uint64(), "parentHash", hex.EncodeToString(parentHash.Bytes()), "hash", hex.EncodeToString(hash.Bytes()), "key", string(nonceStorageKey), "valueLength", len(nonces), "nonce", hex.EncodeToString(nonce), "err", err)
		return err
	} else {
		vh.db.Put(hash, nonceStorageKey, enValue)
		log.Info("Storage previous nonce Success", "blockNumber", currentBlockNumber.Uint64(), "parentHash", hex.EncodeToString(parentHash.Bytes()), "hash", hex.EncodeToString(hash.Bytes()), "valueLength", len(nonces), "EpochValidatorNum", EpochValidatorNum, "nonce", hex.EncodeToString(nonce), "firstNonce", hex.EncodeToString(nonces[0]), "lastNonce", hex.EncodeToString(nonces[len(nonces)-1]))
	}
	return nil
}

func (vh *vrfHandler) Load(hash common.Hash) ([][]byte, error) {
	if value, err := vh.db.Get(hash, nonceStorageKey); nil != err {
		log.Error("Loading previous nonce failed", "hash", hash, "key", string(nonceStorageKey), "err", err)
		return nil, err
	} else {
		nonces := make([][]byte, 0)
		if err := rlp.DecodeBytes(value, &nonces); nil != err {
			log.Error("rlpDecode previous nonce failed", "hash", hash, "key", string(nonceStorageKey), "err", err)
			return nil, err
		}
		return nonces, nil
	}
}

func (vh *vrfHandler) getParentNonce(currentBlockNumber *big.Int, parentHash common.Hash) ([]byte, error) {
	// If it is the first block, take the random number from the Genesis block.
	log.Debug("Get the nonce of the previous block", "blockNumber", currentBlockNumber.Uint64(), "parentHash", hex.EncodeToString(parentHash.Bytes()))
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
	return nil, fmt.Errorf("nonce of the previous block could not be found, blockNumber：%v, parentHash：%x", currentBlockNumber.Uint64(), parentHash)
}
