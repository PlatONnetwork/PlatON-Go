package plugin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"github.com/stretchr/testify/assert"
)

func initInfo(t *testing.T) (*SlashingPlugin, xcom.StateDB) {
	si := SlashInstance()
	StakingInstance()
	RestrictingInstance()
	chain := mock.NewChain(nil)
	return si, chain.StateDB
}

func buildStakingData(blockHash common.Hash, pri *ecdsa.PrivateKey, blsKey bls.SecretKey, t *testing.T, stateDb xcom.StateDB) {
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

	c1 := &staking.Candidate{
		NodeId:             nodeIdA,
		BlsPubKey:          *blsKey.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[1],
		StakingTxIndex:     uint32(2),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big256,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
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
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
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
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
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

	epochArr := &staking.ValidatorArray{
		Start: 1,
		End:   uint64(xutil.CalcBlocksEachEpoch()),
		Arr:   queue,
	}

	preArr := &staking.ValidatorArray{
		Start: 1,
		End:   xutil.ConsensusSize(),
		Arr:   queue,
	}

	curArr := &staking.ValidatorArray{
		Start: xutil.ConsensusSize() + 1,
		End:   xutil.ConsensusSize() * 2,
		Arr:   queue,
	}

	setVerifierList(blockHash, epochArr)
	setRoundValList(blockHash, preArr)
	setRoundValList(blockHash, curArr)
	stk.storeRoundValidatorAddrs(blockHash, 1, queue)
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
	pri, phash := buildBlock(t, int(startNumber), stateDB)
	startNumber++
	blockNumber := new(big.Int).SetInt64(int64(startNumber))
	if err := snapshotdb.Instance().NewBlock(blockNumber, phash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	buildStakingData(common.ZeroHash, pri, blsKey, t, stateDB)

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
	sk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	sign, err := crypto.Sign(header.SealHash().Bytes(), sk)
	if nil != err {
		t.Fatal(err)
	}
	copy(header.Extra[len(header.Extra)-common.ExtraSeal:], sign[:])
	if err := snapshotdb.Instance().NewBlock(header.Number, phash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	if err := si.BeginBlock(common.ZeroHash, header, stateDB); nil != err {
		t.Fatal(err)
	}
}

func buildBlock(t *testing.T, maxNumber int, stateDb xcom.StateDB) (*ecdsa.PrivateKey, common.Hash) {
	pri, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	pri2, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
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
		if err := db.NewBlock(blockNum, parentHash, common.ZeroHash); err != nil {
			t.Fatal(err)
		}
		if err := SlashInstance().BeginBlock(common.ZeroHash, header, stateDb); nil != err {
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
	var nodeBlsKey bls.SecretKey
	nodeBlsSkByte, err := hex.DecodeString("f91a76833f65af2fe971955fceb147fbf963d8b165942e957204d73944423a2c")
	if nil != err {
		t.Fatalf("ReportDuplicateSign DecodeString byte data fail: %v", err)
	}
	nodeBlsKey.SetLittleEndian(nodeBlsSkByte)
	buildStakingData(common.ZeroHash, crypto.HexMustToECDSA("4cb47bd14a95fa89e40303b56df2e152fd3ece3657db9c76b9f3beb1abd0c301"), nodeBlsKey, t, stateDB)
	if err := snapshotdb.Instance().Flush(chash, blockNumber); nil != err {
		t.Fatal(err)
	}
	if err := snapshotdb.Instance().Commit(chash); nil != err {
		t.Fatal(err)
	}
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	si.SetDecodeEvidenceFun(evidence.NewEvidence)
	GovPluginInstance()
	data := `{
	"prepareA": {
		"epoch": 1,
		"viewNumber": 1,
		"blockHash": "0x74093e065b4ef5eefc6247e22917ac6ad5f38d9214f2217efd3615f8c69f6492",
		"blockNumber": 1,
		"blockIndex": 1,
		"validateNode": {
			"index": 0,
			"address": "0xb8819949231dc7304c5906928751fd87ee489146",
			"nodeId": "a7a369700eaac3ced34e58c5d44fc78e85d09f3c214c521e7371ffbe7bd678a649422b280aea99fe6af118b8f5aac2c76ed3d86534556fd204759f61ae9bfeda",
			"blsPubKey": "3804db9244b172c1e5eb1ff52412bc42f0a8f1572aed04fa994d16fc1e2bafd4957dfc66269f89f8b7e6debbc5f1181902bcf20ea1d2737cac407d3e8e173e49e6eda171ee3be8f7576e3b77395bd059408eac4eb768d3057d083b3ed9287995"
		},
		"signature": "0x9d983b860cfa25356e624632808e4045fbddb4e57e497745c118e69c5ef4bbfda8afae2ad5daea94fa7979778106b48900000000000000000000000000000000"
	},
	"prepareB": {
		"epoch": 1,
		"viewNumber": 1,
		"blockHash": "0xa50d1eb4746d685bf0277d70a7094d0278582ddaa68c354a05f45ffb78fa1a3b",
		"blockNumber": 1,
		"blockIndex": 1,
		"validateNode": {
			"index": 0,
			"address": "0xb8819949231dc7304c5906928751fd87ee489146",
			"nodeId": "a7a369700eaac3ced34e58c5d44fc78e85d09f3c214c521e7371ffbe7bd678a649422b280aea99fe6af118b8f5aac2c76ed3d86534556fd204759f61ae9bfeda",
			"blsPubKey": "3804db9244b172c1e5eb1ff52412bc42f0a8f1572aed04fa994d16fc1e2bafd4957dfc66269f89f8b7e6debbc5f1181902bcf20ea1d2737cac407d3e8e173e49e6eda171ee3be8f7576e3b77395bd059408eac4eb768d3057d083b3ed9287995"
		},
		"signature": "0xc730346898c50361f705526dfb4bd8e3e8df8af7bfb152d9fd10e697433c29961859f0658c2dc68a32b2c3dd3841858700000000000000000000000000000000"
	}
}`
	data2 := `{
           "prepareA": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0x86c86e7ddb977fbd2f1d0b5cb92510c230775deef02b60d161c3912244473b54",
            "blockNumber": 1,
            "blockIndex": 1,
            "validateNode": {
             "index": 0,
             "address": "0x076c72c53c569df9998448832a61371ac76d0d05",
             "nodeId": "b68b23496b820f4133e42b747f1d4f17b7fd1cb6b065c613254a5717d856f7a56dabdb0e30657f18fb9074c7cb60eb62a6b35ad61898da407dae2cb8efe68511",
             "blsPubKey": "6021741b867202a3e60b91452d80e98f148aefadbb5ff1860f1fec5a8af14be20ca81fd73c231d6f67d4c9d2d516ac1297c8126ed7c441e476c0623c157638ea3b5b2189f3a20a78b2fd5fb32e5d7de055e4d2a0c181d05892be59cf01f8ab88"
            },
            "signature": "0x8c77b2178239fd525b774845cc7437ecdf5e6175ab4cc49dcb93eae6df288fd978e5290f59420f93bba22effd768f38900000000000000000000000000000000"
           },
           "prepareB": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0xeccd7a0b7793a74615721e883ab5223de30c5cf4d2ced9ab9dfc782e8604d416",
            "blockNumber": 1,
            "blockIndex": 1,
            "validateNode": {
             "index": 0,
             "address": "0x076c72c53c569df9998448832a61371ac76d0d05",
             "nodeId": "b68b23496b820f4133e42b747f1d4f17b7fd1cb6b065c613254a5717d856f7a56dabdb0e30657f18fb9074c7cb60eb62a6b35ad61898da407dae2cb8efe68511",
             "blsPubKey": "6021741b867202a3e60b91452d80e98f148aefadbb5ff1860f1fec5a8af14be20ca81fd73c231d6f67d4c9d2d516ac1297c8126ed7c441e476c0623c157638ea3b5b2189f3a20a78b2fd5fb32e5d7de055e4d2a0c181d05892be59cf01f8ab88"
            },
            "signature": "0x5213b4122f8f86874f537fa9eda702bba2e47a7b8ecc0ff997101d675a174ee5884ec85e8ea5155c5a6ad6b55326670d00000000000000000000000000000000"
           }
          }`
	blockNumber = new(big.Int).Add(blockNumber, common.Big1)
	addr := common.HexToAddress("0x076c72c53c569df9998448832a61371ac76d0d05")
	nodeId, err := discover.HexID("b68b23496b820f4133e42b747f1d4f17b7fd1cb6b065c613254a5717d856f7a56dabdb0e30657f18fb9074c7cb60eb62a6b35ad61898da407dae2cb8efe68511")
	if nil != err {
		t.Fatal(err)
	}
	var blsKey bls.SecretKey
	skbyte, err := hex.DecodeString("155b9a6f5575b9b5a4d8658f616660a549674b36c858e6c606d08ec5c20c4637")
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

		Released:           common.Big256,
		ReleasedHes:        common.Big0,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
	}
	stateDB.CreateAccount(addr)
	stateDB.AddBalance(addr, new(big.Int).SetUint64(1000000000000000000))
	if err := snapshotdb.Instance().NewBlock(blockNumber, chash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	if err := StakingInstance().CreateCandidate(stateDB, common.ZeroHash, blockNumber, can.Shares, 0, addr, can); nil != err {
		t.Fatal(err)
	}
	evidence1, err := si.DecodeEvidence(1, data)
	if nil != err {
		t.Fatal(err)
	}
	err = si.Slash(evidence1, common.ZeroHash, blockNumber.Uint64(), stateDB, sender)
	assert.NotNil(t, err)
	if err := si.Slash(evidence1, common.ZeroHash, blockNumber.Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800")); nil != err {
		t.Fatal(err)
	}
	if value, err := si.CheckDuplicateSign(common.HexToAddress("0xb8819949231dc7304c5906928751fd87ee489146"), common.Big1.Uint64(), 1, stateDB); nil != err || len(value) == 0 {
		t.Fatal(err)
	}
	evidence2, err := si.DecodeEvidence(1, data2)
	if nil != err {
		t.Fatal(err)
	}
	err = si.Slash(evidence2, common.ZeroHash, blockNumber.Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800"))
	assert.NotNil(t, err)

	err = si.Slash(evidence1, common.ZeroHash, blockNumber.Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800"))
	assert.NotNil(t, err)

	err = si.Slash(evidence1, common.ZeroHash, new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch()*uint64(xcom.EvidenceValidEpoch())*2).Uint64(), stateDB, common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52800"))
	assert.NotNil(t, err)
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
