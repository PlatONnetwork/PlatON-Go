package vm

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/node"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

const (
	CreateStakingEvent     = "1000"
	EditorCandidateEvent   = "1001"
	IncreaseStakingEvent   = "1002"
	WithdrewCandidateEvent = "1003"
	DelegateEvent          = "1004"
	WithdrewDelegateEvent  = "1005"
	BLSPUBKEYLEN           = 96 //  the bls public key length must be 96 byte
	BLSPROOFLEN            = 64 // the bls proof length must be 64 byte
)

type StakingContract struct {
	Plugin   *plugin.StakingPlugin
	Contract *Contract
	Evm      *EVM
}

func (stkc *StakingContract) RequiredGas(input []byte) uint64 {
	return params.StakingGas
}

func (stkc *StakingContract) Run(input []byte) ([]byte, error) {
	return execPlatonContract(input, stkc.FnSigns())
}

func (stkc *StakingContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	return nil
}

func (stkc *StakingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		1000: stkc.createStaking,
		1001: stkc.editCandidate,
		1002: stkc.increaseStaking,
		1003: stkc.withdrewStaking,
		1004: stkc.delegate,
		1005: stkc.withdrewDelegate,

		// Get
		1100: stkc.getVerifierList,
		1101: stkc.getValidatorList,
		1102: stkc.getCandidateList,
		1103: stkc.getRelatedListByDelAddr,
		1104: stkc.getDelegateInfo,
		1105: stkc.getCandidateInfo,
	}
}

