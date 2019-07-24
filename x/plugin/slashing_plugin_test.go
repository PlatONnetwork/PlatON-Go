package plugin_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
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
		t.Error(err)
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

	nodeId_A := discover.PubkeyID(&pri.PublicKey)
	addr_A, _ := xutil.NodeId2Addr(nodeId_A)

	nodeId_B := nodeIdArr[1]
	addr_B, _ := xutil.NodeId2Addr(nodeId_B)

	nodeId_C := nodeIdArr[2]
	addr_C, _ := xutil.NodeId2Addr(nodeId_C)

	//canArr := make(staking.CandidateQueue, 0)

	c1 := &staking.Candidate{
		NodeId:             nodeId_A,
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

	c2 := &staking.Candidate{
		NodeId:             nodeId_B,
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

	c3 := &staking.Candidate{
		NodeId:             nodeId_C,
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

	//canArr = append(canArr, c1)
	//canArr = append(canArr, c2)
	//canArr = append(canArr, c3)

	stakingDB.SetCanPowerStore(blockHash, addr_A, c1)
	stakingDB.SetCanPowerStore(blockHash, addr_B, c2)
	stakingDB.SetCanPowerStore(blockHash, addr_C, c3)

	stakingDB.SetCandidateStore(blockHash, addr_A, c1)
	stakingDB.SetCandidateStore(blockHash, addr_B, c2)
	stakingDB.SetCandidateStore(blockHash, addr_C, c3)

	fmt.Println("addr_A", hex.EncodeToString(addr_A.Bytes()), "addr_B", hex.EncodeToString(addr_B.Bytes()), "addr_C", hex.EncodeToString(addr_C.Bytes()))

	queue := make(staking.ValidatorQueue, 0)

	v1 := &staking.Validator{
		NodeAddress:   addr_A,
		NodeId:        c1.NodeId,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c1.StakingBlockNum), fmt.Sprint(c1.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v2 := &staking.Validator{
		NodeAddress:   addr_B,
		NodeId:        c2.NodeId,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c2.StakingBlockNum), fmt.Sprint(c2.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v3 := &staking.Validator{
		NodeAddress:   addr_C,
		NodeId:        c3.NodeId,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c3.StakingBlockNum), fmt.Sprint(c3.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	queue = append(queue, v1)
	queue = append(queue, v2)
	queue = append(queue, v3)

	val_Arr := &staking.Validator_array{
		Start: 1,
		End:   22000,
		Arr:   queue,
	}

	setVerifierList(blockHash, val_Arr)
	setRoundValList(blockHash, val_Arr)
	setRoundValList(blockHash, val_Arr)
}

func TestSlashingPlugin_BeginBlock(t *testing.T) {
	_, _, _ = newChainState()
	si, stateDB := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	pri, phash := confirmBlock(t, 478)
	blockNumber := new(big.Int).SetInt64(479)
	if err := snapshotdb.Instance().NewBlock(blockNumber, phash, common.ZeroHash); err != nil {
		t.Error(err)
		return
	}
	buildStakingData(common.ZeroHash, pri, t, stateDB)

	phash = common.HexToHash("0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802130406")
	if err := snapshotdb.Instance().Flush(phash, blockNumber); err != nil {
		t.Error(err)
		return
	}
	if err := snapshotdb.Instance().Commit(phash); err != nil {
		t.Error(err)
		return
	}
	header := &types.Header{
		Number: new(big.Int).SetUint64(480),
		Extra:  make([]byte, 97),
	}
	if err := snapshotdb.Instance().NewBlock(header.Number, phash, common.ZeroHash); nil != err {
		t.Error(err)
		return
	}
	if err := si.BeginBlock(common.ZeroHash, header, stateDB); nil != err {
		t.Error(err)
		return
	}
}

func TestSlashingPlugin_Confirmed(t *testing.T) {
	si, _ := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	confirmBlock(t, 251)
	result, err := si.GetPreNodeAmount()
	if nil != err {
		t.Error(err)
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

		if i == 7 {
			sk = pri2
		}
		header := &types.Header{
			Number: blockNum,
			Extra:  make([]byte, 97),
		}
		sign, err := crypto.Sign(header.SealHash().Bytes(), sk)
		if nil != err {
			t.Error(err)
		}
		copy(header.Extra[len(header.Extra)-common.ExtraSeal:], sign[:])
		block := types.NewBlock(header, nil, nil)
		if err := plugin.SlashInstance().Confirmed(block); nil != err {
			t.Error(err)
		}
		if err := db.NewBlock(blockNum, parentHash, common.ZeroHash); err != nil {
			panic(err)
		}
		//hash = crypto.Keccak256Hash(common.Int32ToBytes(int32(i)))
		if err := db.Flush(header.Hash(), blockNum); err != nil {
			panic(err)
		}
		if err := db.Commit(header.Hash()); err != nil {
			panic(err)
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
	snapshotdb.Instance().NewBlock(blockNumber, genesis.Hash(), common.ZeroHash)
	buildStakingData(common.ZeroHash, nil, t, stateDB)
	snapshotdb.Instance().Flush(chash, blockNumber)
	snapshotdb.Instance().Commit(chash)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	plugin.GovPluginInstance()
	si.SetDecodeEvidenceFun(cbft.NewEvidences)
	data := `{
          "duplicate_prepare": [
            {
              "VoteA": {
                "timestamp": 0,
                "block_hash": "0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802130407",
                "block_number": 1,
                "validator_index": 1,
                "validator_address": "0x120b77ab712589ebd42d69003893ef962cc52832",
                "signature": "0xa65e16b3bc4862fdd893eaaaaecf1e415cdc2c8a08e4bbb1f6b2a1f4bf4e2d0c0ec27857da86a5f3150b32bee75322073cec320e51fe0a123cc4238ee4155bf001"
              },
              "VoteB": {
                "timestamp": 0,
                "block_hash": "0x18030d1e01071b1d071a12151e100a091f060801031917161e0a0d0f02161d0e",
                "block_number": 1,
                "validator_index": 1,
                "validator_address": "0x120b77ab712589ebd42d69003893ef962cc52832",
                "signature": "0x9126f9a339c8c4a873efc397062d67e9e9109895cd9da0d09a010d5f5ebbc6e76d285f7d87f801850c8552234101b651c8b7601b4ea077328c27e4f86d66a1bf00"
              }
            }
          ],
          "duplicate_viewchange": [],
          "timestamp_viewchange": []
        }`
	blockNumber = new(big.Int).Add(blockNumber, common.Big1)
	addr := common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52832")
	nodeId, err := discover.HexID("0x38e2724b366d66a5acb271dba36bc45e2161e868d961ee299f4e331927feb5e9373f35229ef7fe7e84c083b0fbf24264faef01faaf388df5f459b87638aa620b")
	if nil != err {
		t.Error(err)
	}
	can := &staking.Candidate{
		NodeId:          nodeId,
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
		panic(err)
	}
	if err := plugin.StakingInstance().CreateCandidate(stateDB, common.ZeroHash, blockNumber, can.Shares, 0, addr, can); nil != err {
		t.Error(err)
	}
	evidence, err := si.DecodeEvidence(data)
	if nil != err {
		t.Error(err)
		return
	}
	if err := si.Slash(evidence, common.ZeroHash, blockNumber.Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800")); nil != err {
		t.Error(err)
	}
	if value, err := si.CheckDuplicateSign(addr, common.Big1.Uint64(), 1, stateDB); nil != err || len(value) == 0 {
		t.Error(err)
	}
	evidence, err = si.DecodeEvidence(data)
	if nil != err {
		t.Error(err)
		return
	}
	err = si.Slash(evidence, common.ZeroHash, blockNumber.Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800"))
	assert.NotNil(t, err)
	data = `{
          "duplicate_prepare": [
            {
              "VoteA": {
                "timestamp": 0,
                "block_hash": "0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802130407",
                "block_number": 2,
                "validator_index": 1,
                "validator_address": "0x120b77ab712589ebd42d69003893ef962cc52832",
                "signature": "0xa65e16b3bc4862fdd893eaaaaecf1e415cdc2c8a08e4bbb1f6b2a1f4bf4e2d0c0ec27857da86a5f3150b32bee75322073cec320e51fe0a123cc4238ee4155bf001"
              },
              "VoteB": {
                "timestamp": 0,
                "block_hash": "0x18030d1e01071b1d071a12151e100a091f060801031917161e0a0d0f02161d0e",
                "block_number": 2,
                "validator_index": 1,
                "validator_address": "0x120b77ab712589ebd42d69003893ef962cc52832",
                "signature": "0x9126f9a339c8c4a873efc397062d67e9e9109895cd9da0d09a010d5f5ebbc6e76d285f7d87f801850c8552234101b651c8b7601b4ea077328c27e4f86d66a1bf00"
              }
            }
          ],
          "duplicate_viewchange": [],
          "timestamp_viewchange": []
        }`
	evidence, err = si.DecodeEvidence(data)
	if nil != err {
		t.Error(err)
		return
	}
	err = si.Slash(evidence, common.ZeroHash, blockNumber.Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800"))
	assert.NotNil(t, err)
}

func TestSlashingPlugin_CheckMutiSign(t *testing.T) {
	si, stateDB := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	addr := common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52832")
	if _, err := si.CheckDuplicateSign(addr, 1, 1, stateDB); nil != err {
		t.Error(err)
	}
}
