// Copyright 2021 The PlatON Network Authors
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

package vm

import (
	"fmt"
	"math"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"

	"github.com/PlatONnetwork/PlatON-Go/x/reward"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/node"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

const (
	TxCreateStaking      = 1000
	TxEditorCandidate    = 1001
	TxIncreaseStaking    = 1002
	TxWithdrewCandidate  = 1003
	TxDelegate           = 1004
	TxWithdrewDelegation = 1005
	TxRedeemDelegation   = 1006
	QueryVerifierList    = 1100
	QueryValidatorList   = 1101
	QueryCandidateList   = 1102
	QueryRelateList      = 1103
	QueryDelegateInfo    = 1104
	QueryCandidateInfo   = 1105
	QueryDelegationLock  = 1106
	GetPackageReward     = 1200
	GetStakingReward     = 1201
	GetAvgPackTime       = 1202
)

const (
	BLSPUBKEYLEN = 96 //  the bls public key length must be 96 byte
	BLSPROOFLEN  = 64 // the bls proof length must be 64 byte
)

type StakingContract struct {
	Plugin   *plugin.StakingPlugin
	Contract *Contract
	Evm      *EVM
}

func (stkc *StakingContract) RequiredGas(input []byte) uint64 {
	if checkInputEmpty(input) {
		return 0
	}
	return params.StakingGas
}

func (stkc *StakingContract) Run(input []byte) ([]byte, error) {
	if checkInputEmpty(input) {
		return nil, nil
	}
	if stkc.Evm.chainRules.IsEinstein {
		return execPlatonContract(input, stkc.FnSigns())
	}
	return execPlatonContract(input, stkc.FnSignsV1())
}

func (stkc *StakingContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	return nil
}

func (stkc *StakingContract) FnSignsV1() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		TxCreateStaking:      stkc.createStaking,
		TxEditorCandidate:    stkc.editCandidate,
		TxIncreaseStaking:    stkc.increaseStaking,
		TxWithdrewCandidate:  stkc.withdrewStaking,
		TxDelegate:           stkc.delegate,
		TxWithdrewDelegation: stkc.withdrewDelegation,

		// Get
		QueryVerifierList:  stkc.getVerifierList,
		QueryValidatorList: stkc.getValidatorList,
		QueryCandidateList: stkc.getCandidateList,
		QueryRelateList:    stkc.getRelatedListByDelAddr,
		QueryDelegateInfo:  stkc.getDelegateInfo,
		QueryCandidateInfo: stkc.getCandidateInfo,

		GetPackageReward: stkc.getPackageReward,
		GetStakingReward: stkc.getStakingReward,
		GetAvgPackTime:   stkc.getAvgPackTime,
	}
}

func (stkc *StakingContract) FnSigns() map[uint16]interface{} {
	fnSigns := stkc.FnSignsV1()
	fnSigns[TxRedeemDelegation] = stkc.redeemDelegation
	fnSigns[QueryDelegationLock] = stkc.getDelegateLock
	return fnSigns
}

