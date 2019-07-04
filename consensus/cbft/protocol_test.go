package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"math/big"
	"reflect"
	"testing"
	"time"
)

func TestErrorCode(t *testing.T) {
	var code errCode
	var str string
	// ErrMsgTooLarge
	code = ErrMsgTooLarge
	str = code.String()
	if str != "Message too long" {
		t.Error("error string")
	}

	// ErrDecode
	code = ErrDecode
	str = code.String()
	if str != "Invalid message" {
		t.Error("error string")
	}

	// ErrInvalidMsgCode
	code = ErrInvalidMsgCode
	str = code.String()
	if str != "Invalid message code" {
		t.Error("error string")
	}

	// ErrProtocolVersionMismatch
	code = ErrProtocolVersionMismatch
	str = code.String()
	if str != "Protocol version mismatch" {
		t.Error("error string")
	}

	// ErrNetworkIdMismatch
	code = ErrNetworkIdMismatch
	str = code.String()
	if str != "NetworkId mismatch" {
		t.Error("error string")
	}

	// ErrGenesisBlockMismatch
	code = ErrGenesisBlockMismatch
	str = code.String()
	if str != "Genesis block mismatch" {
		t.Error("error string")
	}

	// ErrNoStatusMsg
	code = ErrNoStatusMsg
	str = code.String()
	if str != "No status message" {
		t.Error("error string")
	}

	// ErrExtraStatusMsg
	code = ErrExtraStatusMsg
	str = code.String()
	if str != "Extra status message" {
		t.Error("error string")
	}

	// ErrSuspendedPeer
	code = ErrSuspendedPeer
	str = code.String()
	if str != "Suspended peer" {
		t.Error("error string")
	}
}

func TestPrepareBlock(t *testing.T) {
	// define block
	block := types.NewBlockWithHeader(&types.Header{
		GasLimit: uint64(3141592),
		GasUsed:  uint64(21000),
		Coinbase: common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"),
		Root:     common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017"),
		//Hash: common.HexToHash("0a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9e"),
		Nonce: types.EncodeNonce(RandBytes(81)),
		Time:  big.NewInt(1426516743),
		Extra: make([]byte, 100),
	})
	pb := &prepareBlock{
		Timestamp:     uint64(time.Now().Unix()),
		Block:         block,
		ProposalIndex: 1,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address")),
		View:          &viewChange{},
	}
	//
	var consensusMsg ConsensusMsg = pb

	// check sign
	consensusMsg.Sign()
	var empty *prepareBlock
	check(t, empty, pb)
}

func TestPrepareBlockHash(t *testing.T) {
	pbh := &prepareBlockHash{
		Hash:   common.BytesToHash([]byte("I'm hash")),
		Number: 1,
	}
	var empty *prepareBlockHash
	check(t, empty, pbh)
}

func TestMessageType(t *testing.T) {
	testCases := []struct {
		msgType interface{}
		want    uint64
	}{
		{msgType: &prepareBlock{}, want: PrepareBlockMsg},
		{msgType: &prepareVote{}, want: PrepareVoteMsg},
		{msgType: &viewChange{}, want: ViewChangeMsg},
		{msgType: &viewChangeVote{}, want: ViewChangeVoteMsg},
		{msgType: &confirmedPrepareBlock{}, want: ConfirmedPrepareBlockMsg},
		{msgType: &getPrepareVote{}, want: GetPrepareVoteMsg},
		{msgType: &prepareVotes{}, want: PrepareVotesMsg},
		{msgType: &getPrepareBlock{}, want: GetPrepareBlockMsg},
		{msgType: &getHighestPrepareBlock{}, want: GetHighestPrepareBlockMsg},
		{msgType: &highestPrepareBlock{}, want: HighestPrepareBlockMsg},
		{msgType: &cbftStatusData{}, want: CBFTStatusMsg},
		{msgType: &prepareBlockHash{}, want: PrepareBlockHashMsg},
	}
	for _, v := range testCases {
		if MessageType(v.msgType) != uint64(v.want) {
			t.Errorf("MessageType error, want: %v, is: %v", uint64(v.want), MessageType(v.msgType))
		}
	}
}

