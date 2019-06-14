package cbft

import (
	"context"
	"encoding/json"
	"fmt"
	sort2 "sort"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/pkg/errors"
)

var (
	errEmptyRootBlock             = errors.New("empty root block")
	errExistViewChange            = errors.New("had viewchange")
	errNotExistViewChange         = errors.New("not exist viewchange")
	errTimestampNotEqual          = errors.New("timestamp not equal")
	errBlockHashNotEqual          = errors.New("block hash not equal")
	errViewChangeBlockNumTooLower = errors.New("block number too lower")
	errInvalidProposalAddr        = errors.New("invalid proposal address")
	errRecvViewTimeout            = errors.New("receive viewchange timeout")
	errTimestamp                  = errors.New("viewchange timestamp too low")
	errInvalidViewChangeVote      = errors.New("invalid viewchange vote")
	errInvalidConfirmNumTooLow    = errors.New("confirm block number lower than local prepare")
	errViewChangeForked           = errors.New("view change's baseblock is forked")
	errViewChangeBaseTooLow       = errors.New("local PrepareVote that have been signed higher than BaseBlockNum of view")
	emptyAddr                     = common.Address{}
)

type AcceptStatus int

const (
	Accept = iota
	Discard
	Cache
)

type PendingVote map[common.Hash]*prepareVote
type PendingBlock map[common.Hash]*prepareBlock
type ProcessingVote map[common.Hash]map[discover.NodeID]*prepareVote
type ViewChangeVotes map[common.Address]*viewChangeVote

type RoundState struct {
	//mux    sync.RWMutex // all about consensus round state run in same goroutine, mux is use on read
	master bool

	viewChange          *viewChange //current viewchange message
	viewChangeResp      *viewChangeVote
	viewChangeVotes     ViewChangeVotes //current viewchange response from other consensus nodes
	lastViewChange      *viewChange     //last viewchange message
	lastViewChangeVotes ViewChangeVotes //last viewchange response from other consensus nodes

	pendingVotes    PendingVote    //pending prepareVotes sign
	pendingBlocks   PendingBlock   //pending blocks
	processingVotes ProcessingVote //to be processed prepareVotes sign

	producerBlocks *ProducerBlocks
	blockExtMap    *BlockExtMap

	localHighestPrepareVoteNum uint64
}

func (vv ViewChangeVotes) String() string {
	if vv == nil {
		return ""
	}
	s := "["
	for k, v := range vv {
		s += fmt.Sprintf("[addr:%s, vote:%s]", k.String(), v.String())
	}
	s += "]"
	return s
}

func (vv ViewChangeVotes) Bits(cnt int) string {
	bitArray := NewBitArray(uint32(cnt))
	for _, v := range vv {
		bitArray.SetIndex(v.ValidatorIndex, true)
	}
	return bitArray.String()
}

func (vv ViewChangeVotes) MarshalJSON() ([]byte, error) {
	type Vote struct {
		Address common.Address  `json:"address"`
		Vote    *viewChangeVote `json:"vote"`
	}
	type Votes struct {
		Votes []*Vote `json:"votes"`
	}

	votes := &Votes{
		Votes: make([]*Vote, 0),
	}
	for k, v := range vv {
		votes.Votes = append(votes.Votes, &Vote{
			Address: k,
			Vote:    v,
		})
	}
	return json.Marshal(votes)
}

func (vv ViewChangeVotes) UnmarshalJSON(b []byte) error {
	type Vote struct {
		Address common.Address  `json:"address"`
		Vote    *viewChangeVote `json:"vote"`
	}
	type Votes struct {
		Votes []*Vote `json:"votes"`
	}
	var votes Votes
	err := json.Unmarshal(b, &votes)
	if err != nil {
		return err
	}

	for _, vote := range votes.Votes {
		vv[vote.Address] = vote.Vote
	}
	return nil
}

func (rs RoundState) String() string {

	return fmt.Sprintf("[ master:%v, viewChange:%s, viewChangeResp:%s, viewChangeVotes:%s, lastViewChange:%s, lastViewChangeVotes:%s, pendingVotes:%s, pendingBlocks:%s, processingVotes:%s, localHighestPrepareVoteNum:%d, blockExtMap:%s",
		rs.master, rs.viewChange.String(), rs.viewChangeResp.String(), rs.viewChangeVotes.String(), rs.lastViewChange.String(), rs.lastViewChangeVotes.String(),
		rs.pendingVotes.String(), rs.pendingBlocks.String(), rs.processingVotes.String(), rs.localHighestPrepareVoteNum, rs.blockExtMap.BlockString())
}

func (pv PendingVote) Add(hash common.Hash, vote *prepareVote) {
	pv[hash] = vote
}
func (pv PendingVote) String() string {
	if pv == nil {
		return ""
	}
	s := "["
	for k, v := range pv {
		s += fmt.Sprintf("[hash:%s, vote:%s]", k.TerminalString(), v.String())
	}
	s += "]"
	return s
}

func (pv *PendingVote) Clear() {
	*pv = make(map[common.Hash]*prepareVote)
}

func (pb PendingBlock) Add(hash common.Hash, ext *prepareBlock) {
	pb[hash] = ext
}

func (pb PendingBlock) String() string {
	if pb == nil {
		return ""
	}
	s := "["
	for _, v := range pb {
		s += fmt.Sprintf("[block:%s]", v.String())
	}
	s += "]"
	return s
}

func (pb *PendingBlock) Clear() {
	*pb = make(map[common.Hash]*prepareBlock)
}

func (pv ProcessingVote) Add(hash common.Hash, nodeId discover.NodeID, vote *prepareVote) {
	var votes map[discover.NodeID]*prepareVote
	if v := pv[hash]; v != nil {
		votes = v
	} else {
		votes = make(map[discover.NodeID]*prepareVote)
	}
	votes[nodeId] = vote
	pv[hash] = votes
}

func (pv ProcessingVote) String() string {
	if pv == nil {
		return ""
	}
	s := "["
	for h, votes := range pv {
		s += fmt.Sprintf("hash:%s,[", h.TerminalString())
		for k, v := range votes {
			s += fmt.Sprintf("[nodeId:%s, vote:%s]", k.TerminalString(), v.String())
		}
		s += "]"
	}
	s += "]"
	return s
}

func (pv *ProcessingVote) Clear() {
	*pv = make(map[common.Hash]map[discover.NodeID]*prepareVote)
}

func (cbft *Cbft) init() {
	log.Info("Current block", "number", cbft.blockChain.CurrentBlock().Number())
	cbft.clear()
}

