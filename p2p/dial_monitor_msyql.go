package p2p

import (
	"strings"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	monitorDbOnce sync.Once
	db            *gorm.DB
)

type TbEpoch struct {
	Epoch      uint64 `gorm:"primaryKey"` //从1开始
	NodeId     string `gorm:"primaryKey"`
	CreateTime int64  `gorm:"autoCreateTime"`
	UpdateTime int64  `gorm:"autoUpdateTime"`
}

func (TbEpoch) TableName() string {
	return "Tb_Epoch" //指定表名。缺省表名是tb_epoches(复数形式的）
}

type TbConsensus struct {
	ConsensusNo  uint64 `gorm:"primaryKey"` //从1开始
	NodeId       string `gorm:"primaryKey"`
	StatBlockQty uint64
	CreateTime   int64 `gorm:"autoCreateTime"`
	UpdateTime   int64 `gorm:"autoUpdateTime"`
}

func (TbConsensus) TableName() string {
	return "Tb_Consensus" //指定表名。
}

type TbNodePing struct {
	NodeId     string `gorm:"primaryKey"`
	Status     int8
	ReplyTime  int64
	ReplyBlock uint64
	Ip         string
	Port       string
	CreateTime int64 `gorm:"autoCreateTime"`
	UpdateTime int64 `gorm:"autoUpdateTime"`
}

func (TbNodePing) TableName() string {
	return "Tb_Node_Ping" //指定表名。
}

/*func MonitorDb() *gorm.DB {
	monitorDbOnce.Do(func() {
		sqlDb, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
		sqlDb.SetMaxOpenConns(8)
		sqlDb.SetMaxIdleConns(4)
		gormDB, _ = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDb,
		}), &gorm.Config{})
	})
	return gormDB
}*/

func MonitorDB() *gorm.DB {
	monitorDbOnce.Do(func() {
		// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
		datasource := "uatdevelop:uat@develop@tcp(al-sz-polardb-uat.rwlb.rds.aliyuncs.com:3306)/sun-test-plt-tracking"
		//datasource := "user:pass@tcp(127.0.0.1:3306)/al-sz-polardb-uat.rwlb.rds.aliyuncs.com?charset=utf8mb4&parseTime=True&loc=Local"
		db, _ = gorm.Open(mysql.Open(datasource), &gorm.Config{})
		mysqlDb, _ := db.DB()
		mysqlDb.SetMaxIdleConns(10) //设置最大连接数
		mysqlDb.SetMaxOpenConns(10) //设置最大的空闲连接数
	})
	return db
}

func SaveEpochElection(epoch uint64, nodeIdList []discover.NodeID) {
	epochList := make([]TbEpoch, len(nodeIdList))
	for idx, nodeId := range nodeIdList {
		epochList[idx] = TbEpoch{Epoch: epoch, NodeId: nodeId.String()}
	}
	MonitorDB().Create(&epochList)
}

func SaveConsensusElection(consensusNo uint64, nodeIdList []discover.NodeID) {
	consensusList := make([]TbConsensus, len(nodeIdList))
	for idx, nodeId := range nodeIdList {
		consensusList[idx] = TbConsensus{ConsensusNo: consensusNo, NodeId: nodeId.String(), StatBlockQty: 0}
	}
	MonitorDB().Create(&consensusList)
}

func InitNodePing(nodeIdList []discover.NodeID) {
	for _, nodeId := range nodeIdList {
		var nodePing TbNodePing
		MonitorDB().Find(&nodePing, "node_id=?", nodeId.String())
		if nodePing.NodeId == "" {
			nodePing = TbNodePing{NodeId: nodeId.String(), Status: 0, ReplyTime: time.Now().Unix()}
			MonitorDB().Create(&nodePing)
		} else {
			nodePing.Status = 0
			MonitorDB().Save(&nodePing)
		}
	}
}

func SaveNodePingResult(nodeId string, ip string, port string, status int8) {
	var nodePing TbNodePing
	MonitorDB().Find(&nodePing, "node_id=?", nodeId)
	if strings.TrimSpace(nodePing.NodeId) != "" {
		nodePing.Ip = ip
		nodePing.Port = port
		nodePing.Status = status
		if status == 1 {
			nodePing.ReplyTime = time.Now().Unix()
		}
		MonitorDB().Save(&nodePing)
	}

	/*var nodePing = TbNodePing{NodeId: nodeId, Ip: ip, Port: port, Status: status, ReplyTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	MonitorDB().Save(&nodePing)*/
}
