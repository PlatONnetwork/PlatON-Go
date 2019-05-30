package cbft

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"sync"
	"time"
)

var (
	SyncInterval     = time.Duration(500)
	errInvalidAuthor = errors.New("invalid author")
)

type ProducerBlocks struct {
	author       discover.NodeID
	baseBlockNum uint64
	blocks       map[uint64]*types.Block
	lock         sync.Mutex
}

func NewProducerBlocks(author discover.NodeID, blockNum uint64) *ProducerBlocks {
	return &ProducerBlocks{
		author:       author,
		baseBlockNum: blockNum,
		blocks:       make(map[uint64]*types.Block),
	}
}
func (pb ProducerBlocks) String() string {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	return fmt.Sprintf("author:%s base block:%d, total:%d", pb.author.String(), pb.baseBlockNum, len(pb.blocks))
}
func (pb *ProducerBlocks) SetAuthor(author discover.NodeID) {
	pb.lock.Lock()
	pb.author = author
	pb.blocks = make(map[uint64]*types.Block)
	pb.lock.Unlock()
}

func (pb *ProducerBlocks) AddBlock(block *types.Block) {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	if pb.baseBlockNum < block.NumberU64() {
		pb.blocks[block.NumberU64()] = block
	}
}

func (pb *ProducerBlocks) ExistBlock(block *types.Block) bool {
	pb.lock.Lock()
	defer pb.lock.Unlock()
	if v, ok := pb.blocks[block.NumberU64()]; ok && v.Hash() == block.Hash() {
		return true
	}
	return false
}
func (pb *ProducerBlocks) Author() discover.NodeID {
	pb.lock.Lock()
	pb.lock.Unlock()
	return pb.author
}

func (pb *ProducerBlocks) MaxSequenceBlockNum() uint64 {
	pb.lock.Lock()
	pb.lock.Unlock()
	num := pb.baseBlockNum
	for {
		_, ok := pb.blocks[num]
		if !ok {
			break
		}
		num++
	}
	return num
}

func (pb *ProducerBlocks) MaxSequenceBlock() *types.Block {
	pb.lock.Lock()
	pb.lock.Unlock()
	var block *types.Block
	num := pb.baseBlockNum + 1
	for {
		var ok bool
		block, ok = pb.blocks[num]
		if !ok {
			break
		}
		num++
	}
	return block
}

func (pb *ProducerBlocks) Len() int {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	return len(pb.blocks)
}

type Syncing struct {
}

