package plugin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/consensus"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/life/utils"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

var (
	// Identifies the prefix of the current round
	curAbnormalPrefix = []byte("SlashCb")
	// Identifies the prefix of the previous round
	preAbnormalPrefix = []byte("SlashPb")

	errDuplicateSignVerify = errors.New("duplicate signature verification failed")
	errSlashExist          = errors.New("punishment has been implemented")

	once = sync.Once{}
)

type SlashingPlugin struct {
	db             snapshotdb.DB
	decodeEvidence func(data string) (consensus.Evidences, error)
	privateKey     *ecdsa.PrivateKey
}

var slsh *SlashingPlugin

func SlashInstance() *SlashingPlugin {
	once.Do(func() {
		log.Info("Init Slashing plugin ...")
		slsh = &SlashingPlugin{
			db: snapshotdb.Instance(),
		}
	})
	return slsh
}

//func ClearSlashPlugin() error {
//	if nil == slsh {
//		return common.NewSysError("the SlashPlugin already be nil")
//	}
//	slsh = nil
//	return nil
//}

func (sp *SlashingPlugin) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	sp.privateKey = privateKey
}

func (sp *SlashingPlugin) SetDecodeEvidenceFun(f func(data string) (consensus.Evidences, error)) {
	sp.decodeEvidence = f
}

func (sp *SlashingPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	// If it is the first block in each round, switch the number of blocks in the upper and lower rounds.
	log.Debug("slashingPlugin BeginBlock", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()),
		"consensusSize", xutil.ConsensusSize())
	if (header.Number.Uint64()%xutil.ConsensusSize() == 1) && header.Number.Uint64() > 1 {
		if err := sp.switchEpoch(header.Number.Uint64(), blockHash); nil != err {
			log.Error("slashingPlugin switchEpoch fail", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "err", err)
			return err
		}
	}
	if err := sp.setPackAmount(blockHash, header); nil != err {
		log.Error("slashingPlugin setPackAmount fail", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "err", err)
		return err
	}
	// If it is the 230th block of each round, it will punish the node with abnormal block rate.
	if header.Number.Uint64() > xutil.ConsensusSize() && xutil.IsElection(header.Number.Uint64()) {
		log.Debug("slashingPlugin Ranking block amount", "blockNumber", header.Number.Uint64(), "blockHash",
			hex.EncodeToString(blockHash.Bytes()), "consensusSize", xutil.ConsensusSize(),
			"electionDistance", xcom.ElectionDistance())
		if result, err := sp.GetPreNodeAmount(header.ParentHash); nil != err {
			return err
		} else {
			if nil == result {
				log.Error("slashingPlugin GetPreNodeAmount is nil", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()))
				return common.NewBizError("block rate data not found")
			}
			validatorList, err := stk.GetCandidateONRound(blockHash, header.Number.Uint64(), PreviousRound, QueryStartIrr)
			if nil != err {
				return err
			}
			for _, validator := range validatorList {
				nodeId := validator.NodeId
				amount, success := result[nodeId]
				isSlash := false
				var rate uint32
				isDelete := false
				if success {
					// Start to punish nodes with abnormal block rate
					log.Debug("slashingPlugin node block amount", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "nodeId", hex.EncodeToString(nodeId.Bytes()), "amount", amount)
					if isAbnormal(amount) {
						if amount <= xcom.PackAmountAbnormal() && amount > xcom.PackAmountHighAbnormal() {
							isSlash = true
							rate = xcom.PackAmountLowSlashRate()
						} else if amount <= xcom.PackAmountHighAbnormal() {
							isSlash = true
							isDelete = true
							rate = xcom.PackAmountHighSlashRate()
						}
					}
				} else {
					isSlash = true
					isDelete = true
					rate = xcom.PackAmountHighSlashRate()
				}
				if isSlash && rate > 0 {
					slashAmount, sumAmount := calcSlashAmount(validator, rate, header.Number.Uint64())
					log.Info("Call SlashCandidates anomalous nodes", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()),
						"nodeId", hex.EncodeToString(nodeId.Bytes()), "packAmount", amount, "isDelete", isDelete, "sumAmount", sumAmount, "slash balance rate of remain", rate, "slashAmount", slashAmount)
					// If there is no record of the node, it means that there is no block, then the penalty is directly
					if err := stk.SlashCandidates(state, blockHash, header.Number.Uint64(), nodeId, slashAmount, isDelete, staking.LowRatio, common.ZeroAddr); nil != err {
						log.Error("slashingPlugin SlashCandidates failed", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "nodeId", hex.EncodeToString(nodeId.Bytes()), "err", err)
						return err
					}
				}
			}
		}
	}
	return nil
}

