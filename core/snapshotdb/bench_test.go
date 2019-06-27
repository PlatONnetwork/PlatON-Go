package snapshotdb

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/big"
	"math/rand"
	"os"
	"runtime"
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

func openDBBench(b *testing.B) *dbBench {
	_, err := os.Stat(dbpath)
	if err == nil {
		err = os.RemoveAll(dbpath)
		if err != nil {
			b.Fatal("cannot remove old db: ", err)
		}
	}
	if err := initDB(); err != nil {
		b.Fatal("init db fail", err)
	}
	return &dbBench{
		b:  b,
		db: dbInstance,
	}
}

func (p *dbBench) populate(n int) {
	p.keys, p.values = make([][]byte, n), make([][]byte, n)
	v := newValueGen(0.5)
	for i := range p.keys {
		p.keys[i], p.values[i] = []byte(fmt.Sprintf("%016d", i)), v.get(100)
	}
}

func (p *dbBench) puts() {
	b := p.b
	db := p.db
	err := p.db.NewBlock(big.NewInt(1), rlpHash("a"), common.ZeroHash)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.StartTimer()
	for i := range p.keys {
		err := db.Put(common.ZeroHash, p.keys[i], p.values[i])
		if err != nil {
			b.Fatal("put failed: ", err)
		}
	}
	b.StopTimer()
	b.SetBytes(116)
}

func (p *dbBench) close() {
	p.db.Clear()
	p.db = nil
	p.keys = nil
	p.values = nil
	runtime.GC()
}

func (p *dbBench) fill() {
	b := p.b
	db := p.db
	err := db.NewBlock(big.NewInt(1), rlpHash("a"), common.ZeroHash)
	if err != nil {
		b.Fatal(err)
	}
	perBatch := 10000
	for i, n := 0, len(p.keys); i < n; {
		first := true
		for ; i < n && ((i+1)%perBatch != 0 || first); i++ {
			first = false
			err := db.Put(common.ZeroHash, p.keys[i], p.values[i])
			if err != nil {
				b.Fatal("write failed: ", err)
			}
		}
	}
}

func BenchmarkDBPut(b *testing.B) {
	p := openDBBench(b)
	p.populate(b.N)
	p.puts()
	p.close()
}

func BenchmarkDBGet(b *testing.B) {
	p := openDBBench(b)
	p.populate(b.N)
	p.fill()
	//p.gets()
	p.close()
}
