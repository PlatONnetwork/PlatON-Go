package plugin

import (
	"crypto/ecdsa"
	"encoding/hex"
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

	errDuplicateSignVerify = common.NewBizError(401, "duplicate signature verification failed")
	errSlashExist          = common.NewBizError(402, "punishment has been implemented")

	once = sync.Once{}
)

type SlashingPlugin struct {
	db             snapshotdb.DB
	decodeEvidence func(dupType consensus.EvidenceType, data string) (consensus.Evidence, error)
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

func (sp *SlashingPlugin) SetDecodeEvidenceFun(f func(dupType consensus.EvidenceType, data string) (consensus.Evidence, error)) {
	sp.decodeEvidence = f
}

func (sp *SlashingPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	// If it is the first block in each round, switch the number of blocks in the upper and lower rounds.
	if (header.Number.Uint64()%xutil.ConsensusSize() == 1) && header.Number.Uint64() > 1 {
		if err := sp.switchEpoch(header.Number.Uint64(), blockHash); nil != err {
			log.Error("Failed to slashingPlugin switchEpoch fail", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
			return err
		}
	}
	if err := sp.setPackAmount(blockHash, header); nil != err {
		log.Error("slashingPlugin setPackAmount fail", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
		return err
	}
	// If it is the 230th block of each round, it will punish the node with abnormal block rate.
	if header.Number.Uint64() > xutil.ConsensusSize() && xutil.IsElection(header.Number.Uint64()) {
		log.Debug("slashingPlugin Ranking block amount", "blockNumber", header.Number.Uint64(), "blockHash",
			blockHash.TerminalString(), "consensusSize", xutil.ConsensusSize(), "electionDistance", xcom.ElectionDistance())
		if result, err := sp.GetPreNodeAmount(header.ParentHash); nil != err {
			return err
		} else {
			if nil == result {
				log.Error("Failed to slashingPlugin GetPreNodeAmount is nil", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString())
				return common.InternalError.Wrap("packAmount data not found")
			}
			validatorList, err := stk.GetCandidateONRound(blockHash, header.Number.Uint64(), PreviousRound, QueryStartIrr)
			if nil != err {
				log.Error("Failed to slashingPlugin BeginBlock, call GetCandidateONRound is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString())
				return err
			}
			for _, validator := range validatorList {
				nodeId := validator.NodeId
				amount, success := result[nodeId]
				isSlash := false
				var rate uint32

				slashType := staking.LowRatio

				if success {
					// Start to punish nodes with abnormal block rate
					log.Debug("slashingPlugin node block amount", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(),
						"nodeId", nodeId.TerminalString(), "amount", amount)
					if isAbnormal(amount) {
						if amount <= xcom.PackAmountAbnormal() && amount > xcom.PackAmountHighAbnormal() {
							isSlash = true
							rate = xcom.PackAmountLowSlashRate()
						} else if amount <= xcom.PackAmountHighAbnormal() {
							isSlash = true
							slashType = staking.LowRatioDel
							rate = xcom.PackAmountHighSlashRate()
						}
					}
				} else {
					isSlash = true
					slashType = staking.LowRatioDel
					rate = xcom.PackAmountHighSlashRate()
				}
				if isSlash && rate > 0 {
					slashAmount, sumAmount := calcSlashAmount(validator, rate, header.Number.Uint64())
					log.Info("Call SlashCandidates anomalous nodes", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(),
						"nodeId", nodeId.TerminalString(), "packAmount", amount, "slashType", slashType, "sumAmount", sumAmount, "slash balance rate of remain", rate, "slashAmount", slashAmount)
					// If there is no record of the node, it means that there is no block, then the penalty is directly
					if err := stk.SlashCandidates(state, blockHash, header.Number.Uint64(), nodeId, slashAmount, slashType, common.ZeroAddr); nil != err {
						log.Error("Failed to slashingPlugin SlashCandidates failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "nodeId", nodeId.TerminalString(), "err", err)
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

func (sp *SlashingPlugin) getPackAmount(blockHash common.Hash, nodeId discover.NodeID) (uint32, error) {
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
	nodeId, err := parseNodeId(header)
	if nil != err {
		return err
	}
	if value, err := sp.getPackAmount(blockHash, nodeId); nil != err {
		return err
	} else {
		value++
		if err := sp.db.Put(blockHash, curKey(nodeId.Bytes()), common.Uint32ToBytes(value)); nil != err {
			return err
		}
		log.Debug("slashingPlugin setPackAmount success", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "nodeId", nodeId.TerminalString(), "value", value)
	}
	return nil
}

func (sp *SlashingPlugin) switchEpoch(blockNumber uint64, blockHash common.Hash) error {
	preCount := 0
	preIterator := sp.db.Ranking(blockHash, preAbnormalPrefix, 0)
	for preIterator.Next() {
		key := preIterator.Key()
		value := preIterator.Value()
		log.Debug("slashingPlugin switchEpoch ranking pre", "blockNumber", blockNumber, "key", hex.EncodeToString(key), "value", common.BytesToUint32(value))
		if err := sp.db.Del(blockHash, key); nil != err {
			return err
		}
		preCount++
	}
	preIterator.Release()
	curCount := 0
	curIterator := sp.db.Ranking(blockHash, curAbnormalPrefix, 0)
	for curIterator.Next() {
		key := curIterator.Key()
		value := curIterator.Value()
		log.Debug("slashingPlugin switchEpoch ranking cur", "blockNumber", blockNumber, "key", hex.EncodeToString(key), "value", common.BytesToUint32(value))
		if err := sp.db.Del(blockHash, key); nil != err {
			return err
		}
		key = preKey(key[len(curAbnormalPrefix):])
		log.Debug("slashingPlugin switchEpoch ranking change pre", "blockNumber", blockNumber, "key", hex.EncodeToString(key), "value", common.BytesToUint32(value))
		if err := sp.db.Put(blockHash, key, value); nil != err {
			return err
		}
		curCount++
	}
	curIterator.Release()
	log.Info("slashingPlugin switchEpoch success", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "preCount", preCount, "curCount", curCount)
	return nil
}

// Get the consensus rate of all nodes in the previous round
func (sp *SlashingPlugin) GetPreNodeAmount(parentHash common.Hash) (map[discover.NodeID]uint32, error) {
	result := make(map[discover.NodeID]uint32)
	iter := sp.db.Ranking(parentHash, preAbnormalPrefix, 0)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		amount := common.BytesToUint32(value)
		nodeId, err := getNodeId(preAbnormalPrefix, key)
		if nil != err {
			return nil, err
		}
		log.Debug("slashingPlugin GetPreNodeAmount", "parentHash", parentHash.Hex(), "nodeId", nodeId.TerminalString(), "value", amount)
		result[nodeId] = amount
	}
	return result, nil
}

func (sp *SlashingPlugin) DecodeEvidence(dupType consensus.EvidenceType, data string) (consensus.Evidence, error) {
	if sp.decodeEvidence == nil {
		return nil, common.InternalError.Wrap("decodeEvidence function is nil")
	}
	return sp.decodeEvidence(dupType, data)
}

func (sp *SlashingPlugin) Slash(evidence consensus.Evidence, blockHash common.Hash, blockNumber uint64, stateDB xcom.StateDB, caller common.Address) error {
	if err := evidence.Validate(); nil != err {
		log.Warn("slashing evidence validate failed", "blockNumber", blockNumber, "err", err)
		return common.InternalError.Wrap(err.Error())
	}
	if value := sp.getSlashResult(evidence.Address(), evidence.BlockNumber(), evidence.Type(), stateDB); len(value) > 0 {
		log.Error("slashing failed", "evidenceBlockNumber", evidence.BlockNumber(), "evidenceHash", hex.EncodeToString(evidence.Hash()), "addr", hex.EncodeToString(evidence.Address().Bytes()), "type", evidence.Type(), "err", errSlashExist.Error())
		return errSlashExist
	}
	if candidate, err := stk.GetCandidateInfo(blockHash, evidence.Address()); nil != err {
		log.Error("slashing failed", "blockNumber", evidence.BlockNumber(), "blockHash", blockHash.TerminalString(), "addr", hex.EncodeToString(evidence.Address().Bytes()), "err", err)
		return common.InternalError.Wrap(err.Error())
	} else {
		if nil == candidate {
			log.Error("slashing failed GetCandidateInfo is nil", "blockNumber", blockNumber, "blockHash", hex.EncodeToString(blockHash.Bytes()), "addr", hex.EncodeToString(evidence.Address().Bytes()), "type", evidence.Type())
			return errDuplicateSignVerify
		}
		slashAmount, sumAmount := calcSlashAmount(candidate, xcom.DuplicateSignHighSlash(), blockNumber)
		log.Info("Call SlashCandidates on executeSlash", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"nodeId", candidate.NodeId.String(), "sumAmount", sumAmount, "rate", xcom.DuplicateSignHighSlash(), "slashAmount", slashAmount, "reporter", caller.Hex())
		if err := stk.SlashCandidates(stateDB, blockHash, blockNumber, candidate.NodeId, slashAmount, staking.DuplicateSign, caller); nil != err {
			log.Error("slashing failed SlashCandidates failed", "blockNumber", blockNumber, "blockHash", hex.EncodeToString(blockHash.Bytes()), "nodeId", hex.EncodeToString(candidate.NodeId.Bytes()), "err", err)
			return err
		}
		sp.putSlashResult(evidence.Address(), evidence.BlockNumber(), evidence.Type(), stateDB)
		log.Info("slashing duplicate signature success", "currentBlockNumber", blockNumber, "signBlockNumber", evidence.BlockNumber(), "blockHash", blockHash.TerminalString(),
			"nodeId", candidate.NodeId.TerminalString(), "dupType", evidence.Type(), "txHash", stateDB.TxHash().TerminalString())
	}
	return nil
}

func (sp *SlashingPlugin) CheckDuplicateSign(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType, stateDB xcom.StateDB) ([]byte, error) {
	if value := sp.getSlashResult(addr, blockNumber, dupType, stateDB); len(value) > 0 {
		log.Info("CheckDuplicateSign exist", "blockNumber", blockNumber, "addr", hex.EncodeToString(addr.Bytes()), "dupType", dupType, "txHash", hex.EncodeToString(value))
		return value, nil
	}
	return nil, nil
}

func (sp *SlashingPlugin) putSlashResult(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType, stateDB xcom.StateDB) {
	stateDB.SetState(vm.SlashingContractAddr, duplicateSignKey(addr, blockNumber, dupType), stateDB.TxHash().Bytes())
}

func (sp *SlashingPlugin) getSlashResult(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType, stateDB xcom.StateDB) []byte {
	return stateDB.GetState(vm.SlashingContractAddr, duplicateSignKey(addr, blockNumber, dupType))
}

// duplicate signature result key format addr+blockNumber+_+etype
func duplicateSignKey(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType) []byte {
	value := append(addr.Bytes(), utils.Uint64ToBytes(blockNumber)...)
	value = append(value, []byte("_")...)
	value = append(value, common.Uint16ToBytes(uint16(dupType))...)
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
