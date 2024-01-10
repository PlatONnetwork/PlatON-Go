package p2p

import (
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"strings"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

type Downloading interface {
	HighestBlock() uint64
}

type monitorDialTasks struct {
	queue []*dialTask
}

func NewMonitorDialedTasks() *monitorDialTasks {
	tasks := &monitorDialTasks{
		queue: make([]*dialTask, 0),
	}
	return tasks
}

func (tasks *monitorDialTasks) listTask() []*dialTask {
	return tasks.queue
}

// adding new task to the end of the queue
func (tasks *monitorDialTasks) offer(task *dialTask) {
	tasks.queue = append(tasks.queue, task)
}

func (tasks *monitorDialTasks) removeTask(NodeID enode.ID) {

	log.Info("[before remove]monitor dial task list before removeTask operation", "task queue", tasks.description())
	if !tasks.isEmpty() {
		for i, t := range tasks.queue {
			if t.dest.ID() == NodeID {
				tasks.queue = append(tasks.queue[:i], tasks.queue[i+1:]...)
				break
			}
		}
	}
	log.Info("[after remove]monitor dial task list after removeTask operation", "task queue", tasks.description())
}

// remove the first task in the queue
func (tasks *monitorDialTasks) poll() *dialTask {
	if tasks.isEmpty() {
		log.Info("dialedTasks is empty!")
		return nil
	}

	pollTask := tasks.queue[0]
	tasks.queue = tasks.queue[1:]
	return pollTask
}

// remove the specify index task in the queue
func (tasks *monitorDialTasks) pollIndex(index int) *dialTask {
	if tasks.isEmpty() {
		log.Info("dialedTasks is empty!")
		return nil
	}

	pollTask := tasks.queue[index]
	tasks.queue = append(tasks.queue[:index], tasks.queue[index+1:]...)
	return pollTask
}

// index of task in the queue
func (tasks *monitorDialTasks) index(task *dialTask) int {
	for i, t := range tasks.queue {
		if t.dest.ID() == task.dest.ID() {
			return i
		}
	}
	return -1
}

// queue size
func (tasks *monitorDialTasks) size() int {
	return len(tasks.queue)
}

// clear queue
func (tasks *monitorDialTasks) clear() bool {
	if tasks.isEmpty() {
		log.Info("queue is empty!")
		return false
	}
	for i := 0; i < tasks.size(); i++ {
		tasks.queue[i] = nil
	}
	tasks.queue = nil
	return true
}

// whether the queue is empty
func (tasks *monitorDialTasks) isEmpty() bool {
	if len(tasks.queue) == 0 {
		return true
	}
	return false
}

func (tasks *monitorDialTasks) description() []string {
	var description []string
	for _, t := range tasks.queue {
		description = append(description, t.dest.ID().TerminalString())
	}
	return description
}

func saveNodePingResult(node *enode.Node, addr string, status int8) {
	log.Info("SaveNodePingResult", "idV0", node.IDv0(), "addr", addr, "status", status)

	var nodePing TbNodePing
	if result := MonitorDB().Find(&nodePing, "node_id=?", node.IDv0()); result.Error != nil {
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

// param: eventMux
// param: blockNumber, 选举块高，结算周期末，选出下一个结算周期的备选101节点
// param: epoch 结算周期。从1开始计算。创世块可以认为是0
// param: verifierList
// param: downloading
func PostMonitorNodeEvent(eventMux *event.TypeMux, blockNumber uint64, epoch uint64, enodeList []*enode.Node, downloading Downloading) {
	//nodeIdList := ConvertToNodeIdList(verifierList)
	//nodeIdStringList := xcom.ConvertToNodeIdStringList(verifierList)
	//MONITOR，保存这一轮结算周期的新101名单
	log.Info("PostMonitorNodeEvent", "blockNumber", blockNumber, "epoch", epoch, "enodeList", enodeList)

	//SaveEpochElection(epoch, nodeIdList)

	if blockNumber == 0 {
		log.Info("current block is genesis block")

		InitNodePing(enodeList)
		if err := eventMux.Post(cbfttypes.ElectNextEpochVerifierEvent{NodeList: enodeList}); err != nil {
			log.Error("post ElectNextEpochVerifierEvent failed", "enodeList", enodeList, "err", err)
		}

	} else {
		if blockNumber > downloading.HighestBlock() {
			log.Info("current block is consensus block")
			//说明区块是共识协议得到的，需要执行monitor任务
			//MONITOR，保存新101名单
			InitNodePing(enodeList)
			if err := eventMux.Post(cbfttypes.ElectNextEpochVerifierEvent{NodeList: enodeList}); err != nil {
				log.Error("post ElectNextEpochVerifierEvent failed", "enodeList", enodeList, "err", err)
			}
		} else {
			//此次需要同步的最高块，可以认为是链上当前块，如果当前结算周期和链上结算周期一致，则开始执行monitor任务
			chainEpoch := xutil.CalculateEpoch(downloading.HighestBlock())
			if chainEpoch == epoch && epoch != 1 { //epoch=1 第一个epoch的 监控信息，已经有创世块写入。
				log.Info("current block is downloaded block and epoch is same as chain")
				InitNodePing(enodeList)
				if err := eventMux.Post(cbfttypes.ElectNextEpochVerifierEvent{NodeList: enodeList}); err != nil {
					log.Error("post ElectNextEpochVerifierEvent failed", "enodeList", enodeList, "err", err)
				}
			} else {
				log.Info("current block is downloaded block but far way from chain")
			}
		}
	}
}

func ConvertToNodeIdStringList(verifierList []*staking.Validator) []string {
	nodeIdStringList := make([]string, len(verifierList))
	for i, verifier := range verifierList {
		nodeIdStringList[i] = verifier.NodeId.String()
	}
	return nodeIdStringList
}

func ConvertToCommonNodeIdList(verifierList []*staking.Validator) []common.NodeID {
	nodeIdList := make([]common.NodeID, len(verifierList))
	for i, verifier := range verifierList {
		nodeIdList[i] = common.NodeID(verifier.NodeId)
	}
	return nodeIdList
}

func ConvertToENodeList(verifierList []*staking.Validator) []*enode.Node {
	enodeList := make([]*enode.Node, len(verifierList))
	for i, verifier := range verifierList {
		pub, err := verifier.NodeId.Pubkey()
		if err != nil {
			panic(err)
		}
		enodeList[i] = enode.NewV4(pub, nil, 0, 0)
	}
	return enodeList
}
