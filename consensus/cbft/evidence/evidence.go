// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

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

// Add tries to add prepare block to PrepareBlockEvidence.
// if the return error is DuplicatePrepareBlockEvidence instructions the prepare is duplicated
func (pbe PrepareBlockEvidence) Add(pb *EvidencePrepare, id Identity) error {
	var l NumberOrderPrepareBlock

	if l = pbe[id]; l == nil {
		l = make(NumberOrderPrepareBlock, 0)
	}
	err := l.Add(pb)
	pbe[id] = l
	return err
}

// Clear tries to clear stale intermediate prepare
func (pbe PrepareBlockEvidence) Clear(epoch uint64, viewNumber uint64) {
	for k := range pbe {
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

// find tries to find the same prepare evidence
// the same epoch,viewNumber,blockNumber and different blockHash
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

// Add tries to add prepare vote to PrepareVoteEvidence.
// if the return error is DuplicatePrepareVoteEvidence instructions the vote is duplicated
func (pve PrepareVoteEvidence) Add(pv *EvidenceVote, id Identity) error {
	var l NumberOrderPrepareVote

	if l = pve[id]; l == nil {
		l = make(NumberOrderPrepareVote, 0)
	}
	err := l.Add(pv)
	pve[id] = l
	return err
}

// Clear tries to clear stale intermediate vote
func (pve PrepareVoteEvidence) Clear(epoch uint64, viewNumber uint64) {
	for k := range pve {
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

// find tries to find the same vote evidence
// the same epoch,viewNumber,blockNumber and different blockHash
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

// Add tries to add view to ViewChangeEvidence.
// if the return error is DuplicateViewChangeEvidence instructions the view is duplicated
func (vce ViewChangeEvidence) Add(vc *EvidenceView, id Identity) error {
	var l NumberOrderViewChange

	if l = vce[id]; l == nil {
		l = make(NumberOrderViewChange, 0)
	}
	err := l.Add(vc)
	vce[id] = l
	return err
}

// Clear tries to clear stale intermediate view
func (vce ViewChangeEvidence) Clear(epoch uint64, viewNumber uint64) {
	for k := range vce {
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
	if ev := ovc.find(vc.Epoch, vc.ViewNumber); ev != nil {
		if ev.BlockNumber != vc.BlockNumber || ev.BlockHash != vc.BlockHash {
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

// find tries to find the same view evidence
// the same epoch,viewNumber,blockNumber and different blockHash
func (ovc NumberOrderViewChange) find(epoch uint64, viewNumber uint64) *EvidenceView {
	for _, v := range ovc {
		if v.Epoch == epoch && v.ViewNumber == viewNumber {
			return v
		}
	}
	return nil
}

//func (ovc NumberOrderViewChange) find(epoch uint64, viewNumber uint64, blockNumber uint64) *EvidenceView {
//	for _, v := range ovc {
//		if v.Epoch == epoch && v.ViewNumber == viewNumber && v.BlockNumber == blockNumber {
//			return v
//		}
//	}
//	return nil
//}

func (ovc NumberOrderViewChange) Len() int {
	return len(ovc)
}

func (ovc NumberOrderViewChange) Less(i, j int) bool {
	return ovc[i].ViewNumber < ovc[j].ViewNumber
}

func (ovc NumberOrderViewChange) Swap(i, j int) {
	ovc[i], ovc[j] = ovc[j], ovc[i]
}
