package cbft

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/deckarep/golang-set"
	"sort"
	"strings"
	"time"
)

const (
	maxReceiveRecords = 80
	maxSendRecords    = 80
)

type receiveRecord struct {
	Time    time.Time `json:"time"`
	Order   int64     `json:"order"`
	SelfId  string    `json:"self_id"`
	FromId  string    `json:"from_id"`
	MsgHash string    `json:"msg_hash"`
	Type    string    `json:"type"`
}

type sendRecord struct {
	Time      time.Time `json:"time"`
	Order     int64     `json:"order"`
	SelfId    string    `json:"self_id"`
	TargetIds []string  `json:"target_ids"`
	MsgHash   string    `json:"msg_hash"`
	Type      string    `json:"type"`
}

func (rr *receiveRecord) ToJSON() ([]byte, error) {
	return json.Marshal(rr)
}

func (rr *receiveRecord) ParseFromJSON(input []byte) error {
	var r receiveRecord
	if err := json.Unmarshal(input, &r); err != nil {
		return err
	}
	return nil
}

func (rr *receiveRecord) String() string {
	j, err := rr.ToJSON()
	if err != nil {
		return ""
	}
	return string(j)
}

func (sr *sendRecord) ToJSON() ([]byte, error) {
	return json.Marshal(sr)
}

func (sr *sendRecord) ParseFromJSON(input []byte) error {
	var s sendRecord
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	return nil
}

func (rr *sendRecord) String() string {
	j, err := rr.ToJSON()
	if err != nil {
		return ""
	}
	return string(j)
}

type tracing struct {
	receiveRQueue mapset.Set
	sendRQueue    mapset.Set
	isRecord      bool
	quite         chan struct{}
}

func NewTracing() *tracing {
	t := &tracing{
		receiveRQueue: mapset.NewSet(),
		sendRQueue:    mapset.NewSet(),
		isRecord:      false,
		quite:         make(chan struct{}, 1),
	}
	return t
}

func (t *tracing) On() {
	t.isRecord = true
	go t.start()
}

func (t *tracing) Off() {
	t.isRecord = false
	t.receiveRQueue.Clear()
	t.sendRQueue.Clear()
	go func() {
		t.quite <- struct{}{}
	}()
}

func (t *tracing) RecordReceive(selfId, fromId, msgHash, msgType string) {
	if !t.isRecord {
		return
	}
	record := new(receiveRecord)
	record.Time = time.Now()
	record.MsgHash = msgHash
	record.Order = time.Now().Unix()
	record.SelfId = selfId
	record.FromId = fromId
	record.Type = msgType
	t.recordReceiveAction(record)
}

func (t *tracing) recordReceiveAction(record *receiveRecord) {
	for t.receiveRQueue.Cardinality() >= maxReceiveRecords {
		t.receiveRQueue.Pop()
	}
	t.receiveRQueue.Add(record)
}

func (t *tracing) RecordSend(selfId, msgHash, msgType string, targetIds string) {
	if !t.isRecord {
		return
	}
	record := new(sendRecord)
	record.Time = time.Now()
	record.MsgHash = msgHash
	record.Order = time.Now().Unix()
	record.SelfId = selfId
	targetIdArr := strings.Split(targetIds, ",")
	record.TargetIds = targetIdArr
	record.Type = msgType
	t.recordSendAction(record)
}

func (t *tracing) recordSendAction(record *sendRecord) {
	if !t.isRecord {
		return
	}
	for t.sendRQueue.Cardinality() >= maxSendRecords {
		t.sendRQueue.Pop()
	}
	t.sendRQueue.Add(record)
}

func (t *tracing) String() string {
	// receive msg / send msg / 2s/per.

	// receive
	receives := t.receiveRQueue.Clone()
	var receiveHeap receiveTimeHeap
	heap.Init(&receiveHeap)
	for v := range receives.Iter() {
		heap.Push(&receiveHeap, v)
	}
	sort.Sort(receiveHeap)

	// send
	sends := t.sendRQueue.Clone()
	var sendHeap sendTimeHeap
	heap.Init(&sendHeap)
	for v := range sends.Iter() {
		heap.Push(&sendHeap, v)
	}
	sort.Sort(sendHeap)
	return fmt.Sprintf("{\"receive\": %v, \"send\": %v }", receiveHeap, sendHeap)
}

func (t *tracing) start() {
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			log.Debug("tracing data", "data", t)
		case <-t.quite:
			log.Warn("tracing stop")
			return
		}
	}
}

type receiveTimeHeap []*receiveRecord

func (h receiveTimeHeap) Len() int      { return len(h) }
func (h receiveTimeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h receiveTimeHeap) Less(i, j int) bool {
	return h[i].Order < h[j].Order
}

func (h *receiveTimeHeap) Push(x interface{}) {
	*h = append(*h, x.(*receiveRecord))
}

func (h *receiveTimeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h receiveTimeHeap) String() string {
	j, err := json.Marshal(&h)
	if err != nil {
		return ""
	}
	return string(j)
}

type sendTimeHeap []*sendRecord

func (h sendTimeHeap) Len() int      { return len(h) }
func (h sendTimeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h sendTimeHeap) Less(i, j int) bool {
	return h[i].Order < h[j].Order
}

func (h *sendTimeHeap) Push(x interface{}) {
	*h = append(*h, x.(*sendRecord))
}

func (h *sendTimeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h sendTimeHeap) String() string {
	j, err := json.Marshal(&h)
	if err != nil {
		return ""
	}
	return string(j)
}
