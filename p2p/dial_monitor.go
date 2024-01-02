package p2p

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

var (
	monitorSchedulerOnce sync.Once
	monitorSchedulerRef  *monitorScheduler
)

type monitorScheduler struct {
	queue                    []*monitorTask
	monitorTaskDoneFurtherFn monitorTaskDoneFurtherFn
}

type monitorTask struct {
	staticPoolIndex int
	flags           connFlag
	dest            *enode.Node
	lastResolved    time.Time
	resolveDelay    time.Duration
	err             error
}

func (t *monitorTask) String() string {
	return fmt.Sprintf("{flags: %v, dest: %v, lastResolved: %v}", t.flags, t.dest, t.lastResolved)
}

func MonitorScheduler() *monitorScheduler {
	monitorSchedulerOnce.Do(func() {
		log.Info("Init node monitor scheduler ...")
		monitorSchedulerRef = &monitorScheduler{}
	})
	return monitorSchedulerRef
}

func (tasks *monitorScheduler) InitMonitorTaskDoneFurtherFn(monitorTaskDoneFurtherFn monitorTaskDoneFurtherFn) {
	tasks.monitorTaskDoneFurtherFn = monitorTaskDoneFurtherFn
}

// whether the queue is empty
func (tasks *monitorScheduler) isEmpty() bool {
	if len(tasks.queue) == 0 {
		return true
	}
	return false
}

func (tasks *monitorScheduler) ClearMonitorScheduler() {
	tasks.queue = []*monitorTask{}
}

func (tasks *monitorScheduler) AddMonitorTask(task *monitorTask) {
	tasks.queue = append(tasks.queue, task)
}
func (tasks *monitorScheduler) RemoveMonitorTask(nodeId enode.Node) {
	log.Info("before RemoveMonitorTask", "nodeId", nodeId, "task queue length", len(tasks.queue), "task queue", tasks.description())
	if !tasks.isEmpty() {
		for i, t := range tasks.queue {
			if t.dest.IDv0() == nodeId.IDv0() {
				tasks.queue = append(tasks.queue[:i], tasks.queue[i+1:]...)
				break
			}
		}
	}
	log.Info("after RemoveMonitorTask", "nodeId", nodeId, "task queue length", len(tasks.queue), "task queue", tasks.description())
}

func (tasks *monitorScheduler) ListTask() []*monitorTask {
	return tasks.queue
}

func (tasks *monitorScheduler) RemoveTask(NodeID enode.Node) {
	if !tasks.isEmpty() {
		for i, t := range tasks.queue {
			if t.dest.IDv0() == NodeID.IDv0() {
				tasks.queue = append(tasks.queue[:i], tasks.queue[i+1:]...)
				break
			}
		}
	}
}

func (tasks *monitorScheduler) description() string {
	var description []string
	for _, t := range tasks.queue {
		description = append(description, fmt.Sprintf("%x", t.dest.IDv0()))
	}
	return strings.Join(description, ",")
}

func (t *monitorTask) Do(srv *Server) {
	log.Info("monitorTask.Do", "id", t.dest.ID)
	if t.dest.Incomplete() {
		if !t.resolve(srv) {
			return
		}
	}
	err := t.dial(srv, t.dest)
	if err != nil {
		log.Trace("Dial error", "task", t, "err", err)
		// Try resolving the ID of static nodes if dialing failed.
		if _, ok := err.(*dialError); ok && t.flags&staticDialedConn != 0 {
			if t.resolve(srv) {
				t.dial(srv, t.dest)
			}
		}
	}
	t.err = err
}

// monitorTask任务结束后的后续操作（保存NodePing结果，从监控任务列表删除任务）
func (t *monitorTask) MonitorTaskDoneFurther() bool {
	log.Info("monitorTask.MonitorTaskDoneFurther", "id", t.dest)
	return MonitorScheduler().monitorTaskDoneFurtherFn(t.dest)
}