func (cbft *Cbft) SetLocalHighestPrepareNum(num uint64) {
	if cbft.localHighestPrepareVoteNum < num {
		cbft.localHighestPrepareVoteNum = num
	}
	cbft.log.Debug("SetLocalHighestPrepareNum", "local", cbft.localHighestPrepareVoteNum, "number", num)
}

func (cbft *Cbft) checkViewChangeVotes(votes []*viewChangeVote) error {
	if cbft.viewChange == nil {
		log.Error("ViewChange is nil, check prepareVotes failed")
		return errNotExistViewChange
	}

	for _, vote := range votes {
		if vote.EqualViewChange(cbft.viewChange) {
			if err := cbft.verifyValidatorSign(cbft.nextRoundValidator(cbft.viewChange.BaseBlockNum), vote.ValidatorIndex, vote.ValidatorAddr, vote, vote.Signature[:]); err != nil {
				log.Error("Verify validator failed", "vote", vote.String(), "err", err)
				return errInvalidViewChangeVotes
			}
		} else {
			log.Error("Invalid viewchange vote", "vote", vote.String(), "view", cbft.viewChange.String())
			return errInvalidViewChangeVote
		}
	}

	return nil
}

func (cbft *Cbft) verifyValidatorSign(blockNumber uint64, validatorIndex uint32, validatorAddr common.Address, msg ConsensusMsg, signature []byte) error {
	vds, err := cbft.agency.GetValidator(blockNumber)
	if err != nil {
		return err
	}
	if vn, err := vds.AddressIndex(validatorAddr); err == nil && uint32(vn.Index) == validatorIndex {
		buf, err := msg.CannibalizeBytes()
		if err != nil {
			return err
		}
		if !vn.Verify(buf, signature) {
			return errSign
		}
	} else {
		return err
	}
	return nil
}
func (cbft *Cbft) AgreeViewChange() bool {
	//cbft.mux.Lock()
	//defer cbft.mux.Unlock()
	return cbft.agreeViewChange()
}

func (cbft *Cbft) addPrepareBlockVote(pbd *prepareBlock) {
	if cbft.viewChange == nil {
		return
	}
	pbd.Timestamp = cbft.viewChange.Timestamp
	except := cbft.viewChange.BaseBlockNum + 1
	log.Info("add prepare block", "number", pbd.Block.NumberU64(), "except", except, "irr", cbft.viewChange.BaseBlockNum)
	log.Info(fmt.Sprintf("master:%v prepareVotes:%d ", cbft.master, len(cbft.viewChangeVotes)))
	if cbft.master && cbft.agreeViewChange() && except == pbd.Block.NumberU64() {
		pbd.View = cbft.viewChange
		pbd.ViewChangeVotes = make([]*viewChangeVote, 0)
		for _, v := range cbft.viewChangeVotes {
			pbd.ViewChangeVotes = append(pbd.ViewChangeVotes, v)
		}
		log.Debug("add prepareVotes in prepare block", "hash", pbd.Block.Hash(), "number", pbd.Block.NumberU64())
	} else {
		pbd.View = cbft.viewChange.CopyWithoutVotes()
	}
}
func (cbft *Cbft) agreeViewChange() bool {
	return len(cbft.viewChangeVotes) >= cbft.getThreshold()
}

//
//func (cbft *Cbft) AgreeReceive(block *types.Block) bool {
//	//cbft.mux.Lock()
//	//defer cbft.mux.Unlock()
//	return cbft.agreeReceive(block)
//}

//func (cbft *Cbft) agreeReceive(block *types.Block) bool {
//	if cbft.viewChange == nil || cbft.hadSendViewChange() {
//		return true
//	}
//
//	if cbft.viewChange.BaseBlockNum < block.NumberU64() || cbft.viewChange.BaseBlockHash == block.Hash() {
//		return true
//	} else {
//		if cbft.producerBlocks != nil && cbft.producerBlocks.ExistBlock(block) {
//			return true
//		}
//	}
//
//	cbft.log.Warn("refuse receive block", "hash", block.Hash(), "number", block.NumberU64(), "view", cbft.viewChange.String())
//	return false
//}

func (cbft *Cbft) viewVoteState() string {
	var state string
	for k, v := range cbft.viewChangeVotes {
		state += k.String() + "=" + v.String()
	}
	return state
}

func (cbft *Cbft) HadSendViewChange() bool {
	//cbft.mux.Lock()
	//defer cbft.mux.Unlock()
	return cbft.hadSendViewChange()
}

func (cbft *Cbft) hadSendViewChange() bool {
	return cbft.viewChanging() && cbft.master
}

func (cbft *Cbft) validViewChange() bool {
	//check current timestamp match view's timestamp
	maxViewProducerBlocksLimit := uint64(cbft.config.Duration) / cbft.config.Period
	now := time.Now().Unix()
	return (now-int64(cbft.viewChange.Timestamp) < cbft.config.Duration) && (cbft.producerBlocks == nil || uint64(cbft.producerBlocks.Len()) < maxViewProducerBlocksLimit)
}

func (cbft *Cbft) viewChanging() bool {
	return cbft.viewChange != nil
}

//func (cbft *Cbft) AcceptBlock(hash common.Hash, number uint64) bool {
//	//cbft.mux.Lock()
//	//defer cbft.mux.Unlock()
//
//	//1. not in viewchanging
//	if cbft.viewChanging() && !cbft.agreeViewChange() {
//		// changing
//		if number <= cbft.viewChange.BaseBlockNum {
//			log.Debug("accept true", "number", number, "hash", hash, "irr num", cbft.viewChange.BaseBlockNum)
//			return true
//		}
//	} else {
//		log.Warn("accept true", "view", cbft.viewChange.String(), "view prepareVotes", cbft.viewVoteState())
//		return true
//	}
//
//	if cbft.viewChange != nil && cbft.viewChangeVotes != nil {
//		log.Warn("accept false", "view", cbft.viewChange.String(), "view prepareVotes", cbft.viewVoteState())
//	} else {
//		log.Warn("accept false", "has view change", cbft.viewChange != nil, "has view vote", cbft.viewChangeVotes != nil)
//	}
//
//	return false
//}

