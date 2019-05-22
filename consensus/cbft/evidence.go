package cbft

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb"
	"sort"
)

var (
	//prefix +[timestamp/number] + address
	viewTimestampPrefix = "vt"
	viewDualPrefix      = "vh"
	prepareDualPrefix   = "ph"
)

type Evidence interface {
	Verify(ecdsa.PublicKey) error
	Equal(Evidence) error
	//return lowest number
	BlockNumber() uint64
}

type DuplicatePrepareVoteEvidence struct {
	voteA *prepareVote
	voteB *prepareVote
}

type DuplicateViewChangeVoteEvidence struct {
	voteA *viewChangeVote
	voteB *viewChangeVote
}

type TimestampViewChangeVoteEvidence struct {
	voteA *viewChangeVote
	voteB *viewChangeVote
}

type TimeOrderViewChange []*viewChangeVote
type NumberOrderViewChange []*viewChangeVote
type NumberOrderPrepare []*prepareVote

type ViewTimeEvidence map[common.Address]TimeOrderViewChange
type ViewNumberEvidence map[common.Address]NumberOrderViewChange
type PrepareEvidence map[common.Address]NumberOrderPrepare

func (vt TimeOrderViewChange) Len() int {
	return len(vt)
}

func (vt TimeOrderViewChange) Less(i, j int) bool {
	return vt[i].Timestamp > vt[j].Timestamp
}

func (vt TimeOrderViewChange) Swap(i, j int) {
	vt[i], vt[j] = vt[j], vt[i]
}

func (vt TimeOrderViewChange) findFrontAndBack(v *viewChangeVote) (*viewChangeVote, *viewChangeVote) {
	i := 0
	var front *viewChangeVote
	var back *viewChangeVote
	for i < len(vt) {
		if (vt)[i].Timestamp > v.Timestamp {
			break
		}
		i++
	}

	if i != 0 {
		front = nil
	} else {
		front = vt[i-1]
	}

	if len(vt) != i {
		back = vt[i]
	}
	return front, back
}
func (vt *TimeOrderViewChange) Add(v *viewChangeVote) error {
	front, back := vt.findFrontAndBack(v)
	if front != nil && front.Timestamp != v.Timestamp {

	}

	if back.Timestamp != v.Timestamp {

	}

	*vt = append(*vt, v)
	sort.Sort(*vt)

	return nil
}

func (vt TimeOrderViewChange) Remove(timestamp uint64) {

}

func (vt NumberOrderViewChange) Len() int {
	return len(vt)
}

func (vt NumberOrderViewChange) Less(i, j int) bool {
	return vt[i].BlockNum > vt[j].BlockNum
}

func (vt NumberOrderViewChange) Swap(i, j int) {
	vt[i], vt[j] = vt[j], vt[i]
}

func (vt NumberOrderViewChange) Add(v *viewChangeVote) error {
	return nil
}

func (vt NumberOrderViewChange) Remove(blockNum uint64) {

}

func (vt NumberOrderPrepare) Len() int {
	return len(vt)
}

func (vt NumberOrderPrepare) Less(i, j int) bool {
	return vt[i].Number > vt[j].Number
}

func (vt NumberOrderPrepare) Swap(i, j int) {
	vt[i], vt[j] = vt[j], vt[i]
}

func (vt NumberOrderPrepare) Add(v *prepareVote) error {
	return nil
}

func (vt NumberOrderPrepare) Remove(blockNum uint64) {

}

func (vt ViewTimeEvidence) Add(v *viewChangeVote) error {
	return nil
}

func (vt ViewTimeEvidence) Clear() error {
	return nil
}

func (vt ViewNumberEvidence) Add(v *viewChangeVote) error {
	return nil
}

func (vt ViewNumberEvidence) Clear() error {
	return nil
}

func (vt PrepareEvidence) Add(v *prepareVote) error {
	return nil
}

func (vt PrepareEvidence) Clear(v *prepareVote) error {
	return nil
}

type EvidencePool struct {
	exitCh chan struct{}
	db     *leveldb.DB
}

func NewEvidencePool(path string) (*EvidencePool, error) {
	return nil, nil
}

func (ev *EvidencePool) AddViewChangeVote(v *viewChangeVote) error {
	return nil
}

func (ev *EvidencePool) AddPrepareVote(p *prepareVote) error {
	return nil
}

func (ev *EvidencePool) Clear(timestamp, blockNum uint64) error {
	return nil
}

func (ev *EvidencePool) ClearDB() error {
	return nil
}

func (ev *EvidencePool) Close() error {
	return nil
}

func (ev *EvidencePool) Evidences() {

}
