package p2p

import (
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
)

var (
	enodeArr = RandomEnode(10)
)

func RandomEnode(count int) []*enode.Node {
	enodeList := make([]*enode.Node, 10)
	for i := 0; i <= count; i++ {
		nodeKey := newkey()
		enodeList[i] = enode.NewV4(&nodeKey.PublicKey, nil, 0, 0)
	}
	return enodeList
}
func TestEvent(t *testing.T) {

	event := cbfttypes.ElectNextEpochVerifierEvent{NodeList: enodeArr}
	t.Log("TestMonitorTask", event)
}

func TestNodeIDList(t *testing.T) {
	t.Log(enodeArr)
}

func TestInitNodePing(t *testing.T) {
	InitNodePing(enodeArr)
}

func TestFindNodePing(t *testing.T) {
	var nodePing TbNodePing
	MonitorDB().Find(&nodePing, "node_id=?", "11a315747ce79cdf3d6aaf87ff2b6897950a20bda281838f922ea9407736fec9029d85f6202fd059a57a9119d05895402e7570948ae759cb093a54c3da9e0a4a")
	t.Log(nodePing)
}

func TestListNodePing(t *testing.T) {
	var nodePings []TbNodePing
	MonitorDB().Find(&nodePings, "status=?", 0)
	t.Log(len(nodePings))
}

func TestSaveNodePingResult(t *testing.T) {
	SaveNodePingResult(enodeArr[0], "127.0.0.1:8080", 1)
}

type TbBlock struct {
	Id          string `gorm:"primaryKey"`
	BlockNumber uint64
	BlockHash   string
	Epoch       uint64
	Additional  uint64
	Reward      uint64
	BlockFee    uint64
	TxCount     uint64
	NodeID      string
	ConsensusNo uint64
	CreateTime  int64 `gorm:"autoCreateTime"`
}

func (TbBlock) TableName() string {
	return "Tb_Block" //指定表名。
}
func TestInsertBlock(t *testing.T) {
	//b := types.NewSimplifiedBlock(124, common.HexToHash("499987a73fa100f582328c92c1239262edf5c0a3479face652c89f60314aa805"))
	tbBlock := TbBlock{BlockNumber: uint64(123), BlockHash: "499987a73fa100f582328c92c1239262edf5c0a3479face652c89f60314aa805", Epoch: 2, ConsensusNo: 31323, NodeID: "11a315747ce79cdf3d6aaf87ff2b6897950a20bda281838f922ea9407736fec9029d85f6202fd059a57a9119d05895402e7570948ae759cb093a54c3da9e0a4a"}
	MonitorDB().Create(&tbBlock)
}

func InsertConsensusElect(t *testing.T) {
}
