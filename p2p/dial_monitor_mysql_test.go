package p2p

import (
	"net"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	nodeIdArr = []discover.NodeID{
		discover.MustHexID("0x362003c50ed3a523cdede37a001803b8f0fed27cb402b3d6127a1a96661ec202318f68f4c76d9b0bfbabfd551a178d4335eaeaa9b7981a4df30dfc8c0bfe3384"),
		discover.MustHexID("0xced880d4769331f47af07a8d1b79de1e40c95a37ea1890bb9d3f0da8349e1a7c0ea4cadbb9c5bf185b051061eef8e5eadca251c24e1db1d9faf0fb24cbd06f9a"),
		discover.MustHexID("0xda56501a77fc1dfe0399b81f3909061d9a176cb9433fab4d3dfb1a10344c243274e38155e18878c7a0b3fcdd6182000c7784a95e2c4d9e0691ce67798624786e"),
		discover.MustHexID("0x89a4409abe1ace8b77c4497c2073a8a2046dbdabb58c8bb58fe73926bbdc572fb848d739b1d2d09dd0796abcc1ed8d9a33bb3ef0a6c2e106e408090df179b041"),
		discover.MustHexID("0x65e2ab09161e32e6d07d82adaa416ee6d41d617c52db20e3145a4d1b7d396af38d095c87508ad5bb35df741513bdc4bf12fec215e58450e255f05d194d41d089"),
		discover.MustHexID("0x9bfacd628f3adb0f94e8b3968064d5248fa18efa75c680fdffea3af2575406461f3395817dd2a1be07a79bd81ffa00f57ad82286061d4a6caceece048e352380"),
		discover.MustHexID("0x1e07d66b56bbc931ddce7cc5b9f55672d7fe4e19897a42f19d4ad7c969435cad652d720401d68f5769e245ec0f4e23362c8b1b062771d614876fdbb875ba9d44"),
		discover.MustHexID("0x11a315747ce79cdf3d6aaf87ff2b6897950a20bda281838f922ea9407736fec9029d85f6202fd059a57a9119d05895402e7570948ae759cb093a54c3da9e0a4a"),
		discover.MustHexID("0x248af08a775ff63a47a5970e4928bcccd1a8cef984fd4142ea7f89cd13015bdab9ca4a8c5e1070dc00fa81a047542f53ca596f553c4acfb7abe75a8fb5019057"),
		discover.MustHexID("0xfd790ff5dc48baccb9418ce5cfac6a10c3646f20a3fe32d9502c4edce3a77fa90bfee0361d8a72093b7994f8cbc28ee537bdda2b634c5966b1a9253d9d270145"),
		discover.MustHexID("0x56d243db84a521cb204f582ee84bca7f4af29437dd447a6e36d17f4853888e05343844bd64294b99b835ca7f72ef5b1325ef1c89b0c5c2744154cdadf7c4e9fa"),
		discover.MustHexID("0x8796a6fcefd9037d8433e3a959ff8f3c4552a482ce727b00a90bfd1ec365ce2faa33e19aa6a172b5c186b51f5a875b5acd35063171f0d9501a9c8f1c98513825"),
		discover.MustHexID("0x547b876036165d66274ce31692165c8acb6f140a65cab0e0e12f1f09d1c7d8d53decf997830919e4f5cacb2df1adfe914c53d22e3ab284730b78f5c63a273b8c"),
		discover.MustHexID("0x9fdbeb873bea2557752eabd2c96419b8a700b680716081472601ddf7498f0db9b8a40797b677f2fac541031f742c2bbd110ff264ae3400bf177c456a76a93d42"),
		discover.MustHexID("0xc553783799bfef7c34a84b2737f2c77f8f2c5cfedc3fd7af2d944da6ece90aa94cf621e6de5c4495881fbfc9beec655ffb10e39cb4ca9be7768d284409040f32"),
		discover.MustHexID("0x75ad2ee8ca77619c3ba0ddcec5dab1375fe4fa90bab9e751caef3996ce082dfed32fe4c137401ee05e501c079b2e4400397b09de14b08b09c9e7f9698e9e4f0a"),
		discover.MustHexID("0xdb18af9be2af9dff2347c3d06db4b1bada0598d099a210275251b68fa7b5a863d47fcdd382cc4b3ea01e5b55e9dd0bdbce654133b7f58928ce74629d5e68b974"),
		discover.MustHexID("0x472d19e5e9888368c02f24ebbbe0f2132096e7183d213ab65d96b8c03205f88398924af8876f3c615e08aa0f9a26c38911fda26d51c602c8d4f8f3cb866808d7"),
		discover.MustHexID("4f1f036e5e18cc812347d5073cbec2a8da7930de323063c39b0d4413a396e088bfa90e8c28174313d8d82e9a14bc0884b13a48fc28e619e44c48a49b4fd9f107"),
		discover.MustHexID("f18c596232d637409c6295abb1e720db99ffc12363a1eb8123d6f54af80423a5edd06f91115115a1dca1377e97b9031e2ddb864d34d9b3491d6fa07e8d9b951b"),
		discover.MustHexID("7a8f7a28ac1c4eaf98b2be890f372e5abc58ebe6d3aab47aedcb0076e34eb42882e926676ebab327a4ef4e2ea5c4296e9c7bc0991360cb44f52672631012db1b"),
		discover.MustHexID("9eeb448babf9e93449e831b91f98d9cbc0c2324fe8c43baac69d090717454f3f930713084713fe3a9f01e4ca59b80a0f2b41dbd6d531f414650bab0363e3691a"),
		discover.MustHexID("cc1d7314c15e30dc5587f675eb5f803b1a2d88bfe76cec591cec1ff678bc6abce98f40054325bdcb44fb83174f27d38a54fbce4846af8f027b333868bc5144a4"),
		discover.MustHexID("e4d99694be2fc8a53d8c2446f947aec1c7de3ee26f7cd43f4f6f77371f56f11156218dec32b51ddce470e97127624d330bb7a3237ba5f0d87d2d3166faf1035e"),
		discover.MustHexID("9c61f59f70296b6d494e7230888e58f19b13c5c6c85562e57e1fe02d0ff872b4957238c73559d017c8770b999891056aa6329dbf628bc19028d8f4d35ec35823"),
	}
)

func TestEvent(t *testing.T) {
	event := cbfttypes.ElectNextEpochVerifierEvent{NodeIdList: nodeIdArr}
	t.Log("TestMonitorTask", event)
}

func TestMonitorTask(t *testing.T) {
	ip := net.IPv4(123, 1, 2, 3)
	node := &discover.Node{ID: randomID(), IP: ip, TCP: 8080}
	task := &monitorTask{flags: monitorConn, dest: node}
	t.Log("TestMonitorTask", task)
}

func TestNodeIDList(t *testing.T) {
	t.Log(nodeIdArr)
}

func TestNodeIDToString(t *testing.T) {
	t.Log(nodeIdArr[0].String())
	t.Log(hexutil.Encode(nodeIdArr[0].Bytes()))
	t.Log(nodeIdArr[0].HexPrefixString())
}

func TestSaveEpochElection(t *testing.T) {
	SaveEpochElection(1, nodeIdArr)
}

func TestSaveConsensusElection(t *testing.T) {
	SaveConsensusElection(31323, nodeIdArr)
}

func TestInitNodePing(t *testing.T) {
	InitNodePing(nodeIdArr)
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
	SaveNodePingResult(nodeIdArr[0], "127.0.0.1:8080", 1)
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