func (cbft *Cbft) AcceptPrepareBlock(request *prepareBlock) AcceptStatus {
	if cbft.viewChange == nil {
		//todo need check prepareblock is belong to last accepted viewchange
		cbft.log.Debug("Accept block, viewchange is empty")
		return Accept
	}

	if cbft.viewChange.Timestamp > request.Timestamp && cbft.viewChange.BaseBlockNum < request.Block.NumberU64() {
		cbft.log.Debug("Cache block", "view", cbft.viewChange.String(), "prepareBlock", request.String())
		return Cache
	}

	if cbft.viewChange.BaseBlockNum > request.Block.NumberU64() {
		return Accept
	}

	if request.ProposalIndex != cbft.viewChange.ProposalIndex {
		if cbft.lastViewChange == nil {
			cbft.log.Debug("Discard block, lastViewChange is empty")
			return Discard
		}
		if request.ProposalIndex == cbft.lastViewChange.ProposalIndex {
			cbft.log.Debug("Cache block, ProposalIndex is lastviewchange")
			return Cache
		} else {
			cbft.log.Debug("Discard block, unknown block")
			return Discard
		}
	}

	if cbft.viewChanging() && !cbft.agreeViewChange() {
		cbft.log.Debug("Cache block, is viewchanging")
		return Cache
	}
	return Accept
}

func (cbft *Cbft) AcceptPrepareVote(vote *prepareVote) AcceptStatus {
	if vote.Number < cbft.getHighestConfirmed().number {
		cbft.log.Debug("Discard prepare vote, vote's number lower than local confirmed")
		return Discard
	}
	if (cbft.lastViewChange != nil && vote.Number < cbft.lastViewChange.BaseBlockNum) ||
		(cbft.viewChange != nil && vote.Number < cbft.viewChange.BaseBlockNum) {
		return Accept
	}
	//1. not in viewchanging
	if cbft.viewChanging() && !cbft.agreeViewChange() {
		// changing, if vote's timestamp equal viewchanging's timestamp, local is too slower than other,need to accept vote
		if vote.Number <= cbft.viewChange.BaseBlockNum || vote.Timestamp == cbft.viewChange.Timestamp {
			log.Debug("Accept vote", "hash", vote.Hash, "number", vote.Number, "irr num", cbft.viewChange.BaseBlockNum)
			return Accept
		}
	} else {
		//log.Warn("Accept vote", "view", cbft.viewChange.String(), "view prepareVotes", cbft.viewVoteState())
		if cbft.viewChange == nil || (cbft.viewChange != nil && vote.Timestamp == cbft.viewChange.Timestamp) {
			return Accept
		}
	}

	if cbft.viewChange != nil && cbft.viewChangeVotes != nil {
		log.Warn("Cache vote", "view", cbft.viewChange.String(), "view prepareVotes", cbft.viewVoteState())
	} else {
		log.Warn("Cache vote", "viewchange", cbft.viewChange != nil, "has view vote", cbft.viewChangeVotes != nil)
	}

	return Cache
}

func (cbft *Cbft) ClearPending() {
	cbft.clearPending()
}

func (cbft *Cbft) clearPending() {
	cbft.pendingVotes.Clear()
	cbft.pendingBlocks.Clear()
	cbft.processingVotes.Clear()
}

func (cbft *Cbft) ClearViewChange() {
	cbft.clearViewChange()
}

func (cbft *Cbft) clearViewChange() {
	cbft.viewChange = nil
	cbft.viewChangeVotes = make(map[common.Address]*viewChangeVote)
	cbft.lastViewChange = nil
	cbft.lastViewChangeVotes = make(map[common.Address]*viewChangeVote)
}

func (cbft *Cbft) Clear() {
	cbft.clear()
}

func (cbft *Cbft) clear() {
	cbft.clearPending()
	cbft.clearViewChange()
}

func (cbft *Cbft) handleCache() {
	votes := cbft.processingVotes
	cbft.processingVotes.Clear()
	go cbft.processing(votes)
	cbft.pendingProcess()
}

func (cbft *Cbft) processing(votes ProcessingVote) {
	for _, v := range votes {
		for k, v := range v {
			cbft.peerMsgCh <- &MsgInfo{
				Msg:    v,
				PeerID: k,
			}
		}
	}
}

func (cbft *Cbft) pendingProcess() {
	var pendingVote PendingVote
	var pendingBlock PendingBlock

	if len(cbft.pendingVotes) != 0 {
		pendingVote = cbft.pendingVotes
		cbft.pendingVotes = make(map[common.Hash]*prepareVote)
		cbft.log.Trace("Process pending vote", "total", len(pendingVote))

	}
	if len(cbft.pendingBlocks) != 0 {
		pendingBlock = cbft.pendingBlocks
		cbft.log.Trace("Process pending block", "total", len(pendingBlock))
		cbft.pendingBlocks = make(map[common.Hash]*prepareBlock)
	}

	for _, pv := range pendingVote {
		cbft.log.Debug("Handle cache pending votes", "hash", pv.Hash, "number", pv.Number)
		cbft.SetLocalHighestPrepareNum(pv.Number)
		ext := cbft.blockExtMap.findBlock(pv.Hash, pv.Number)
		if ext != nil {
			ext.prepareVotes.Add(pv)
			cbft.blockExtMap.Add(pv.Hash, pv.Number, ext)
		}
		cbft.handler.SendAllConsensusPeer(pv)
	}

	for _, v := range pendingBlock {
		cbft.log.Debug("Handle cache pending votes", "hash", v.Block.Hash(), "number", v.Block.Number())
		cbft.handler.SendAllConsensusPeer(v)
	}

}

func (cbft *Cbft) AddProcessingVote(nodeId discover.NodeID, vote *prepareVote) {
	//cbft.mux.Lock()
	//defer cbft.mux.Unlock()
	cbft.processingVotes.Add(vote.Hash, nodeId, vote)
}

func (cbft *Cbft) newViewChange() (*viewChange, error) {

	ext := cbft.getHighestConfirmed()

	if ext.number < cbft.localHighestPrepareVoteNum {
		//todo ask prepare vote to other, need optimize
		cbft.handler.SendAllConsensusPeer(&getHighestPrepareBlock{Lowest: ext.number})

		return nil, errInvalidConfirmNumTooLow
	}
	validator, err := cbft.getValidators().NodeIndexAddress(cbft.config.NodeID)
	if err != nil {
		return nil, errInvalidatorCandidateAddress
	}
	view := &viewChange{
		Timestamp:     uint64(time.Now().Unix()),
		BaseBlockNum:  ext.block.NumberU64(),
		BaseBlockHash: ext.block.Hash(),
		ProposalIndex: uint32(validator.Index),
		ProposalAddr:  validator.Address,
	}

	sign, err := cbft.signMsg(view)
	if err != nil {
		return nil, err
	}
	view.Signature.SetBytes(sign)
	view.BaseBlockPrepareVote = ext.Votes()
	cbft.resetViewChange()
	cbft.viewChange = view
	cbft.master = true
	log.Debug("Make new view change", "view", view.String(), "msgHash", view.MsgHash().TerminalString())
	return view, nil
}