func (stkc *StakingContract) createStaking(typ uint16, benefitAddress common.Address, nodeId enode.IDv0,
	externalId, nodeName, website, details string, amount *big.Int, rewardPer uint16, programVersion uint32,
	programVersionSign common.VersionSign, blsPubKey bls.PublicKeyHex, blsProof bls.SchnorrProofHex) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	txIndex := stkc.Evm.StateDB.TxIdx()
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call createStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "typ", typ,
		"benefitAddress", benefitAddress.String(), "nodeId", nodeId.String(), "externalId", externalId,
		"nodeName", nodeName, "website", website, "details", details, "amount", amount, "rewardPer", rewardPer,
		"programVersion", programVersion, "programVersionSign", programVersionSign.Hex(),
		"from", from, "blsPubKey", blsPubKey, "blsProof", blsProof)

	if !stkc.Contract.UseGas(params.CreateStakeGas) {
		return nil, ErrOutOfGas
	}

	if !verifyRewardPer(rewardPer) {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("invalid param rewardPer: %d", rewardPer),
			TxCreateStaking, staking.ErrInvalidRewardPer)
	}

	if len(blsPubKey) != BLSPUBKEYLEN {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("got blsKey length: %d, must be: %d", len(blsPubKey), BLSPUBKEYLEN),
			TxCreateStaking, staking.ErrWrongBlsPubKey)
	}

	if len(blsProof) != BLSPROOFLEN {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("got blsProof length: %d, must be: %d", len(blsProof), BLSPROOFLEN),
			TxCreateStaking, staking.ErrWrongBlsPubKeyProof)
	}

	// parse bls publickey
	blsPk, err := blsPubKey.ParseBlsPubKey()
	if nil != err {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("failed to parse blspubkey: %s", err.Error()),
			TxCreateStaking, staking.ErrWrongBlsPubKey)
	}

	// verify bls proof
	if err := verifyBlsProof(blsProof, blsPk); nil != err {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("failed to verify bls proof: %s", err.Error()),
			TxCreateStaking, staking.ErrWrongBlsPubKeyProof)

	}

	// validate programVersion sign
	if !node.GetCryptoHandler().IsSignedByNodeID(programVersion, programVersionSign.Bytes(), nodeId) {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			"call IsSignedByNodeID is failed",
			TxCreateStaking, staking.ErrWrongProgramVersionSign)
	}

	if ok, threshold := plugin.CheckStakeThreshold(blockNumber.Uint64(), blockHash, amount); !ok {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("staking threshold: %d, deposit: %d", threshold, amount),
			TxCreateStaking, staking.ErrStakeVonTooLow)
	}

	// check Description length
	desc := &staking.Description{
		NodeName:   nodeName,
		ExternalId: externalId,
		Website:    website,
		Details:    details,
	}
	if err := desc.CheckLength(); nil != err {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			staking.ErrDescriptionLen.Msg+":"+err.Error(),
			TxCreateStaking, staking.ErrDescriptionLen)
	}

	// Query current active version
	originVersion := gov.GetVersionForStaking(blockHash, state)
	currVersion := xutil.CalcVersion(originVersion)
	inputVersion := xutil.CalcVersion(programVersion)

	var isDeclareVersion bool

	realVersion := programVersion
	// Compare version
	// Just like that:
	// eg: 2.1.x == 2.1.x; 2.1.x > 2.0.x
	if inputVersion < currVersion {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("input Version: %s, current valid Version: %s",
				xutil.ProgramVersion2Str(programVersion), xutil.ProgramVersion2Str(originVersion)),
			TxCreateStaking, staking.ErrProgramVersionTooLow)

	} else if inputVersion > currVersion {
		isDeclareVersion = true
		//If the node version is higher than the current governance version, temporarily use the governance version,  wait for the version to pass the governance proposal, and then replace it
		realVersion = originVersion
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to createStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("nodeid %s to address fail: %s",
				nodeId.String(), err.Error()),
			TxCreateStaking, staking.ErrNodeID2Addr)
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to createStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if canOld.IsNotEmpty() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			"can is not nil",
			TxCreateStaking, staking.ErrCanAlreadyExist)
	}
	if txHash == common.ZeroHash {
		return nil, nil
	}
	/**
	init candidate info
	*/
	canBase := &staking.CandidateBase{
		NodeId:          nodeId,
		BlsPubKey:       blsPubKey,
		StakingAddress:  from,
		BenefitAddress:  benefitAddress,
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  txIndex,
		ProgramVersion:  realVersion,
		Description:     *desc,
	}

	canMutable := &staking.CandidateMutable{
		Shares:               amount,
		Released:             new(big.Int).SetInt64(0),
		ReleasedHes:          new(big.Int).SetInt64(0),
		RestrictingPlan:      new(big.Int).SetInt64(0),
		RestrictingPlanHes:   new(big.Int).SetInt64(0),
		RewardPer:            rewardPer,
		NextRewardPer:        rewardPer,
		RewardPerChangeEpoch: uint32(xutil.CalculateEpoch(blockNumber.Uint64())),
		DelegateRewardTotal:  new(big.Int).SetInt64(0),
	}

	can := &staking.Candidate{}
	can.CandidateBase = canBase
	can.CandidateMutable = canMutable

	err = stkc.Plugin.CreateCandidate(state, blockHash, blockNumber, amount, typ, canAddr, can)

	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {

			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
				bizErr.Error(), TxCreateStaking, bizErr)

		} else {
			log.Error("Failed to createStaking by CreateCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	// Because we must need to staking before we declare the version information.
	if isDeclareVersion {
		// Declare new Version
		err := gov.DeclareVersion(can.StakingAddress, can.NodeId,
			programVersion, programVersionSign, blockHash, blockNumber.Uint64(), stkc.Plugin, state)
		if nil != err {
			log.Error("Failed to CreateCandidate with govplugin DelareVersion failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)

			// the snapshot db can roll back ,so rollBack here no need
			/*if er := stkc.Plugin.RollBackStaking(state, blockHash, blockNumber, canAddr, typ); nil != er {
				log.Error("Failed to createStaking by RollBackStaking", "txHash", txHash,
					"blockNumber", blockNumber, "err", er)
			}*/

			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
				err.Error(), TxCreateStaking, staking.ErrDeclVsFialedCreateCan)
		}
	}

	return txResultHandler(vm.StakingContractAddr, stkc.Evm, "",
		"", TxCreateStaking, common.NoErr)
}

func verifyBlsProof(proofHex bls.SchnorrProofHex, pubKey *bls.PublicKey) error {
	proofByte, err := proofHex.MarshalText()
	if nil != err {
		return err
	}
	// proofHex to proof
	proof := new(bls.SchnorrProof)
	if err = proof.UnmarshalText(proofByte); nil != err {
		return err
	}
	// verify proof
	return proof.VerifySchnorrNIZK(*pubKey)
}

func verifyRewardPer(rewardPer uint16) bool {
	return rewardPer <= 10000 //	1BP(BasePoint)=0.01%
}

func (stkc *StakingContract) editCandidate(benefitAddress *common.Address, nodeId enode.IDv0, rewardPer *uint16,
	externalId, nodeName, website, details *string) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash
	from := stkc.Contract.CallerAddress

	log.Debug("Call editCandidate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
		"benefitAddress", benefitAddress, "nodeId", nodeId.String(), "rewardPer", rewardPer,
		"externalId", externalId, "nodeName", nodeName, "website", website,
		"details", details, "from", from)

	if !stkc.Contract.UseGas(params.EditCandidateGas) {
		return nil, ErrOutOfGas
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to editCandidate by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("nodeid %s to address fail: %s",
				nodeId.String(), err.Error()),
			TxCreateStaking, staking.ErrNodeID2Addr)
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to editCandidate by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return nil, err
	}

	if canOld.IsEmpty() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "editCandidate",
			"can is nil", TxEditorCandidate, staking.ErrCanNoExist)
	}

	if canOld.IsInvalid() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "editCandidate",
			fmt.Sprintf("can status is: %d", canOld.Status),
			TxEditorCandidate, staking.ErrCanStatusInvalid)
	}

	//发起修改交易的钱包地址，必须和发起质押的钱包地址一致
	if from != canOld.StakingAddress {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "editCandidate",
			fmt.Sprintf("contract sender: %s, can stake addr: %s", from, canOld.StakingAddress),
			TxEditorCandidate, staking.ErrNoSameStakingAddr)
	}

	//修改收益地址
	if benefitAddress != nil && canOld.BenefitAddress != vm.RewardManagerPoolAddr {
		canOld.BenefitAddress = *benefitAddress
	}

	if nodeName != nil {
		canOld.Description.NodeName = *nodeName
	}
	if externalId != nil {
		canOld.Description.ExternalId = *externalId
	}
	if website != nil {
		canOld.Description.Website = *website
	}
	if details != nil {
		canOld.Description.Details = *details
	}
	if err := canOld.Description.CheckLength(); nil != err {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "editCandidate",
			staking.ErrDescriptionLen.Msg+":"+err.Error(),
			TxEditorCandidate, staking.ErrDescriptionLen)
	}

	currentEpoch := uint32(xutil.CalculateEpoch(blockNumber.Uint64()))
	if currentEpoch > canOld.RewardPerChangeEpoch && canOld.NextRewardPer != canOld.RewardPer {
		canOld.RewardPer = canOld.NextRewardPer
	}

	if rewardPer != nil && *rewardPer != canOld.NextRewardPer {
		//分红比例修改时，和原有比例不能变化太大
		if !verifyRewardPer(*rewardPer) {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "editCandidate",
				fmt.Sprintf("invalid rewardPer: %d", rewardPer),
				TxEditorCandidate, staking.ErrInvalidRewardPer)
		}

		rewardPerMaxChangeRange, err := gov.GovernRewardPerMaxChangeRange(blockNumber.Uint64(), blockHash)
		if nil != err {
			log.Error("Failed to editCandidate, call GovernRewardPerMaxChangeRange is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"err", err)
			return nil, err
		}
		//分红比例修改时，不能太频繁。要和上次修改间隔一定的epoch
		rewardPerChangeInterval, err := gov.GovernRewardPerChangeInterval(blockNumber.Uint64(), blockHash)
		if nil != err {
			log.Error("Failed to editCandidate, call GovernRewardPerChangeInterval is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"err", err)
			return nil, err
		}

		if uint32(rewardPerChangeInterval) > currentEpoch-canOld.RewardPerChangeEpoch {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "editCandidate",
				fmt.Sprintf("needs interval [%d] epoch to modify", rewardPerChangeInterval),
				TxEditorCandidate, staking.ErrRewardPerInterval)
		}

		canOld.NextRewardPer = *rewardPer
		difference := uint16(math.Abs(float64(canOld.NextRewardPer) - float64(canOld.RewardPer)))
		//分红比例修改时，和原有比例不能变化太大
		if difference > rewardPerMaxChangeRange {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "editCandidate",
				fmt.Sprintf("invalid rewardPer: %d, modified by more than: %d", rewardPer, rewardPerMaxChangeRange),
				TxEditorCandidate, staking.ErrRewardPerChangeRange)
		}
		canOld.RewardPerChangeEpoch = currentEpoch
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}
	err = stkc.Plugin.EditCandidate(blockHash, blockNumber, canAddr, canOld)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "editCandidate",
				bizErr.Error(), TxEditorCandidate, bizErr)
		} else {
			log.Error("Failed to editCandidate by EditCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	return txResultHandler(vm.StakingContractAddr, stkc.Evm, "",
		"", TxEditorCandidate, common.NoErr)
}

