package evidence

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	prepareDualPrefix = byte(0x1)
	voteDualPrefix    = byte(0x2)
	viewDualPrefix    = byte(0x3)
)

type EvidencePool interface {
	consensus.EvidencePool
	AddPrepareBlock(pb *protocols.PrepareBlock) error
	AddPrepareVote(pv *protocols.PrepareVote) error
	AddViewChange(vc *protocols.ViewChange) error
}

type emptyEvidencePool struct {
}

func (pool emptyEvidencePool) AddPrepareBlock(pb *protocols.PrepareBlock) error {
	return nil
}

func (pool emptyEvidencePool) AddPrepareVote(pv *protocols.PrepareVote) error {
	return nil
}

func (pool emptyEvidencePool) AddViewChange(vc *protocols.ViewChange) error {
	return nil
}

func (pool emptyEvidencePool) Evidences() consensus.Evidences {
	return nil
}

func (pool emptyEvidencePool) UnmarshalEvidence(data []byte) (consensus.Evidences, error) {
	return nil, nil
}

func (pool emptyEvidencePool) Clear(epoch uint64, viewNumber uint64) {
}

func (pool emptyEvidencePool) Close() {
}

type baseEvidencePool struct {
	pb PrepareBlockEvidence
	pv PrepareVoteEvidence
	vc ViewChangeEvidence
	db *leveldb.DB
}

func NewEvidencePool(ctx *node.ServiceContext, evidenceDir string) (EvidencePool, error) {
	path := ""
	if ctx != nil {
		path = ctx.ResolvePath(evidenceDir)
	}
	if len(path) == 0 {
		return &emptyEvidencePool{}, nil
	}
	return NewBaseEvidencePool(path)
}

func NewBaseEvidencePool(path string) (*baseEvidencePool, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &baseEvidencePool{
		pb: make(PrepareBlockEvidence),
		pv: make(PrepareVoteEvidence),
		vc: make(ViewChangeEvidence),
		db: db,
	}, nil
}

func (pool baseEvidencePool) AddPrepareBlock(pb *protocols.PrepareBlock) (err error) {
	id := verify(pb)
	if err := pool.pb.Add(pb, id); err != nil {
		if evidence, ok := err.(*DuplicatePrepareBlockEvidence); ok {
			if err := pool.commit(evidence, id); err != nil {
				return err
			}
			return err
		}
	}
	return nil
}

func (pool baseEvidencePool) AddPrepareVote(pv *protocols.PrepareVote) (err error) {
	id := verify(pv)
	if err := pool.pv.Add(pv, id); err != nil {
		if evidence, ok := err.(*DuplicatePrepareVoteEvidence); ok {
			if err := pool.commit(evidence, id); err != nil {
				return err
			}
			return err
		}
	}
	return nil
}

func (pool baseEvidencePool) AddViewChange(vc *protocols.ViewChange) (err error) {
	id := verify(vc)
	if err := pool.vc.Add(vc, id); err != nil {
		if evidence, ok := err.(*DuplicateViewChangeEvidence); ok {
			if err := pool.commit(evidence, id); err != nil {
				return err
			}
			return err
		}
	}
	return nil
}

func (pool baseEvidencePool) Evidences() consensus.Evidences {
	var evds consensus.Evidences
	it := pool.db.NewIterator(nil, nil)
	for it.Next() {
		flag := it.Key()[0]
		switch flag {
		case prepareDualPrefix:
			var e DuplicatePrepareBlockEvidence
			if err := rlp.DecodeBytes(it.Value(), &e); err == nil {
				evds = append(evds, e)
			}
		case voteDualPrefix:
			var e DuplicatePrepareVoteEvidence
			if err := rlp.DecodeBytes(it.Value(), &e); err == nil {
				evds = append(evds, e)
			}
		case viewDualPrefix:
			var e DuplicateViewChangeEvidence
			if err := rlp.DecodeBytes(it.Value(), &e); err == nil {
				evds = append(evds, e)
			}
		}
	}

	it.Release()
	return evds
}

func (pool baseEvidencePool) UnmarshalEvidence(data []byte) (consensus.Evidences, error) {
	var ed EvidenceData
	if err := json.Unmarshal(data, &ed); err != nil {
		return nil, err
	}
	evds := make(consensus.Evidences, 0)
	for _, e := range ed.DP {
		evds = append(evds, e)
	}
	for _, e := range ed.DV {
		evds = append(evds, e)
	}
	for _, e := range ed.DC {
		evds = append(evds, e)
	}
	return evds, nil
}

func (pool baseEvidencePool) Clear(epoch uint64, viewNumber uint64) {
	pool.pb.Clear(epoch, viewNumber)
	pool.pv.Clear(epoch, viewNumber)
	pool.vc.Clear(epoch, viewNumber)
}

func (pool baseEvidencePool) Close() {
	pool.db.Close()
}

func verify(msg types.ConsensusMsg) Identity {
	return Identity(msg.IdentityMsg())
}

func encodeKey(e consensus.Evidence, id Identity) []byte {
	buf := bytes.NewBuffer(nil)
	switch e.(type) {
	case *DuplicatePrepareBlockEvidence:
		buf.WriteByte(prepareDualPrefix)
	case *DuplicatePrepareVoteEvidence:
		buf.WriteByte(voteDualPrefix)
	case *DuplicateViewChangeEvidence:
		buf.WriteByte(viewDualPrefix)
	}

	// epoch
	epoch := [8]byte{}
	binary.BigEndian.PutUint64(epoch[:], e.Epoch())
	buf.Write(epoch[:])
	// viewNumber
	viewNum := [8]byte{}
	binary.BigEndian.PutUint64(viewNum[:], e.ViewNumber())
	buf.Write(viewNum[:])
	// blockNumber
	num := [8]byte{}
	binary.BigEndian.PutUint64(num[:], e.BlockNumber())
	buf.Write(num[:])
	// node identity
	buf.Write(id.Bytes())
	// Evidence hash
	buf.Write(e.Hash())
	return buf.Bytes()
}

func (ev *baseEvidencePool) commit(e consensus.Evidence, id Identity) error {
	key := encodeKey(e, id)
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