func (t *monitorTask) resolve(srv *Server) bool {
	if srv.ntab == nil {
		log.Debug("Can't resolve node", "id", t.dest.ID, "err", "discovery is disabled")
		return false
	}
	if t.resolveDelay == 0 {
		t.resolveDelay = initialResolveDelay
	}
	if time.Since(t.lastResolved) < t.resolveDelay {
		return false
	}
	resolved := srv.ntab.Resolve(t.dest)
	t.lastResolved = time.Now()
	if resolved == nil {
		t.resolveDelay *= 2
		if t.resolveDelay > maxResolveDelay {
			t.resolveDelay = maxResolveDelay
		}
		log.Debug("Resolving node failed", "id", t.dest.ID, "newdelay", t.resolveDelay)
		return false
	}
	// The node was found.
	t.resolveDelay = initialResolveDelay
	t.dest = resolved
	log.Debug("Resolved node", "id", t.dest.ID, "addr", &net.TCPAddr{IP: t.dest.IP(), Port: t.dest.TCP()})

	return true
}

// dial performs the actual connection attempt.
func (t *monitorTask) dial(srv *Server, dest *enode.Node) error {
	fd, err := srv.Dialer.Dial(srv.dialsched.ctx, dest)
	if err != nil {
		return &dialError{err}
	}
	mfd := newMeteredConn(fd, false, &net.TCPAddr{IP: dest.IP(), Port: dest.TCP()})
	return srv.SetupConn(mfd, t.flags, dest)
}

type Downloading interface {
	HighestBlock() uint64
}

// param: eventMux
// param: blockNumber, 选举块高，结算周期末，选出下一个结算周期的备选101节点
// param: epoch 结算周期。从1开始计算。创世块可以认为是0
// param: verifierList
// param: downloading
func PostMonitorNodeEvent(eventMux *event.TypeMux, blockNumber uint64, epoch uint64, nodeIdList []common.NodeID, downloading Downloading) {
	//nodeIdList := ConvertToNodeIdList(verifierList)
	//nodeIdStringList := xcom.ConvertToNodeIdStringList(verifierList)
	//MONITOR，保存这一轮结算周期的新101名单
	log.Info("PostMonitorNodeEvent", "blockNumber", blockNumber, "epoch", epoch, "nodeIdList", nodeIdList)

	//SaveEpochElection(epoch, nodeIdList)

	if blockNumber == 0 {
		log.Info("current block is genesis block")

		InitNodePing(nodeIdList)
		if err := eventMux.Post(cbfttypes.ElectNextEpochVerifierEvent{NodeIdList: nodeIdList}); err != nil {
			log.Error("post ElectNextEpochVerifierEvent failed", "nodeIdList", nodeIdList, "err", err)
		}

	} else {
		if blockNumber > downloading.HighestBlock() {
			log.Info("current block is consensus block")
			//说明区块是共识协议得到的，需要执行monitor任务
			//MONITOR，保存新101名单
			InitNodePing(nodeIdList)
			if err := eventMux.Post(cbfttypes.ElectNextEpochVerifierEvent{NodeIdList: nodeIdList}); err != nil {
				log.Error("post ElectNextEpochVerifierEvent failed", "nodeIdList", nodeIdList, "err", err)
			}
		} else {
			//此次需要同步的最高块，可以认为是链上当前块，如果当前结算周期和链上结算周期一致，则开始执行monitor任务
			chainEpoch := xutil.CalculateEpoch(downloading.HighestBlock())
			if chainEpoch == epoch && epoch != 1 { //epoch=1 第一个epoch的 监控信息，已经有创世块写入。
				log.Info("current block is downloaded block and epoch is same as chain")
				InitNodePing(nodeIdList)
				if err := eventMux.Post(cbfttypes.ElectNextEpochVerifierEvent{NodeIdList: nodeIdList}); err != nil {
					log.Error("post ElectNextEpochVerifierEvent failed", "nodeIdList", nodeIdList, "err", err)
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

//func ConvertToNodeIdList(verifierList []*staking.Validator) []enode.Node {
//	nodeIdList := make([]enode.Node, len(verifierList))
//	for i, verifier := range verifierList {
//		nodeIdList[i] = verifier.NodeId
//	}
//	return nodeIdList
//}

func ConvertToCommonNodeIdList(verifierList []*staking.Validator) []common.NodeID {
	nodeIdList := make([]common.NodeID, len(verifierList))
	for i, verifier := range verifierList {
		nodeIdList[i] = common.NodeID(verifier.NodeId)
	}
	return nodeIdList
}

/*func (t *monitorTask) dialAndClose(srv *Server, dest *discover.Node) error {
	fd, err := srv.Dialer.Dial(dest)
	if err != nil {
		return &dialError{err}
	}
	fd.Close()
	return nil
}
*/