func (stkc *StakingContract) increaseStaking(nodeId enode.IDv0, typ uint16, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call increaseStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String(), "typ", typ,
		"amount", amount, "from", from)

	if !stkc.Contract.UseGas(params.IncStakeGas) {
		return nil, ErrOutOfGas
	}

	if ok, threshold := plugin.CheckOperatingThreshold(blockNumber.Uint64(), blockHash, amount); !ok {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "increaseStaking",
			fmt.Sprintf("increase staking threshold: %d, deposit: %d", threshold, amount),
			TxIncreaseStaking, staking.ErrIncreaseStakeVonTooLow)
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to increaseStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("nodeid %s to address fail: %s",
				nodeId.String(), err.Error()),
			TxCreateStaking, staking.ErrNodeID2Addr)
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to increaseStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if canOld.IsEmpty() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "increaseStaking",
			"can is nil", TxIncreaseStaking, staking.ErrCanNoExist)
	}

	if canOld.IsInvalid() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "increaseStaking",
			fmt.Sprintf("can status is: %d", canOld.Status),
			TxIncreaseStaking, staking.ErrCanStatusInvalid)
	}

	if from != canOld.StakingAddress {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "increaseStaking",
			fmt.Sprintf("contract sender: %s, can stake addr: %s", from, canOld.StakingAddress),
			TxIncreaseStaking, staking.ErrNoSameStakingAddr)
	}
	if txHash == common.ZeroHash {
		return nil, nil
	}

	err = stkc.Plugin.IncreaseStaking(state, blockHash, blockNumber, amount, typ, canAddr, canOld)

	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "increaseStaking",
				bizErr.Error(), TxIncreaseStaking, bizErr)

		} else {
			log.Error("Failed to increaseStaking by EditCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}
	return txResultHandler(vm.StakingContractAddr, stkc.Evm, "",
		"", TxIncreaseStaking, common.NoErr)
}

func (stkc *StakingContract) withdrewStaking(nodeId enode.IDv0) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call withdrewStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String(), "from", from)

	if !stkc.Contract.UseGas(params.WithdrewStakeGas) {
		return nil, ErrOutOfGas
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to withdrewStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("nodeid %s to address fail: %s",
				nodeId.String(), err.Error()),
			TxCreateStaking, staking.ErrNodeID2Addr)
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to withdrewStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	if canOld.IsEmpty() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "withdrewStaking",
			"can is nil", TxWithdrewCandidate, staking.ErrCanNoExist)
	}

	if canOld.IsInvalid() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "withdrewStaking",
			fmt.Sprintf("can status is: %d", canOld.Status),
			TxWithdrewCandidate, staking.ErrCanStatusInvalid)
	}

	if from != canOld.StakingAddress {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "withdrewStaking",
			fmt.Sprintf("contract sender: %s, can stake addr: %s", from, canOld.StakingAddress),
			TxWithdrewCandidate, staking.ErrNoSameStakingAddr)
	}
	if txHash == common.ZeroHash {
		return nil, nil
	}
	err = stkc.Plugin.WithdrewStaking(state, blockHash, blockNumber, canAddr, canOld)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "withdrewStaking",
				bizErr.Error(), TxWithdrewCandidate, bizErr)
		} else {
			log.Error("Failed to withdrewStaking by WithdrewStaking", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}

	return txResultHandler(vm.StakingContractAddr, stkc.Evm, "",
		"", TxWithdrewCandidate, common.NoErr)
}

