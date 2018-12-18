package cbft

import (
	"Platon-go/common"
)

type rotating struct {
	dpos         *dpos
	rotaList     []common.Address // This round of cyclic block node order list
	startTime    int64            // The current cycle start timestamp, in milliseconds
	endTime      int64            // The current cycle end timestamp, in milliseconds
	timeInterval int64            // Block time per unit, in milliseconds
}

func newRotating(dpos *dpos, timeInterval int64) *rotating {
	rotating := &rotating{
		dpos:         dpos,
		timeInterval: timeInterval,
	}
	return rotating
}

func sort() {
	// New round of consensus sorting function
	// xor(Last block last block hash + node public key address)
}

func (r *rotating) IsRotating(common.Address) bool {
	// Determine whether the current node is out of order
	// Sort by consensus and time window
	return false
}

/*func (r *rotating) inturn(number uint64, signer common.Address) bool {
	sort.Sort(signerOrderingRule(r.rotaList))
	offset :=  0
	for offset < len(r.rotaList) && r.rotaList[offset] != signer {
		offset++
	}
	return (number % uint64(len(r.rotaList))) == uint64(offset)
}

type signerOrderingRule []common.Address
func (s signerOrderingRule) Len() int           { return len(s) }
func (s signerOrderingRule) Less(i, j int) bool { return bytes.Compare(s[i][:], s[j][:]) < 0 }
func (s signerOrderingRule) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }*/