func TestPrepareVote(t *testing.T) {
	privateHex := "e4eb3e58ab7810984a0c77d432b07fe9f9897158dd4bb4f63d0a4366e6d949fa"
	pri, _ := crypto.HexToECDSA(privateHex)
	pv := &prepareVote{
		Timestamp:      uint64(time.Now().Unix()),
		Hash:           common.BytesToHash([]byte("I'm hash")),
		Number:         1,
		ValidatorIndex: 0,
		ValidatorAddr:  common.BytesToAddress([]byte("I'm address")),
	}

	var consensusMsg ConsensusMsg = pv
	cb, _ := consensusMsg.CannibalizeBytes()
	sign, _ := crypto.Sign(cb, pri)
	pv.Signature.SetBytes(sign)

	// check sign
	signRes := consensusMsg.Sign()
	if !reflect.DeepEqual(signRes, pv.Signature.Bytes()) {
		t.Error("sign not equal.")
	}
	var empty *prepareVote
	check(t, empty, pv)
}

func TestViewChange(t *testing.T) {
	privateHex := "e4eb3e58ab7810984a0c77d432b07fe9f9897158dd4bb4f63d0a4366e6d949fa"
	pri, _ := crypto.HexToECDSA(privateHex)
	pv := &viewChange{
		Timestamp:     uint64(time.Now().Unix()),
		ProposalIndex: 0,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address")),
		BaseBlockHash: common.BytesToHash([]byte("I'm hash")),
		BaseBlockNum:  1,
		Extra:         make([]byte, 100),
	}

	var consensusMsg ConsensusMsg = pv
	cb, _ := consensusMsg.CannibalizeBytes()
	sign, _ := crypto.Sign(cb, pri)
	pv.Signature.SetBytes(sign)

	// check sign
	signRes := consensusMsg.Sign()
	if !reflect.DeepEqual(signRes, sign) {
		t.Error("sign not equal.")
	}
	var empty *viewChange
	check(t, empty, pv)
}

func TestViewChange_Equal(t *testing.T) {
	now := time.Now().Unix()
	vc1 := &viewChange{
		Timestamp:     uint64(now),
		ProposalIndex: 0,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address")),
		BaseBlockHash: common.BytesToHash([]byte("I'm hash")),
		BaseBlockNum:  1,
	}
	vc2 := &viewChange{
		Timestamp:     uint64(now),
		ProposalIndex: 0,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address")),
		BaseBlockHash: common.BytesToHash([]byte("I'm hash")),
		BaseBlockNum:  1,
	}
	vc3 := &viewChange{
		Timestamp:     uint64(now),
		ProposalIndex: 0,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address 3")),
		BaseBlockHash: common.BytesToHash([]byte("I'm hash")),
		BaseBlockNum:  1,
	}
	if !reflect.DeepEqual(vc1, vc2) {
		t.Error("must be equal")
	}
	if reflect.DeepEqual(vc1, vc3) {
		t.Error("must not be equal")
	}

	duplicate_vc4 := vc3.CopyWithoutVotes()
	pvs := []*prepareVote{
		{Timestamp: uint64(now), Hash: common.BytesToHash([]byte("v3 pv01"))},
		{Timestamp: uint64(now), Hash: common.BytesToHash([]byte("v3 pv02"))},
	}
	vc3.BaseBlockPrepareVote = pvs
	duplicate_vc5 := vc3.Copy()
	if duplicate_vc4 == duplicate_vc5 {
		t.Error("shoud be equal")
	}
}

func TestViewChangeVote(t *testing.T) {
	privateHex := "e4eb3e58ab7810984a0c77d432b07fe9f9897158dd4bb4f63d0a4366e6d949fa"
	pri, _ := crypto.HexToECDSA(privateHex)
	pcv := &viewChangeVote{
		Timestamp:      uint64(time.Now().Unix()),
		ProposalIndex:  0,
		ProposalAddr:   common.BytesToAddress([]byte("I'm address")),
		BlockHash:      common.BytesToHash([]byte("I'm hash")),
		BlockNum:       1,
		ValidatorAddr:  common.BytesToAddress([]byte("I'm validtor address")),
		ValidatorIndex: 1,
		Extra:          make([]byte, 100),
	}

	var consensusMsg ConsensusMsg = pcv
	cb, _ := consensusMsg.CannibalizeBytes()
	sign, _ := crypto.Sign(cb, pri)
	pcv.Signature.SetBytes(sign)

	// check sign
	signRes := consensusMsg.Sign()
	if !reflect.DeepEqual(signRes, sign) {
		t.Error("sign not equal.")
	}
	var empty *viewChangeVote
	check(t, empty, pcv)
}