func (stkc *StakingContract) delegate(typ uint16, nodeId enode.IDv0, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call delegate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from, "typ", typ,
		"nodeId", nodeId.String(), "amount", amount)

	if !stkc.Contract.UseGas(params.DelegateGas) {
		return nil, ErrOutOfGas
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to delegate by parse nodeId", "txHash", txHash, "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "createStaking",
			fmt.Sprintf("nodeid %s to address fail: %s",
				nodeId.String(), err.Error()),
			TxCreateStaking, staking.ErrNodeID2Addr)
	}

	canMutable, err := stkc.Plugin.GetCanMutable(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to delegate by GetCandidateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if canMutable.IsEmpty() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "delegate",
			"can is nil", TxDelegate, staking.ErrCanNoExist)
	}

	if canMutable.IsInvalid() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "delegate",
			fmt.Sprintf("can status is: %d", canMutable.Status),
			TxDelegate, staking.ErrCanStatusInvalid)
	}

	// the can base must exist if canMutable is exist,so no need check if canBase==nil
	canBase, err := stkc.Plugin.GetCanBase(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to delegate by GetCandidateBase", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if canBase.StakingBlockNum == blockNumber.Uint64() {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "delegate",
			fmt.Sprintf("delegate fail,can't not delgate in the staking block:%d", blockNumber.Uint64()),
			TxDelegate, staking.ErrCanNoExist)
	}

	del, err := stkc.Plugin.GetDelegateInfo(blockHash, from, nodeId, canBase.StakingBlockNum)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to delegate by GetDelegateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}
	if del.IsEmpty() {
		// build delegate
		del = staking.NewDelegation()
	}
	var delegateRewardPerList []*reward.DelegateRewardPer
	if del.DelegateEpoch > 0 {
		delegateRewardPerList, err = plugin.RewardMgrInstance().GetDelegateRewardPerList(blockHash, canBase.NodeId, canBase.StakingBlockNum, uint64(del.DelegateEpoch), xutil.CalculateEpoch(blockNumber.Uint64())-1)
		if snapshotdb.NonDbNotFoundErr(err) {
			log.Error("Failed to delegate by GetDelegateRewardPerList", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
		result, err := stkc.calcRewardPerUseGas(delegateRewardPerList, del)
		if nil != err {
			return result, err
		}
	}

	if ok, threshold := plugin.CheckOperatingThreshold(blockNumber.Uint64(), blockHash, amount); !ok {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "delegate",
			fmt.Sprintf("delegate threshold: %d, deposit: %d", threshold, amount),
			TxDelegate, staking.ErrDelegateVonTooLow)
	}

	// check account
	hasStake, err := stkc.Plugin.HasStake(blockHash, from)
	if nil != err {
		return nil, err
	}

	if hasStake {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "delegate",
			fmt.Sprintf("'%s' has staking, so don't allow to delegate", from),
			TxDelegate, staking.ErrAccountNoAllowToDelegate)
	}

	// If the candidate’s benefitaAddress is the RewardManagerPoolAddr, no delegation is allowed
	if canBase.BenefitAddress == vm.RewardManagerPoolAddr {
		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "delegate",
			"the can benefitAddr is reward addr",
			TxDelegate, staking.ErrCanNoAllowDelegate)
	}
	if txHash == common.ZeroHash {
		return nil, nil
	}

	can := &staking.Candidate{}
	can.CandidateBase = canBase
	can.CandidateMutable = canMutable

	err = stkc.Plugin.Delegate(state, blockHash, blockNumber, from, del, canAddr, can, typ, amount, delegateRewardPerList, stkc.Evm.chainRules.IsEinstein)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "delegate",
				bizErr.Error(), TxDelegate, bizErr)
		} else {
			log.Error("Failed to delegate by Delegate", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	return txResultHandler(vm.StakingContractAddr, stkc.Evm, "",
		"", TxDelegate, common.NoErr)
}

// 撤消委托，当撤消某个节点的全部委托时，委托奖励将立刻发放到委托用户账户；当撤消某个接的部分委托时，只是计算委托奖励。
//	stakingBlockNum + nodeId 确定一个质押的节点
// param: stakingBlockNum	代表着某个node的某次质押的唯一标示
// param: nodeId			被质押的节点的NodeId
// param: amount			减持的金额
// return:

func (stkc *StakingContract) withdrewDelegation(stakingBlockNum uint64, nodeId enode.IDv0, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call withdrewDelegation of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from, "nodeId", nodeId.String(),
		"stakingNum", stakingBlockNum, "amount", amount)

	if !stkc.Contract.UseGas(params.WithdrewDelegationGas) {
		return nil, ErrOutOfGas
	}

	del, err := stkc.Plugin.GetDelegateInfo(blockHash, from, nodeId, stakingBlockNum)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to withdrewDelegation by GetDelegateInfo",
			"txHash", txHash.Hex(), "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if del.IsEmpty() {
		if txHash == common.ZeroHash {
			return nil, nil
		} else {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "withdrewDelegation",
				"del is nil", TxWithdrewDelegation, staking.ErrDelegateNoExist)
		}
	}

	delegateRewardPerList, err := plugin.RewardMgrInstance().GetDelegateRewardPerList(blockHash, nodeId, stakingBlockNum, uint64(del.DelegateEpoch), xutil.CalculateEpoch(blockNumber.Uint64())-1)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to delegate by GetDelegateRewardPerList", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	result, err := stkc.calcRewardPerUseGas(delegateRewardPerList, del)
	if nil != err {
		return result, err
	}

	if ok, threshold := plugin.CheckOperatingThreshold(blockNumber.Uint64(), blockHash, amount); !ok {

		return txResultHandler(vm.StakingContractAddr, stkc.Evm, "withdrewDelegation",
			fmt.Sprintf("withdrewDelegation threshold: %d, deposit: %d", threshold, amount),
			TxWithdrewDelegation, staking.ErrWithdrewDelegationVonTooLow)
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	issueIncome, released, restrictingPlan, lockReleased, lockRestrictingPlan, err := stkc.Plugin.WithdrewDelegation(state, blockHash, blockNumber, amount, txHash, from, nodeId, stakingBlockNum, del, delegateRewardPerList, stkc.Evm.chainRules.IsEinstein)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {

			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "withdrewDelegation",
				bizErr.Error(), TxWithdrewDelegation, bizErr)

		} else {
			log.Error("Failed to withdrewDelegation by WithdrewDelegation", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}
	if stkc.Evm.chainRules.IsEinstein {
		return txResultHandlerWithRes(vm.StakingContractAddr, stkc.Evm, "",
			"", TxWithdrewDelegation, int(common.NoErr.Code), issueIncome, released, restrictingPlan, lockReleased, lockRestrictingPlan), nil
	} else {
		return txResultHandlerWithRes(vm.StakingContractAddr, stkc.Evm, "",
			"", TxWithdrewDelegation, int(common.NoErr.Code), issueIncome), nil
	}
}

