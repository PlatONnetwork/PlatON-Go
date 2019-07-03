package cbft

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/stretchr/testify/assert"
)

func createPath(dirName string) string {
	name, err := ioutil.TempDir(os.TempDir(), dirName)
	if err != nil {
		panic(err)
	}
	return name
}

func createRandomPaths(basename string, num int) []string {
	var paths []string
	for i := 0; i < num; i++ {
		dirName := fmt.Sprintf("%s_%d", basename, i)
		path := createPath(dirName)
		paths = append(paths, path)
	}
	return paths
}

func removePaths(paths []string) {
	for _, path := range paths {
		os.RemoveAll(path)
	}
}

func mockCbft(paths []string) ([]*Cbft, []*testBackend, *testValidator) {
	validators := createTestValidator(createAccount(len(paths)))
	var engines []*Cbft
	var backends []*testBackend
	for i, path := range paths {
		engine := CreateCBFT(path, validators.validator(uint32(i)).privateKey)
		engines = append(engines, engine)
		backend := CreateBackend(engine, validators.Nodes())
		backends = append(backends, backend)
	}
	return engines, backends, validators
}

func closeAllCbft(engines []*Cbft) {
	for _, engine := range engines {
		engine.Close()
	}
}

func notProposalNodeIndex(nodeNum, nodeIndex uint32) uint32 {
	var notNodeIndex uint32
	if nodeIndex == nodeNum-1 {
		notNodeIndex = nodeIndex - 2
	} else {
		notNodeIndex = nodeIndex + 1
	}
	return notNodeIndex
}

func nextProposalNodeIndex(nodeNum, nodeIndex uint32) uint32 {
	var nextNodeIndex uint32
	if nodeIndex == nodeNum-1 {
		nextNodeIndex = 0
	} else {
		nextNodeIndex = nodeIndex + 1
	}
	return nextNodeIndex
}

func TestCbft_OnViewChange(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, validators := randomCBFT(path, 4)
	defer engine.Close()
	node := nodeIndexNow(validators, engine.startTimeOfEpoch)
	testCases := []*viewChange{
		// The timestamp is greater than the end time of the window period
		makeViewChange(node.privateKey, uint64(time.Now().Unix()+20000), 0, engine.blockChain.Genesis().Hash(), uint32(node.index), node.address, nil),
		// The timestamp is less than the start time of the window period
		makeViewChange(node.privateKey, uint64(time.Now().Unix()-20000), 0, engine.blockChain.Genesis().Hash(), uint32(node.index), node.address, nil),
		// Block hash is empty
		func() *viewChange {
			p := &viewChange{
				Timestamp:            uint64(time.Now().Unix()),
				BaseBlockNum:         0,
				ProposalIndex:        uint32(node.index),
				ProposalAddr:         node.address,
				BaseBlockPrepareVote: nil,
			}
			cb, _ := p.CannibalizeBytes()
			sign, _ := crypto.Sign(cb, node.privateKey)
			p.Signature.SetBytes(sign)
			return p
		}(),
		// The proposed person is empty
		func() *viewChange {
			p := &viewChange{
				Timestamp:            uint64(time.Now().Unix()),
				BaseBlockNum:         0,
				BaseBlockHash:        engine.blockChain.Genesis().Hash(),
				ProposalIndex:        uint32(node.index),
				BaseBlockPrepareVote: nil,
			}
			cb, _ := p.CannibalizeBytes()
			sign, _ := crypto.Sign(cb, node.privateKey)
			p.Signature.SetBytes(sign)
			return p
		}(),
		// Proposed non-window node
		func() *viewChange {
			nodeIndex := notProposalNodeIndex(4, uint32(node.index))
			p := &viewChange{
				Timestamp:            uint64(time.Now().Unix()),
				BaseBlockNum:         0,
				BaseBlockHash:        engine.blockChain.Genesis().Hash(),
				ProposalIndex:        nodeIndex,
				ProposalAddr:         validators.validator(nodeIndex).address,
				BaseBlockPrepareVote: nil,
			}
			cb, _ := p.CannibalizeBytes()
			sign, _ := crypto.Sign(cb, validators.validator(nodeIndex).privateKey)
			p.Signature.SetBytes(sign)
			return p
		}(),
		// Proposal and signature error
		func() *viewChange {
			p := &viewChange{
				Timestamp:            uint64(time.Now().Unix()),
				BaseBlockNum:         0,
				BaseBlockHash:        engine.blockChain.Genesis().Hash(),
				ProposalIndex:        uint32(node.index),
				ProposalAddr:         node.address,
				BaseBlockPrepareVote: nil,
			}
			cb, _ := p.CannibalizeBytes()
			errPri := createAccount(1)[0]
			sign, _ := crypto.Sign(cb, errPri)
			p.Signature.SetBytes(sign)
			return p
		}(),
		// Message not signed
		func() *viewChange {
			p := &viewChange{
				Timestamp:            uint64(time.Now().Unix()),
				BaseBlockNum:         0,
				BaseBlockHash:        engine.blockChain.Genesis().Hash(),
				ProposalIndex:        uint32(node.index),
				ProposalAddr:         node.address,
				BaseBlockPrepareVote: nil,
			}
			return p
		}(),
	}
	for i, view := range testCases {
		err := engine.OnViewChange(node.nodeID, view)
		assert.NotNil(t, err, "case:%d is fail", i)
		engine.viewChange = nil
	}
}

func TestCbft_Status(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 4)
	defer engine.Close()
	status := engine.Status()
	reg := regexp.MustCompile(`master:false`)
	if len(reg.FindAllString(status, -1)) == 0 {
		t.Errorf("status err,expected false,but now is true")
	}
	engine.master = true
	status = engine.Status()
	reg = regexp.MustCompile(`master:true`)
	if len(reg.FindAllString(status, -1)) == 0 {
		t.Errorf("status err,expected true,but now is false")
	}
}

func TestCbft_TracingSwitch(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 1)
	defer engine.Close()
	engine.TracingSwitch(1)
	assert.Equal(t, true, engine.tracing.isRecord)
	engine.TracingSwitch(10)
	assert.Equal(t, false, engine.tracing.isRecord)
}

func TestCbft_ViewChangeVote(t *testing.T) {
	paths := createRandomPaths("platon_test", 4)
	defer removePaths(paths)
	engines, _, validators := mockCbft(paths)
	defer closeAllCbft(engines)
	viewNode := nodeIndexNow(validators, time.Now().Unix())
	view, err := engines[viewNode.index].newViewChange()
	assert.Nil(t, err)
	nextNode := nextProposalNodeIndex(4, uint32(viewNode.index))
	engines[nextNode].OnViewChange(viewNode.nodeID, view)
	viewVote := makeViewChangeVote(validators.validator(nextNode).privateKey, view.Timestamp, view.BaseBlockNum,
		view.BaseBlockHash, view.ProposalIndex, view.ProposalAddr, nextNode, validators.validator(nextNode).address)
	err = engines[viewNode.index].OnViewChangeVote(validators.validator(nextNode).nodeID, viewVote)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(engines[viewNode.index].viewChangeVotes))
}
