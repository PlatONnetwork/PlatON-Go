package snapshotdb

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
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
	b                                    *testing.B
	db                                   *snapshotDB
	baseDBkeys, baseDBvalues             [][]byte
	recognizedkeys, recognizedvalues     [][]byte
	unrecognizedkeys, unrecognizedvalues [][]byte
	committedkeys, committedvalues       [][]byte
	hashs                                []common.Hash
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
	p.baseDBkeys, p.baseDBvalues = make([][]byte, n), make([][]byte, n)
	p.recognizedkeys, p.recognizedvalues = make([][]byte, n), make([][]byte, n)
	p.unrecognizedkeys, p.unrecognizedvalues = make([][]byte, n), make([][]byte, n)
	p.committedkeys, p.committedvalues = make([][]byte, n), make([][]byte, n)
	p.hashs = make([]common.Hash, n)
	v := newValueGen(0.5)
	for i := range p.baseDBkeys {
		p.baseDBkeys[i], p.baseDBvalues[i] = []byte(fmt.Sprintf("%016d", i)), v.get(100)
	}

	v2 := newValueGen(0.6)
	for i := range p.baseDBkeys {
		p.recognizedkeys[i], p.recognizedvalues[i] = []byte(fmt.Sprintf("%016d", i)), v2.get(100)
	}

	v3 := newValueGen(0.7)
	for i := range p.baseDBkeys {
		p.unrecognizedkeys[i], p.unrecognizedvalues[i] = []byte(fmt.Sprintf("%016d", i)), v3.get(100)
	}

	v4 := newValueGen(0.8)
	for i := range p.baseDBkeys {
		p.committedkeys[i], p.committedvalues[i] = []byte(fmt.Sprintf("%016d", i)), v4.get(100)
	}
	for i := range p.baseDBkeys {
		p.hashs[i] = generateHash(fmt.Sprint(i))
	}
}

func (p *dbBench) putsUnrecognized() {
	b := p.b
	db := p.db
	recognizedHash := generateHash("recognizedHash")
	parentHash := generateHash("parentHash")
	num := big.NewInt(1)
	b.ResetTimer()
	b.StartTimer()
	if err := p.db.NewBlock(num, parentHash, common.ZeroHash); err != nil {
		b.Fatal(err)
	}
	for i := range p.unrecognizedkeys {
		err := db.Put(common.ZeroHash, p.unrecognizedkeys[i], p.unrecognizedvalues[i])
		if err != nil {
			b.Fatal("put failed: ", err)
		}
	}
	if err := db.Flush(recognizedHash, num); err != nil {
		b.Fatal("put failed: ", err)
	}
	if err := db.Commit(recognizedHash); err != nil {
		b.Fatal(err)
	}
	if err := db.Compaction(); err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
	b.SetBytes(116)
}

func (p *dbBench) putsRecognized() {
	b := p.b
	db := p.db
	recognizedHash := generateHash("recognizedHash")
	parentHash := generateHash("parentHash")
	num := big.NewInt(1)
	b.ResetTimer()
	b.StartTimer()
	if err := p.db.NewBlock(num, parentHash, recognizedHash); err != nil {
		b.Fatal(err)
	}
	for i := range p.recognizedkeys {
		db.GetLastKVHash(recognizedHash)
		err := db.Put(recognizedHash, p.recognizedkeys[i], p.recognizedvalues[i])
		if err != nil {
			b.Fatal("put failed: ", err)
		}
	}
	if err := db.Commit(recognizedHash); err != nil {
		b.Fatal(err)
	}
	if err := db.Compaction(); err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
	b.SetBytes(116)
}

func (p *dbBench) close() {
	p.db.Clear()
	p.db = nil
	p.unrecognizedvalues = nil
	p.unrecognizedkeys = nil
	p.recognizedkeys = nil
	p.recognizedvalues = nil
	p.committedkeys = nil
	p.committedvalues = nil
	p.baseDBkeys = nil
	p.baseDBvalues = nil
	runtime.GC()
}

func (p *dbBench) fill() {

	//commit 100
	//recongized 100
	//unreconized 1
}

func BenchmarkDBPutUnrecognized(b *testing.B) {
	logger.SetHandler(log.DiscardHandler())
	p := openDBBench(b)
	p.populate(b.N)
	p.putsUnrecognized()
	p.close()
}

func BenchmarkDBPutRecognized(b *testing.B) {
	logger.SetHandler(log.DiscardHandler())
	p := openDBBench(b)
	p.populate(b.N)
	p.putsRecognized()
	p.close()
}