func (stkc *StakingContract) redeemDelegation() ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call redeemDelegation of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from)

	if !stkc.Contract.UseGas(params.RedeemDelegationGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	released, restrictingPlan, err := stkc.Plugin.RedeemDelegation(state, blockHash, blockNumber, from)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			return txResultHandler(vm.StakingContractAddr, stkc.Evm, "redeemDelegation",
				bizErr.Error(), TxRedeemDelegation, bizErr)
		} else {
			log.Error("Failed to redeemDelegation by RedeemDelegation", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}
	log.Debug("Call redeemDelegation of stakingContract finished", "restrictingPlan", restrictingPlan, "released", released)

	return txResultHandlerWithRes(vm.StakingContractAddr, stkc.Evm, "",
		"", TxRedeemDelegation, int(common.NoErr.Code), released, restrictingPlan), nil
}

func (stkc *StakingContract) calcRewardPerUseGas(delegateRewardPerList []*reward.DelegateRewardPer, del *staking.Delegation) ([]byte, error) {
	unCalcEpoch := len(delegateRewardPerList)
	if unCalcEpoch > 0 {
		if delegateRewardPerList[0].Epoch == uint64(del.DelegateEpoch) {
			if del.Released.Cmp(common.Big0) == 0 && del.RestrictingPlan.Cmp(common.Big0) == 0 {
				unCalcEpoch -= 1
			}
		}
		if !stkc.Contract.UseGas(params.WithdrawDelegateEpochGas * uint64(unCalcEpoch)) {
			return nil, ErrOutOfGas
		}
	}
	return nil, nil
}