func (cbft *Cbft) VerifyAndViewChange(view *viewChange) error {

	now := time.Now().UnixNano() / 1e6
	if !cbft.isLegal(now, view.ProposalAddr) {
		cbft.log.Error("Receive view change timeout", "current", now, "remote", view.Timestamp)
		return errRecvViewTimeout
	}

	if cbft.viewChange != nil && cbft.viewChange.Timestamp > view.Timestamp {
		cbft.log.Error("Verify view change failed", "local timestamp", cbft.viewChange.Timestamp, "remote", view.Timestamp)
		return errTimestamp
	}

	if cbft.localHighestPrepareVoteNum > view.BaseBlockNum {
		cbft.log.Error(fmt.Sprintf("Local highest PrepareVote's blocknum higher than view's BaseBlockNum hash:%s, number:%d localHighest:%d",
			view.BaseBlockHash.TerminalString(), view.BaseBlockNum, cbft.localHighestPrepareVoteNum))
		return errViewChangeBaseTooLow
	}

	ext := cbft.getHighestConfirmed()

	if ext.number == view.BaseBlockNum {
		if ext.block.Hash() != view.BaseBlockHash {
			cbft.log.Error(fmt.Sprintf("View's block is forked hash:%s, number:%d confirmed hash:%s, number %d localHighest:%d",
				view.BaseBlockHash.TerminalString(), view.BaseBlockNum, ext.block.Hash(), ext.number, cbft.localHighestPrepareVoteNum))
			return errViewChangeForked
		}
	} else if ext.number > view.BaseBlockNum {
		cbft.log.Warn(fmt.Sprintf("View change block too lower 2/3 block hash:%s num:%d, view irr block hash:%s num:%d",
			ext.block.Hash().TerminalString(),
			ext.number, view.BaseBlockHash.TerminalString(), view.BaseBlockNum))
		return errViewChangeBlockNumTooLower
	} else {
		cbft.log.Error(fmt.Sprintf("View's block is not found hash:%s, number:%d confirmed:%d, localHighest:%d",
			view.BaseBlockHash.TerminalString(), view.BaseBlockNum, ext.number, cbft.localHighestPrepareVoteNum))
		return errNotFoundViewBlock
	}

	if view.BaseBlockNum != 0 && len(view.BaseBlockPrepareVote) < cbft.getThreshold() {
		cbft.log.Error("View's prepare vote < 2f", "view", view.String())
		return errTwoThirdPrepareVotes
	}

	for _, vote := range view.BaseBlockPrepareVote {
		if err := cbft.verifyValidatorSign(view.BaseBlockNum, vote.ValidatorIndex, vote.ValidatorAddr, vote, vote.Signature[:]); err != nil {
			cbft.log.Error("Verify validator failed", "vote", vote.String(), "err", err)
			return errInvalidPrepareVotes
		}
	}

	return nil
}

func (cbft *Cbft) setViewChange(view *viewChange) {
	log.Info("Make viewchange vote", "vote", view.String())

	cbft.resetViewChange()
	cbft.viewChange = view
	cbft.master = false
}

func (cbft *Cbft) afterUpdateValidator() {
	cbft.master = false
}

func (cbft *Cbft) nextRoundValidator(blockNumber uint64) uint64 {
	return blockNumber + 1
}

func (cbft *Cbft) OnViewChangeVote(peerID discover.NodeID, vote *viewChangeVote) error {
	log.Debug("Receive view change vote", "peer", peerID, "vote", vote.String(), "view", cbft.viewChange.String())
	bpCtx := context.WithValue(context.Background(), "peer", peerID)
	cbft.bp.ViewChangeBP().ReceiveViewChangeVote(bpCtx, vote, cbft)
	if cbft.needBroadcast(peerID, vote) {
		go cbft.handler.SendBroadcast(vote)
	}
	//cbft.mux.Lock()
	//defer cbft.mux.Unlock()
	hadAgree := cbft.agreeViewChange()
	if cbft.viewChange != nil && vote.EqualViewChange(cbft.viewChange) {
		if err := cbft.verifyValidatorSign(cbft.nextRoundValidator(cbft.viewChange.BaseBlockNum), vote.ValidatorIndex, vote.ValidatorAddr, vote, vote.Signature[:]); err == nil {
			cbft.viewChangeVotes[vote.ValidatorAddr] = vote
			log.Info("Agree receive view change response", "peer", peerID, "viewChangeVotes", len(cbft.viewChangeVotes))
		} else {
			cbft.log.Warn("Verify sign failed", "peer", peerID, "vote", vote.String())
			return err
		}
	} else {
		switch {
		case cbft.viewChange == nil:
			cbft.bp.ViewChangeBP().InvalidViewChangeVote(bpCtx, vote, errNotExistViewChange, cbft)
			return errNotExistViewChange
		case vote.Timestamp != cbft.viewChange.Timestamp:
			cbft.bp.ViewChangeBP().InvalidViewChangeVote(bpCtx, vote, errTimestampNotEqual, cbft)
			return errTimestampNotEqual
		case vote.BlockHash != cbft.viewChange.BaseBlockHash:
			cbft.bp.ViewChangeBP().InvalidViewChangeVote(bpCtx, vote, errBlockHashNotEqual, cbft)
			return errBlockHashNotEqual
		case vote.ProposalAddr != cbft.viewChange.ProposalAddr:
			cbft.bp.ViewChangeBP().InvalidViewChangeVote(bpCtx, vote, errInvalidProposalAddr, cbft)
			return errInvalidProposalAddr
		default:
			return errInvalidViewChangeVotes
		}
	}

	if err := cbft.evPool.AddViewChangeVote(vote); err != nil {
		switch err.(type) {
		case *DuplicateViewChangeVoteEvidence:
		case *TimestampViewChangeVoteEvidence:
			cbft.log.Warn("Receive TimestampViewChangeVoteEvidence msg", "err", err.Error())
			return err
		}
	}

	if !hadAgree && cbft.agreeViewChange() {
		viewChangeVoteFulfillTimer.UpdateSince(time.Unix(int64(cbft.viewChange.Timestamp), 0))
		cbft.wal.UpdateViewChange(&ViewChangeMessage{
			Hash:   vote.BlockHash,
			Number: vote.BlockNum,
		})
		cbft.bp.ViewChangeBP().TwoThirdViewChangeVotes(bpCtx, cbft.viewChange, cbft.viewChangeVotes, cbft)
		cbft.flushReadyBlock()
		cbft.producerBlocks = NewProducerBlocks(cbft.config.NodeID, cbft.viewChange.BaseBlockNum)
		cbft.clearPending()
		cbft.ClearChildren(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum, cbft.viewChange.Timestamp)

		cbft.log.Info("Previous round state",
			"logicalNum", cbft.getHighestLogical().number,
			"logicalHash", cbft.getHighestLogical().block.Hash(),
			"logicalTimestamp", cbft.getHighestLogical().block.Time(),
			"logicalVoteBits", cbft.getHighestLogical().prepareVotes.voteBits.String(),
			"confirmedNum", cbft.getHighestConfirmed().number,
			"confirmedHash", cbft.getHighestConfirmed().block.Hash(),
			"confirmedTimestamp", cbft.getHighestConfirmed().block.Time(),
			"confirmedVoteBits", cbft.getHighestConfirmed().prepareVotes.voteBits.String(),
			"view", cbft.getHighestLogical().view.String(),
		)
	}

	log.Info("Receive viewchange vote", "msg", vote.String(), "had votes", len(cbft.viewChangeVotes), "voteBits", cbft.viewChangeVotes.Bits(cbft.getValidators().Len()))
	return nil
}

