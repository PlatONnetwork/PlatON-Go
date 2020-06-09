package common

import (
	"bytes"
	"strconv"
)

//type NodeID [64]byte

type GenesisData struct {
	AllocItemList []*AllocItem
}
type AllocItem struct {
	Address Address
	Amount  uint64
}

func (g *GenesisData) AddAllocItem(address Address, amount uint64) {
	//todo: test
	g.AllocItemList = append(g.AllocItemList, &AllocItem{Address: address, Amount: amount})
}

type AdditionalIssuanceData struct {
	AdditionalNo     uint64          //增发周期
	AdditionalBase   uint64          //增发基数
	AdditionalRate   uint16          //增发比例 单位：万分之一
	AdditionalAmount uint64          //增发金额
	IssuanceItemList []*IssuanceItem //增发分配
}

type IssuanceItem struct {
	Address Address //增发金额分配地址
	Amount  uint64  //增发金额
}

func (d *AdditionalIssuanceData) AddIssuanceItem(address Address, amount uint64) {
	//todo: test
	d.IssuanceItemList = append(d.IssuanceItemList, &IssuanceItem{Address: address, Amount: amount})
}

type RewardData struct {
	BlockRewardAmount   uint64           //出块奖励
	StakingRewardAmount uint64           //一结算周期内所有101节点的质押奖励
	CandidateInfoList   []*CandidateInfo `rlp:"nil"` //备选节点信息
}

type CandidateInfo struct {
	NodeID       [64]byte //备选节点ID
	MinerAddress Address  //备选节点的矿工地址（收益地址）
}

type ZeroSlashingItem struct {
	NodeID         [64]byte //备选节点ID
	SlashingAmount uint64   //0出块处罚金(从质押金扣)
}

type DuplicatedSignSlashingSetting struct {
	PenaltyRatioByValidStakings uint32 //unit:1%%		//罚金 = 有效质押 & PenaltyRatioByValidStakings / 10000
	RewardRatioByPenalties      uint32 //unit:1%		//给举报人的赏金=罚金 * RewardRatioByPenalties / 100
}

type UnstakingRefundItem struct {
	NodeID        [64]byte    //备选节点ID
	NodeAddress   NodeAddress //备选节点地址
	RefundEpochNo uint64      //解除质押,资金真正退回的结算周期（此周期最后一个块的endBlocker里
}

type RestrictingReleaseItem struct {
	DestAddress   Address //释放地址
	ReleaseAmount uint64  //释放金额
}

var ExeBlockDataCollector map[uint64]*ExeBlockData

func GetExeBlockData(blockNumber uint64) *ExeBlockData {
	return ExeBlockDataCollector[blockNumber]
}

func InitExeBlockData(blockNumber uint64) *ExeBlockData {
	exeBlockData := &ExeBlockData{
		ZeroSlashingItemList:       make([]*ZeroSlashingItem, 0),
		UnstakingRefundItemList:    make([]*UnstakingRefundItem, 0),
		RestrictingReleaseItemList: make([]*RestrictingReleaseItem, 0),
	}

	ExeBlockDataCollector[blockNumber] = exeBlockData
	return exeBlockData
}

type ExeBlockData struct {
	AdditionalIssuanceData        *AdditionalIssuanceData        `rlp:"nil"`
	RewardData                    *RewardData                    `rlp:"nil"`
	ZeroSlashingItemList          []*ZeroSlashingItem            `rlp:"nil"`
	DuplicatedSignSlashingSetting *DuplicatedSignSlashingSetting `rlp:"nil"`
	UnstakingRefundItemList       []*UnstakingRefundItem         `rlp:"nil"`
	RestrictingReleaseItemList    []*RestrictingReleaseItem      `rlp:"nil"`
}

func (d *ExeBlockData) AddUnstakingRefundItem(nodeId [64]byte, nodeAddress NodeAddress, refundEpochNo uint64) {
	//todo: test
	d.UnstakingRefundItemList = append(d.UnstakingRefundItemList, &UnstakingRefundItem{NodeID: nodeId, NodeAddress: nodeAddress, RefundEpochNo: refundEpochNo})
}

func (d *ExeBlockData) AddRestrictingReleaseItem(destAddress Address, releaseAmount uint64) {
	//todo: test
	d.RestrictingReleaseItemList = append(d.RestrictingReleaseItemList, &RestrictingReleaseItem{DestAddress: destAddress, ReleaseAmount: releaseAmount})
}

func (d *ExeBlockData) AttachDuplicatedSignSlashingSetting(penaltyRatioByValidStakings, rewardRatioByPenalties uint32) {
	d.DuplicatedSignSlashingSetting = &DuplicatedSignSlashingSetting{PenaltyRatioByValidStakings: penaltyRatioByValidStakings, RewardRatioByPenalties: rewardRatioByPenalties}
}

var KeyDelimiter = []byte(":")
var ExeBlockDataKey = []byte("ExeBlockDataKey")

//todo: refer to blockchian.go's writeHeader(header *types.Header) error {
func ExeBlockDataKeyValue(blockNumber uint64) []byte {
	return bytes.Join([][]byte{
		ExeBlockDataKey,
		[]byte(strconv.FormatUint(blockNumber, 10)),
	}, KeyDelimiter)
}

/*
func AddUnstakingRefundItem(blockHash Hash, blockNumber uint64, nodeID discover.NodeID, nodeAddress NodeAddress, refundEpochNo uint64) error {
	item := &UnstakingRefundItem{NodeID: nodeID, NodeAddress: nodeAddress, RefundEpochNo: refundEpochNo}
	if unstakingRefundData, err := GetUnstakingRefundData(blockHash, blockNumber); err != nil {
		return err
	} else {
		if unstakingRefundData == nil {
			list := []*UnstakingRefundItem{item}
			unstakingRefundData = &UnstakingRefundData{UnstakingRefundItemList: list}
		} else {
			unstakingRefundData.UnstakingRefundItemList = append(unstakingRefundData.UnstakingRefundItemList, item)
		}
		return snapshotdb.Instance().Put(blockHash, UnstakingRefundItemListKeyValue(blockNumber), MustRlpEncode(unstakingRefundData))
	}
}

func GetUnstakingRefundData(blockHash Hash, blockNumber uint64) (*UnstakingRefundData, error) {
	key := UnstakingRefundItemListKeyValue(blockNumber)
	bytes, err := snapshotdb.Instance().Get(blockHash, key)
	if err != nil && err != snapshotdb.ErrNotFound {
		return nil, err
	}
	var unstakingRefundData *UnstakingRefundData
	if len(bytes) > 0 {
		if err = rlp.DecodeBytes(bytes, unstakingRefundData); err != nil {
			return nil, err
		}
	}
	return unstakingRefundData, nil
}
func RemoveUnstakingRefundData(blockHash Hash, blockNumber uint64) error {
	key := UnstakingRefundItemListKeyValue(blockNumber)
	return snapshotdb.Instance().Del(blockHash, key)
}
*/