func (stkc *StakingContract) getVerifierList() ([]byte, error) {

	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash

	arr, err := stkc.Plugin.GetVerifierList(blockHash, blockNumber.Uint64(), plugin.QueryStartNotIrr)

	if snapshotdb.NonDbNotFoundErr(err) {
		return callResultHandler(stkc.Evm, "getVerifierList",
			arr, staking.ErrGetVerifierList.Wrap(err.Error())), nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || arr.IsEmpty() {
		return callResultHandler(stkc.Evm, "getVerifierList",
			arr, staking.ErrGetVerifierList.Wrap("VerifierList info is not found")), nil
	}

	return callResultHandler(stkc.Evm, "getVerifierList",
		arr, nil), nil
}

func (stkc *StakingContract) getValidatorList() ([]byte, error) {

	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash

	arr, err := stkc.Plugin.GetValidatorList(blockHash, blockNumber.Uint64(), plugin.CurrentRound, plugin.QueryStartNotIrr)
	if snapshotdb.NonDbNotFoundErr(err) {

		return callResultHandler(stkc.Evm, "getValidatorList",
			arr, staking.ErrGetValidatorList.Wrap(err.Error())), nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || arr.IsEmpty() {
		return callResultHandler(stkc.Evm, "getValidatorList",
			arr, staking.ErrGetValidatorList.Wrap("ValidatorList info is not found")), nil
	}

	return callResultHandler(stkc.Evm, "getValidatorList",
		arr, nil), nil
}

func (stkc *StakingContract) getCandidateList() ([]byte, error) {

	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash

	arr, err := stkc.Plugin.GetCandidateList(blockHash, blockNumber.Uint64())
	if snapshotdb.NonDbNotFoundErr(err) {
		return callResultHandler(stkc.Evm, "getCandidateList",
			arr, staking.ErrGetCandidateList.Wrap(err.Error())), nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || arr.IsEmpty() {
		return callResultHandler(stkc.Evm, "getCandidateList",
			arr, staking.ErrGetCandidateList.Wrap("CandidateList info is not found")), nil
	}

	return callResultHandler(stkc.Evm, "getCandidateList",
		arr, nil), nil
}

func (stkc *StakingContract) getRelatedListByDelAddr(addr common.Address) ([]byte, error) {

	blockHash := stkc.Evm.Context.BlockHash
	arr, err := stkc.Plugin.GetRelatedListByDelAddr(blockHash, addr)
	if snapshotdb.NonDbNotFoundErr(err) {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getRelatedListByDelAddr, delAddr: %s", addr),
			arr, staking.ErrGetDelegateRelated.Wrap(err.Error())), nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || arr.IsEmpty() {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getRelatedListByDelAddr, delAddr: %s", addr),
			arr, staking.ErrGetDelegateRelated.Wrap("RelatedList info is not found")), nil
	}

	return callResultHandler(stkc.Evm, fmt.Sprintf("getRelatedListByDelAddr, delAddr: %s", addr),
		arr, nil), nil
}

func (stkc *StakingContract) getDelegateInfo(stakingBlockNum uint64, delAddr common.Address,
	nodeId enode.IDv0) ([]byte, error) {

	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash

	del, err := stkc.Plugin.GetDelegateExCompactInfo(blockHash, blockNumber.Uint64(), delAddr, nodeId, stakingBlockNum)
	if snapshotdb.NonDbNotFoundErr(err) {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getDelegateInfo, delAddr: %s, nodeId: %s, stakingBlockNumber: %d",
			delAddr, nodeId, stakingBlockNum),
			del, staking.ErrQueryDelegateInfo.Wrap(err.Error())), nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || del.IsEmpty() {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getDelegateInfo, delAddr: %s, nodeId: %s, stakingBlockNumber: %d",
			delAddr, nodeId, stakingBlockNum),
			del, staking.ErrQueryDelegateInfo.Wrap("Delegate info is not found")), nil
	}
	if stkc.Evm.chainRules.IsEinstein {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getDelegateInfo, delAddr: %s, nodeId: %s, stakingBlockNumber: %d",
			delAddr, nodeId, stakingBlockNum),
			del, nil), nil
	}
	return callResultHandler(stkc.Evm, fmt.Sprintf("getDelegateInfo, delAddr: %s, nodeId: %s, stakingBlockNumber: %d",
		delAddr, nodeId, stakingBlockNum),
		del.V1(), nil), nil
}

