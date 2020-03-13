// Copyright 2018-2019 The PlatON Network Authors
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

package plugin

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/x/slashing"

	"github.com/pkg/errors"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/consensus"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

const (
	HundredDenominator     = 100
	TenThousandDenominator = 10000
)

var (
	// The prefix key of the number of blocks packed in the recording node
	packAmountPrefix = []byte("nodePackAmount")
	once             sync.Once
	slash            *SlashingPlugin
)

type SlashingPlugin struct {
	db             snapshotdb.DB
	decodeEvidence func(dupType consensus.EvidenceType, data string) (consensus.Evidence, error)
	privateKey     *ecdsa.PrivateKey
}

func SlashInstance() *SlashingPlugin {
	once.Do(func() {
		log.Info("Init Slashing plugin ...")
		slash = &SlashingPlugin{
			db: snapshotdb.Instance(),
		}
	})
	return slash
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
			log.Error("Failed to BeginBlock,  call switchEpoch is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
			return err
		}
	}
	if err := sp.setPackAmount(blockHash, header); nil != err {
		log.Error("Failed to BeginBlock, call setPackAmount is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
		return err
	}
	// If it is the 230th block of each round,
	// it will punish the node with abnormal block rate.
	// Do this from the second consensus round
	if header.Number.Uint64() > xutil.ConsensusSize() && xutil.IsElection(header.Number.Uint64()) {
		log.Debug("Call GetPrePackAmount", "blockNumber", header.Number.Uint64(), "blockHash",
			blockHash.TerminalString(), "consensusSize", xutil.ConsensusSize(), "electionDistance", xcom.ElectionDistance())
		if result, err := sp.GetPrePackAmount(header.Number.Uint64(), header.ParentHash); nil != err {
			return err
		} else {
			if nil == result {
				log.Error("Failed to BeginBlock, call GetPrePackAmount is failed, the result is nil", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString())
				return errors.New("packAmount data not found")
			}

			preRoundVal, err := stk.getPreValList(blockHash, header.Number.Uint64(), QueryStartIrr)
			if nil != err {
				log.Error("Failed to BeginBlock, query previous round validators is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
				return err
			}

			slashQueue := make(staking.SlashQueue, 0)

			blockReward, err := gov.GovernSlashBlocksReward(header.Number.Uint64(), blockHash)
			if nil != err {
				log.Error("Failed to BeginBlock, query GovernSlashBlocksReward is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
				return err
			}

			for _, validator := range preRoundVal.Arr {
				nodeId := validator.NodeId
				count := result[nodeId]
				if count > 0 {
					continue
				}
				slashType := staking.LowRatioDel
				slashAmount := common.Big0

				canMutable, err := stk.GetCanMutableByIrr(validator.NodeAddress)
				if nil != err {
					log.Error("Failed to BeginBlock, call candidate mutable info is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
					if err == snapshotdb.ErrNotFound {
						continue
					}
					return err
				}
				totalBalance := calcCanTotalBalance(header.Number.Uint64(), canMutable)
				if blockReward > 0 {
					slashAmount, err = calcSlashBlockRewards(sp.db, blockHash, uint64(blockReward))
					if nil != err {
						log.Error("Failed to BeginBlock, call calcSlashBlockRewards fail", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
					}
					if slashAmount.Cmp(totalBalance) > 0 {
						slashAmount = totalBalance
					}
				}
				log.Info("Need to call SlashCandidates anomalous nodes", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "nodeId", nodeId.TerminalString(),
					"packBlockCount", count, "slashType", slashType, "totalBalance", totalBalance, "slashAmount", slashAmount, "SlashBlocksReward", blockReward)

				slashItem := &staking.SlashNodeItem{
					NodeId:      nodeId,
					Amount:      slashAmount,
					SlashType:   slashType,
					BenefitAddr: vm.RewardManagerPoolAddr,
				}

				slashQueue = append(slashQueue, slashItem)
			}

			// Real to slash the node
			// If there is no record of the node,
			// it means that there is no block,
			// then the penalty is directly
			if len(slashQueue) != 0 {
				if err := stk.SlashCandidates(state, blockHash, header.Number.Uint64(), slashQueue...); nil != err {
					log.Error("Failed to BeginBlock, call SlashCandidates is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
					return err
				}
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
	if snapshotdb.NonDbNotFoundErr(err) {
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
		log.Debug("Call setPackAmount finished", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString(), "nodeId", nodeId.TerminalString(), "value", value)
	}
	return nil
}

func (sp *SlashingPlugin) switchEpoch(blockNumber uint64, blockHash common.Hash) error {

	iter := sp.db.Ranking(blockHash, buildPrefixByRound(xutil.CalculateRound(blockNumber)-2), 0)
	if err := iter.Error(); nil != err {
		return err
	}
	defer iter.Release()
	count := 0
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		log.Debug("Call switchEpoch ranking old", "blockNumber", blockNumber, "key", hex.EncodeToString(key), "value", common.BytesToUint32(value))
		if err := sp.db.Del(blockHash, key); nil != err {
			return err
		}
		count++
	}
	log.Info("Call switchEpoch finished", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "count", count)
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
		log.Error("Failed to Slash, evidence validate is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "err", err)
		return slashing.ErrDuplicateSignVerify
	}
	if evidence.BlockNumber() > blockNumber {
		log.Error("Failed to Slash, Evidence is higher than the current block height",
			"blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "evidenceBlockNumber", evidence.BlockNumber())
		return slashing.ErrBlockNumberTooHigh
	}
	evidenceEpoch := xutil.CalculateEpoch(evidence.BlockNumber())
	blocksOfEpoch := xutil.CalcBlocksEachEpoch()
	invalidNum := evidenceEpoch * blocksOfEpoch
	if invalidNum < blockNumber {

		evidenceAge, err := gov.GovernMaxEvidenceAge(blockNumber, blockHash)
		if nil != err {
			log.Error("Failed to Slash, query Gov SlashFractionDuplicateSign is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"err", err)
			return err
		}

		if validSize := blocksOfEpoch * uint64(evidenceAge); blockNumber-invalidNum > validSize {
			log.Warn("Failed to Slash, Evidence time expired", "blockNumber", blockNumber,
				"blockHash", blockHash.TerminalString(), "evidenceBlockNum", evidence.BlockNumber(),
				"blocksOfEpoch", blocksOfEpoch, "the end blockNum of evidenceEpoch", invalidNum)
			return slashing.ErrIntervalTooLong
		}
	}
	if slashTxHash := sp.getSlashTxHash(evidence.Address(), evidence.BlockNumber(), evidence.Type(), stateDB); len(slashTxHash) > 0 {
		log.Error("Failed to Slash, the evidence had slashed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"evidenceBlockNumber", evidence.BlockNumber(), "evidenceHash", hex.EncodeToString(evidence.Hash()),
			"canAddr", evidence.Address().String(), "evidence type", evidence.Type(), "err", slashing.ErrSlashingExist)
		return slashing.ErrSlashingExist
	}

	canBase, err := stk.GetCanBase(blockHash, evidence.Address())
	if nil != err {
		log.Error("Failed to Slash, query CandidateBase info is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"evidenceBlockNumber", evidence.BlockNumber(), "canAddr", evidence.Address().String(), "err", err)
		return slashing.ErrGetCandidate
	}

	if canBase.IsEmpty() {
		log.Error("Failed to Slash, the candidate info is nil", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"canAddr", hex.EncodeToString(evidence.Address().Bytes()), "evidence type", evidence.Type())
		return slashing.ErrGetCandidate
	}

	if caller == canBase.StakingAddress {
		log.Error("Failed to Slash, can't report self", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"nodeId", canBase.NodeId.TerminalString(), "stakingAddress", caller.String(), "evidence type", evidence.Type())
		return slashing.ErrSameAddr
	}

	if canBase.NodeId != evidence.NodeID() {
		log.Error("Failed to Slash, Mismatch nodeId", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"can nodeId", canBase.NodeId.TerminalString(), "evidence nodeId", evidence.NodeID().TerminalString(), "evidence type", evidence.Type())
		return slashing.ErrNodeIdMismatch
	}

	pk, err := canBase.NodeId.Pubkey()
	if nil != err {
		log.Error("slashing failed candidate nodeId parse fail", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"nodeId", canBase.NodeId.TerminalString(), "type", evidence.Type(), "err", err)
		return slashing.ErrDuplicateSignVerify
	}
	addr := crypto.PubkeyToAddress(*pk)
	if addr != evidence.Address() {
		log.Error("slashing failed Mismatch addr", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "candidateNodeId", canBase.NodeId.TerminalString(),
			"candidateAddr", addr.Hex(), "evidenceAddr", evidence.Address().Hex(), "type", evidence.Type())
		return slashing.ErrAddrMismatch
	}

	blsKey, _ := canBase.BlsPubKey.ParseBlsPubKey()
	if !bytes.Equal(blsKey.Serialize(), evidence.BlsPubKey().Serialize()) {
		log.Error("Failed to Slash, Mismatch blsPubKey", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"nodeId", canBase.NodeId.TerminalString(), "can blsKey", hex.EncodeToString(blsKey.Serialize()),
			"evidence blsKey", hex.EncodeToString(evidence.BlsPubKey().Serialize()), "evidence type", evidence.Type())
		return slashing.ErrBlsPubKeyMismatch
	}

	if has, err := stk.checkRoundValidatorAddr(blockHash, evidence.BlockNumber(), evidence.Address()); nil != err {
		log.Error("Failed to Slash, checkRoundValidatorAddr is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"evidenceBlockNum", evidence.BlockNumber(), "canAddr", evidence.Address().Hex(), "err", err)
		return slashing.ErrDuplicateSignVerify
	} else if !has {
		log.Error("Failed to Slash, this node is not a validator, maybe!", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"evidenceBlockNum", evidence.BlockNumber(), "canAddr", evidence.Address().Hex())
		return slashing.ErrNotValidator
	}

	canMutable, err := stk.GetCanMutable(blockHash, evidence.Address())
	if nil != err {
		log.Error("Failed to Slash, query CandidateMutable info is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"evidenceBlockNumber", evidence.BlockNumber(), "canAddr", evidence.Address().String(), "err", err)
		return slashing.ErrGetCandidate
	}

	fraction, err := gov.GovernSlashFractionDuplicateSign(blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to Slash, query Gov SlashFractionDuplicateSign is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"err", err)
		return err
	}

	rewardFraction, err := gov.GovernDuplicateSignReportReward(blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to Slash, query Gov DuplicateSignReportReward is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"err", err)
		return err
	}

	totalBalance := calcCanTotalBalance(blockNumber, canMutable)
	slashAmount := calcAmountByRate(totalBalance, uint64(fraction), TenThousandDenominator)

	log.Info("Call SlashCandidates on executeSlash", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
		"nodeId", canBase.NodeId.TerminalString(), "totalBalance", totalBalance, "fraction", fraction, "rewardFraction", rewardFraction,
		"slashAmount", slashAmount, "reporter", caller.Hex())

	toCallerAmount := calcAmountByRate(slashAmount, uint64(rewardFraction), HundredDenominator)
	toCallerItem := &staking.SlashNodeItem{
		NodeId:      canBase.NodeId,
		Amount:      toCallerAmount,
		SlashType:   staking.DuplicateSign,
		BenefitAddr: caller,
	}

	toRewardPoolAmount := new(big.Int).Sub(slashAmount, toCallerAmount)
	toRewardPoolItem := &staking.SlashNodeItem{
		NodeId:      canBase.NodeId,
		Amount:      toRewardPoolAmount,
		SlashType:   staking.DuplicateSign,
		BenefitAddr: vm.RewardManagerPoolAddr,
	}

	if err := stk.SlashCandidates(stateDB, blockHash, blockNumber, toCallerItem, toRewardPoolItem); nil != err {
		log.Error("Failed to Slash, call SlashCandidates is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"nodeId", canBase.NodeId.TerminalString(), "err", err)
		return slashing.ErrSlashingFail
	}
	sp.putSlashTxHash(evidence.Address(), evidence.BlockNumber(), evidence.Type(), stateDB)
	log.Info("Call Slash finished", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
		"evidenceBlockNum", evidence.BlockNumber(), "nodeId", canBase.NodeId.TerminalString(), "evidenceType", evidence.Type(),
		"the txHash", stateDB.TxHash().TerminalString())

	return nil
}

func (sp *SlashingPlugin) CheckDuplicateSign(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType, stateDB xcom.StateDB) ([]byte, error) {
	if value := sp.getSlashTxHash(addr, blockNumber, dupType, stateDB); len(value) > 0 {
		return value, nil
	}
	return nil, nil
}

func (sp *SlashingPlugin) putSlashTxHash(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType, stateDB xcom.StateDB) {
	stateDB.SetState(vm.SlashingContractAddr, duplicateSignKey(addr, blockNumber, dupType), stateDB.TxHash().Bytes())
}

func (sp *SlashingPlugin) getSlashTxHash(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType, stateDB xcom.StateDB) []byte {
	return stateDB.GetState(vm.SlashingContractAddr, duplicateSignKey(addr, blockNumber, dupType))
}

// duplicate signature result key format addr+blockNumber+_+type
func duplicateSignKey(addr common.Address, blockNumber uint64, dupType consensus.EvidenceType) []byte {
	return append(append(addr.Bytes(), common.Uint64ToBytes(blockNumber)...), common.Uint16ToBytes(uint16(dupType))...)
}

func buildKey(blockNumber uint64, key []byte) []byte {
	return append(buildPrefix(blockNumber), key...)
}

func buildPrefix(blockNumber uint64) []byte {
	return buildPrefixByRound(xutil.CalculateRound(blockNumber))
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

func calcCanTotalBalance(blockNumber uint64, candidate *staking.CandidateMutable) *big.Int {
	// Recalculate the quality deposit
	lazyCalcStakeAmount(xutil.CalculateEpoch(blockNumber), candidate)
	return new(big.Int).Add(candidate.Released, candidate.RestrictingPlan)
}

func calcAmountByRate(balance *big.Int, numerator, denominator uint64) *big.Int {
	if balance.Cmp(common.Big0) > 0 {
		amount := new(big.Int).Mul(balance, new(big.Int).SetUint64(numerator))
		return amount.Div(amount, new(big.Int).SetUint64(denominator))
	}
	return new(big.Int).SetInt64(0)
}

func calcSlashBlockRewards(db snapshotdb.DB, hash common.Hash, blockRewardAmount uint64) (*big.Int, error) {
	newBlockReward, err := LoadNewBlockReward(hash, db)
	if nil != err {
		return nil, err
	}
	return new(big.Int).Mul(newBlockReward, new(big.Int).SetUint64(blockRewardAmount)), nil
}
