package p2p

import (
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"strings"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/log"

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
	Addr       string
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

func InitMonitorDB(dataSource string) {
	monitorDbOnce.Do(func() {
		datasource := dataSource
		//datasource := "user:pass@tcp(127.0.0.1:3306)/al-sz-polardb-uat.rwlb.rds.aliyuncs.com?charset=utf8mb4&parseTime=True&loc=Local"
		db, _ = gorm.Open(mysql.Open(datasource), &gorm.Config{})
		mysqlDb, _ := db.DB()
		mysqlDb.SetMaxIdleConns(10) //设置最大连接数
		mysqlDb.SetMaxOpenConns(10) //设置最大的空闲连接数
	})
}
func MonitorDB() *gorm.DB {
	return db
}

func SaveEpochElection(epoch uint64, nodeIdList []enode.IDv0) {
	log.Info("SaveEpochElection", "epoch", epoch, "nodeIdList", nodeIdList)
	epochList := make([]TbEpoch, len(nodeIdList))
	for idx, nodeId := range nodeIdList {
		epochList[idx] = TbEpoch{Epoch: epoch, NodeId: nodeId.String()}
	}
	if result := MonitorDB().Create(&epochList); result.Error != nil {
		log.Error("failed to insert into tb_epoch", "err", result.Error)
	}
}

func SaveConsensusElection(consensusNo uint64, nodeIdList []enode.IDv0) {
	log.Info("SaveConsensusElection", "consensusNo", consensusNo, "nodeIdList", nodeIdList)
	consensusList := make([]TbConsensus, len(nodeIdList))
	for idx, nodeId := range nodeIdList {
		consensusList[idx] = TbConsensus{ConsensusNo: consensusNo, NodeId: nodeId.String(), StatBlockQty: 0}
	}
	if result := MonitorDB().Create(&consensusList); result.Error != nil {
		log.Error("failed to insert into tb_consensus", "err", result.Error)
	}

}

func InitNodePing(nodeIdList []common.NodeID) {
	log.Info("InitNodePing", "nodeIdList", nodeIdList)
	for _, nodeId := range nodeIdList {
		var nodePing TbNodePing
		if result := MonitorDB().Find(&nodePing, "node_id=?", nodeId.String()); result.Error != nil {
			log.Error("failed to query tb_node_ping", "err", result.Error)
		}
		if nodePing.NodeId == "" {
			nodePing = TbNodePing{NodeId: nodeId.String(), Status: 0}
			if result := MonitorDB().Create(&nodePing); result.Error != nil {
				log.Error("failed to insert into tb_node_ping", "err", result.Error)
			}
		} else {
			//nodePing.Status = 0
			if result := MonitorDB().Save(&nodePing); result.Error != nil {
				log.Error("failed to update tb_node_ping", "err", result.Error)
			}
		}
	}
}

func SaveNodePingResult(nodeId enode.IDv0, addr string, status int8) {
	log.Info("SaveNodePingResult", "nodeId", nodeId.String(), "addr", addr, "status", status)

	var nodePing TbNodePing
	if result := MonitorDB().Find(&nodePing, "node_id=?", nodeId.String()); result.Error != nil {
		log.Error("failed to query tb_node_ping", "err", result.Error)
	}
	if strings.TrimSpace(nodePing.NodeId) != "" {
		nodePing.Addr = addr
		nodePing.Status = status
		if status == 1 {
			nodePing.ReplyTime = time.Now().Unix()
		}
		if result := MonitorDB().Save(&nodePing); result.Error != nil {
			log.Error("failed to update tb_node_ping", "err", result.Error)
		}
	}

	/*var nodePing = TbNodePing{NodeId: nodeId, Ip: ip, Port: port, Status: status, ReplyTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
	MonitorDB().Save(&nodePing)*/
}