func (stkc *StakingContract) getDelegateLock(delAddr common.Address) ([]byte, error) {
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash
	locks, err := stkc.Plugin.GetGetDelegationLockCompactInfo(blockHash, blockNumber.Uint64(), delAddr)

	if err != nil {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getDelegateLock, delAddr: %s", delAddr),
			nil, staking.ErrQueryDelegationLockInfo.Wrap(err.Error())), nil
	}

	return callResultHandler(stkc.Evm, fmt.Sprintf("getDelegateLock, delAddr: %s",
		delAddr),
		locks, nil), nil
}

func (stkc *StakingContract) getCandidateInfo(nodeId enode.IDv0) ([]byte, error) {
	blockNumber := stkc.Evm.Context.BlockNumber
	blockHash := stkc.Evm.Context.BlockHash

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getCandidateInfo, nodeId: %s",
			nodeId), nil, staking.ErrQueryCandidateInfo.Wrap(err.Error())), nil
	}
	can, err := stkc.Plugin.GetCandidateCompactInfo(blockHash, blockNumber.Uint64(), canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getCandidateInfo, nodeId: %s",
			nodeId), can, staking.ErrQueryCandidateInfo.Wrap(err.Error())), nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || can.IsEmpty() {
		return callResultHandler(stkc.Evm, fmt.Sprintf("getCandidateInfo, nodeId: %s",
			nodeId), can, staking.ErrQueryCandidateInfo.Wrap("Candidate info is not found")), nil
	}

	return callResultHandler(stkc.Evm, fmt.Sprintf("getCandidateInfo, nodeId: %s",
		nodeId), can, nil), nil
}

