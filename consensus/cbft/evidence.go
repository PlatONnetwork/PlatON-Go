package cbft

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"math/big"
	"sort"
)

var (
	//prefix +[number] + address + hash
	viewTimestampPrefix                = byte(0x1)
	viewDualPrefix                     = byte(0x2)
	prepareDualPrefix                  = byte(0x3)
	errDuplicatePrepareVoteEvidence    = errors.New("duplicate prepare vote")
	errDuplicateViewChangeVoteEvidence = errors.New("duplicate view change")
	errTimestampViewChangeVoteEvidence = errors.New("view change timestamp out of order")

	evidenceDir = "evidenceDir"
)

type Evidence interface {
	Verify(ecdsa.PublicKey) error
	Equal(Evidence) bool
	//return lowest number
	BlockNumber() uint64
	Hash() []byte
	Address() common.Address
	Validate() error
}

type EvidenceData struct {
	DP []*DuplicatePrepareVoteEvidence    `json:"duplicate_prepare"`
	DV []*DuplicateViewChangeVoteEvidence `json:"duplicate_viewchange"`
	TV []*TimestampViewChangeVoteEvidence `json:"timestamp_viewchange"`
}

func NewEvidenceData() *EvidenceData {
	return &EvidenceData{
		DP: make([]*DuplicatePrepareVoteEvidence, 0),
		DV: make([]*DuplicateViewChangeVoteEvidence, 0),
		TV: make([]*TimestampViewChangeVoteEvidence, 0),
	}
}
func ClassifyEvidence(evds []Evidence) *EvidenceData {
	ed := NewEvidenceData()
	for _, e := range evds {
		switch e.(type) {
		case *DuplicatePrepareVoteEvidence:
			ed.DP = append(ed.DP, e.(*DuplicatePrepareVoteEvidence))
		case *DuplicateViewChangeVoteEvidence:
			ed.DV = append(ed.DV, e.(*DuplicateViewChangeVoteEvidence))
		case *TimestampViewChangeVoteEvidence:
			ed.TV = append(ed.TV, e.(*TimestampViewChangeVoteEvidence))
		}
	}
	return ed
}

//Evidence A.Number == B.Number but A.Hash != B.Hash
type DuplicatePrepareVoteEvidence struct {
	VoteA *prepareVote
	VoteB *prepareVote
}

func (d DuplicatePrepareVoteEvidence) Verify(pub ecdsa.PublicKey) error {
	addr := crypto.PubkeyToAddress(pub)
	if err := verifyAddr(d.VoteA, addr); err != nil {
		return err
	}
	return verifyAddr(d.VoteB, addr)
}
func (d DuplicatePrepareVoteEvidence) Equal(ev Evidence) bool {
	_, ok := ev.(*DuplicatePrepareVoteEvidence)
	if !ok {
		return false
	}
	dh := d.Hash()
	eh := ev.Hash()
	return bytes.Equal(dh, eh)
}
func (d DuplicatePrepareVoteEvidence) BlockNumber() uint64 {
	return d.VoteA.Number
}

func (d DuplicatePrepareVoteEvidence) Address() common.Address {
	return d.VoteA.ValidatorAddr
}

func (d DuplicatePrepareVoteEvidence) Hash() []byte {
	var buf []byte
	if ac, err := d.VoteA.CannibalizeBytes(); err == nil {
		if bc, err := d.VoteB.CannibalizeBytes(); err == nil {
			buf, err = rlp.EncodeToBytes([]interface{}{
				ac,
				d.VoteA.Sign(),
				bc,
				d.VoteB.Sign(),
			})
		}
	}
	return crypto.Keccak256(buf)
}