func (sp *SlashingPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	return nil
}

func (sp *SlashingPlugin) Confirmed(block *types.Block) error {
	return nil
}

func (sp *SlashingPlugin) getPackAmount(blockHash common.Hash, header *types.Header) (uint32, error) {
	log.Debug("slashingPlugin getPackAmount", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()))
	nodeId, err := parseNodeId(header)
	if nil != err {
		return 0, err
	}
	value, err := sp.db.Get(blockHash, curKey(nodeId.Bytes()))
	if nil != err && err != snapshotdb.ErrNotFound {
		return 0, err
	}
	var amount uint32
	if err == snapshotdb.ErrNotFound {
		amount = 0
	} else {
		amount = common.BytesToUint32(value)
	}
	return amount, nil
}

func (sp *SlashingPlugin) setPackAmount(blockHash common.Hash, header *types.Header) error {
	log.Debug("slashingPlugin setPackAmount", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()))
	nodeId, err := parseNodeId(header)
	if nil != err {
		return err
	}
	if value, err := sp.getPackAmount(blockHash, header); nil != err {
		return err
	} else {
		value++
		if err := sp.db.Put(blockHash, curKey(nodeId.Bytes()), common.Uint32ToBytes(value)); nil != err {
			return err
		}
		log.Debug("slashingPlugin setPackAmount success", "blockNumber", header.Number.Uint64(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "nodeId", hex.EncodeToString(nodeId.Bytes()), "value", value)
	}
	return nil
}

func (sp *SlashingPlugin) switchEpoch(blockNumber uint64, blockHash common.Hash) error {
	log.Debug("slashingPlugin switchEpoch", "blockNumber", blockNumber, "blockHash", hex.EncodeToString(blockHash.Bytes()))
	preCount := 0
	iter := sp.db.Ranking(blockHash, preAbnormalPrefix, 0)
	for iter.Next() {
		log.Debug("slashingPlugin switchEpoch ranking pre", "blockNumber", blockNumber, "key", hex.EncodeToString(iter.Key()), "value", iter.Value())
		if err := sp.db.Del(blockHash, iter.Key()); nil != err {
			return err
		}
		preCount++
	}

	curCount := 0
	iter = sp.db.Ranking(blockHash, curAbnormalPrefix, 0)
	for iter.Next() {
		key := iter.Key()
		log.Debug("slashingPlugin switchEpoch ranking cur", "blockNumber", blockNumber, "key", hex.EncodeToString(iter.Key()), "value", iter.Value())
		if err := sp.db.DelBaseDB(key); nil != err {
			return err
		}
		key = preKey(key[len(curAbnormalPrefix):])
		log.Debug("slashingPlugin switchEpoch ranking change pre", "blockNumber", blockNumber, "key", hex.EncodeToString(iter.Key()), "value", iter.Value())
		if err := sp.db.PutBaseDB(key, iter.Value()); nil != err {
			return err
		}
		curCount++
	}
	log.Info("slashingPlugin switchEpoch success", "blockNumber", blockNumber, "blockHash", hex.EncodeToString(blockHash.Bytes()), "preCount", preCount, "curCount", curCount)
	return nil
}