//func (cbft *Cbft) handleSync(p *peer, msg p2p.Msg) error {
//	switch {
//	case msg.Code == ConsensusStateMsg:
//		var request consensusState
//		if err := msg.Decode(&request); err != nil {
//			return errResp(ErrDecode, "%v: %v", msg, err)
//		}
//
//		for _, sa := range request.UnCommitted {
//			subSigns := cbft.SubVote(sa.BlockHash, sa.SignBits)
//			if len(subSigns) > 0 {
//				for _, s := range subSigns {
//					if err := p2p.Send(p.rw, PrepareVoteMsg, s); err != nil {
//						log.Error("send BlockSyncSignatureMsg failed", "peer", p.id, "err", err)
//						return err
//					}
//				}
//			}
//		}
//		if cbft.producerBlocks.MaxSequenceBlockNum() > request.MemMaxBlockNum {
//			num := request.MemMaxBlockNum
//			if request.MemMaxBlockNum < cbft.producerBlocks.baseBlockNum {
//				log.Warn("peer too low", "peer", p.id)
//				return nil
//			}
//			for {
//				num += 1
//				if b, ok := cbft.producerBlocks.blocks[num]; ok {
//					var signs []*blockSyncSignature
//					if vs, ok := cbft.votesState[b.Hash()]; ok {
//						signs = vs.Signs()
//					}
//					if err := p2p.Send(p.rw,
//						BlockSyncMsg,
//						&blockSync{
//							Block: b,
//							sign:  signs,
//						}); err != nil {
//						log.Error("send BlockSyncMsg failed", "peer", p.id, "err", err)
//						return err
//					}
//				}
//			}
//		}
//		if cbft.getRootIrreversible().block.NumberU64() < request.IrreversibleBlockNum {
//			cbft.sendConsensusState()
//		}
//	case msg.Code == BlockSyncSignatureMsg:
//		var bs blockSyncSignature
//		if err := msg.Decode(&bs); err != nil {
//			return errResp(ErrDecode, "%v: %v", msg, err)
//		}
//		cbft.OnBlockSignature(cbft.blockChain, bs.ID, bs.Sign)
//	case msg.Code == BlockSyncMsg:
//		var bs blockSync
//		if err := msg.Decode(&bs); err != nil {
//			return errResp(ErrDecode, "%v: %v", msg, err)
//		}
//		cbft.OnNewBlock(cbft.blockChain, bs.Block)
//		for _, sign := range bs.sign {
//			cbft.OnBlockSignature(cbft.blockChain, sign.ID, sign.Sign)
//		}
//	}
//
//	return nil
//}
func (cbft *Cbft) sendConsensusState() {
	//cbft.mux.Lock()
	//defer cbft.mux.Unlock()
	//maxBlockNum := cbft.producerBlocks.MaxSequenceBlockNum()
	//rootBlockNum := cbft.getRootIrreversible().block.Number()
	//if rootBlockNum.Add(rootBlockNum, big.NewInt(1)).Cmp(maxBlockNum) == 0 {
	//	return
	//}
	//cs := &consensusState{
	//	UnCommitted:          cbft.BlockVoteBitArray(),
	//	IrreversibleBlockNum: cbft.getRootIrreversible().block.Number(),
	//	MemMaxBlockNum:       maxBlockNum,
	//}
	//cbft.baseHandler.sendConsensusState(cs)
}

//func (cbft *Cbft) startSync() {
//	timer := time.NewTimer(time.Millisecond * SyncInterval)
//
//)
//
//type ConsensusState struct {
//	UnCommitted          []SignBitArray
//	IrreversibleBlockNum *big.Int
//	MemMaxBlockNum       *big.Int
//}
//
//type SignBitArray struct {
//	BlockHash common.Hash
//	BlockNum  *big.Int
//	SignBits  []byte
//}
//
//type CommittedSigns struct {
//	BlockHash common.Hash
//	BlockNum  *big.Int
//	SignBits  []byte
//}

//const (
//	ConsensusStateCode = 0
//	SignBitArrayCode   = 1
//	CommittedSignsCode = 2
//)
//
//var (
//	messages = []interface{}{
//		ConsensusState{},
//		SignBitArray{},
//		CommittedSigns{},
//	}
//)
//
//type Syncing struct {
//	cbft  *Cbft
//	timer *time.Timer
//}
//
//func (s *Syncing) Protocols() []p2p.Protocol {
//	return []p2p.Protocol{
//		{
//			Name:    "cbft",
//			Version: 1,
//			Length:  uint64(len(messages)),
//			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
//				return s.baseHandler(p, rw)
//			},
//		},
//	}
//}
//
//func (s *Syncing) baseHandler(p *p2p.Peer, rw p2p.MsgReadWriter) error {
//	s.cbft.dpos.NodeIndex(p.ID())
//	for {
//		msg, err := rw.ReadMsg()
//		if err != nil {
//			log.Error("read peer message error", "err", err)
//			return err
//		}
//		switch msg.Code {
//		case ConsensusStateCode:
//		case SignBitArrayCode:
//		case CommittedSignsCode:
//
//		}
//	}
//
//	return nil
//}
//
//func (s *Syncing) Run() {
//	timer := time.NewTimer(time.Millisecond * 500)
//
//	for {
//		select {
//		case <-timer.C:
//			//cbft.sendConsensusState()
//
//		case <-cbft.exitCh:
//			cbft.log.Debug("consensus engine exit")
//			return
//		}
//
//	}
//}
//
//		}
//	}
//}
//
//func (s *Syncing) OnNewBlock() {
//
//}
//
//func (s *Syncing) OnBlockSign() {
//
//}
