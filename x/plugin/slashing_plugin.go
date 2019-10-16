package plugin

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"sync"

	"github.com/pkg/errors"

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
	// The prefix key of the number of blocks packed in the recording node
	packAmountPrefix = []byte("nodePackAmount")

	errDuplicateSignVerify = common.NewBizError(303000, "duplicate signature verification failed")
	errSlashingExist       = common.NewBizError(303001, "punishment has been implemented")
	errBlockNumberTooHigh  = common.NewBizError(303002, "blockNumber too high")
	errIntervalTooLong     = common.NewBizError(303003, "evidence interval is too long")
	errGetCandidate        = common.NewBizError(303004, "failed to get certifier information")
	errAddrMismatch        = common.NewBizError(303005, "address does not match")
	errNodeIdMismatch      = common.NewBizError(303006, "nodeId does not match")
	errBlsPubKeyMismatch   = common.NewBizError(303007, "blsPubKey does not match")
	errSlashingFail        = common.NewBizError(303008, "slashing node fail")
	errNotValidator        = common.NewBizError(303009, "This node is not a validator")
	errSameAddr            = common.NewBizError(303010, "Can't report yourself")

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

func (sp *SlashingPlugin) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	sp.privateKey = privateKey
}

func (sp *SlashingPlugin) SetDecodeEvidenceFun(f func(dupType consensus.EvidenceType, data string) (consensus.Evidence, error)) {
	sp.decodeEvidence = f
}