func (stkc *StakingContract) getPackageReward() ([]byte, error) {
	packageReward, err := plugin.LoadNewBlockReward(common.ZeroHash, stkc.Evm.SnapshotDB)
	if nil != err {
		return callResultHandler(stkc.Evm, "getPackageReward", nil, common.NotFound.Wrap(err.Error())), nil
	}
	return callResultHandler(stkc.Evm, "getPackageReward", (*hexutil.Big)(packageReward), nil), nil
}

func (stkc *StakingContract) getStakingReward() ([]byte, error) {
	stakingReward, err := plugin.LoadStakingReward(common.ZeroHash, stkc.Evm.SnapshotDB)
	if nil != err {
		return callResultHandler(stkc.Evm, "getStakingReward", nil, common.NotFound.Wrap(err.Error())), nil
	}
	return callResultHandler(stkc.Evm, "getStakingReward", (*hexutil.Big)(stakingReward), nil), nil
}

func (stkc *StakingContract) getAvgPackTime() ([]byte, error) {
	avgPackTime, err := xcom.LoadCurrentAvgPackTime()
	if nil != err {
		return callResultHandler(stkc.Evm, "getAvgPackTime", nil, common.InternalError.Wrap(err.Error())), nil
	}
	return callResultHandler(stkc.Evm, "getAvgPackTime", avgPackTime, nil), nil
}
