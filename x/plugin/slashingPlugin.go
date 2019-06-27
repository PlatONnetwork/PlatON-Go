package plugin

import (
	"encoding/hex"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/go-errors/errors"
)

const (
	MutiSignPrepare		uint8 = iota+1
	MutiSignViewChange
)

var (
	curAbnormalPrefix = []byte("SlashCb")
	preAbnormalPrefix = []byte("SlashPb")

	// The number of low exceptions per consensus round
	blockAmountLow 				uint16 	= 8
	// The number of high exceptions per consensus round
	blockAmountHigh 			uint16 	= 5
	//
	blockAmountLowSlashing		uint32	= 10
	blockAmountHighSlashing		uint32	= 20
	duplicateSignNum			uint32	= 2
	duplicateSignLowSlashing	uint32	= 10
	duplicateSignHighSlashing	uint32	= 10

	errSetBlockAmount	= errors.New("set block amount fail")
	errMutiSignVerify	= errors.New("Multi-sign verification failed")
)

type slashingPlugin struct {
	db		snapshotdb.DB
}

var slashPlugin *slashingPlugin

func SlashInstance(db snapshotdb.DB) *slashingPlugin {
	if slashPlugin == nil {
		slashPlugin = &slashingPlugin{
			db:db,
		}
	}
	return slashPlugin
}

func (sp *slashingPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {
	return true, nil
}

func (sp *slashingPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {
	// If it is the 230th block of each round, it will punish the node with abnormal block rate.
	if header.Number.Uint64() % (xcom.ConsensusSize - xcom.ElectionDistance) == 0 {
		log.Debug("slashingPlugin Ranking block amount", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "consensusSize", xcom.ConsensusSize, "electionDistance", xcom.ElectionDistance)
		it := sp.db.Ranking(blockHash, curAbnormalPrefix, 0)
		for end := true; end; end = it.Next() {
			key := it.Key()
			value := it.Value()
			var amount uint16
			if err := rlp.DecodeBytes(value, amount); nil != err {
				log.Error("slashingPlugin rlp block amount fail", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "value", value, "err", err)
				return false, err
			}
			// Start to punish nodes with abnormal block rate
			log.Debug("slashingPlugin node block amount", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "key", hex.EncodeToString(key), "value", amount)
			if uint64(amount) < (xcom.ConsensusSize / xcom.ConsValidatorNum) {
				key = key[len(curAbnormalPrefix):]
				nodeId, err := discover.BytesID(key)
				if nil != err {
					return false, err
				}
				log.Debug("Slashing anomalous nodes", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "nodeId", hex.EncodeToString(nodeId.Bytes()))
				if amount <= blockAmountLow && amount > blockAmountHigh {

				} else if amount <= blockAmountHigh {

				}
			}
		}
	}
	return true, nil
}

func (sp *slashingPlugin) Confirmed(block *types.Block) error {
	// If it is the first block in each round, switch the number of blocks in the upper and lower rounds.
	log.Debug("slashingPlugin Confirmed", "blockNumber", block.NumberU64(), "blockHash", hex.EncodeToString(block.Hash().Bytes()), "consensusSize", xcom.ConsensusSize)
	if (block.NumberU64() % xcom.ConsensusSize == 1) && block.NumberU64() / xcom.ConsensusSize >= 2 {
		if err := sp.switchEpoch(block.Hash()); nil != err {
			return err
		}
	}
	if err := sp.setBlockAmount(block.Hash(), block.Header()); nil != err {
		log.Error("slashingPlugin setBlockAmount fail", "blockNumber", block.NumberU64(), "blockHash", hex.EncodeToString(block.Hash().Bytes()), "err", err)
		return err
	}
	return nil
}

func (sp *slashingPlugin) getBlockAmount(blockHash common.Hash, header *types.Header) (uint16, error) {
	log.Debug("slashingPlugin getBlockAmount", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()))
	nodeId, err := getNodeId(header)
	if nil != err {
		return 0, err
	}
	value, err := sp.db.Get(blockHash, curKey(nodeId.Bytes()))
	if nil != err && err != snapshotdb.ErrNotFound {
		return 0, err
	}
	var amount uint16
	if err == snapshotdb.ErrNotFound {
		amount = 0
	} else {
		if err := rlp.DecodeBytes(value, &amount); nil != err {
			return 0, err
		}
	}
	return amount, nil
}

func (sp *slashingPlugin) setBlockAmount(blockHash common.Hash, header *types.Header) error {
	log.Debug("slashingPlugin setBlockAmount", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()))
	nodeId, err := getNodeId(header)
	if nil != err {
		return err
	}
	if value, err := sp.getBlockAmount(blockHash, header); nil != err {
		return err
	} else {
		value++
		if enValue, err := rlp.EncodeToBytes(value); nil != err {
			return err
		} else {
			if success, err := sp.db.NewBlock(header.Number, header.ParentHash, blockHash); nil != err {
				log.Error("slashingPlugin setBlockAmount NewBlock failed", "blockNumber", header.Number.Uint64(), "parentHash", hex.EncodeToString(header.ParentHash.Bytes()), "hash", string(blockHash.Bytes()), "err", err)
				return err
			} else if !success {
				log.Error("slashingPlugin setBlockAmount NewBlock failed", "blockNumber", header.Number.Uint64(), "parentHash", hex.EncodeToString(header.ParentHash.Bytes()), "hash", string(blockHash.Bytes()))
				return errSetBlockAmount
			}
			if _, err := sp.db.Put(blockHash, curKey(nodeId.Bytes()), enValue); nil != err {
				return err
			}
			log.Debug("slashingPlugin setBlockAmount success", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "value", value)
		}
	}
	return nil
}

func (sp *slashingPlugin) switchEpoch(blockHash common.Hash) error {
	log.Debug("slashingPlugin switchEpoch", "blockHash", hex.EncodeToString(blockHash.Bytes()))
	it := sp.db.Ranking(blockHash, preAbnormalPrefix, 0)
	preCount := 0
	for end := true; end; end = it.Next() {
		if _, err := sp.db.Del(blockHash, it.Key()); nil != err {
			return err
		}
		preCount++
	}
	curCount := 0
	it = sp.db.Ranking(blockHash, curAbnormalPrefix, 0)
	for end := true; end; end = it.Next() {
		key := it.Key()
		if _, err := sp.db.Del(blockHash, key); nil != err {
			return err
		}
		key = preKey(key[len(curAbnormalPrefix):])
		if _, err := sp.db.Put(blockHash, key, it.Value()); nil != err {
			return err
		}
		curCount++
	}
	log.Info("slashingPlugin switchEpoch success", "blockHash", hex.EncodeToString(blockHash.Bytes()), "preCount", preCount, "curCount", curCount)
	return nil
}

func (sp *slashingPlugin) Slash(mutiSignType uint8, evidence xcom.Evidence) error {
	if err := evidence.Validate(); nil != err {
		return err
	}

	return nil
}

func curKey(key []byte) []byte {
	return append(curAbnormalPrefix, key...)
}

func preKey(key []byte) []byte {
	return append(preAbnormalPrefix, key...)
}

func getNodeId(header *types.Header) (discover.NodeID, error) {
	sign := header.Extra[32:97]
	pk, err := crypto.SigToPub(header.SealHash().Bytes(), sign)
	if nil != err {
		return discover.NodeID{}, err
	}
	return discover.PubkeyID(pk), nil
}