func (cbft *Cbft) ClearChildren(baseBlockHash common.Hash, baseBlockNum uint64, Timestamp uint64) {
	if cbft.getHighestLogical().number > baseBlockNum {
		if ext := cbft.blockExtMap.findBlock(baseBlockHash, baseBlockNum); ext != nil {
			cbft.highestLogical.Store(ext)
		}
	}
	cbft.blockExtMap.ClearChildren(cbft.viewChange.BaseBlockHash, cbft.viewChange.BaseBlockNum, cbft.viewChange.Timestamp)
}

func (cbft *Cbft) resetViewChange() {
	cbft.lastViewChange, cbft.lastViewChangeVotes = cbft.viewChange, cbft.viewChangeVotes
	cbft.viewChange, cbft.viewChangeVotes = nil, make(map[common.Address]*viewChangeVote)
}

func (cbft *Cbft) broadcastBlock(ext *BlockExt) {
	validator, err := cbft.getValidators().NodeIndexAddress(cbft.config.NodeID)
	if err != nil {
		return
	}
	ext.proposalIndex, ext.proposalAddr = uint32(validator.Index), validator.Address
	p := &prepareBlock{Block: ext.block, ProposalIndex: uint32(validator.Index), ProposalAddr: validator.Address}

	cbft.addPrepareBlockVote(p)
	ext.prepareBlock = p
	if cbft.viewChange != nil && !cbft.agreeViewChange() && cbft.viewChange.BaseBlockNum < ext.block.NumberU64() {
		log.Debug("Pending block", "number", ext.block.Number())
		cbft.pendingBlocks[ext.block.Hash()] = p
		return
	} else {
		log.Debug("Send block", "nodeID", cbft.config.NodeID, "number", ext.block.Number(), "hash", ext.block.Hash())
		cbft.bp.PrepareBP().SendBlock(context.TODO(), p, cbft)

		cbft.handler.SendAllConsensusPeer(p)
	}
}

func (cbft *Cbft) AddPrepareBlock(block *types.Block) {
	//cbft.mux.Lock()
	//defer cbft.mux.Unlock()
	if cbft.hadSendViewChange() {
		if block.NumberU64() > cbft.viewChange.BaseBlockNum {
			if cbft.producerBlocks == nil {
				cbft.producerBlocks = NewProducerBlocks(cbft.config.NodeID, cbft.viewChange.BaseBlockNum)
			}
			cbft.producerBlocks.AddBlock(block)
		}
	}
}

type GetBlock struct {
	hash   common.Hash
	number uint64
	ch     chan *types.Block
}

type prepareVoteSet struct {
	votes    map[uint32]*prepareVote
	voteBits *BitArray
}

func NewPrepareVoteSet(threshold uint32) *prepareVoteSet {
	return &prepareVoteSet{
		votes:    make(map[uint32]*prepareVote),
		voteBits: NewBitArray(threshold),
	}
}

func (pv *prepareVoteSet) Add(vote *prepareVote) {
	pv.votes[vote.ValidatorIndex] = vote
	pv.voteBits.setIndex(vote.ValidatorIndex, true)
}
func (pv *prepareVoteSet) Get(index uint32) *prepareVote {
	if pv.voteBits.GetIndex(index) {
		return pv.votes[index]
	}
	return nil
}
func (pv *prepareVoteSet) Merge(vs *prepareVoteSet) {
	for k, v := range vs.votes {
		pv.votes[k] = v
		pv.voteBits.setIndex(k, true)
	}
}

func (pv *prepareVoteSet) IsMaj23() bool {
	if pv == nil {
		return false
	}
	return uint32(len(pv.votes)) >= pv.voteBits.Size()
}

func (pv *prepareVoteSet) Signs() []common.BlockConfirmSign {
	signs := make([]common.BlockConfirmSign, 0)

	for _, v := range pv.votes {
		signs = append(signs, v.Signature)
	}
	return signs
}

func (pv *prepareVoteSet) Votes() []*prepareVote {
	votes := make([]*prepareVote, 0)
	for _, v := range pv.votes {
		votes = append(votes, v)
	}
	return votes
}

func (pv *prepareVoteSet) Len() int {
	return len(pv.votes)
}

// BlockExt is an extension from Block
type BlockExt struct {
	block           *types.Block
	inTree          bool
	inTurn          bool
	executing       bool
	isExecuted      bool
	isSigned        bool
	isConfirmed     bool
	number          uint64
	rcvTime         int64
	timestamp       uint64
	proposalIndex   uint32
	proposalAddr    common.Address
	prepareVotes    *prepareVoteSet //all prepareVotes for block
	prepareBlock    *prepareBlock
	view            *viewChange
	viewChangeVotes []*viewChangeVote
	parent          *BlockExt
	children        map[common.Hash]*BlockExt
	syncState       chan error
}

