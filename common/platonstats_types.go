package common

//type NodeID [64]byte

var PlatONStatsServiceRunning bool = false

type BlockType uint8

const (
	GenesisBlock BlockType = iota
	GeneralBlock
	ConsensusBeginBlock
	ConsensusElectionBlock
	ConsensusEndBlock
	EpochBeginBlock
	EpochEndBlock
	EndOfYear
)

type EmbedTransferTx struct {
	TxHash Hash    `rlp:"nil"`
	From   Address `rlp:"nil"`
	To     Address `rlp:"nil"`
	Amount uint64
}

type EmbedContractTx struct {
	TxHash          Hash    `rlp:"nil"`
	From            Address `rlp:"nil"`
	ContractAddress Address `rlp:"nil"`
	Input           []byte  `rlp:"nil"`
}

type GenesisData struct {
	AllocItemList []*AllocItem `rlp:"nil"`
}
type AllocItem struct {
	Address Address `rlp:"nil"`
	Amount  uint64
}

func (g *GenesisData) AddAllocItem(address Address, amount uint64) {
	//todo: test
	g.AllocItemList = append(g.AllocItemList, &AllocItem{Address: address, Amount: amount})
}

type AdditionalIssuanceData struct {
	AdditionalNo     uint32          //增发周期
	AdditionalBase   uint64          //增发基数
	AdditionalRate   uint16          //增发比例 单位：万分之一
	AdditionalAmount uint64          //增发金额
	IssuanceItemList []*IssuanceItem `rlp:"nil"` //增发分配
}

type IssuanceItem struct {
	Address Address `rlp:"nil"` //增发金额分配地址
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
	NodeID        [64]byte    `rlp:"nil"` //备选节点ID
	NodeAddress   NodeAddress `rlp:"nil"` //备选节点地址
	RefundEpochNo uint64      //解除质押,资金真正退回的结算周期（此周期最后一个块的endBlocker里
}

type RestrictingReleaseItem struct {
	DestAddress   Address `rlp:"nil"` //释放地址
	ReleaseAmount uint64  //释放金额
}

var ExeBlockDataCollector = make(map[uint64]*ExeBlockData)

func PopExeBlockData(blockNumber uint64) *ExeBlockData {
	exeBlockData, ok := ExeBlockDataCollector[blockNumber]
	if ok {
		delete(ExeBlockDataCollector, blockNumber)
		return exeBlockData
	}
	return nil
}

func InitExeBlockData(blockNumber uint64) {
	if PlatONStatsServiceRunning {
		exeBlockData := &ExeBlockData{
			ZeroSlashingItemList:       make([]*ZeroSlashingItem, 0),
			UnstakingRefundItemList:    make([]*UnstakingRefundItem, 0),
			RestrictingReleaseItemList: make([]*RestrictingReleaseItem, 0),
		}

		ExeBlockDataCollector[blockNumber] = exeBlockData
	}
}

func GetExeBlockData(blockNumber uint64) *ExeBlockData {
	return ExeBlockDataCollector[blockNumber]
}

type ExeBlockData struct {
	AdditionalIssuanceData        *AdditionalIssuanceData        `rlp:"nil"`
	RewardData                    *RewardData                    `rlp:"nil"`
	ZeroSlashingItemList          []*ZeroSlashingItem            `rlp:"nil"`
	DuplicatedSignSlashingSetting *DuplicatedSignSlashingSetting `rlp:"nil"`
	UnstakingRefundItemList       []*UnstakingRefundItem         `rlp:"nil"`
	RestrictingReleaseItemList    []*RestrictingReleaseItem      `rlp:"nil"`
	EmbedTransferTxList           []*EmbedTransferTx             `rlp:"nil"`
	EmbedContractTxList           []*EmbedContractTx             `rlp:"nil"`
}

func CollectUnstakingRefundItem(blockNumber uint64, nodeId [64]byte, nodeAddress NodeAddress, refundEpochNo uint64) {
	if PlatONStatsServiceRunning && ExeBlockDataCollector[blockNumber] != nil {
		d := ExeBlockDataCollector[blockNumber]
		d.UnstakingRefundItemList = append(d.UnstakingRefundItemList, &UnstakingRefundItem{NodeID: nodeId, NodeAddress: nodeAddress, RefundEpochNo: refundEpochNo})
	}
}

func CollectRestrictingReleaseItem(blockNumber uint64, destAddress Address, releaseAmount uint64) {
	if PlatONStatsServiceRunning && ExeBlockDataCollector[blockNumber] != nil {
		d := ExeBlockDataCollector[blockNumber]
		d.RestrictingReleaseItemList = append(d.RestrictingReleaseItemList, &RestrictingReleaseItem{DestAddress: destAddress, ReleaseAmount: releaseAmount})
	}
}

func CollectRewardData(blockNumber uint64, rewardData *RewardData) {
	if PlatONStatsServiceRunning && ExeBlockDataCollector[blockNumber] != nil {
		d := ExeBlockDataCollector[blockNumber]
		d.RewardData = rewardData
	}
}

func CollectDuplicatedSignSlashingSetting(blockNumber uint64, penaltyRatioByValidStakings, rewardRatioByPenalties uint32) {
	if PlatONStatsServiceRunning && ExeBlockDataCollector[blockNumber] != nil {
		d := ExeBlockDataCollector[blockNumber]
		if d.DuplicatedSignSlashingSetting != nil {
			//在同一个区块中，只要设置一次即可
			d.DuplicatedSignSlashingSetting = &DuplicatedSignSlashingSetting{PenaltyRatioByValidStakings: penaltyRatioByValidStakings, RewardRatioByPenalties: rewardRatioByPenalties}
		}
	}
}

func CollectZeroSlashingItem(blockNumber uint64, zeroSlashingItemList []*ZeroSlashingItem) {
	if PlatONStatsServiceRunning && ExeBlockDataCollector[blockNumber] != nil {
		d := ExeBlockDataCollector[blockNumber]
		d.ZeroSlashingItemList = zeroSlashingItemList
	}
}

func CollectEmbedTransferTx(blockNumber uint64, txHash Hash, from, to Address, amount uint64) {
	if PlatONStatsServiceRunning && ExeBlockDataCollector[blockNumber] != nil {
		d := ExeBlockDataCollector[blockNumber]
		d.EmbedTransferTxList = append(d.EmbedTransferTxList, &EmbedTransferTx{TxHash: txHash, From: from, To: to, Amount: amount})
	}
}

func CollectEmbedContractTx(blockNumber uint64, txHash Hash, from, contractAddress Address, input []byte) {
	if PlatONStatsServiceRunning && ExeBlockDataCollector[blockNumber] != nil {
		d := ExeBlockDataCollector[blockNumber]
		d.EmbedContractTxList = append(d.EmbedContractTxList, &EmbedContractTx{TxHash: txHash, From: from, ContractAddress: contractAddress, Input: input})
	}
}