func (stkc *StakingContract) createStaking(typ uint16, benefitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string, amount *big.Int, programVersion uint32,
	programVersionSign common.VersionSign, blsPubKey bls.PublicKeyHex, blsProof bls.SchnorrProofHex) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	txIndex := stkc.Evm.StateDB.TxIdx()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call createStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "typ", typ,
		"benefitAddress", benefitAddress.String(), "nodeId", nodeId.String(), "externalId", externalId,
		"nodeName", nodeName, "website", website, "details", details, "amount", amount,
		"programVersion", programVersion, "programVersionSign", programVersionSign.Hex(),
		"from", from.Hex(), "blsPubKey", blsPubKey, "blsProof", blsProof)

	if !stkc.Contract.UseGas(params.CreateStakeGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	if len(blsPubKey) != BLSPUBKEYLEN {
		receipt := strconv.Itoa(int(staking.ErrWrongBlsPubKey.Code))
		stkc.badLog(CreateStakingEvent, receipt,
			fmt.Sprintf("got length: %d, must be: %d", len(blsPubKey), BLSPUBKEYLEN), "createStaking")
		return []byte(receipt), nil
	}

	if len(blsProof) != BLSPROOFLEN {
		receipt := strconv.Itoa(int(staking.ErrWrongBlsPubKeyProof.Code))
		stkc.badLog(CreateStakingEvent, receipt,
			fmt.Sprintf("got length: %d, must be: %d", len(blsPubKey), BLSPUBKEYLEN), "createStaking")
		return []byte(receipt), nil
	}

	// parse bls publickey
	blsPk, err := blsPubKey.ParseBlsPubKey()
	if nil != err {
		receipt := strconv.Itoa(int(staking.ErrWrongBlsPubKey.Code))
		stkc.badLog(CreateStakingEvent, receipt,
			fmt.Sprintf("failed to parse blspubkey: %s", err.Error()), "createStaking")
		return []byte(receipt), nil
	}

	// verify bls proof
	if err := verifyBlsProof(blsProof, blsPk); nil != err {
		receipt := strconv.Itoa(int(staking.ErrWrongBlsPubKeyProof.Code))
		stkc.badLog(CreateStakingEvent, receipt,
			fmt.Sprintf("failed to verify bls proof: %s", err.Error()), "createStaking")
		return []byte(receipt), nil
	}

	// validate programVersion sign
	if !node.GetCryptoHandler().IsSignedByNodeID(programVersion, programVersionSign.Bytes(), nodeId) {
		receipt := strconv.Itoa(int(staking.ErrWrongProgramVersionSign.Code))
		stkc.badLog(CreateStakingEvent, receipt,
			"call IsSignedByNodeID is failed", "createStaking")
		return []byte(receipt), nil
	}

	if ok, threshold := plugin.CheckStakeThreshold(blockNumber.Uint64(), blockHash, amount); !ok {
		receipt := strconv.Itoa(int(staking.ErrStakeVonTooLow.Code))
		stkc.badLog(CreateStakingEvent, receipt,
			fmt.Sprintf("staking threshold: %d, deposit: %d", threshold, amount), "createStaking")
		return []byte(receipt), nil
	}

	// check Description length
	desc := &staking.Description{
		NodeName:   nodeName,
		ExternalId: externalId,
		Website:    website,
		Details:    details,
	}
	if err := desc.CheckLength(); nil != err {
		receipt := strconv.Itoa(int(staking.ErrDescriptionLen.Code))
		stkc.badLog(CreateStakingEvent, receipt,
			staking.ErrDescriptionLen.Msg+":"+err.Error(), "createStaking")
		return []byte(receipt), nil
	}

	// Query current active version
	originVersion := gov.GetVersionForStaking(state)
	currVersion := xutil.CalcVersion(originVersion)
	inputVersion := xutil.CalcVersion(programVersion)

	var isDeclareVersion bool

	// Compare version
	// Just like that:
	// eg: 2.1.x == 2.1.x; 2.1.x > 2.0.x
	if inputVersion < currVersion {
		err := fmt.Sprintf("input Version: %s, current valid Version: %s",
			xutil.ProgramVersion2Str(programVersion), xutil.ProgramVersion2Str(originVersion))
		receipt := strconv.Itoa(int(staking.ErrProgramVersionTooLow.Code))
		stkc.badLog(CreateStakingEvent, receipt, err, "createStaking")
		return []byte(receipt), nil

	} else if inputVersion > currVersion {
		isDeclareVersion = true
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to createStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to createStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if canOld.IsNotEmpty() {
		receipt := strconv.Itoa(int(staking.ErrCanAlreadyExist.Code))
		stkc.badLog(CreateStakingEvent, receipt,
			"can is not nil", "createStaking")
		return []byte(receipt), nil
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
		ProgramVersion:  currVersion,
		Description:     *desc,
	}

	canMutable := &staking.CandidateMutable{
		Shares:             amount,
		Released:           new(big.Int).SetInt64(0),
		ReleasedHes:        new(big.Int).SetInt64(0),
		RestrictingPlan:    new(big.Int).SetInt64(0),
		RestrictingPlanHes: new(big.Int).SetInt64(0),
	}

	can := &staking.Candidate{}
	can.CandidateBase = canBase
	can.CandidateMutable = canMutable

	err = stkc.Plugin.CreateCandidate(state, blockHash, blockNumber, amount, typ, canAddr, can)

	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			receipt := strconv.Itoa(int(bizErr.Code))
			stkc.badLog(CreateStakingEvent, receipt,
				fmt.Sprintf("failed to createStaking: %s", bizErr.Error()), "createStaking")
			return []byte(receipt), nil
		} else {
			log.Error("Failed to createStaking by CreateCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	if isDeclareVersion {
		// Declare new Version
		err := gov.DeclareVersion(can.StakingAddress, can.NodeId,
			programVersion, programVersionSign, blockHash, blockNumber.Uint64(), stkc.Plugin, state)
		if nil != err {
			log.Error("Failed to CreateCandidate with govplugin DelareVersion failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)

			if er := stkc.Plugin.RollBackStaking(state, blockHash, blockNumber, canAddr, typ); nil != er {
				log.Error("Failed to createStaking by RollBackStaking", "txHash", txHash,
					"blockNumber", blockNumber, "err", er)
			}

			receipt := strconv.Itoa(int(staking.ErrDeclVsFialedCreateCan.Code))
			stkc.badLog(CreateStakingEvent, receipt, err.Error(), "createStaking")
			return []byte(receipt), nil
		}
	}
	receipt := strconv.Itoa(int(common.NoErr.Code))
	stkc.goodLog(CreateStakingEvent, receipt, "createStaking")
	return []byte(receipt), nil
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

func (stkc *StakingContract) editCandidate(benefitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash
	from := stkc.Contract.CallerAddress

	log.Debug("Call editCandidate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
		"benefitAddress", benefitAddress.String(), "nodeId", nodeId.String(),
		"externalId", externalId, "nodeName", nodeName, "website", website,
		"details", details, "from", from.Hex())

	if !stkc.Contract.UseGas(params.EditCandidatGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to editCandidate by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to editCandidate by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return nil, err
	}

	if canOld.IsEmpty() {
		receipt := strconv.Itoa(int(staking.ErrCanNoExist.Code))
		stkc.badLog(EditorCandidateEvent, receipt,
			"can is nil", "editCandidate")
		return []byte(receipt), nil
	}

	if canOld.Is_Invalid() {
		receipt := strconv.Itoa(int(staking.ErrCanStatusInvalid.Code))
		stkc.badLog(EditorCandidateEvent, receipt,
			fmt.Sprintf("can status is: %d", canOld.Status), "editCandidate")
		return []byte(receipt), nil
	}

	if from != canOld.StakingAddress {
		receipt := strconv.Itoa(int(staking.ErrNoSameStakingAddr.Code))
		stkc.badLog(EditorCandidateEvent, receipt,
			fmt.Sprintf("contract sender: %s, can stake addr: %s", from.Hex(), canOld.StakingAddress.Hex()),
			"editCandidate")
		return []byte(receipt), nil
	}

	if canOld.BenefitAddress != vm.RewardManagerPoolAddr {
		canOld.BenefitAddress = benefitAddress
	}

	// check Description length
	desc := &staking.Description{
		NodeName:   nodeName,
		ExternalId: externalId,
		Website:    website,
		Details:    details,
	}
	if err := desc.CheckLength(); nil != err {
		err := common.NewBizError(staking.ErrDescriptionLen.Code, staking.ErrDescriptionLen.Msg+":"+err.Error())
		receipt := strconv.Itoa(int(staking.ErrDescriptionLen.Code))
		stkc.badLog(EditorCandidateEvent, receipt,
			err.Error(), "editCandidate")
		return []byte(receipt), nil
	}

	canOld.Description = *desc

	err = stkc.Plugin.EditCandidate(blockHash, blockNumber, canAddr, canOld)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			receipt := strconv.Itoa(int(bizErr.Code))
			stkc.badLog(EditorCandidateEvent, receipt,
				fmt.Sprintf("failed to editCandidate: %s", bizErr.Error()), "editCandidate")
			return []byte(receipt), nil
		} else {
			log.Error("Failed to editCandidate by EditCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}

	receipt := strconv.Itoa(int(common.NoErr.Code))
	stkc.goodLog(EditorCandidateEvent, receipt, "editCandidate")
	return []byte(receipt), nil
}

func (stkc *StakingContract) increaseStaking(nodeId discover.NodeID, typ uint16, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call increaseStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String(), "typ", typ,
		"amount", amount, "from", from.Hex())

	if !stkc.Contract.UseGas(params.IncStakeGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	if ok, threshold := plugin.CheckOperatingThreshold(blockNumber.Uint64(), blockHash, amount); !ok {
		receipt := strconv.Itoa(int(staking.ErrIncreaseStakeVonTooLow.Code))
		stkc.badLog(IncreaseStakingEvent, receipt,
			fmt.Sprintf("increase staking threshold: %d, deposit: %d", threshold, amount), "increaseStaking")
		return []byte(receipt), nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to increaseStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to increaseStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if canOld.IsEmpty() {
		receipt := strconv.Itoa(int(staking.ErrCanNoExist.Code))
		stkc.badLog(IncreaseStakingEvent, receipt,
			"can is nil", "increaseStaking")
		return []byte(receipt), nil
	}

	if canOld.Is_Invalid() {
		receipt := strconv.Itoa(int(staking.ErrCanStatusInvalid.Code))
		stkc.badLog(IncreaseStakingEvent, receipt,
			fmt.Sprintf("can status is: %d", canOld.Status), "increaseStaking")
		return []byte(receipt), nil
	}

	if from != canOld.StakingAddress {
		receipt := strconv.Itoa(int(staking.ErrNoSameStakingAddr.Code))
		stkc.badLog(IncreaseStakingEvent, receipt,
			fmt.Sprintf("contract sender: %s, can stake addr: %s", from.Hex(), canOld.StakingAddress.Hex()),
			"increaseStaking")
		return []byte(receipt), nil
	}

	err = stkc.Plugin.IncreaseStaking(state, blockHash, blockNumber, amount, typ, canAddr, canOld)

	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			receipt := strconv.Itoa(int(bizErr.Code))
			stkc.badLog(IncreaseStakingEvent, receipt,
				fmt.Sprintf("failed to increaseStaking: %s", bizErr.Error()), "increaseStaking")
			return []byte(receipt), nil
		} else {
			log.Error("Failed to increaseStaking by EditCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}
	receipt := strconv.Itoa(int(common.NoErr.Code))
	stkc.goodLog(IncreaseStakingEvent, receipt, "increaseStaking")
	return []byte(receipt), nil
}

func (stkc *StakingContract) withdrewStaking(nodeId discover.NodeID) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call withdrewStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String(), "from", from.Hex())

	if !stkc.Contract.UseGas(params.WithdrewStakeGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to withdrewStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to withdrewStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	if canOld.IsEmpty() {
		receipt := strconv.Itoa(int(staking.ErrCanNoExist.Code))
		stkc.badLog(WithdrewCandidateEvent, receipt,
			"can is nil", "withdrewStaking")
		return []byte(receipt), nil
	}

	if canOld.Is_Invalid() {
		receipt := strconv.Itoa(int(staking.ErrCanStatusInvalid.Code))
		stkc.badLog(WithdrewCandidateEvent, receipt,
			fmt.Sprintf("can status is: %d", canOld.Status), "withdrewStaking")
		return []byte(receipt), nil
	}

	if from != canOld.StakingAddress {
		receipt := strconv.Itoa(int(staking.ErrNoSameStakingAddr.Code))
		stkc.badLog(WithdrewCandidateEvent, receipt,
			fmt.Sprintf("contract sender: %s, can stake addr: %s", from.Hex(), canOld.StakingAddress.Hex()),
			"withdrewStaking")
		return []byte(receipt), nil
	}

	err = stkc.Plugin.WithdrewStaking(state, blockHash, blockNumber, canAddr, canOld)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			receipt := strconv.Itoa(int(bizErr.Code))
			stkc.badLog(WithdrewCandidateEvent, receipt,
				fmt.Sprintf("failed to withdrewStaking: %s", bizErr.Error()), "withdrewStaking")
			return []byte(receipt), nil
		} else {
			log.Error("Failed to withdrewStaking by WithdrewStaking", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}

	receipt := strconv.Itoa(int(common.NoErr.Code))
	stkc.goodLog(WithdrewCandidateEvent, receipt, "withdrewStaking")
	return []byte(receipt), nil
}

func (stkc *StakingContract) delegate(typ uint16, nodeId discover.NodeID, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call delegate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from.Hex(), "typ", typ,
		"nodeId", nodeId.String(), "amount", amount)

	if !stkc.Contract.UseGas(params.DelegateGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	if ok, threshold := plugin.CheckOperatingThreshold(blockNumber.Uint64(), blockHash, amount); !ok {
		receipt := strconv.Itoa(int(staking.ErrDelegateVonTooLow.Code))
		stkc.badLog(DelegateEvent, receipt,
			fmt.Sprintf("delegate threshold: %d, deposit: %d", threshold, amount), "delegate")
		return []byte(receipt), nil
	}

	// check account
	hasStake, err := stkc.Plugin.HasStake(blockHash, from)
	if nil != err {
		return nil, err
	}

	if hasStake {
		receipt := strconv.Itoa(int(staking.ErrAccountNoAllowToDelegate.Code))
		stkc.badLog(DelegateEvent, receipt,
			fmt.Sprintf("'%s' has staking, so don't allow to delegate", from.Hex()), "delegate")
		return []byte(receipt), nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to delegate by parse nodeId", "txHash", txHash, "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canMutable, err := stkc.Plugin.GetCanMutable(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to delegate by GetCandidateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if canMutable.IsEmpty() {
		receipt := strconv.Itoa(int(staking.ErrCanNoExist.Code))
		stkc.badLog(DelegateEvent, receipt,
			"can is nil", "delegate")
		return []byte(receipt), nil
	}

	if canMutable.Is_Invalid() {
		receipt := strconv.Itoa(int(staking.ErrCanStatusInvalid.Code))
		stkc.badLog(DelegateEvent, receipt,
			fmt.Sprintf("can status is: %d", canMutable.Status), "delegate")
		return []byte(receipt), nil
	}

	canBase, err := stkc.Plugin.GetCanBase(blockHash, canAddr)

	// If the candidate’s benefitaAddress is the RewardManagerPoolAddr, no delegation is allowed
	if canBase.BenefitAddress == vm.RewardManagerPoolAddr {
		receipt := strconv.Itoa(int(staking.ErrCanNoAllowDelegate.Code))
		stkc.badLog(DelegateEvent, receipt,
			"the can benefitAddr is reward addr", "delegate")
		return []byte(receipt), nil
	}

	del, err := stkc.Plugin.GetDelegateInfo(blockHash, from, nodeId, canBase.StakingBlockNum)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to delegate by GetDelegateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if del.IsEmpty() {
		// build delegate
		del = new(staking.Delegation)
		// Prevent null pointer initialization
		del.Released = new(big.Int).SetInt64(0)
		del.RestrictingPlan = new(big.Int).SetInt64(0)
		del.ReleasedHes = new(big.Int).SetInt64(0)
		del.RestrictingPlanHes = new(big.Int).SetInt64(0)
	}
	can := &staking.Candidate{}
	can.CandidateBase = canBase
	can.CandidateMutable = canMutable

	err = stkc.Plugin.Delegate(state, blockHash, blockNumber, from, del, canAddr, can, typ, amount)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			receipt := strconv.Itoa(int(bizErr.Code))
			stkc.badLog(DelegateEvent, receipt,
				fmt.Sprintf("failed to delegate: %s", bizErr.Error()), "delegate")
			return []byte(receipt), nil
		} else {
			log.Error("Failed to delegate by Delegate", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	receipt := strconv.Itoa(int(common.NoErr.Code))
	stkc.goodLog(DelegateEvent, receipt, "delegate")
	return []byte(receipt), nil
}

func (stkc *StakingContract) withdrewDelegate(stakingBlockNum uint64, nodeId discover.NodeID, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash
	from := stkc.Contract.CallerAddress
	state := stkc.Evm.StateDB

	log.Debug("Call withdrewDelegate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from.Hex(), "nodeId", nodeId.String(),
		"stakingNum", stakingBlockNum, "amount", amount)

	if !stkc.Contract.UseGas(params.WithdrewDelegateGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	if ok, threshold := plugin.CheckOperatingThreshold(blockNumber.Uint64(), blockHash, amount); !ok {
		receipt := strconv.Itoa(int(staking.ErrWithdrewDelegateVonTooLow.Code))
		stkc.badLog(WithdrewDelegateEvent, receipt,
			fmt.Sprintf("withdrewDelegate threshold: %d, deposit: %d", threshold, amount), "withdrewDelegate")
		return []byte(receipt), nil
	}

	del, err := stkc.Plugin.GetDelegateInfo(blockHash, from, nodeId, stakingBlockNum)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to withdrewDelegate by GetDelegateInfo",
			"txHash", txHash.Hex(), "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if del.IsEmpty() {
		receipt := strconv.Itoa(int(staking.ErrDelegateNoExist.Code))
		stkc.badLog(WithdrewDelegateEvent, receipt,
			"del is nil", "withdrewDelegate")
		return []byte(receipt), nil
	}

	err = stkc.Plugin.WithdrewDelegate(state, blockHash, blockNumber, amount, from, nodeId, stakingBlockNum, del)
	if nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			receipt := strconv.Itoa(int(bizErr.Code))
			stkc.badLog(WithdrewDelegateEvent, receipt,
				fmt.Sprintf("failed to withdrewDelegate: %s", bizErr.Error()), "withdrewDelegate")
			return []byte(receipt), nil
		} else {
			log.Error("Failed to withdrewDelegate by WithdrewDelegate", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}
	receipt := strconv.Itoa(int(common.NoErr.Code))
	stkc.goodLog(WithdrewDelegateEvent, receipt, "withdrewDelegate")
	return []byte(receipt), nil
}

func (stkc *StakingContract) getVerifierList() ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetVerifierList(blockHash, blockNumber.Uint64(), plugin.QueryStartIrr)

	if snapshotdb.NonDbNotFoundErr(err) {
		data := xcom.NewFailedResult(staking.ErrGetVerifierList.Wrap(err.Error()))
		log.Error("Failed to getVerifierList: Query VerifierList is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || arr.IsEmpty() {
		data := xcom.NewFailedResult(staking.ErrGetVerifierList.Wrap("VerifierList info is not found"))
		log.Error("Failed to getVerifierList: VerifierList info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex())
		return data, nil
	}

	arrByte, err := json.Marshal(arr)
	if nil != err {
		data := xcom.NewFailedResult(staking.ErrGetVerifierList.Wrap(err.Error()))
		log.Error("Failed to getVerifierList: VerifierList Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}
	data := xcom.NewOkResult(string(arrByte))
	log.Debug("getVerifierList", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "verArr", string(arrByte))
	return data, nil
}

func (stkc *StakingContract) getValidatorList() ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetValidatorList(blockHash, blockNumber.Uint64(), plugin.CurrentRound, plugin.QueryStartIrr)
	if snapshotdb.NonDbNotFoundErr(err) {
		data := xcom.NewFailedResult(staking.ErrGetValidatorList.Wrap(err.Error()))
		log.Error("Failed to getValidatorList: Query ValidatorList is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || arr.IsEmpty() {
		data := xcom.NewFailedResult(staking.ErrGetValidatorList.Wrap("ValidatorList info is not found"))
		log.Error("Failed to getValidatorList: ValidatorList info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex())
		return data, nil
	}

	arrByte, err := json.Marshal(arr)
	if nil != err {
		data := xcom.NewFailedResult(staking.ErrGetValidatorList.Wrap(err.Error()))
		log.Error("Failed to getValidatorList: ValidatorList Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}
	data := xcom.NewOkResult(string(arrByte))
	log.Debug("getValidatorList", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "valArr", string(arrByte))
	return data, nil
}

func (stkc *StakingContract) getCandidateList() ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetCandidateList(blockHash, blockNumber.Uint64())
	if snapshotdb.NonDbNotFoundErr(err) {
		data := xcom.NewFailedResult(staking.ErrGetCandidateList.Wrap(err.Error()))
		log.Error("Failed to getCandidateList: Query CandidateList is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || arr.IsEmpty() {
		data := xcom.NewFailedResult(staking.ErrGetCandidateList.Wrap("CandidateList info is not found"))
		log.Error("Failed to getCandidateList: CandidateList info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex())
		return data, nil
	}

	arrByte, err := json.Marshal(arr)
	if nil != err {
		data := xcom.NewFailedResult(staking.ErrGetCandidateList.Wrap(err.Error()))
		log.Error("Failed to getCandidateList: CandidateList Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}
	data := xcom.NewOkResult(string(arrByte))
	log.Debug("getCandidateList", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "canArr", string(arrByte))
	return data, nil
}

func (stkc *StakingContract) getRelatedListByDelAddr(addr common.Address) ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetRelatedListByDelAddr(blockHash, addr)
	if snapshotdb.NonDbNotFoundErr(err) {
		data := xcom.NewFailedResult(staking.ErrGetDelegateRelated.Wrap(err.Error()))
		log.Error("Failed to getRelatedListByDelAddr: Query RelatedList is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", addr.Hex(), "err", err)
		return data, nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || arr.IsEmpty() {
		data := xcom.NewFailedResult(staking.ErrGetDelegateRelated.Wrap("RelatedList info is not found"))
		log.Error("Failed to getRelatedListByDelAddr: RelatedList info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", addr.Hex())
		return data, nil
	}

	jsonByte, err := json.Marshal(arr)
	if nil != err {
		data := xcom.NewFailedResult(staking.ErrGetDelegateRelated.Wrap(err.Error()))
		log.Error("Failed to getRelatedListByDelAddr: RelatedList Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", addr.Hex(), "err", err)
		return data, nil
	}
	data := xcom.NewOkResult(string(jsonByte))
	log.Debug("getRelatedListByDelAddr", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"delAddr", addr.Hex(), "relateArr", string(jsonByte))
	return data, nil
}

func (stkc *StakingContract) getDelegateInfo(stakingBlockNum uint64, delAddr common.Address,
	nodeId discover.NodeID) ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	del, err := stkc.Plugin.GetDelegateExCompactInfo(blockHash, blockNumber.Uint64(), delAddr, nodeId, stakingBlockNum)
	if snapshotdb.NonDbNotFoundErr(err) {
		data := xcom.NewFailedResult(staking.ErrQueryDelegateInfo.Wrap(err.Error()))
		log.Error("Failed to getDelegateInfo: Query Delegate info is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNumber", stakingBlockNum, "err", err)
		return data, nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || del.IsEmpty() {
		data := xcom.NewFailedResult(staking.ErrQueryDelegateInfo.Wrap("Delegate info is not found"))
		log.Error("Failed to getDelegateInfo: Delegate info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNumber", stakingBlockNum)
		return data, nil
	}

	jsonByte, err := json.Marshal(del)
	if nil != err {
		data := xcom.NewFailedResult(staking.ErrQueryDelegateInfo.Wrap(err.Error()))
		log.Error("Failed to getDelegateInfo: Delegate Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNumber", stakingBlockNum, "err", err)
		return data, nil
	}
	data := xcom.NewOkResult(string(jsonByte))
	log.Debug("getDelegateInfo", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNumber", stakingBlockNum, "delinfo", string(jsonByte))
	return data, nil
}

func (stkc *StakingContract) getCandidateInfo(nodeId discover.NodeID) ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		data := xcom.NewFailedResult(staking.ErrQueryCandidateInfo.Wrap(err.Error()))
		log.Error("Failed to getCandidateInfo: Parse NodeId to Address is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return data, nil
	}
	can, err := stkc.Plugin.GetCandidateCompactInfo(blockHash, blockNumber.Uint64(), canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		data := xcom.NewFailedResult(staking.ErrQueryCandidateInfo.Wrap(err.Error()))
		log.Error("Failed to getCandidateInfo: Query Candidate info is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return data, nil
	}

	if snapshotdb.IsDbNotFoundErr(err) || can.IsEmpty() {
		data := xcom.NewFailedResult(staking.ErrQueryCandidateInfo.Wrap("Candidate info is not found"))
		log.Error("Failed to getCandidateInfo: Candidate info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return data, nil
	}

	jsonByte, err := json.Marshal(can)
	if nil != err {
		data := xcom.NewFailedResult(staking.ErrQueryCandidateInfo.Wrap(err.Error()))
		log.Error("Failed to getCandidateInfo: Candidate Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return data, nil
	}
	data := xcom.NewOkResult(string(jsonByte))
	log.Debug("getCandidateInfo", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"nodeId", nodeId.String(), "caninfo", string(jsonByte))
	return data, nil
}

func (stkc *StakingContract) goodLog(eventType, eventData, callFn string) {

	blockNumber := stkc.Evm.BlockNumber.Uint64()
	xcom.AddLog(stkc.Evm.StateDB, blockNumber, vm.StakingContractAddr, eventType, eventData)
}

func (stkc *StakingContract) badLog(eventType, eventData, reason, callFn string) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber.Uint64()
	xcom.AddLog(stkc.Evm.StateDB, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Error("Failed to "+callFn+" of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber, "receipt: ", eventData, "the reason", reason)
}