func (b BlockExt) MarshalJSON() ([]byte, error) {
	type BlockExt struct {
		Timestamp       uint64      `json:"timestamp"`
		InTree          bool        `json:"in_tree"`
		InTurn          bool        `json:"in_turn"`
		Executing       bool        `json:"executing"`
		IsExecuted      bool        `json:"is_executed"`
		IsSigned        bool        `json:"is_signed"`
		IsConfirmed     bool        `json:"is_confirmed"`
		Number          uint64      `json:"block_number"`
		RcvTime         int64       `json:"receive_time"`
		Hash            common.Hash `json:"block_hash"`
		Parent          common.Hash `json:"parent_hash"`
		ViewChangeVotes int         `json:"viewchange_votes"`
		PrepareVotes    int         `json:"prepare_votes"`
	}
	ext := BlockExt{
		Timestamp:       b.timestamp,
		InTree:          b.inTree,
		InTurn:          b.inTurn,
		Executing:       b.executing,
		IsExecuted:      b.isExecuted,
		IsSigned:        b.isSigned,
		IsConfirmed:     b.isConfirmed,
		Number:          b.number,
		RcvTime:         b.rcvTime,
		ViewChangeVotes: len(b.viewChangeVotes),
		PrepareVotes:    b.prepareVotes.Len(),
	}
	if b.block != nil {
		ext.Hash = b.block.Hash()
		ext.Parent = b.block.ParentHash()
	}

	return json.Marshal(&ext)
}

func (b BlockExt) String() string {
	if b.block == nil {
		return fmt.Sprintf("number:%d inTree:%v inTurn:%v isExecuted:%v, isSigned:%v isConfirmed:%v timestamp:%d rcvTime:%d prepareVotes:%d children:%d",
			b.number, b.inTree, b.inTurn, b.isExecuted, b.isSigned, b.isConfirmed, b.timestamp, b.rcvTime, b.prepareVotes.Len(), len(b.children))
	}
	return fmt.Sprintf("hash:%s number:%d inTree:%v inTurn:%v isExecuted:%v, isSigned:%v isConfirmed:%v timestamp:%d rcvTime:%d prepareVotes:%d children:%d",
		b.block.Hash().TerminalString(), b.block.NumberU64(), b.inTree, b.inTurn, b.isExecuted, b.isSigned, b.isConfirmed, b.timestamp, b.rcvTime, b.prepareVotes.Len(), len(b.children))
}

func (b *BlockExt) SetSyncState(err error) {
	if b.syncState != nil {
		b.syncState <- err
		b.syncState = nil
	}
}
func (b *BlockExt) PrepareBlock() (*prepareBlock, error) {

	if b.prepareBlock == nil {
		return nil, errors.Errorf("empty block")
	}
	return b.prepareBlock, nil
}

func (b *BlockExt) IsParent(hash common.Hash) bool {
	return b.block.Hash() == hash
}
func (b *BlockExt) Merge(ext *BlockExt) {
	if b != ext && b.number == ext.number {
		if b.block == nil && ext.block != nil {
			//receive PrepareVote before receive PrepareBlock, so view is need set
			b.block = ext.block
			b.view = ext.view
		}
		if b.proposalAddr == emptyAddr {
			b.proposalAddr = ext.proposalAddr
			b.proposalIndex = ext.proposalIndex
		}
		if b.prepareBlock == nil {
			b.prepareBlock = ext.prepareBlock
		}
		b.prepareVotes.Merge(ext.prepareVotes)

		if ext.syncState != nil && b.syncState != nil {
			panic("invalid syncState: double state channel")
		}

		if ext.syncState != nil {
			b.syncState = ext.syncState
		}
	}
}
func (b BlockExt) Signs() []common.BlockConfirmSign {
	return b.prepareVotes.Signs()
}

func (b BlockExt) Votes() []*prepareVote {
	return b.prepareVotes.Votes()
}

func (b BlockExt) BlockExtra() *BlockExtra {
	return &BlockExtra{
		Prepare:         b.Votes(),
		ViewChange:      b.view,
		ViewChangeVotes: b.viewChangeVotes,
	}
}

// New creates a BlockExt object
func NewBlockExt(block *types.Block, blockNum uint64, threshold int) *BlockExt {
	return &BlockExt{
		block:        block,
		number:       blockNum,
		prepareVotes: NewPrepareVoteSet(uint32(threshold)),
		children:     make(map[common.Hash]*BlockExt),
	}
}

func NewBlockExtBySeal(block *types.Block, blockNum uint64, threshold int) *BlockExt {
	return &BlockExt{
		block:        block,
		number:       blockNum,
		prepareVotes: NewPrepareVoteSet(uint32(threshold)),
		children:     make(map[common.Hash]*BlockExt),
		rcvTime:      common.Millis(time.Now()),
		inTree:       true,
		executing:    true,
		isExecuted:   true,
		isSigned:     true,
		isConfirmed:  true,
	}
}

func NewBlockExtByPeer(block *types.Block, blockNum uint64, threshold int) *BlockExt {
	return &BlockExt{
		block:        block,
		number:       blockNum,
		prepareVotes: NewPrepareVoteSet(uint32(threshold)),
		children:     make(map[common.Hash]*BlockExt),
		rcvTime:      common.Millis(time.Now()),
		inTree:       false,
		executing:    false,
		isExecuted:   false,
		isSigned:     false,
		isConfirmed:  false,
	}
}

func NewBlockExtByPrepareBlock(pb *prepareBlock, threshold int) *BlockExt {
	return &BlockExt{
		block:         pb.Block,
		view:          pb.View,
		number:        pb.Block.NumberU64(),
		prepareVotes:  NewPrepareVoteSet(uint32(threshold)),
		children:      make(map[common.Hash]*BlockExt),
		rcvTime:       common.Millis(time.Now()),
		timestamp:     pb.Timestamp,
		proposalIndex: pb.ProposalIndex,
		proposalAddr:  pb.ProposalAddr,
		prepareBlock:  pb,
		inTree:        false,
		executing:     false,
		isExecuted:    false,
		isSigned:      false,
		isConfirmed:   false,
	}
}

type BlockExtra struct {
	Prepare         []*prepareVote
	ViewChange      *viewChange
	ViewChangeVotes []*viewChangeVote
}

type ExecuteBlockStatus struct {
	block *BlockExt
	err   error
}

type SealBlock struct {
	block        *types.Block
	sealResultCh chan<- *types.Block
	stopCh       <-chan struct{}
}

func NewSealBlock(block *types.Block, sealResultCh chan<- *types.Block, stopCh <-chan struct{}) *SealBlock {
	return &SealBlock{
		block:        block,
		sealResultCh: sealResultCh,
		stopCh:       stopCh,
	}
}

