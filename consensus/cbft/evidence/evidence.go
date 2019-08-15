package evidence

import (
	"math/big"
	"sort"
	"strconv"
	"strings"
)

const (
	// IdentityLength is the expected length of the identity
	IdentityLength = 20
)

type NumberOrderPrepareBlock []*EvidencePrepare
type NumberOrderPrepareVote []*EvidenceVote
type NumberOrderViewChange []*EvidenceView

type Identity string

type PrepareBlockEvidence map[Identity]NumberOrderPrepareBlock
type PrepareVoteEvidence map[Identity]NumberOrderPrepareVote
type ViewChangeEvidence map[Identity]NumberOrderViewChange

// Bytes gets the string representation of the underlying identity.
func (id Identity) Bytes() []byte { return []byte(id) }

func (e PrepareBlockEvidence) Add(pb *EvidencePrepare, id Identity) error {
	var l NumberOrderPrepareBlock

	if l = e[id]; l == nil {
		l = make(NumberOrderPrepareBlock, 0)
	}
	err := l.Add(pb)
	e[id] = l
	return err
}

func (pbe PrepareBlockEvidence) Clear(epoch uint64, viewNumber uint64) {
	for k, _ := range pbe {
		s := strings.Split(string(k), "|")
		e, _ := strconv.ParseUint(s[0], 10, 64)
		v, _ := strconv.ParseUint(s[1], 10, 64)

		if e < epoch || e == epoch && v <= viewNumber {
			delete(pbe, k)
		}
	}
}

func (pbe PrepareBlockEvidence) Size() int {
	size := 0
	for _, v := range pbe {
		size = size + v.Len()
	}
	return size
}

func (opb *NumberOrderPrepareBlock) Add(pb *EvidencePrepare) error {
	if ev := opb.find(pb.Epoch, pb.ViewNumber, pb.BlockNumber); ev != nil {
		if ev.BlockHash != pb.BlockHash {
			a, b := pb, ev
			ha := new(big.Int).SetBytes(pb.BlockHash.Bytes())
			hb := new(big.Int).SetBytes(ev.BlockHash.Bytes())

			if ha.Cmp(hb) > 0 {
				a, b = ev, pb
			}
			return &DuplicatePrepareBlockEvidence{
				PrepareA: a,
				PrepareB: b,
			}
		}
	} else {
		*opb = append(*opb, pb)
		sort.Sort(*opb)
	}
	return nil
}

func (opb NumberOrderPrepareBlock) find(epoch uint64, viewNumber uint64, blockNumber uint64) *EvidencePrepare {
	for _, v := range opb {
		if v.Epoch == epoch && v.ViewNumber == viewNumber && v.BlockNumber == blockNumber {
			return v
		}
	}
	return nil
}

func (opb NumberOrderPrepareBlock) Len() int {
	return len(opb)
}

func (opb NumberOrderPrepareBlock) Less(i, j int) bool {
	return opb[i].ViewNumber < opb[j].ViewNumber
}

func (opb NumberOrderPrepareBlock) Swap(i, j int) {
	opb[i], opb[j] = opb[j], opb[i]
}

func (e PrepareVoteEvidence) Add(pv *EvidenceVote, id Identity) error {
	var l NumberOrderPrepareVote

	if l = e[id]; l == nil {
		l = make(NumberOrderPrepareVote, 0)
	}
	err := l.Add(pv)
	e[id] = l
	return err
}

func (pve PrepareVoteEvidence) Clear(epoch uint64, viewNumber uint64) {
	for k, _ := range pve {
		s := strings.Split(string(k), "|")
		e, _ := strconv.ParseUint(s[0], 10, 64)
		v, _ := strconv.ParseUint(s[1], 10, 64)

		if e < epoch || e == epoch && v <= viewNumber {
			delete(pve, k)
		}
	}
}

func (pve PrepareVoteEvidence) Size() int {
	size := 0
	for _, v := range pve {
		size = size + v.Len()
	}
	return size
}

func (opv *NumberOrderPrepareVote) Add(pv *EvidenceVote) error {
	if ev := opv.find(pv.Epoch, pv.ViewNumber, pv.BlockNumber); ev != nil {
		if ev.BlockHash != pv.BlockHash {
			a, b := pv, ev
			ha := new(big.Int).SetBytes(pv.BlockHash.Bytes())
			hb := new(big.Int).SetBytes(ev.BlockHash.Bytes())

			if ha.Cmp(hb) > 0 {
				a, b = ev, pv
			}
			return &DuplicatePrepareVoteEvidence{
				VoteA: a,
				VoteB: b,
			}
		}
	} else {
		*opv = append(*opv, pv)
		sort.Sort(*opv)
	}
	return nil
}

func (opv NumberOrderPrepareVote) find(epoch uint64, viewNumber uint64, blockNumber uint64) *EvidenceVote {
	for _, v := range opv {
		if v.Epoch == epoch && v.ViewNumber == viewNumber && v.BlockNumber == blockNumber {
			return v
		}
	}
	return nil
}

func (opv NumberOrderPrepareVote) Len() int {
	return len(opv)
}

func (opv NumberOrderPrepareVote) Less(i, j int) bool {
	return opv[i].ViewNumber < opv[j].ViewNumber
}

func (opv NumberOrderPrepareVote) Swap(i, j int) {
	opv[i], opv[j] = opv[j], opv[i]
}

func (e ViewChangeEvidence) Add(vc *EvidenceView, id Identity) error {
	var l NumberOrderViewChange

	if l = e[id]; l == nil {
		l = make(NumberOrderViewChange, 0)
	}
	err := l.Add(vc)
	e[id] = l
	return err
}

func (vce ViewChangeEvidence) Clear(epoch uint64, viewNumber uint64) {
	for k, _ := range vce {
		s := strings.Split(string(k), "|")
		e, _ := strconv.ParseUint(s[0], 10, 64)
		v, _ := strconv.ParseUint(s[1], 10, 64)

		if e < epoch || e == epoch && v <= viewNumber {
			delete(vce, k)
		}
	}
}

func (vce ViewChangeEvidence) Size() int {
	size := 0
	for _, v := range vce {
		size = size + v.Len()
	}
	return size
}

func (ovc *NumberOrderViewChange) Add(vc *EvidenceView) error {
	if ev := ovc.find(vc.Epoch, vc.ViewNumber, vc.BlockNumber); ev != nil {
		if ev.BlockHash != vc.BlockHash {
			a, b := vc, ev
			ha := new(big.Int).SetBytes(vc.BlockHash.Bytes())
			hb := new(big.Int).SetBytes(ev.BlockHash.Bytes())

			if ha.Cmp(hb) > 0 {
				a, b = ev, vc
			}
			return &DuplicateViewChangeEvidence{
				ViewA: a,
				ViewB: b,
			}
		}
	} else {
		*ovc = append(*ovc, vc)
		sort.Sort(*ovc)
	}
	return nil
}

func (ovc NumberOrderViewChange) find(epoch uint64, viewNumber uint64, blockNumber uint64) *EvidenceView {
	for _, v := range ovc {
		if v.Epoch == epoch && v.ViewNumber == viewNumber && v.BlockNumber == blockNumber {
			return v
		}
	}
	return nil
}

func (ovc NumberOrderViewChange) Len() int {
	return len(ovc)
}

func (ovc NumberOrderViewChange) Less(i, j int) bool {
	return ovc[i].ViewNumber < ovc[j].ViewNumber
}

func (ovc NumberOrderViewChange) Swap(i, j int) {
	ovc[i], ovc[j] = ovc[j], ovc[i]
}
