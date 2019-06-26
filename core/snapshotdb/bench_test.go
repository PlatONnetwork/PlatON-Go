package snapshotdb

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

type valueGen struct {
	src []byte
	pos int
}

func (v *valueGen) get(n int) []byte {
	if v.pos+n > len(v.src) {
		v.pos = 0
	}
	v.pos += n
	return v.src[v.pos-n : v.pos]
}

func newValueGen(frac float32) *valueGen {
	v := new(valueGen)
	r := rand.New(rand.NewSource(301))
	v.src = make([]byte, 0, 1048576+100)
	for len(v.src) < 1048576 {
		v.src = append(v.src, compressibleStr(r, frac, 100)...)
	}
	return v
}

func compressibleStr(r *rand.Rand, frac float32, n int) []byte {
	nn := int(float32(n) * frac)
	rb := randomString(r, nn)
	b := make([]byte, 0, n+nn)
	for len(b) < n {
		b = append(b, rb...)
	}
	return b[:n]
}

func randomString(r *rand.Rand, n int) []byte {
	b := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		b.WriteByte(' ' + byte(r.Intn(95)))
	}
	return b.Bytes()
}

type dbBench struct {
	b            *testing.B
	db           *snapshotDB
	keys, values [][]byte
}

//
//func openDBBench(b *testing.B) *dbBench {
//	_, err := os.Stat(dbpath)
//	if err == nil {
//		err = os.RemoveAll(dbpath)
//		if err != nil {
//			b.Fatal("cannot remove old db: ", err)
//		}
//	}
//	if err := initDB(); err != nil {
//		b.Fatal("init db fail", err)
//	}
//	p := &dbBench{
//		b:  b,
//		db: dbInstance,
//	}
//}

func (p *dbBench) populate(n int) {
	p.keys, p.values = make([][]byte, n), make([][]byte, n)
	v := newValueGen(0.5)
	for i := range p.keys {
		p.keys[i], p.values[i] = []byte(fmt.Sprintf("%016d", i)), v.get(100)
	}
}

//
//func (p *dbBench) writes(perBatch int) {
//	b := p.b
//	db := p.db
//
//	n := len(p.keys)
//	m := n / perBatch
//	if n%perBatch > 0 {
//		m++
//	}
//	//batches := make([]Batch, m)
//	//j := 0
//	//for i := range batches {
//	//	first := true
//	//	for ; j < n && ((j+1)%perBatch != 0 || first); j++ {
//	//		first = false
//	//		batches[i].Put(p.keys[j], p.values[j])
//	//	}cd
//	//}
//	//runtime.GC()
//	//
//	//b.ResetTimer()
//	//b.StartTimer()
//	//for i := range batches {
//	//	err := db.Write(&(batches[i]), p.wo)
//	//	if err != nil {
//	//		b.Fatal("write failed: ", err)
//	//	}
//	//}
//	//b.StopTimer()
//	//b.SetBytes(116)
//}