type BlockExtMap struct {
	baseBlockNum uint64
	head         *BlockExt
	blocks       map[uint64]map[common.Hash]*BlockExt
	threshold    int
}

func NewBlockExtMap(baseBlock *BlockExt, threshold int) *BlockExtMap {
	extMap := &BlockExtMap{
		head:      baseBlock,
		blocks:    make(map[uint64]map[common.Hash]*BlockExt),
		threshold: threshold,
	}
	extMap.Add(baseBlock.block.Hash(), baseBlock.number, baseBlock)
	return extMap
}

func (bm *BlockExtMap) Add(hash common.Hash, number uint64, blockExt *BlockExt) {
	if extMap, ok := bm.blocks[number]; ok {
		log.Debug(fmt.Sprintf("hash:%s, number:%d", hash.TerminalString(), number))
		if ext, ok := extMap[hash]; ok {
			log.Debug(fmt.Sprintf("hash:%s, number:%d", hash.TerminalString(), number))
			ext.Merge(blockExt)
			if ext.prepareVotes.IsMaj23() {
				bm.removeFork(number, hash)
			}
			if ext.block != nil {
				bm.fixChain(ext)
			}
		} else {
			log.Debug(fmt.Sprintf("hash:%s, number:%d", hash.TerminalString(), number))
			if blockExt.prepareVotes.IsMaj23() {
				bm.removeFork(number, hash)
			}
			extMap[hash] = blockExt
			if blockExt.block != nil {
				bm.fixChain(blockExt)
			}
		}
	} else {
		log.Debug(fmt.Sprintf("hash:%s, number:%d", hash.TerminalString(), number))

		extMap := make(map[common.Hash]*BlockExt)
		extMap[hash] = blockExt
		bm.blocks[number] = extMap
		if blockExt.prepareVotes.IsMaj23() {
			bm.removeFork(number, hash)
		}
		if blockExt.block != nil {
			bm.fixChain(blockExt)
		}
	}
}

func (bm *BlockExtMap) removeFork(number uint64, hash common.Hash) {
	if extMap, ok := bm.blocks[number]; ok {
		for k, v := range extMap {
			if k != hash {
				if v.prepareVotes.IsMaj23() {
					panic(fmt.Sprintf("forked block has 2f+1 prepare votes:%s", k.TerminalString()))
				}
				if v.parent != nil {
					delete(v.parent.children, k)
				}
				if v.children != nil {
					for _, p := range v.children {
						p.parent = nil
					}
				}
			}
		}
	}
}

func (bm *BlockExtMap) fixChain(blockExt *BlockExt) {
	if bm.baseBlockNum > blockExt.number {
		panic(fmt.Sprintf("Insert invalid block base number is %d, insert chain is %d", bm.baseBlockNum, blockExt.number))
		return
	}

	if blockExt.prepareVotes.Len() >= bm.threshold {
		log.Debug("Block is confirmed", "hash", blockExt.block.Hash(), "number", blockExt.number)
		blockExt.isConfirmed = true
		blockMinedTimer.UpdateSince(common.MillisToTime(blockExt.rcvTime))
	}

	parent := bm.findParent(blockExt.block.ParentHash(), blockExt.number)
	child := bm.findChild(blockExt.block.Hash(), blockExt.number)

	if parent != nil {
		parent.children[blockExt.block.Hash()] = blockExt
		blockExt.parent = parent
	}

	if child != nil {
		child.parent = blockExt
		blockExt.children[child.block.Hash()] = child
	}

}
func (bm *BlockExtMap) BlockString() string {
	var blockStr string

	var keys []float64
	for k, _ := range bm.blocks {
		keys = append(keys, float64(k))
	}
	sort2.Float64s(keys)

	for _, k := range keys {
		for _, ext := range bm.blocks[uint64(k)] {
			if ext.block != nil {
				blockStr += fmt.Sprintf("[Hash:%s, Number:%d PrepareVotes:%d, Execute:%v, %d ", ext.block.Hash().TerminalString(), ext.block.NumberU64(), ext.prepareVotes.Len(), ext.isExecuted, ext.timestamp)
				for _, v := range ext.children {
					blockStr += fmt.Sprintf("child[%s,%d]", v.block.Hash().TerminalString(), v.block.NumberU64())
				}
				blockStr += fmt.Sprintf("]")
			} else {
				blockStr += fmt.Sprintf("[Hash:{}, Number:%d PrepareVotes:%d, Execute:%v, %d", ext.number, ext.prepareVotes.Len(), ext.isExecuted, ext.timestamp)
				for _, v := range ext.children {
					blockStr += fmt.Sprintf("child[%s,%d,%d]", v.block.Hash().TerminalString(), v.block.NumberU64(), v.timestamp)
				}
				blockStr += fmt.Sprintf("]")
			}
		}
	}
	return blockStr
}

func (bm *BlockExtMap) Len() int {
	return len(bm.blocks)
}

func (bm *BlockExtMap) Total() int {
	total := 0
	for _, v := range bm.blocks {
		total += len(v)
	}
	return total
}
func (bm *BlockExtMap) GetSubChainWithTwoThirdVotes(hash common.Hash, number uint64) []*BlockExt {
	base := bm.findBlock(hash, number)
	if base == nil || base.prepareVotes.Len() < bm.threshold {
		return nil
	}

	blockExts := make([]*BlockExt, 0)

	hash = bm.head.block.Hash()
	number = bm.head.number

	for be := bm.findChild(hash, number); be != nil && be.prepareVotes.Len() >= bm.threshold && be.isExecuted && be.number <= base.number; be = bm.findChild(hash, number) {
		blockExts = append(blockExts, be)
		hash = be.block.Hash()
		number = be.number
		log.Debug("GetSubChainWithTwoThirdVotes", "hash", hash.TerminalString(), "number", number)
	}
	if number != base.number {
		log.Error("GetSubChainWithTwoThirdVotesError", "number", number, "basenumber", base.number)
		return nil
	}

	return blockExts
}

func (bm *BlockExtMap) GetHasVoteWithoutBlock(highest uint64) []*HashNumberBits {
	wb := make([]*HashNumberBits, 0)
	for i := uint64(bm.head.number); i < highest; i++ {
		blocks := bm.blocks[i]
		if blocks != nil {
			for h, b := range blocks {
				if b.block == nil {
					wb = append(wb, &HashNumberBits{hash: h, number: b.number, bits: nil})
				}
			}
		}
	}
	return wb
}