func (d DuplicatePrepareVoteEvidence) Validate() error {
	if d.VoteA.Number != d.VoteB.Number {
		return fmt.Errorf("DuplicatePrepareVoteEvidence BlockNum is different, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}
	if d.VoteA.Hash == d.VoteB.Hash {
		return fmt.Errorf("DuplicatePrepareVoteEvidence BlockHash is equal, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}

	if d.VoteA.ValidatorIndex != d.VoteB.ValidatorIndex ||
		d.VoteA.ValidatorAddr != d.VoteB.ValidatorAddr {
		return fmt.Errorf("DuplicatePrepareVoteEvidence Validator do not match, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}

	if err := verifyAddr(d.VoteA, d.VoteA.ValidatorAddr); err != nil {
		return fmt.Errorf("DuplicatePrepareVoteEvidence Vote verify failed, VoteA:%s", d.VoteA.String())
	}
	if err := verifyAddr(d.VoteB, d.VoteB.ValidatorAddr); err != nil {
		return fmt.Errorf("DuplicatePrepareVoteEvidence Vote verify failed, VoteA:%s", d.VoteA.String())
	}
	return nil
}

func (d DuplicatePrepareVoteEvidence) Error() string {
	return fmt.Sprintf("DuplicatePrepareVoteEvidence ValidatorIndex:%d, ValidatorAddr:%s, blockNum:%d blockHashA:%s, blockHashB:%s",
		d.VoteA.ValidatorIndex, d.VoteB.ValidatorAddr, d.VoteA.Number, d.VoteA.Hash.String(), d.VoteB.Hash.String())
}

//Evidence A.BlockNum == B.BlockNum but A.BlockHash != B.BlockHash
type DuplicateViewChangeVoteEvidence struct {
	VoteA *viewChangeVote
	VoteB *viewChangeVote
}

func (d DuplicateViewChangeVoteEvidence) Verify(pub ecdsa.PublicKey) error {
	addr := crypto.PubkeyToAddress(pub)
	if err := verifyAddr(d.VoteA, addr); err != nil {
		return err
	}
	return verifyAddr(d.VoteB, addr)
}

func (d DuplicateViewChangeVoteEvidence) Equal(ev Evidence) bool {
	_, ok := ev.(*DuplicateViewChangeVoteEvidence)
	if !ok {
		return false
	}
	dh := d.Hash()
	eh := ev.Hash()
	return bytes.Equal(dh, eh)
}

func (d DuplicateViewChangeVoteEvidence) BlockNumber() uint64 {
	return d.VoteA.BlockNum
}

func (d DuplicateViewChangeVoteEvidence) Hash() []byte {
	var buf []byte
	if ac, err := d.VoteA.CannibalizeBytes(); err == nil {
		if bc, err := d.VoteB.CannibalizeBytes(); err == nil {
			buf, err = rlp.EncodeToBytes([]interface{}{
				ac,
				d.VoteA.Sign(),
				bc,
				d.VoteB.Sign(),
			})
		}
	}
	return crypto.Keccak256(buf)
}

func (d DuplicateViewChangeVoteEvidence) Address() common.Address {
	return d.VoteA.ValidatorAddr
}

func (d DuplicateViewChangeVoteEvidence) Validate() error {
	ba := new(big.Int).SetBytes(d.VoteA.BlockHash.Bytes())
	bb := new(big.Int).SetBytes(d.VoteB.BlockHash.Bytes())

	if ba.Cmp(bb) >= 0 {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence BlockHash do not match, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}

	if d.VoteA.BlockNum != d.VoteB.BlockNum {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence BlockNum is not equal, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}

	if d.VoteA.ValidatorIndex != d.VoteB.ValidatorIndex ||
		d.VoteA.ValidatorAddr != d.VoteB.ValidatorAddr {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence Validator do not match, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}

	if err := verifyAddr(d.VoteA, d.VoteA.ValidatorAddr); err != nil {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence Vote verify failed, VoteA:%s", d.VoteA.String())
	}
	if err := verifyAddr(d.VoteB, d.VoteB.ValidatorAddr); err != nil {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence Vote verify failed, VoteA:%s", d.VoteA.String())
	}
	return nil
}

func (d DuplicateViewChangeVoteEvidence) Error() string {
	return fmt.Sprintf("DuplicateViewChangeVoteEvidence timestamp:%d blockNumberA:%d, blockHashA:%s, blockNumberB:%d, blockHashB:%s",
		d.VoteA.Timestamp, d.VoteA.BlockNum, d.VoteA.BlockHash.String(), d.VoteB.BlockNum, d.VoteB.BlockHash.String())
}

//Evidence A.Timestamp < B.Timestamp but A.BlockNum > B.BlockNum
type TimestampViewChangeVoteEvidence struct {
	VoteA *viewChangeVote
	VoteB *viewChangeVote
}

func (d TimestampViewChangeVoteEvidence) Verify(pub ecdsa.PublicKey) error {
	addr := crypto.PubkeyToAddress(pub)
	if err := verifyAddr(d.VoteA, addr); err != nil {
		return err
	}
	return verifyAddr(d.VoteB, addr)
}

func (d TimestampViewChangeVoteEvidence) Equal(ev Evidence) bool {
	_, ok := ev.(*TimestampViewChangeVoteEvidence)
	if !ok {
		return false
	}
	dh := d.Hash()
	eh := ev.Hash()
	return bytes.Equal(dh, eh)
}

func (d TimestampViewChangeVoteEvidence) BlockNumber() uint64 {
	return d.VoteA.BlockNum
}

func (d TimestampViewChangeVoteEvidence) Hash() []byte {
	var buf []byte
	if ac, err := d.VoteA.CannibalizeBytes(); err == nil {
		if bc, err := d.VoteB.CannibalizeBytes(); err == nil {
			buf, err = rlp.EncodeToBytes([]interface{}{
				ac,
				d.VoteA.Sign(),
				bc,
				d.VoteB.Sign(),
			})
		}
	}
	return crypto.Keccak256(buf)
}

func (d TimestampViewChangeVoteEvidence) Address() common.Address {
	return d.VoteA.ValidatorAddr
}

func (d TimestampViewChangeVoteEvidence) Validate() error {
	if d.VoteA.Timestamp > d.VoteB.Timestamp {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence Timestamp do not match, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}

	if d.VoteA.BlockNum <= d.VoteB.BlockNum {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence BlockNum do not match, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}

	if d.VoteA.ValidatorIndex != d.VoteB.ValidatorIndex ||
		d.VoteA.ValidatorAddr != d.VoteB.ValidatorAddr {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence Validator do not match, VoteA:%s, VoteB:%s", d.VoteA.String(), d.VoteB.String())
	}

	if err := verifyAddr(d.VoteA, d.VoteA.ValidatorAddr); err != nil {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence Vote verify failed, VoteA:%s", d.VoteA.String())
	}
	if err := verifyAddr(d.VoteB, d.VoteB.ValidatorAddr); err != nil {
		return fmt.Errorf("DuplicateViewChangeVoteEvidence Vote verify failed, VoteA:%s", d.VoteA.String())
	}
	return nil
}

func (d TimestampViewChangeVoteEvidence) Error() string {
	return fmt.Sprintf("TimestampViewChangeVoteEvidence timestamp:%d blockNumberA:%d, blockHashA:%s, blockNumberB:%d, blockHashB:%s",
		d.VoteA.Timestamp, d.VoteA.BlockNum, d.VoteA.BlockHash.String(), d.VoteB.BlockNum, d.VoteB.BlockHash.String())
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

	if i == 0 {
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
		return &TimestampViewChangeVoteEvidence{
			VoteA: front,
			VoteB: v,
		}
	}

	if back != nil && back.BlockNum < v.BlockNum {
		return &TimestampViewChangeVoteEvidence{
			VoteA: v,
			VoteB: back,
		}
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
		(*vt)[i] = nil
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
	if ev := vt.find(v.BlockNum); ev != nil {

		if ev.BlockHash != v.BlockHash {
			a, b := v, ev
			ha := new(big.Int).SetBytes(v.BlockHash.Bytes())
			hb := new(big.Int).SetBytes(ev.BlockHash.Bytes())

			if ha.Cmp(hb) > 0 {
				a, b = ev, v
			}
			return &DuplicateViewChangeVoteEvidence{
				VoteA: a,
				VoteB: b,
			}
		}
	} else {
		*vt = append(*vt, v)
		sort.Sort(*vt)
	}
	return nil
}

func (vt *NumberOrderViewChange) Remove(blockNum uint64) {
	i := 0

	for i < len(*vt) {
		if (*vt)[i].BlockNum > blockNum {
			break
		}
		(*vt)[i] = nil
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
	if ev := vt.find(v.Number); ev != nil {
		if ev.Hash != v.Hash {
			a, b := v, ev
			ha := new(big.Int).SetBytes(v.Hash.Bytes())
			hb := new(big.Int).SetBytes(ev.Hash.Bytes())

			if ha.Cmp(hb) > 0 {
				a, b = ev, v
			}
			return &DuplicatePrepareVoteEvidence{
				VoteA: a,
				VoteB: b,
			}

		}
	} else {
		*vt = append(*vt, v)
		sort.Sort(*vt)
	}
	return nil
}

func (vt *NumberOrderPrepare) Remove(blockNum uint64) {
	i := 0

	for i < len(*vt) {
		if (*vt)[i].Number > blockNum {
			break
		}
		(*vt)[i] = nil
		i++
	}
	if i == len(*vt) {
		*vt = (*vt)[:0]
	} else {
		*vt = append((*vt)[:0], (*vt)[i:]...)
	}
}

func (vt ViewTimeEvidence) Add(v *viewChangeVote) error {
	var l TimeOrderViewChange
	if l = vt[v.ValidatorAddr]; l == nil {
		l = make(TimeOrderViewChange, 0)
	}
	err := l.Add(v)
	vt[v.ValidatorAddr] = l
	return err
}

func (vt ViewTimeEvidence) Clear(timestamp uint64) {
	for k, v := range vt {
		v.Remove(timestamp)
		if v.Len() == 0 {
			delete(vt, k)
		}
	}
}

func (vt ViewNumberEvidence) Add(v *viewChangeVote) error {
	var l NumberOrderViewChange

	if l = vt[v.ValidatorAddr]; l == nil {
		l = make(NumberOrderViewChange, 0)
	}
	err := l.Add(v)
	vt[v.ValidatorAddr] = l
	return err
}

func (vt ViewNumberEvidence) Clear(number uint64) {
	for k, v := range vt {
		v.Remove(number)
		if v.Len() == 0 {
			delete(vt, k)
		}
	}
}

func (vt PrepareEvidence) Add(v *prepareVote) error {
	var l NumberOrderPrepare

	if l = vt[v.ValidatorAddr]; l == nil {
		l = make(NumberOrderPrepare, 0)
	}
	err := l.Add(v)
	vt[v.ValidatorAddr] = l
	return err
}

func (vt PrepareEvidence) Clear(number uint64) {
	for k, v := range vt {
		v.Remove(number)
		if v.Len() == 0 {
			delete(vt, k)
		}
	}
}

type EvidencePool struct {
	vt ViewTimeEvidence
	vn ViewNumberEvidence
	pe PrepareEvidence
	db *leveldb.DB
}

func NewEvidencePool(path string) (*EvidencePool, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &EvidencePool{
		vt: make(ViewTimeEvidence),
		vn: make(ViewNumberEvidence),
		pe: make(PrepareEvidence),
		db: db,
	}, nil
}

func (ev *EvidencePool) AddViewChangeVote(v *viewChangeVote) error {
	if err := verifyAddr(v, v.ValidatorAddr); err != nil {
		return err
	}
	if err := ev.vt.Add(v); err != nil {
		if evidence, ok := err.(*TimestampViewChangeVoteEvidence); ok {
			if err := ev.commit(evidence); err != nil {
				return err
			}
			return err
		}
	}
	if err := ev.vn.Add(v); err != nil {
		if evidence, ok := err.(*DuplicateViewChangeVoteEvidence); ok {
			if err := ev.commit(evidence); err != nil {
				return err
			}
			return err
		}
	}
	return nil
}

func (ev *EvidencePool) AddPrepareVote(p *prepareVote) error {
	if err := verifyAddr(p, p.ValidatorAddr); err != nil {
		return err
	}
	if err := ev.pe.Add(p); err != nil {
		if evidence, ok := err.(*DuplicatePrepareVoteEvidence); ok {
			if err := ev.commit(evidence); err != nil {
				return err
			}
			return err
		}
	}
	return nil
}

func encodeKey(e Evidence) []byte {
	buf := bytes.NewBuffer(nil)
	switch e.(type) {
	case *DuplicatePrepareVoteEvidence:
		buf.WriteByte(prepareDualPrefix)
	case *DuplicateViewChangeVoteEvidence:
		buf.WriteByte(viewDualPrefix)
	case *TimestampViewChangeVoteEvidence:
		buf.WriteByte(viewTimestampPrefix)
	}

	num := [8]byte{}
	binary.BigEndian.PutUint64(num[:], e.BlockNumber())
	buf.Write(num[:])
	buf.Write(e.Address().Bytes())
	buf.Write(e.Hash())
	return buf.Bytes()
}
func (ev *EvidencePool) commit(e Evidence) error {
	key := encodeKey(e)
	var buf []byte
	var err error
	ok := false
	if ok, err = ev.db.Has(key, nil); !ok {
		if buf, err = rlp.EncodeToBytes(e); err == nil {
			err = ev.db.Put(key, buf, &opt.WriteOptions{Sync: true})
		}
	}
	return err
}

func (ev *EvidencePool) Clear(timestamp, blockNum uint64) {
	ev.vt.Clear(timestamp)
	ev.vn.Clear(blockNum)
	ev.pe.Clear(blockNum)
}

func (ev *EvidencePool) Close() {
	ev.db.Close()
}

func (ev *EvidencePool) Evidences() []Evidence {
	var evds []Evidence
	it := ev.db.NewIterator(nil, nil)
	for it.Next() {
		flag := it.Key()[0]
		switch flag {
		case prepareDualPrefix:
			var e DuplicatePrepareVoteEvidence
			if err := rlp.DecodeBytes(it.Value(), &e); err == nil {
				evds = append(evds, e)
			}
		case viewDualPrefix:
			var e DuplicateViewChangeVoteEvidence
			if err := rlp.DecodeBytes(it.Value(), &e); err == nil {
				evds = append(evds, e)
			}
		case viewTimestampPrefix:
			var e TimestampViewChangeVoteEvidence
			if err := rlp.DecodeBytes(it.Value(), &e); err == nil {
				evds = append(evds, e)
			}
		}
	}

	it.Release()
	return evds
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
	return nil
}
