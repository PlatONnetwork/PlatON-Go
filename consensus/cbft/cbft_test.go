package cbft

import (
	"container/list"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
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
var forkRootHash common.Hash
var rootBlock *types.Block
var rootNumber uint64
var hash23 common.Hash
var lastMainBlockHash common.Hash

var discreteBlock *types.Block
var lastMainBlock *types.Block
var lastButOneMainBlock *types.Block

func TestMarshalJson(t *testing.T){
	obj := NewBlockExt(lastButOneMainBlock, 0)
	obj.inTree = true
	obj.isExecuted = true
	obj.isSigned = true
	obj.isConfirmed = true
	obj.parent = nil

	t.Log("cbft.highestConfirmed is null", "hash:", obj.Hash)
	fmt.Println("cbft.highestConfirmed:" + obj.Hash)

	jsons, errs := json.Marshal(obj) //byte[] for json
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("obj:" + string(jsons)) //byte[] to string
}

func TestBlockSynced_discreteBlock(t *testing.T) {
	doBlockSynced(discreteBlock, t)
}

func TestBlockSynced_lastMainBlock(t *testing.T) {
	doBlockSynced(lastMainBlock, t)
}

func TestBlockSynced_lastButOneMainBlock(t *testing.T) {
	doBlockSynced(lastButOneMainBlock, t)
}

func TestBlockSynced_hash23(t *testing.T) {
	doBlockSynced(cbft.findBlockExt(hash23).block, t)
}


func doBlockSynced(currentBlock *types.Block, t *testing.T){
	if currentBlock.NumberU64() > cbft.getRootIrreversible().Number {
		log.Debug("chain has a higher irreversible block", "hash", currentBlock.Hash(), "number", currentBlock.NumberU64())
		newRoot := cbft.findBlockExt(currentBlock.Hash())
		if newRoot == nil || newRoot.block == nil {
			log.Debug("higher irreversible block is not existing in memory", "newTree", newRoot.toJson())
			//the block synced from other peer is a new block in local peer
			//remove all blocks referenced in old tree after being cut off
			cbft.cleanByTailoredTree(cbft.getRootIrreversible())

			newRoot = NewBlockExt(currentBlock, currentBlock.NumberU64())
			newRoot.inTree = true
			newRoot.isExecuted = true
			newRoot.isSigned = true
			newRoot.isConfirmed = true
			newRoot.Number = currentBlock.NumberU64()
			newRoot.parent = nil
			cbft.saveBlockExt(newRoot.block.Hash(), newRoot)

			cbft.buildChildNode(newRoot)

			if len(newRoot.Children) > 0{
				//the new root's children should re-execute base on new state
				for _, child := range newRoot.Children {
					if err := cbft.executeBlockAndDescendantMock(child, newRoot); err != nil {
						log.Error("execute the block error", "err", err)
						break
					}
				}
				//there are some redundancy code for newRoot, but these codes are necessary for other logical blocks
				cbft.signLogicalAndDescendantMock(newRoot)
			}

		} else if newRoot.block != nil {
			log.Debug("higher irreversible block is existing in memory","newTree", newRoot.toJson())

			//the block synced from other peer exists in local peer
			newRoot.isExecuted = true
			newRoot.isSigned = true
			newRoot.isConfirmed = true
			newRoot.Number = currentBlock.NumberU64()

			if newRoot.inTree == false {
				newRoot.inTree = true

				cbft.setDescendantInTree(newRoot)

				if len(newRoot.Children) > 0{
					//the new root's children should re-execute base on new state
					for _, child := range newRoot.Children {
						if err := cbft.executeBlockAndDescendantMock(child, newRoot); err != nil {
							log.Error("execute the block error", "err", err)
							break
						}
					}
					//there are some redundancy code for newRoot, but these codes are necessary for other logical blocks
					cbft.signLogicalAndDescendantMock(newRoot)
				}
			}else{
				//cut off old tree from new root,
				tailorTree(newRoot)
			}

			//remove all blocks referenced in old tree after being cut off
			cbft.cleanByTailoredTree(cbft.getRootIrreversible())
		}

		//remove all other blocks those their numbers are too low
		cbft.cleanByNumber(newRoot.Number)

		log.Debug("the cleared new tree in memory", "json", newRoot.toJson())

		//reset the new root irreversible
		cbft.rootIrreversible.Store(newRoot)

		log.Debug("reset the new root irreversible by synced", "hash", newRoot.block.Hash(), "number", newRoot.block.NumberU64())


		//reset logical path
		highestLogical := cbft.findHighestLogical(newRoot)
		//cbft.setHighestLogical(highestLogical)
		cbft.highestLogical.Store(highestLogical)

		t.Log("reset the highestLogical by synced", "number", highestLogical.block.NumberU64(),  "newRoot.number", newRoot.block.NumberU64() )

		//reset highest confirmed block
		highestConfirmed :=cbft.findLastClosestConfirmedIncludingSelf(newRoot)
		cbft.highestConfirmed.Store(highestConfirmed)
		if highestConfirmed != nil {
			t.Log("cbft.highestConfirmed", "number", highestConfirmed.block.NumberU64(),"rcvTime", highestConfirmed.rcvTime)
		} else {

			t.Log("reset the highestConfirmed by synced failure because findLastClosestConfirmedIncludingSelf() returned nil", "newRoot.number", newRoot.block.NumberU64())

		}
	}
}



func TestFindLastClosestConfirmedIncludingSelf(t *testing.T) {

	newRoot := NewBlockExt(rootBlock, 0)
	newRoot.inTree = true
	newRoot.isExecuted = true
	newRoot.isSigned = true
	newRoot.isConfirmed = true
	newRoot.Number = rootBlock.NumberU64()

	//reorg the block tree
	children := cbft.findChildren(newRoot)
	for _, child := range children {
		child.parent = newRoot
		child.inTree = true
	}
	newRoot.Children = children

	//save the root in BlockExtMap
	cbft.saveBlockExt(newRoot.block.Hash(), newRoot)

	//reset the new root irreversible
	cbft.rootIrreversible.Store(newRoot)
	//the new root's children should re-execute base on new state
	for _, child := range newRoot.Children {
		if err := cbft.executeBlockAndDescendant(child, newRoot); err != nil {
			//remove bad block from tree and map
			cbft.removeBadBlock(child)
			break
		}
	}

	//there are some redundancy code for newRoot, but these codes are necessary for other logical blocks
	cbft.signLogicalAndDescendant(newRoot)

	//reset logical path
	//highestLogical := cbft.findHighestLogical(newRoot)
	//cbft.setHighestLogical(highestLogical)
	//reset highest confirmed block
	cbft.highestConfirmed.Store(cbft.findLastClosestConfirmedIncludingSelf(newRoot))

	if cbft.getHighestLogical() != nil {
		t.Log("ok")
	} else {
		t.Log("cbft.highestConfirmed is null")
	}
}
func TestFindClosestConfirmedExcludingSelf(t *testing.T) {
	newRoot := NewBlockExt(rootBlock, 0)
	newRoot.inTree = true
	newRoot.isExecuted = true
	newRoot.isSigned = true
	newRoot.isConfirmed = true
	newRoot.Number = rootBlock.NumberU64()

	closest := cbft.findClosestConfirmedExcludingSelf(newRoot)
	fmt.Println(closest)
}

func TestExecuteBlockAndDescendant(t *testing.T) {
	newRoot := NewBlockExt(rootBlock, 0)
	newRoot.inTree = true
	newRoot.isExecuted = true
	newRoot.isSigned = true
	newRoot.isConfirmed = true
	newRoot.Number = rootBlock.NumberU64()

	//save the root in BlockExtMap
	cbft.saveBlockExt(newRoot.block.Hash(), newRoot)

	//reset the new root irreversible
	cbft.rootIrreversible.Store(newRoot)
	//reorg the block tree
	cbft.buildChildNode(newRoot)

	//the new root's children should re-execute base on new state
	for _, child := range newRoot.Children {
		if err := cbft.executeBlockAndDescendant(child, newRoot); err != nil {
			//remove bad block from tree and map
			cbft.removeBadBlock(child)
			break
		}
	}

	//there are some redundancy code for newRoot, but these codes are necessary for other logical blocks
	cbft.signLogicalAndDescendant(newRoot)
}

func TestBackTrackBlocksIncludingEnd(t *testing.T) {
	testBackTrackBlocks(t, true)
}

func TestBackTrackBlocksExcludingEnd(t *testing.T) {
	testBackTrackBlocks(t, false)
}

func testBackTrackBlocks(t *testing.T, includeEnd bool) {
	end, _ := cbft.blockExtMap.Load(rootBlock.Hash())
	exts := cbft.backTrackBlocks(cbft.getHighestLogical(), end.(*BlockExt), includeEnd)

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
		config: cbftConfig,
		//blockExtMap:   make(map[common.Hash]*BlockExt),
		//signedSet:     make(map[uint64]struct{}),
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

	parentHash := forkRootHash
	for i := uint64(3); i <= 5; i++ {
		header := &types.Header{
			ParentHash: parentHash,
			Number:     big.NewInt(int64(i)),
			TxHash:     hash(2, i),
		}
		block := types.NewBlockWithHeader(header)

		ext := &BlockExt{
			block:       block,
			inTree:      true,
			isExecuted:  true,
			isConfirmed: false,
			rcvTime:    int64(20+i),
			Number:      block.NumberU64(),
			signs:       make([]*common.BlockConfirmSign, 0),
		}
		cbft.blockExtMap.Store(block.Hash(), ext)

		parentHash = block.Hash()

		if i == 3 {
			ext.isSigned=false
			ext.isConfirmed=false
			hash23 = block.Hash()
		} else if i == 4 {
			ext.isSigned=true
			ext.isConfirmed=false
		} else if i == 5 {
			ext.isSigned=false
			ext.isConfirmed=false
		}
		cbft.buildIntoTree(ext)
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

	rootHeader := &types.Header{Number: big.NewInt(0), TxHash: hash(1, 0)}
	rootBlock = types.NewBlockWithHeader(rootHeader)
	rootNumber = 0

	rootExt := &BlockExt{
		block:       rootBlock,
		inTree:      true,
		isExecuted:  true,
		isSigned:    true,
		isConfirmed: true,
		rcvTime:	 0,
		Number:      rootBlock.NumberU64(),
		signs:       make([]*common.BlockConfirmSign, 0),
	}

	//hashSet[uint64(0)] = rootBlock.Hash()
	cbft.blockExtMap.Store(rootBlock.Hash(), rootExt)
	cbft.highestConfirmed.Store(rootExt)
	cbft.rootIrreversible.Store(rootExt)

	parentHash := rootBlock.Hash()
	for i := uint64(1); i <= 5; i++ {
		header := &types.Header{
			ParentHash: parentHash,
			Number:     big.NewInt(int64(i)),
			TxHash:     hash(1, i),
		}
		block := types.NewBlockWithHeader(header)

		ext := &BlockExt{
			block:       block,
			inTree:      true,
			isExecuted:  true,
			isSigned:    true,
			isConfirmed: false,
			rcvTime:    int64(i),
			Number:      block.NumberU64(),
			signs:       make([]*common.BlockConfirmSign, 0),
		}
		//hashSet[i] = rootBlock.Hash()
		cbft.blockExtMap.Store(block.Hash(), ext)

		cbft.highestLogical.Store(ext)

		parentHash = block.Hash()

		if i == 2 {
			forkRootHash = block.Hash()
			ext.isSigned=false
			ext.isConfirmed=true
		} else if i == 3 {
			ext.isSigned=true
			ext.isConfirmed=false
		} else if i == 4 {
			lastButOneMainBlock = block

			ext.isSigned=false
			ext.isConfirmed=false
		} else if i == 5 {
			lastMainBlock = block
			ext.isSigned=false
			ext.isConfirmed=false

		}
		cbft.buildIntoTree(ext)
	}

	notExistHeader := &types.Header{
		ParentHash: parentHash,
		Number:     big.NewInt(int64(20)),
		TxHash:     hash(1, 20),
	}
	notExistBlock := types.NewBlockWithHeader(notExistHeader)


	discreteHeader := &types.Header{
		ParentHash: notExistBlock.Hash(),
		Number:     big.NewInt(int64(6)),
		TxHash:     hash(1, 6),
	}
	discreteBlock = types.NewBlockWithHeader(discreteHeader)

	discreteExt := &BlockExt{
		block:       discreteBlock,
		inTree:      false,
		isExecuted:  false,
		isSigned:    false,
		isConfirmed: false,
		rcvTime:     int64(220),
		Number:      discreteBlock.NumberU64(),
		signs:       make([]*common.BlockConfirmSign, 0),
	}
	//hashSet[i] = rootBlock.Hash()
	cbft.blockExtMap.Store(discreteBlock.Hash(), discreteExt)

	parentHash = discreteBlock.Hash()
	for i := 7; i <= 7; i++ {
		header := &types.Header{
			ParentHash: parentHash,
			Number:     big.NewInt(int64(i)),
			TxHash:     hash(1, uint64(i)),
		}
		block := types.NewBlockWithHeader(header)

		ext := &BlockExt{
			block:       block,
			inTree:      false,
			isExecuted:  false,
			isSigned:    false,
			isConfirmed: false,
			rcvTime:    int64(i),
			Number:      block.NumberU64(),
			signs:       make([]*common.BlockConfirmSign, 0),
		}
		//hashSet[i] = rootBlock.Hash()
		cbft.blockExtMap.Store(block.Hash(), ext)
		parentHash = block.Hash()
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