func (bm *BlockExtMap) GetWithoutTwoThirdVotes(highest uint64) []*HashNumberBits {
	wb := make([]*HashNumberBits, 0)
	for i := uint64(bm.head.number); i < highest; i++ {
		blocks := bm.blocks[i]
		if blocks != nil {
			for h, b := range blocks {
				if b.prepareVotes.Len() < bm.threshold {
					wb = append(wb, &HashNumberBits{hash: h, number: b.number, bits: b.prepareVotes.voteBits})
				}
			}
		}
	}
	return wb
}

func (bm *BlockExtMap) GetSubChainUnExecuted() []*BlockExt {
	blockExts := make([]*BlockExt, 0)
	hash := bm.head.block.Hash()
	number := bm.head.number

	for be := bm.findChild(hash, number); be != nil && be.executing; be = bm.findChild(hash, number) {
		hash = be.block.Hash()
		number = be.number
	}

	for be := bm.findChild(hash, number); be != nil && !be.executing; be = bm.findChild(hash, number) {
		blockExts = append(blockExts, be)
		hash = be.block.Hash()
		number = be.number
	}

	log.Debug("subunexecuted", "head", bm.head.block.Hash(), "number", bm.head.number, "blocks", bm.BlockString())

	return blockExts
}

func (bm *BlockExtMap) ClearParents(hash common.Hash, number uint64) {
	base := bm.findBlock(hash, number)
	if base == nil {
		return
	}

	for n, blocks := range bm.blocks {
		if n < number {
			if blocks != nil {
				for _, b := range blocks {
					if b.children != nil {
						for _, p := range b.children {
							p.parent = nil
						}
					}
					b.children = nil
					b.parent = nil
				}
			}
			delete(bm.blocks, n)
		}
	}

	delete(bm.blocks, bm.head.number)
	//log.Debug("clear block", "number", number)
	//parentHash := base.block.ParentHash()
	//for be := bm.findParent(parentHash, number); be != nil && bm.head.block.Hash() != parentHash && be.prepareVotes.Len() >= bm.threshold; be = bm.findParent(parentHash, number) {
	//	delete(bm.blocks, be.number)
	//	parentHash = be.block.ParentHash()
	//	number = be.number
	//}
}

func (bm *BlockExtMap) ClearChildren(hash common.Hash, number uint64, timestamp uint64) {
	for i := number + 1; bm.blocks[i] != nil; i++ {
		log.Debug("clear block", "number", i)
		for hash, ext := range bm.blocks[i] {
			if ext.parent != nil {
				delete(ext.parent.children, hash)
			}
			if ext.timestamp != timestamp {
				ext.SetSyncState(nil)
				delete(bm.blocks[i], hash)
			}
			if len(bm.blocks[i]) == 0 {
				delete(bm.blocks, i)
			}
		}
	}
}

func (bm *BlockExtMap) RemoveBlock(block *BlockExt) {

}

func (bm *BlockExtMap) FindHighestConfirmed(hash common.Hash, number uint64) *BlockExt {
	var highest *BlockExt
	for be := bm.findChild(hash, number); be != nil && be.prepareVotes.Len() >= bm.threshold && be.isExecuted; be = bm.findChild(hash, number) {
		highest = be
		hash = be.block.Hash()
		number = be.number
	}
	return highest
}

func (bm *BlockExtMap) FindHighestConfirmedWithHeader() *BlockExt {
	var highest *BlockExt
	hash := bm.head.block.Hash()
	number := bm.head.block.NumberU64()
	for be := bm.findChild(hash, number); be != nil && be.prepareVotes.Len() >= bm.threshold && be.isExecuted; be = bm.findChild(hash, number) {
		highest = be
		hash = be.block.Hash()
		number = be.number
	}
	return highest
}

func (bm *BlockExtMap) FindHighestLogical(hash common.Hash, number uint64) *BlockExt {
	var highest *BlockExt
	for be := bm.findChild(hash, number); be != nil && be.block != nil && be.isExecuted; be = bm.findChild(hash, number) {
		highest = be
		hash = be.block.Hash()
		number = be.number
	}
	return highest
}

func (bm *BlockExtMap) BaseBlock(hash common.Hash, number uint64) *BlockExt {
	if b := bm.findBlock(hash, number); b != nil {
		bm.head = b
	} else {
		//it's impossible
		panic("New base block is none")
	}
	return bm.head
}

func (bm *BlockExtMap) findBlock(hash common.Hash, number uint64) *BlockExt {
	if extMap, ok := bm.blocks[number]; ok {
		for _, v := range extMap {
			if v.block != nil {
				if v.block.Hash() == hash {
					return v
				}
			}
		}
	}
	return nil
}

func (bm *BlockExtMap) findParent(hash common.Hash, number uint64) *BlockExt {
	if extMap, ok := bm.blocks[number-1]; ok {
		for _, v := range extMap {
			if v.block != nil {
				if v.block.Hash() == hash {
					return v
				}
			}
		}
	}
	return nil
}

func (bm *BlockExtMap) findChild(hash common.Hash, number uint64) *BlockExt {
	if extMap, ok := bm.blocks[number+1]; ok {
		for _, v := range extMap {
			if v.block != nil {
				if v.block.ParentHash() == hash {
					return v
				}
			}
		}
	}
	return nil
}

func (bm *BlockExtMap) findBlockByNumber(low, high uint64) types.Blocks {
	blocks := make([]*types.Block, 0)
	for i := low; i <= high; i++ {
		if extMap, ok := bm.blocks[i]; ok {
			for _, v := range extMap {
				if v.block != nil {
					blocks = append(blocks, v.block)
				}
			}
		}
	}
	return blocks
}

func (bm *BlockExtMap) findBlockByHash(hash common.Hash) *types.Block {
	for _, extMap := range bm.blocks {
		for existHash, ext := range extMap {
			if existHash == hash {
				return ext.block
			}
		}
	}
	return nil
}

func (bm *BlockExtMap) findBlockExtByNumber(low, high uint64) []*BlockExt {
	blocks := make([]*BlockExt, 0)
	for i := low; i <= high; i++ {
		if extMap, ok := bm.blocks[i]; ok {
			for _, v := range extMap {
				if v.block != nil {
					blocks = append(blocks, v)
				}
			}
		}
	}
	return blocks
}

type HashNumberBits struct {
	hash   common.Hash
	number uint64
	bits   *BitArray
}

type HasBlock struct {
	hash   common.Hash
	number uint64
	hasCh  chan bool
}
