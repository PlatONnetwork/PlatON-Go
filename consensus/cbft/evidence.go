package cbft

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/gogo/protobuf/test/enumstringer"
	"github.com/syndtr/goleveldb/leveldb"
	"sort"
)

var (
	//prefix +[timestamp/number] + address
	viewTimestampPrefix                = "vt"
	viewDualPrefix                     = "vh"
	prepareDualPrefix                  = "ph"
	errDuplicatePrepareVoteEvidence    = errors.New("duplicate prepare vote")
	errDuplicateViewChangeVoteEvidence = errors.New("duplicate view change")
	errTimestampViewChangeVoteEvidence = errors.New("view change timestamp out of order")
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

	if (front != nil && front.Timestamp == v.Timestamp) ||
		(back != nil && back.Timestamp == v.Timestamp) {
		return nil
	}
	if front != nil && front.BlockNum > v.BlockNum {
		return errTimestampViewChangeVoteEvidence
	}

	if back != nil && back.BlockNum < v.BlockNum {
		return errTimestampViewChangeVoteEvidence
	}

	*vt = append(*vt, v)
	sort.Sort(*vt)

	return nil
}

func (vt *TimeOrderViewChange) Remove(timestamp uint64) {
	i := 0

	for i < len(*vt) {
		if (*vt)[i].Timestamp > timestamp {
			break
		}
		i++
	}
	if i == len(*vt) {
		*vt = (*vt)[:0]
	} else {
		*vt = append((*vt)[:0], (*vt)[i:]...)
	}
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

func (vt NumberOrderViewChange) find(number uint64) *viewChangeVote {
	for _, v := range vt {
		if v.BlockNum == number {
			return v
		}
	}
	return nil
}
func (vt *NumberOrderViewChange) Add(v *viewChangeVote) error {
	ev := vt.find(v.BlockNum)
	if ev.BlockHash != ev.BlockHash {
		return errDuplicateViewChangeVoteEvidence
	}
	return nil
}

func (vt *NumberOrderViewChange) Remove(blockNum uint64) {
	i := 0

	for i < len(*vt) {
		if (*vt)[i].BlockNum > blockNum {
			break
		}
		i++
	}
	if i == len(*vt) {
		*vt = (*vt)[:0]
	} else {
		*vt = append((*vt)[:0], (*vt)[i:]...)
	}
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

func (vt NumberOrderPrepare) find(number uint64) *prepareVote {
	for _, v := range vt {
		if v.Number == number {
			return v
		}
	}
	return nil
}
func (vt *NumberOrderPrepare) Add(v *prepareVote) error {
	ev := vt.find(v.Number)
	if ev.Hash != ev.Hash {
		return errDuplicatePrepareVoteEvidence
	}
	return nil
}

func (vt *NumberOrderPrepare) Remove(blockNum uint64) {
	i := 0

	for i < len(*vt) {
		if (*vt)[i].Number > blockNum {
			break
		}
		i++
	}
	if i == len(*vt) {
		*vt = (*vt)[:0]
	} else {
		*vt = append((*vt)[:0], (*vt)[i:]...)
	}
}

func (vt *ViewTimeEvidence) Add(v *viewChangeVote) error {
	if l := (*vt)[v.ValidatorAddr]; l != nil {
		err := l.Add(v)
		(*vt)[v.ValidatorAddr] = l
		return err
	}
	return nil
}

func (vt *ViewTimeEvidence) Clear(timestamp uint64) {
	for k, v := range *vt {
		v.Remove(timestamp)
		if v.Len() == 0 {
			delete(*vt, k)
		}
	}
}

func (vt *ViewNumberEvidence) Add(v *viewChangeVote) error {
	if l := (*vt)[v.ValidatorAddr]; l != nil {
		err := l.Add(v)
		(*vt)[v.ValidatorAddr] = l
		return err
	}
	return nil
}

func (vt *ViewNumberEvidence) Clear(number uint64) {
	for k, v := range *vt {
		v.Remove(number)
		if v.Len() == 0 {
			delete(*vt, k)
		}
	}
}

func (vt *PrepareEvidence) Add(v *prepareVote) error {
	if l := (*vt)[v.ValidatorAddr]; l != nil {
		err := l.Add(v)
		(*vt)[v.ValidatorAddr] = l
		return err
	}
	return nil
}

func (vt *PrepareEvidence) Clear(number uint64) error {
	for k, v := range *vt {
		v.Remove(number)
		if v.Len() == 0 {
			delete(*vt, k)
		}
	}
}

type EvidencePool struct {
	vt     ViewTimeEvidence
	vn     ViewNumberEvidence
	pe     PrepareEvidence
	exitCh chan struct{}
	db     *leveldb.DB
}

func NewEvidencePool(path string) (*EvidencePool, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &EvidencePool{
			vt:     make(ViewTimeEvidence),
			vn:     make(ViewNumberEvidence),
			pe:     make(PrepareEvidence),
			exitCh: make(chan struct{}), db: db},
		nil
}

func (ev *EvidencePool) AddViewChangeVote(v *viewChangeVote) error {
	if err := verifyAddr(v, v.ValidatorAddr); err != nil {
		return nil
	}
	ev.vt.Add(v)
	ev.vn.Add(v)
	return nil
}

func (ev *EvidencePool) AddPrepareVote(p *prepareVote) error {
	if err := verifyAddr(p, p.ValidatorAddr); err != nil {
		return nil
	}
	ev.pe.Add(p)
	return nil
}

func (ev *EvidencePool) Clear(timestamp, blockNum uint64) error {
	ev.vt.Clear(timestamp)
	ev.vn.Clear(blockNum)
	ev.pe.Clear(blockNum)
}

func (ev *EvidencePool) ClearDB() error {
	return nil
}

func (ev *EvidencePool) Close() {
	ev.exitCh <- struct{}{}
}

func (ev *EvidencePool) Evidences() {

}

func verifyAddr(msg ConsensusMsg, addr common.Address) error {
	data, err := msg.CannibalizeBytes()
	recPubKey, err := crypto.Ecrecover(data, msg.Sign())
	if err != nil {
		return err
	}
	pub, err := crypto.UnmarshalPubkey(recPubKey)
	if err != nil {
		return err
	}

	recAddr := crypto.PubkeyToAddress(*pub)

	if !bytes.Equal(recAddr.Bytes(), addr.Bytes()) {
		return fmt.Errorf("validator's address is not match signature")
	}
}