func TestViewChangeVote_View(t *testing.T) {
	pcv := &viewChangeVote{
		Timestamp:      uint64(time.Now().Unix()),
		ProposalIndex:  0,
		ProposalAddr:   common.BytesToAddress([]byte("I'm address")),
		BlockHash:      common.BytesToHash([]byte("I'm hash")),
		BlockNum:       1,
		ValidatorAddr:  common.BytesToAddress([]byte("I'm validtor address")),
		ValidatorIndex: 1,
		Extra:          make([]byte, 100),
	}
	pc := &viewChange{
		Timestamp:     uint64(time.Now().Unix()),
		ProposalIndex: 0,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address")),
		BaseBlockHash: common.BytesToHash([]byte("I'm hash")),
		BaseBlockNum:  1,
	}
	if !pcv.EqualViewChange(pc) {
		t.Error("should equal")
	}

	v := pcv.ViewChangeWithSignature()
	if !pcv.EqualViewChange(v) {
		t.Error("should equal")
	}
}

func check(t *testing.T, empty Message, message Message) {
	str := message.String()
	if str == "" {
		t.Error("error")
	}
	str = empty.String()
	if str != "" {
		t.Error("error")
	}

	msgHash := message.MsgHash()
	emptyCommonHash := common.Hash{}
	if msgHash != emptyCommonHash {
		t.Log(msgHash.String())
	}

	msgHash = empty.MsgHash()
	if msgHash != emptyCommonHash {
		t.Error("error")
	}

	// check bhash method
	bhash := message.BHash()
	bhash = empty.BHash()
	if bhash != emptyCommonHash {
		t.Error("must be empty hash")
	}
}

func emptyConfirmedPrepareBlock() *confirmedPrepareBlock {
	return nil
}

func TestConfirmedPrepareBlock(t *testing.T) {
	pbh := &confirmedPrepareBlock{
		Hash:   common.BytesToHash([]byte("I'm hash")),
		Number: 1,
	}
	check(t, emptyConfirmedPrepareBlock(), pbh)
}

func emptyGetHighestPrepareBlock() *getHighestPrepareBlock {
	return nil
}

func TestGetHighestPrepareBlock(t *testing.T) {
	pbh := &getHighestPrepareBlock{
		Lowest: 1,
	}
	check(t, emptyGetHighestPrepareBlock(), pbh)
}

func TestGetPrepareBlock(t *testing.T) {
	var emptyGetPrepareBlock *getPrepareBlock
	pbh := &getPrepareBlock{
		Hash:   common.BytesToHash([]byte("I'm hash")),
		Number: 1,
	}
	check(t, emptyGetPrepareBlock, pbh)
}

func TestGetPrepareVote(t *testing.T) {
	var empty *getPrepareVote
	pbh := &getPrepareVote{
		Hash:   common.BytesToHash([]byte("I'm hash")),
		Number: 1,
	}
	check(t, empty, pbh)
}

func TestPrepareVotes(t *testing.T) {
	var empty *prepareVotes
	pbh := &prepareVotes{
		Hash:   common.BytesToHash([]byte("I'm hash")),
		Number: 1,
	}
	check(t, empty, pbh)
}

func TestSignBitArray(t *testing.T) {
	var empty *signBitArray
	pbh := &signBitArray{
		BlockHash: common.BytesToHash([]byte("I'm hash")),
		BlockNum:  1,
	}
	check(t, empty, pbh)

	cpy := pbh.Copy()
	if !reflect.DeepEqual(pbh, cpy) {
		t.Error("should equal")
	}
}

func TestCbftStatusData(t *testing.T) {
	var empty *cbftStatusData
	pbh := &cbftStatusData{
		CurrentBlock: common.BytesToHash([]byte("I'm hash")),
		ConfirmedBn:  big.NewInt(1),
		LogicBn:      big.NewInt(1),
	}
	check(t, empty, pbh)
}