// Get the consensus rate of all nodes in the previous round
func (sp *SlashingPlugin) GetPreNodeAmount(parentHash common.Hash) (map[discover.NodeID]uint32, error) {
	result := make(map[discover.NodeID]uint32)
	iter := sp.db.Ranking(parentHash, preAbnormalPrefix, 0)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		amount := common.BytesToUint32(value)
		log.Debug("slashingPlugin GetPreEpochAnomalyNode", "parentHash", parentHash.Hex(), "key", hex.EncodeToString(key), "value", amount)
		nodeId, err := getNodeId(preAbnormalPrefix, key)
		if nil != err {
			return nil, err
		}
		result[nodeId] = amount
	}
	return result, nil
}

func (sp *SlashingPlugin) DecodeEvidence(data string) (consensus.Evidences, error) {
	if sp.decodeEvidence == nil {
		return nil, common.NewBizError("decodeEvidence function is nil")
	}
	return sp.decodeEvidence(data)
}

func (sp *SlashingPlugin) Slash(evidences consensus.Evidences, blockHash common.Hash, blockNumber uint64, stateDB xcom.StateDB, caller common.Address) error {
	log.Debug("slashingPlugin Slash", "blockNumber", blockNumber, "blockHash", hex.EncodeToString(blockHash.Bytes()), "evidencesSize", len(evidences), "caller", hex.EncodeToString(caller.Bytes()))
	for _, evidence := range evidences {
		err := sp.executeSlash(evidence, blockHash, blockNumber, stateDB, caller)
		if nil != err {
			if _, ok := err.(*common.BizError); ok {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func (sp *SlashingPlugin) executeSlash(evidence consensus.Evidence, blockHash common.Hash, blockNumber uint64, stateDB xcom.StateDB, caller common.Address) error {
	if err := evidence.Validate(); nil != err {
		log.Warn("slashing evidence validate failed", "err", err)
		return common.NewBizError(err.Error())
	}
	if value := sp.getSlashResult(evidence.Address(), evidence.BlockNumber(), uint32(evidence.Type()), stateDB); len(value) > 0 {
		log.Error("slashing failed", "blockNumber", evidence.BlockNumber(), "evidenceHash", hex.EncodeToString(evidence.Hash()), "addr", hex.EncodeToString(evidence.Address().Bytes()), "type", evidence.Type(), "err", errSlashExist.Error())
		return common.NewBizError(errSlashExist.Error())
	}
	if candidate, err := stk.GetCandidateInfo(blockHash, evidence.Address()); nil != err {
		log.Error("slashing failed", "blockNumber", evidence.BlockNumber(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "addr", hex.EncodeToString(evidence.Address().Bytes()), "err", err)
		return common.NewBizError(err.Error())
	} else {
		if nil == candidate {
			log.Error("slashing failed GetCandidateInfo is nil", "blockNumber", blockNumber, "blockHash", hex.EncodeToString(blockHash.Bytes()), "addr", hex.EncodeToString(evidence.Address().Bytes()), "type", evidence.Type())
			return common.NewBizError(errDuplicateSignVerify.Error())
		}
		slashAmount, sumAmount := calcSlashAmount(candidate, xcom.DuplicateSignHighSlash(), blockNumber)
		log.Info("Call SlashCandidates on executeSlash", "blockNumber", blockNumber, "blockHash", hex.EncodeToString(blockHash.Bytes()),
			"nodeId", candidate.NodeId.String(), "sumAmount", sumAmount, "rate", xcom.DuplicateSignHighSlash(), "slashAmount", slashAmount, "reporter", caller.Hex())
		if err := stk.SlashCandidates(stateDB, blockHash, blockNumber, candidate.NodeId, slashAmount, true, staking.DuplicateSign, caller); nil != err {
			log.Error("slashing failed SlashCandidates failed", "blockNumber", blockNumber, "blockHash", hex.EncodeToString(blockHash.Bytes()), "nodeId", hex.EncodeToString(candidate.NodeId.Bytes()), "err", err)
			return err
		}
		sp.putSlashResult(evidence.Address(), evidence.BlockNumber(), uint32(evidence.Type()), stateDB)
		log.Info("slashing duplicate signature success", "currentBlockNumber", blockNumber, "signBlockNumber", evidence.BlockNumber(), "blockHash", hex.EncodeToString(blockHash.Bytes()), "nodeId", hex.EncodeToString(candidate.NodeId.Bytes()), "etype", evidence.Type(), "txHash", hex.EncodeToString(stateDB.TxHash().Bytes()))
	}
	return nil
}

func (sp *SlashingPlugin) CheckDuplicateSign(addr common.Address, blockNumber uint64, etype uint32, stateDB xcom.StateDB) ([]byte, error) {
	if value := sp.getSlashResult(addr, blockNumber, etype, stateDB); len(value) > 0 {
		log.Info("CheckDuplicateSign exist", "blockNumber", blockNumber, "addr", hex.EncodeToString(addr.Bytes()), "type", etype, "txHash", hex.EncodeToString(value))
		return value, nil
	}
	return nil, nil
}

func (sp *SlashingPlugin) putSlashResult(addr common.Address, blockNumber uint64, etype uint32, stateDB xcom.StateDB) {
	stateDB.SetState(vm.SlashingContractAddr, duplicateSignKey(addr, blockNumber, etype), stateDB.TxHash().Bytes())
}

func (sp *SlashingPlugin) getSlashResult(addr common.Address, blockNumber uint64, etype uint32, stateDB xcom.StateDB) []byte {
	return stateDB.GetState(vm.SlashingContractAddr, duplicateSignKey(addr, blockNumber, etype))
}

// duplicate signature result key format addr+blockNumber+_+etype
func duplicateSignKey(addr common.Address, blockNumber uint64, etype uint32) []byte {
	value := append(addr.Bytes(), utils.Uint64ToBytes(blockNumber)...)
	value = append(value, []byte("_")...)
	value = append(value, utils.Uint64ToBytes(uint64(etype))...)
	return value
}

func curKey(key []byte) []byte {
	return append(curAbnormalPrefix, key...)
}

func preKey(key []byte) []byte {
	return append(preAbnormalPrefix, key...)
}

func getNodeId(prefix []byte, key []byte) (discover.NodeID, error) {
	key = key[len(prefix):]
	nodeId, err := discover.BytesID(key)
	if nil != err {
		return discover.NodeID{}, err
	}
	return nodeId, nil
}

func isAbnormal(amount uint32) bool {
	return uint64(amount) < (xutil.ConsensusSize() / xcom.ConsValidatorNum())
}

func parseNodeId(header *types.Header) (discover.NodeID, error) {
	log.Debug("extra parseNodeId", "extra", hex.EncodeToString(header.Extra), "sealHash", hex.EncodeToString(header.SealHash().Bytes()))
	if xutil.IsWorker(header.Extra) {
		return discover.PubkeyID(&SlashInstance().privateKey.PublicKey), nil
	} else {
		sign := header.Extra[32:97]
		pk, err := crypto.SigToPub(header.SealHash().Bytes(), sign)
		if nil != err {
			return discover.NodeID{}, err
		}
		return discover.PubkeyID(pk), nil
	}
}

func calcSlashAmount(candidate *staking.Candidate, rate uint32, blockNumber uint64) (*big.Int, *big.Int) {
	// Recalculate the quality deposit
	lazyCalcStakeAmount(xutil.CalculateEpoch(blockNumber), candidate)
	sumAmount := new(big.Int)
	sumAmount.Add(candidate.Released, candidate.RestrictingPlan)
	if sumAmount.Cmp(common.Big0) > 0 {
		amount := new(big.Int).Mul(sumAmount, new(big.Int).SetUint64(uint64(rate)))
		return amount.Div(amount, new(big.Int).SetUint64(100)), sumAmount
	}
	return common.Big0, common.Big0
}
