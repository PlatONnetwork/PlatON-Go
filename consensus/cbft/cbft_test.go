package cbft

import (
	"Platon-go/common"
	"Platon-go/core/types"
	"Platon-go/crypto"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"container/list"
	"crypto/md5"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("test begin ...")
	initTest()

	flag.Parse()
	exitCode := m.Run()

	destroyTest()
	fmt.Println("test end ...")
	// Exit
	os.Exit(exitCode)
}

//var cbftConfig *params.CbftConfig
var forkRoot common.Hash

func TestBackTrackBlocksIncludingEnd(t *testing.T) {
	testBackTrackBlocks(t, true)
}

func TestBackTrackBlocksExcludingEnd(t *testing.T) {
	testBackTrackBlocks(t, false)
}

func testBackTrackBlocks(t *testing.T, includeEnd bool) {
	end := cbft.blockExtMap[forkRoot]
	exts := cbft.backTrackBlocks(cbft.highestLogical, end, includeEnd)

	t.Log("len(exts)", len(exts))
}

func initTest() {
	nodes := initNodes()
	priKey, _ := crypto.HexToECDSA("0x8b54398b67e656dcab213c1b5886845963a9ab0671786eefaf6e241ee9c8074f")

	cbftConfig := &params.CbftConfig{
		Period:       1,
		Epoch:        250000,
		MaxLatency:   600,
		Duration:     10,
		InitialNodes: nodes,
		NodeID:       nodes[0].ID,
		PrivateKey:   priKey,
	}

	cbft = &Cbft{
		config:        cbftConfig,
		blockExtMap:   make(map[common.Hash]*BlockExt),
		signedSet:     make(map[uint64]struct{}),
		netLatencyMap: make(map[discover.NodeID]*list.List),
	}
	buildMain(cbft)

	buildFork(cbft)
}

func destroyTest() {
	cbft.Close()
}

func buildFork(cbft *Cbft) {

	//sealhash := cbft.SealHash(header)

	parentHash := forkRoot
	for i := uint64(4); i <= 9; i++ {
		header := &types.Header{
			ParentHash: parentHash,
			Number:     big.NewInt(int64(i)),
			TxHash:     hash(1, i),
			Difficulty: big.NewInt(int64(i) * 2),
		}
		block := types.NewBlockWithHeader(header)

		ext := &BlockExt{
			block:       block,
			isLinked:    true,
			isSigned:    true,
			isStored:    false,
			isConfirmed: false,
			number:      block.NumberU64(),
			signs:       make([]*common.BlockConfirmSign, 0),
		}
		cbft.blockExtMap[block.Hash()] = ext

		parentHash = block.Hash()
	}
}

func hash(branch uint64, number uint64) (hash common.Hash) {
	s := "branch" + strconv.FormatUint(branch, 10) + "number" + strconv.FormatUint(number, 10)
	signByte := []byte(s)
	hasher := md5.New()
	hasher.Write(signByte)
	return common.BytesToHash(hasher.Sum(nil))
}
func buildMain(cbft *Cbft) {

	//sealhash := cbft.SealHash(header)

	rootHeader := &types.Header{Number: big.NewInt(0), Difficulty: big.NewInt(0), TxHash: hash(1, 0)}
	rootBlock := types.NewBlockWithHeader(rootHeader)
	rootExt := &BlockExt{
		block:       rootBlock,
		isLinked:    true,
		isSigned:    true,
		isStored:    false,
		isConfirmed: false,
		number:      rootBlock.NumberU64(),
		signs:       make([]*common.BlockConfirmSign, 0),
	}

	//hashSet[uint64(0)] = rootBlock.Hash()
	cbft.blockExtMap[rootBlock.Hash()] = rootExt
	cbft.highestConfirmed = rootExt

	parentHash := rootBlock.Hash()
	for i := uint64(1); i <= 10; i++ {
		header := &types.Header{
			ParentHash: parentHash,
			Number:     big.NewInt(int64(i)),
			TxHash:     hash(1, i),
			Difficulty: big.NewInt(int64(i) * 2),
		}
		block := types.NewBlockWithHeader(header)

		ext := &BlockExt{
			block:       block,
			isLinked:    true,
			isSigned:    true,
			isStored:    false,
			isConfirmed: false,
			number:      block.NumberU64(),
			signs:       make([]*common.BlockConfirmSign, 0),
		}
		//hashSet[i] = rootBlock.Hash()
		cbft.blockExtMap[block.Hash()] = ext

		cbft.highestLogical = ext

		parentHash = block.Hash()

		if i == 4 {
			forkRoot = block.Hash()
		}
	}

}

func initNodes() []discover.Node {
	var nodes [3]discover.Node

	initialNodes := [3]string{
		"1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429",
		"751f4f62fccee84fc290d0c68d673e4b0cc6975a5747d2baccb20f954d59ba3315d7bfb6d831523624d003c8c2d33451129e67c3eef3098f711ef3b3e268fd3c",
		"b6c8c9f99bfebfa4fb174df720b9385dbd398de699ec36750af3f38f8e310d4f0b90447acbef64bdf924c4b59280f3d42bb256e6123b53e9a7e99e4c432549d6",
	}
	nodeIDs := convert(initialNodes[:])

	for i, node := range nodes {
		node.ID = nodeIDs[i]
	}

	return nodes[:]
}
func convert(initialNodes []string) []discover.NodeID {
	NodeIDList := make([]discover.NodeID, 0, len(initialNodes))
	for _, value := range initialNodes {
		if nodeID, error := discover.HexID(value); error == nil {
			NodeIDList = append(NodeIDList, nodeID)
		}
	}
	return NodeIDList
}
