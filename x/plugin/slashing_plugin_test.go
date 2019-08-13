package plugin_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"github.com/stretchr/testify/assert"
)

func initInfo(t *testing.T) (*plugin.SlashingPlugin, xcom.StateDB) {
	si := plugin.SlashInstance()
	plugin.StakingInstance()
	db := ethdb.NewMemDatabase()
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(db))
	if nil != err {
		t.Fatal(err)
	}
	return si, stateDB
}

func buildStakingData(blockHash common.Hash, pri *ecdsa.PrivateKey, t *testing.T, stateDb xcom.StateDB) {
	stakingDB := staking.NewStakingDB()

	sender := common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5234")

	buildDbRestrictingPlan(sender, t, stateDb)

	if nil == pri {
		sk, err := crypto.GenerateKey()
		if nil != err {
			panic(err)
		}
		pri = sk
	}

	nodeIdA := discover.PubkeyID(&pri.PublicKey)
	addrA, _ := xutil.NodeId2Addr(nodeIdA)

	nodeIdB := nodeIdArr[1]
	addrB, _ := xutil.NodeId2Addr(nodeIdB)

	nodeIdC := nodeIdArr[2]
	addrC, _ := xutil.NodeId2Addr(nodeIdC)

	var blsKey1 bls.SecretKey
	blsKey1.SetByCSPRNG()
	c1 := &staking.Candidate{
		NodeId:             nodeIdA,
		BlsPubKey:          *blsKey1.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[1],
		StakingTxIndex:     uint32(2),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "xxccccdddddddd",
			NodeName:   "I Am " + fmt.Sprint(1),
			Website:    "www.baidu.com",
			Details:    "this is  baidu ~~",
		},
	}

	var blsKey2 bls.SecretKey
	blsKey2.SetByCSPRNG()
	c2 := &staking.Candidate{
		NodeId:             nodeIdB,
		BlsPubKey:          *blsKey2.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[2],
		StakingTxIndex:     uint32(3),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "SFSFSFSFSFSFSSFS",
			NodeName:   "I Am " + fmt.Sprint(2),
			Website:    "www.JD.com",
			Details:    "this is  JD ~~",
		},
	}

	var blsKey3 bls.SecretKey
	blsKey3.SetByCSPRNG()
	c3 := &staking.Candidate{
		NodeId:             nodeIdC,
		BlsPubKey:          *blsKey3.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[3],
		StakingTxIndex:     uint32(4),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "FWAGGDGDGG",
			NodeName:   "I Am " + fmt.Sprint(3),
			Website:    "www.alibaba.com",
			Details:    "this is  alibaba ~~",
		},
	}

	stakingDB.SetCanPowerStore(blockHash, addrA, c1)
	stakingDB.SetCanPowerStore(blockHash, addrB, c2)
	stakingDB.SetCanPowerStore(blockHash, addrC, c3)

	stakingDB.SetCandidateStore(blockHash, addrA, c1)
	stakingDB.SetCandidateStore(blockHash, addrB, c2)
	stakingDB.SetCandidateStore(blockHash, addrC, c3)

	log.Info("addr_A", hex.EncodeToString(addrA.Bytes()), "addr_B", hex.EncodeToString(addrB.Bytes()), "addr_C", hex.EncodeToString(addrC.Bytes()))

	queue := make(staking.ValidatorQueue, 0)

	v1 := &staking.Validator{
		NodeAddress:   addrA,
		NodeId:        c1.NodeId,
		BlsPubKey:     c1.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c1.StakingBlockNum), fmt.Sprint(c1.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v2 := &staking.Validator{
		NodeAddress:   addrB,
		NodeId:        c2.NodeId,
		BlsPubKey:     c2.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c2.StakingBlockNum), fmt.Sprint(c2.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v3 := &staking.Validator{
		NodeAddress:   addrC,
		NodeId:        c3.NodeId,
		BlsPubKey:     c3.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c3.StakingBlockNum), fmt.Sprint(c3.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	queue = append(queue, v1)
	queue = append(queue, v2)
	queue = append(queue, v3)

	epochArr := &staking.Validator_array{
		Start: 1,
		End:   uint64(xutil.CalcBlocksEachEpoch()),
		Arr:   queue,
	}

	preArr := &staking.Validator_array{
		Start: 1,
		End:   xutil.ConsensusSize(),
		Arr:   queue,
	}

	curArr := &staking.Validator_array{
		Start: xutil.ConsensusSize() + 1,
		End:   xutil.ConsensusSize() * 2,
		Arr:   queue,
	}

	setVerifierList(blockHash, epochArr)
	setRoundValList(blockHash, preArr)
	setRoundValList(blockHash, curArr)
	balance, ok := new(big.Int).SetString("9999999999999999999999999999999999999999999999999", 10)
	if !ok {
		panic("set balance fail")
	}
	stateDb.AddBalance(vm.StakingContractAddr, balance)
}

func TestSlashingPlugin_BeginBlock(t *testing.T) {
	_, _, _ = newChainState()
	si, stateDB := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	startNumber := xutil.ConsensusSize()
	startNumber += xutil.ConsensusSize() - xcom.ElectionDistance() - 2
	pri, phash := confirmBlock(t, int(startNumber))
	startNumber++
	blockNumber := new(big.Int).SetInt64(int64(startNumber))
	if err := snapshotdb.Instance().NewBlock(blockNumber, phash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	buildStakingData(common.ZeroHash, pri, t, stateDB)

	phash = common.HexToHash("0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802130406")
	if err := snapshotdb.Instance().Flush(phash, blockNumber); err != nil {
		t.Fatal(err)
	}
	if err := snapshotdb.Instance().Commit(phash); err != nil {
		t.Fatal(err)
	}
	startNumber++
	header := &types.Header{
		Number: new(big.Int).SetUint64(uint64(startNumber)),
		Extra:  make([]byte, 97),
	}
	if err := snapshotdb.Instance().NewBlock(header.Number, phash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	if err := si.BeginBlock(common.ZeroHash, header, stateDB); nil != err {
		t.Fatal(err)
	}
}

func TestSlashingPlugin_Confirmed(t *testing.T) {
	si, _ := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	startNumber := xutil.ConsensusSize() + 1
	confirmBlock(t, int(startNumber))
	result, err := si.GetPreNodeAmount()
	if nil != err {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func confirmBlock(t *testing.T, maxNumber int) (*ecdsa.PrivateKey, common.Hash) {
	pri, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	pri2, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	//hash := common.HexToHash("0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802111216")
	db := snapshotdb.Instance()

	sk := pri

	_, genesis, _ := newChainState()
	parentHash := genesis.Hash()
	for i := 0; i < maxNumber; i++ {
		blockNum := big.NewInt(int64(i + 1))
		if i == int(xcom.PackAmountAbnormal()) {
			sk = pri2
		}
		header := &types.Header{
			Number: blockNum,
			Extra:  make([]byte, 97),
		}
		sign, err := crypto.Sign(header.SealHash().Bytes(), sk)
		if nil != err {
			t.Fatal(err)
		}
		copy(header.Extra[len(header.Extra)-common.ExtraSeal:], sign[:])
		block := types.NewBlock(header, nil, nil)
		if err := plugin.SlashInstance().Confirmed(block); nil != err {
			t.Fatal(err)
		}
		if err := db.NewBlock(blockNum, parentHash, common.ZeroHash); err != nil {
			t.Fatal(err)
		}
		if err := db.Flush(header.Hash(), blockNum); err != nil {
			t.Fatal(err)
		}
		if err := db.Commit(header.Hash()); err != nil {
			t.Fatal(err)
		}
		parentHash = header.Hash()
	}
	return pri, parentHash
}

func TestSlashingPlugin_Slash(t *testing.T) {
	_, genesis, _ := newChainState()
	si, stateDB := initInfo(t)
	blockNumber := new(big.Int).SetUint64(1)
	chash := common.HexToHash("0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802130406")
	if err := snapshotdb.Instance().NewBlock(blockNumber, genesis.Hash(), common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	buildStakingData(common.ZeroHash, nil, t, stateDB)
	if err := snapshotdb.Instance().Flush(chash, blockNumber); nil != err {
		t.Fatal(err)
	}
	if err := snapshotdb.Instance().Commit(chash); nil != err {
		t.Fatal(err)
	}
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	plugin.GovPluginInstance()
	si.SetDecodeEvidenceFun(evidence.NewEvidences)
	data := `{
         "duplicate_prepare": [
          {
           "PrepareA": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0x00a452c6116ac9df049016437f8a35b4e29c17d63632314f0266df2b0dcd4bef",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257",
             "NodeID": "bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c",
             "blsPubKey": "f3ee4cb60b04358c21460b9dd0832028959a6d0052218d796c96a5eac01b541f88595d62ee52880e0a77ecf8ffde63966a5d0d70028c08dfca622827563df99e"
            },
            "signature": "0x554a2a2f1b0d197730c707b595016b5f735ce0df0a5e9efd28a77764f295af1f"
           },
           "PrepareB": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0x3f643f315a72d54815e3b638e53a1f293834e6d9109c4c0f3e5d9c7171bf1cf2",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257",
             "NodeID": "bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c",
             "blsPubKey": "f3ee4cb60b04358c21460b9dd0832028959a6d0052218d796c96a5eac01b541f88595d62ee52880e0a77ecf8ffde63966a5d0d70028c08dfca622827563df99e"
            },
            "signature": "0x9e626bd0fd19290c7ff23a605259735de216f9c26df51ddaf51f66f0aade4097"
           }
          },
		  {
           "PrepareA": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0x00a452c6116ac9df049016437f8a35b4e29c17d63632314f0266df2b0dcd4bef",
            "block_number": 2,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257",
             "NodeID": "bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c",
             "blsPubKey": "f3ee4cb60b04358c21460b9dd0832028959a6d0052218d796c96a5eac01b541f88595d62ee52880e0a77ecf8ffde63966a5d0d70028c08dfca622827563df99e"
            },
            "signature": "0x554a2a2f1b0d197730c707b595016b5f735ce0df0a5e9efd28a77764f295af1f"
           },
           "PrepareB": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0x3f643f315a72d54815e3b638e53a1f293834e6d9109c4c0f3e5d9c7171bf1cf2",
            "block_number": 2,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257",
             "NodeID": "bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c",
             "blsPubKey": "f3ee4cb60b04358c21460b9dd0832028959a6d0052218d796c96a5eac01b541f88595d62ee52880e0a77ecf8ffde63966a5d0d70028c08dfca622827563df99e"
            },
            "signature": "0x9e626bd0fd19290c7ff23a605259735de216f9c26df51ddaf51f66f0aade4097"
           }
          },
		  {
           "PrepareA": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0x00a452c6116ac9df049016437f8a35b4e29c17d63632314f0266df2b0dcd4bef",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257",
             "NodeID": "bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c",
             "blsPubKey": "f3ee4cb60b04358c21460b9dd0832028959a6d0052218d796c96a5eac01b541f88595d62ee52880e0a77ecf8ffde63966a5d0d70028c08dfca622827563df99e"
            },
            "signature": "0x554a2a2f1b0d197730c707b595016b5f735ce0df0a5e9efd28a77764f295af1f"
           },
           "PrepareB": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0x3f643f315a72d54815e3b638e53a1f293834e6d9109c4c0f3e5d9c7171bf1cf2",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257",
             "NodeID": "bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c",
             "blsPubKey": "f3ee4cb60b04358c21460b9dd0832028959a6d0052218d796c96a5eac01b541f88595d62ee52880e0a77ecf8ffde63966a5d0d70028c08dfca622827563df99e"
            },
            "signature": "0x9e626bd0fd19290c7ff23a605259735de216f9c26df51ddaf51f66f0aade4000"
           }
          }
         ],
         "duplicate_vote": [],
         "duplicate_viewchange": []
        }`
	blockNumber = new(big.Int).Add(blockNumber, common.Big1)
	addr := common.HexToAddress("0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257")
	nodeId, err := discover.HexID("bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c")
	if nil != err {
		t.Fatal(err)
	}
	var blsKey bls.SecretKey
	skbyte, err := hex.DecodeString("d6ba381339988d7984393cd1892969d78eae735c588a9528c834676faf333507")
	if nil != err {
		t.Fatalf("ReportDuplicateSign DecodeString byte data fail: %v", err)
	}
	blsKey.SetLittleEndian(skbyte)
	can := &staking.Candidate{
		NodeId:          nodeId,
		BlsPubKey:       *blsKey.GetPublicKey(),
		StakingAddress:  addr,
		BenefitAddress:  addr,
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  1,
		ProgramVersion:  xutil.CalcVersion(initProgramVersion),
		Shares:          new(big.Int).SetUint64(1000),

		Released:           common.Big0,
		ReleasedHes:        common.Big0,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
	}
	stateDB.CreateAccount(addr)
	stateDB.AddBalance(addr, new(big.Int).SetUint64(1000000000000000000))
	if err := snapshotdb.Instance().NewBlock(blockNumber, chash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	if err := plugin.StakingInstance().CreateCandidate(stateDB, common.ZeroHash, blockNumber, can.Shares, 0, addr, can); nil != err {
		t.Fatal(err)
	}
	evidence, err := si.DecodeEvidence(data)
	if nil != err {
		t.Fatal(err)
	}
	if err := si.Slash(evidence, common.ZeroHash, blockNumber.Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800")); nil != err {
		t.Fatal(err)
	}
	if value, err := si.CheckDuplicateSign(addr, common.Big1.Uint64(), 1, stateDB); nil != err || len(value) == 0 {
		t.Fatal(err)
	}
	err = si.Slash(evidence, common.ZeroHash, blockNumber.Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800"))
	assert.Nil(t, err)
}

func TestSlashingPlugin_CheckMutiSign(t *testing.T) {
	si, stateDB := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	addr := common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52832")
	if _, err := si.CheckDuplicateSign(addr, 1, 1, stateDB); nil != err {
		t.Fatal(err)
	}
}