func (sp *SlashingPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	// If it is the first block in each round, Delete old pack amount record.
	// Do this from the second consensus round
	if xutil.IsBeginOfConsensus(header.Number.Uint64()) && header.Number.Uint64() > 1 {
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
	// Do this from the second consensus round
	if header.Number.Uint64() > xutil.ConsensusSize() && xutil.IsElection(header.Number.Uint64()) {
		log.Debug("slashingPlugin Ranking block amount", "blockNumber", header.Number.Uint64(), "blockHash",
			blockHash.TerminalString(), "consensusSize", xutil.ConsensusSize(), "electionDistance", xcom.ElectionDistance())
		if result, err := sp.GetPrePackAmount(header.Number.Uint64(), header.ParentHash); nil != err {
			return err
		} else {
			if nil == result {
				log.Error("Failed to slashingPlugin GetPrePackAmount is nil", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString())
				return errors.New("packAmount data not found")
			}
			preRoundValArr, err := stk.GetCandidateONRound(blockHash, header.Number.Uint64(), PreviousRound, QueryStartIrr)
			if nil != err {
				log.Error("Failed to slashingPlugin BeginBlock, call GetCandidateONRound is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString())
				return err
			}

			slashQueue := make(staking.SlashQueue, 0)

			for _, validator := range preRoundValArr {
				nodeId := validator.NodeId
				amount := result[nodeId]
				if amount > xcom.PackAmountAbnormal() {
					continue
				}
				var slashType int
				if amount == 0 {
					slashType = staking.LowRatioDel
				} else {
					slashType = staking.LowRatio
				}
				slashAmount := calcEndBlockSlashAmount(header.Number.Uint64(), state)
				sumAmount := calcSumAmount(header.Number.Uint64(), validator)
				if slashAmount.Cmp(sumAmount) > 0 {
					slashAmount = sumAmount
				}
				log.Info("Need to call SlashCandidates anomalous nodes", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "nodeId", nodeId.TerminalString(),
					"packAmount", amount, "slashType", slashType, "slashAmount", slashAmount, "sumAmount", sumAmount, "NumberOfBlockRewardForSlashing", xcom.NumberOfBlockRewardForSlashing())

				slashItem := &staking.SlashNodeItem{
					NodeId:      nodeId,
					Amount:      slashAmount,
					SlashType:   slashType,
					BenefitAddr: vm.RewardManagerPoolAddr,
				}

				slashQueue = append(slashQueue, slashItem)
			}

			// Real to slash the node
			// If there is no record of the node, it means that there is no block, then the penalty is directly
			if err := stk.SlashCandidates(state, blockHash, header.Number.Uint64(), slashQueue...); nil != err {
				log.Error("Failed to slashingPlugin SlashCandidates failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
				return err
			}

		}
	}
	return nil
}

func (sp *SlashingPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	return nil
}

func (sp *SlashingPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

func (sp *SlashingPlugin) getPackAmount(blockNumber uint64, blockHash common.Hash, nodeId discover.NodeID) (uint32, error) {
	value, err := sp.db.Get(blockHash, buildKey(blockNumber, nodeId.Bytes()))
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
	if value, err := sp.getPackAmount(header.Number.Uint64(), blockHash, nodeId); nil != err {
		return err
	} else {
		value++
		if err := sp.db.Put(blockHash, buildKey(header.Number.Uint64(), nodeId.Bytes()), common.Uint32ToBytes(value)); nil != err {
			return err
		}
		log.Debug("slashingPlugin setPackAmount success", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "nodeId", nodeId.TerminalString(), "value", value)
	}
	return nil
}

func (sp *SlashingPlugin) switchEpoch(blockNumber uint64, blockHash common.Hash) error {
	count := 0
	iter := sp.db.Ranking(blockHash, buildPrefixByRound(xutil.CalculateRound(blockNumber)-2), 0)
	if err := iter.Error(); nil != err {
		return err
	}
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		log.Debug("slashingPlugin switchEpoch ranking old", "blockNumber", blockNumber, "key", hex.EncodeToString(key), "value", common.BytesToUint32(value))
		if err := sp.db.Del(blockHash, key); nil != err {
			return err
		}
		count++
	}
	log.Info("slashingPlugin switchEpoch success", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "count", count)
	return nil
}

// Get the consensus rate of all nodes in the previous round
func (sp *SlashingPlugin) GetPrePackAmount(blockNumber uint64, parentHash common.Hash) (map[discover.NodeID]uint32, error) {
	result := make(map[discover.NodeID]uint32)
	prefixKey := buildPrefixByRound(xutil.CalculateRound(blockNumber) - 1)
	iter := sp.db.Ranking(parentHash, prefixKey, 0)

	if err := iter.Error(); nil != err {
		return nil, err
	}
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		amount := common.BytesToUint32(value)
		nodeId, err := getNodeId(prefixKey, key)
		if nil != err {
			return nil, err
		}
		log.Debug("slashingPlugin GetPrePackAmount", "parentHash", parentHash.Hex(), "nodeId", nodeId.TerminalString(), "value", amount)
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
		log.Error("slashing failed evidence validate failed", "blockNumber", blockNumber, "err", err)
		return errDuplicateSignVerify
	}
	if evidence.BlockNumber() > blockNumber {
		log.Warn("slashing failed Evidence is higher than the current block height", "currBlockNumber", blockNumber, "evidenceBlockNumber", evidence.BlockNumber())
		return errBlockNumberTooHigh
	}
	epoch := xutil.CalculateEpoch(evidence.BlockNumber())
	blockAmount := xutil.CalcBlocksEachEpoch()
	evidenceEpochEndBlockNumber := epoch * blockAmount
	if evidenceEpochEndBlockNumber < blockNumber {
		if (blockNumber - evidenceEpochEndBlockNumber) > (blockAmount * uint64(xcom.EvidenceValidEpoch())) {
			log.Warn("slashing failed Evidence time expired", "currBlockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "evidenceBlockNumber", evidence.BlockNumber(),
				"EpochBlockAmount", blockAmount, "evidenceEpochEndBlockNumber", evidenceEpochEndBlockNumber)
			return errIntervalTooLong
		}
	}
	if value := sp.getSlashResult(evidence.Address(), evidence.BlockNumber(), evidence.Type(), stateDB); len(value) > 0 {
		log.Warn("slashing failed", "evidenceBlockNumber", evidence.BlockNumber(), "evidenceHash", hex.EncodeToString(evidence.Hash()),
			"addr", hex.EncodeToString(evidence.Address().Bytes()), "type", evidence.Type(), "err", errSlashingExist.Error())
		return errSlashingExist
	}
	if candidate, err := stk.GetCandidateInfo(blockHash, evidence.Address()); nil != err {
		log.Error("slashing failed", "evidenceBlockNumber", evidence.BlockNumber(), "blockHash", blockHash.TerminalString(), "addr", hex.EncodeToString(evidence.Address().Bytes()), "err", err)
		return errGetCandidate
	} else {
		if nil == candidate {
			log.Error("slashing failed GetCandidateInfo is nil", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"addr", hex.EncodeToString(evidence.Address().Bytes()), "type", evidence.Type())
			return errGetCandidate
		}
		if bytes.Equal(caller.Bytes(), candidate.StakingAddress.Bytes()) {
			log.Error("slashing failed Can't report yourself", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"nodeId", candidate.NodeId.TerminalString(), "stakingAddress", caller.Hex(), "type", evidence.Type())
			return errSameAddr
		}
		pk, err := candidate.NodeId.Pubkey()
		if nil != err {
			log.Error("slashing failed candidate nodeId parse fail", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"nodeId", candidate.NodeId.TerminalString(), "type", evidence.Type(), "err", err)
			return errDuplicateSignVerify
		}
		addr := crypto.PubkeyToAddress(*pk)
		if !bytes.Equal(addr.Bytes(), evidence.Address().Bytes()) {
			log.Error("slashing failed Mismatch addr", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "candidateNodeId", candidate.NodeId.TerminalString(),
				"candidateAddr", addr.Hex(), "evidenceAddr", evidence.Address().Hex(), "type", evidence.Type())
			return errAddrMismatch
		}
		if !bytes.Equal(candidate.NodeId.Bytes(), evidence.NodeID().Bytes()) {
			log.Error("slashing failed Mismatch nodeId", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"candidateNodeId", candidate.NodeId.TerminalString(), "evidenceAddr", evidence.NodeID().TerminalString(), "type", evidence.Type())
			return errNodeIdMismatch
		}
		if !bytes.Equal(candidate.BlsPubKey.Serialize(), evidence.BlsPubKey().Serialize()) {
			log.Error("slashing failed Mismatch blsPubKey", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "candidateNodeId", candidate.NodeId.TerminalString(),
				"candidateBlsPubKey", hex.EncodeToString(candidate.BlsPubKey.Serialize()), "evidenceBlsPubKey", hex.EncodeToString(evidence.BlsPubKey().Serialize()), "type", evidence.Type())
			return errBlsPubKeyMismatch
		}
		if isExists, err := stk.checkRoundValidatorAddr(blockHash, evidence.BlockNumber(), evidence.Address()); nil != err {
			log.Error("slashing failed checkRoundValidatorAddr", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"evidenceBlockNumber", evidence.BlockNumber(), "addr", evidence.Address().Hex(), "err", err)
			return errDuplicateSignVerify
		} else if !isExists {
			log.Warn("slashing failed, This node is not a validator", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"evidenceBlockNumber", evidence.BlockNumber(), "addr", evidence.Address().Hex())
			return errNotValidator
		}
		sumAmount := calcSumAmount(blockNumber, candidate)
		slashAmount := calcSlashAmount(sumAmount, xcom.DuplicateSignHighSlash())
		log.Info("Call SlashCandidates on executeSlash", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"nodeId", candidate.NodeId.String(), "sumAmount", sumAmount, "rate", xcom.DuplicateSignHighSlash(), "slashAmount", slashAmount, "reporter", caller.Hex())

		toCallerAmount := calcSlashAmount(slashAmount, xcom.DuplicateSignReportReward())
		toCallerItem := &staking.SlashNodeItem{
			NodeId:      candidate.NodeId,
			Amount:      toCallerAmount,
			SlashType:   staking.DuplicateSign,
			BenefitAddr: caller,
		}

		toRewardPoolAmount := new(big.Int).Sub(slashAmount, toCallerAmount)
		toRewardPoolItem := &staking.SlashNodeItem{
			NodeId:      candidate.NodeId,
			Amount:      toRewardPoolAmount,
			SlashType:   staking.DuplicateSign,
			BenefitAddr: vm.RewardManagerPoolAddr,
		}

		if err := stk.SlashCandidates(stateDB, blockHash, blockNumber, toCallerItem, toRewardPoolItem); nil != err {
			log.Error("slashing failed SlashCandidates failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"nodeId", hex.EncodeToString(candidate.NodeId.Bytes()), "err", err)
			return errSlashingFail
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

// duplicate signature result key format addr+blockNumber+_+type
func duplicateSignKey(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType) []byte {
	value := append(addr.Bytes(), utils.Uint64ToBytes(blockNumber)...)
	value = append(value, []byte("_")...)
	value = append(value, common.Uint16ToBytes(uint16(dupType))...)
	return value
}

func buildKey(blockNumber uint64, key []byte) []byte {
	return append(buildPrefix(blockNumber), key...)
}

func buildPrefix(blockNumber uint64) []byte {
	round := xutil.CalculateRound(blockNumber)
	return buildPrefixByRound(round)
}

func buildPrefixByRound(round uint64) []byte {
	return append(packAmountPrefix, common.Uint64ToBytes(round)...)
}

func getNodeId(prefix []byte, key []byte) (discover.NodeID, error) {
	key = key[len(prefix):]
	nodeId, err := discover.BytesID(key)
	if nil != err {
		return discover.NodeID{}, err
	}
	return nodeId, nil
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

func calcSumAmount(blockNumber uint64, candidate *staking.Candidate) *big.Int {
	// Recalculate the quality deposit
	lazyCalcStakeAmount(xutil.CalculateEpoch(blockNumber), candidate)
	return new(big.Int).Add(candidate.Released, candidate.RestrictingPlan)
}

func calcSlashAmount(sumAmount *big.Int, rate uint32) *big.Int {
	if sumAmount.Cmp(common.Big0) > 0 {
		amount := new(big.Int).Mul(sumAmount, new(big.Int).SetUint64(uint64(rate)))
		return amount.Div(amount, new(big.Int).SetUint64(100))
	}
	return new(big.Int).SetInt64(0)
}

func calcEndBlockSlashAmount(blockNumber uint64, state xcom.StateDB) *big.Int {
	thisYear := xutil.CalculateYear(blockNumber)
	var lastYear uint32
	if thisYear != 0 {
		lastYear = thisYear - 1
	}
	_, newBlockReward := RewardMgrInstance().calculateExpectReward(thisYear, lastYear, state)
	num := xcom.NumberOfBlockRewardForSlashing()
	return new(big.Int).Mul(newBlockReward, new(big.Int).SetUint64(uint64(num)))
}
