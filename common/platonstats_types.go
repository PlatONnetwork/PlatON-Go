package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

//uint32(年度<<16 | 结算周期<<8 | 共识周期)
// 年度：0-中间；1-开始；2-结束
// 结算周期：0-中间；1-开始；2-结束
// 共识周期：0-中间；1-开始；2-选举；3-结束
//type BlockType uint32

type BlockType uint8

type NodeID [512 / 8]byte

type Input []byte

var nodeIdT = reflect.TypeOf(NodeID{})

// MarshalText returns the hex representation of a.
func (a NodeID) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *NodeID) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("common.NodeID", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *NodeID) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(nodeIdT, input, a[:])
}
func (n NodeID) TerminalString() string {
	return hex.EncodeToString(n[:8])
}

// NodeID prints as a long hexadecimal number.
func (n NodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}

var inputT = reflect.TypeOf(Input{})

// MarshalText returns the hex representation of a.
func (a Input) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Input) UnmarshalText(input []byte) error {
	hexBytes, err := hexutil.Decode(string(input[1 : len(input)-1]))
	if err != nil {
		return err
	}
	aa := make(Input, len(hexBytes))

	err = hexutil.UnmarshalFixedText("common.Input", input, aa)
	if err != nil {
		return err
	}
	a = &aa
	return nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Input) UnmarshalJSON(input []byte) error {
	//string(input)="0x0102030405", so, firstly remove the "", and then, to decode it.
	hexBytes, err := hexutil.Decode(string(input[1 : len(input)-1]))
	if err != nil {
		return err
	}
	aa := make(Input, len(hexBytes))
	err = hexutil.UnmarshalFixedJSON(inputT, input, aa)
	if err != nil {
		return err
	}
	a = &aa
	return nil
}

func MustHexID(in string) NodeID {
	id, err := HexID(in)
	if err != nil {
		panic(err)
	}
	return id
}

func HexID(in string) (NodeID, error) {
	var id NodeID
	b, err := hex.DecodeString(strings.TrimPrefix(in, "0x"))
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("wrong length, want %d hex chars", len(id)*2)
	}
	copy(id[:], b)
	return id, nil
}

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
	TxHash Hash     `json:"txHash,omitempty"`
	From   Address  `json:"from,omitempty"`
	To     Address  `json:"to,omitempty"`
	Amount *big.Int `json:"amount,omitempty"`
}

type EmbedContractTx struct {
	TxHash          Hash    `json:"txHash,omitempty"`
	From            Address `json:"from,omitempty"`
	ContractAddress Address `json:"contractAddress,omitempty"`
	Input           string  `json:"input,omitempty"` //hex string
}

type GenesisData struct {
	AllocItemList             []*AllocItem       `json:"allocItemList,omitempty"`
	StakingItemList           []*StakingItem     `json:"stakingItemList,omitempty"`
	RestrictingCreateItemList []*RestrictingItem `json:"restrictingCreateItemList,omitempty"`
	InitFundItemList          []*InitFundItem    `json:"initFundItemList,omitempty"`
	EpochElection             []NodeID           `json:"epochElection,omitempty"`
	ConsensusElection         []NodeID           `json:"consensusElection,omitempty"`
}
type AllocItem struct {
	Address Address  `json:"address,omitempty"`
	Amount  *big.Int `json:"amount,omitempty"`
}

type StakingItem struct {
	NodeID         NodeID   `json:"nodeID,omitempty"`
	StakingAddress Address  `json:"stakingAddress,omitempty"`
	BenefitAddress Address  `json:"benefitAddress,omitempty"`
	NodeName       string   `json:"nodeName,omitempty"`
	Amount         *big.Int `json:"amount,omitempty"`
}

type RestrictingItem struct {
	From        Address    `json:"from,omitempty"`
	DestAddress Address    `json:"destAddress,omitempty"`
	Plans       []*big.Int `json:"plans,omitempty"`
}

type InitFundItem struct {
	From   Address  `json:"from,omitempty"`
	To     Address  `json:"to,omitempty"`
	Amount *big.Int `json:"amount,omitempty"`
}

type AutoStakingTx struct {
	RestrictingAmount *big.Int `json:"restrictingAmount,omitempty"`
	BalanceAmount     *big.Int `json:"balanceAmount,omitempty"`
}

func (g *GenesisData) AddAllocItem(address Address, amount *big.Int) {
	g.AllocItemList = append(g.AllocItemList, &AllocItem{Address: address, Amount: amount})
}
func (g *GenesisData) AddRestrictingCreateItem(from, to Address, plans []*big.Int) {
	g.RestrictingCreateItemList = append(g.RestrictingCreateItemList, &RestrictingItem{From: from, DestAddress: to, Plans: plans})
}

func (g *GenesisData) AddInitFundItem(from, to Address, initAmount *big.Int) {
	g.InitFundItemList = append(g.InitFundItemList, &InitFundItem{From: from, To: to, Amount: initAmount})
}

func (g *GenesisData) AddStakingItem(nodeID NodeID, nodeName string, stakingAddress, benefitAddress Address, amount *big.Int) {
	g.StakingItemList = append(g.StakingItemList, &StakingItem{NodeID: nodeID, NodeName: nodeName, StakingAddress: stakingAddress, BenefitAddress: benefitAddress, Amount: amount})
}

type AdditionalIssuanceData struct {
	AdditionalNo     uint32          `json:"additionalNo,omitempty"`     //增发周期
	AdditionalBase   *big.Int        `json:"additionalBase,omitempty"`   //增发基数
	AdditionalRate   uint16          `json:"additionalRate,omitempty"`   //增发比例 单位：万分之一
	AdditionalAmount *big.Int        `json:"additionalAmount,omitempty"` //增发金额
	IssuanceItemList []*IssuanceItem `json:"issuanceItemList,omitempty"` //增发分配
}

type IssuanceItem struct {
	Address Address  `json:"address,omitempty"` //增发金额分配地址
	Amount  *big.Int `json:"amount,omitempty"`  //增发金额
}

func (d *AdditionalIssuanceData) AddIssuanceItem(address Address, amount *big.Int) {
	//todo: test
	d.IssuanceItemList = append(d.IssuanceItemList, &IssuanceItem{Address: address, Amount: amount})
}

// 分配奖励，包括出块奖励，质押奖励
//  注意：委托人不一定每次都能参与到出块奖励的分配中（共识论跨结算周期时会出现，此时节点虽然还在出块，但是可能已经不在当前结算周期的101备选人列表里了，那这个出块节点的委托人在当前结算周期，就不参与这个块的出块奖励分配）
type RewardData struct {
	BlockRewardAmount   *big.Int         `json:"blockRewardAmount,omitempty"`   //出块奖励
	DelegatorReward     bool             `json:"delegatorReward"`               //出块奖励中，分配给委托人的奖励
	StakingRewardAmount *big.Int         `json:"stakingRewardAmount,omitempty"` //一结算周期内所有101节点的质押奖励
	CandidateInfoList   []*CandidateInfo `json:"candidateInfoList,omitempty"`   //备选节点信息
}

type CandidateInfo struct {
	NodeID       NodeID  `json:"nodeId,omitempty"`       //备选节点ID
	MinerAddress Address `json:"minerAddress,omitempty"` //备选节点的矿工地址（收益地址）
}

type ZeroSlashingItem struct {
	NodeID         NodeID   `json:"nodeId,omitempty"`         //备选节点ID
	SlashingAmount *big.Int `json:"slashingAmount,omitempty"` //0出块处罚金(从质押金扣)
}

type DuplicatedSignSlashingSetting struct {
	PenaltyRatioByValidStakings uint32 `json:"penaltyRatioByValidStakings,omitempty"` //unit:1%%		//罚金 = 有效质押 * PenaltyRatioByValidStakings / 10000
	RewardRatioByPenalties      uint32 `json:"rewardRatioByPenalties,omitempty"`      //unit:1%		//给举报人的赏金=罚金 * RewardRatioByPenalties / 100
}

type StakingSetting struct {
	OperatingThreshold *big.Int `json:"operatingThreshold,omitempty"` //质押，委托操作，要求的最小数量；当某次操作后，剩余数量小于此值时，这剩余数量将随此次操作一次处理完。
}

type StakingFrozenItem struct {
	NodeID        NodeID  `json:"nodeId,omitempty"`        //备选节点ID
	NodeAddress   Address `json:"nodeAddress,omitempty"`   //备选节点地址
	FrozenEpochNo uint64  `json:"frozenEpochNo,omitempty"` //质押资金，被解冻的结算周期（此周期最后一个块的endBlocker里）
	Recovery      bool    `json:"recovery"`                //Recover=true；表示冻结期结束后，质押将变成有效质押；Recover=false, 表示冻结期结束后，质押将原来退回质押钱包（或者和锁仓合约）
}

type RestrictingReleaseItem struct {
	DestAddress   Address  `json:"destAddress,omitempty,omitempty"` //释放地址
	ReleaseAmount *big.Int `json:"releaseAmount,omitempty"`         //释放金额
	LackingAmount *big.Int `json:"lackingAmount,omitempty"`         //欠释放金额
}

//todo:改名
//撤消委托后领取的奖励（全部减持）
type WithdrawDelegation struct {
	TxHash          Hash     `json:"txHash,omitempty"`                    //委托用户撤销节点的全部委托的交易HASH
	DelegateAddress Address  `json:"delegateAddress,omitempty,omitempty"` //委托用户地址
	NodeID          NodeID   `json:"nodeId,omitempty"`                    //委托用户委托的节点ID
	RewardAmount    *big.Int `json:"rewardAmount,omitempty"`              //委托用户从此节点获取的全部委托奖励
}

var ExeBlockDataCollector = make(map[uint64]*ExeBlockData)

func PopExeBlockData(blockNumber uint64) *ExeBlockData {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		delete(ExeBlockDataCollector, blockNumber)
		return exeBlockData
	}
	return nil
}

func InitExeBlockData(blockNumber uint64) {
	exeBlockData := &ExeBlockData{
		ZeroSlashingItemList:       make([]*ZeroSlashingItem, 0),
		StakingFrozenItemList:      make([]*StakingFrozenItem, 0),
		RestrictingReleaseItemList: make([]*RestrictingReleaseItem, 0),
		EmbedTransferTxList:        make([]*EmbedTransferTx, 0),
		EmbedContractTxList:        make([]*EmbedContractTx, 0),
		AutoStakingMap:             make(map[Hash]*AutoStakingTx),
	}

	ExeBlockDataCollector[blockNumber] = exeBlockData
}

func GetExeBlockData(blockNumber uint64) *ExeBlockData {
	return ExeBlockDataCollector[blockNumber]
}

type ExeBlockData struct {
	ActiveVersion                 string                         `json:"activeVersion,omitempty"` //如果当前块有升级提案生效，则填写新版本,0.14.0
	AdditionalIssuanceData        *AdditionalIssuanceData        `json:"additionalIssuanceData,omitempty"`
	RewardData                    *RewardData                    `json:"rewardData,omitempty"`
	ZeroSlashingItemList          []*ZeroSlashingItem            `json:"zeroSlashingItemList,omitempty"`
	DuplicatedSignSlashingSetting *DuplicatedSignSlashingSetting `json:"duplicatedSignSlashingSetting,omitempty"`
	StakingSetting                *StakingSetting                `json:"stakingSetting,omitempty"`
	StakingFrozenItemList         []*StakingFrozenItem           `json:"stakingFrozenItemList,omitempty"`
	RestrictingReleaseItemList    []*RestrictingReleaseItem      `json:"restrictingReleaseItemList,omitempty"`
	EmbedTransferTxList           []*EmbedTransferTx             `json:"embedTransferTxList,omitempty"`    //一个显式交易引起的内置转账交易：一般有两种情况：1是部署，或者调用合约时，带上了value，则这个value会转账给合约地址；2是调用合约，合约内部调用transfer()函数完成转账
	EmbedContractTxList           []*EmbedContractTx             `json:"embedContractTxList,omitempty"`    //一个显式交易引起的内置合约交易。这个显式交易显然也是个合约交易，在这个合约里，又调用了其他合约（包括内置合约）
	WithdrawDelegationList        []*WithdrawDelegation          `json:"withdrawDelegationList,omitempty"` //当委托用户撤回节点的全部委托时，需要的统计信息（由于Alaya在运行中，只能兼容Alaya的bug）
	AutoStakingMap                map[Hash]*AutoStakingTx        `json:"autoStakingTxMap,omitempty"`
	EpochElection                 []NodeID                       `json:"epochElection,omitempty"`
	ConsensusElection             []NodeID                       `json:"consensusElection,omitempty"`
}

func CollectAdditionalIssuance(blockNumber uint64, additionalIssuanceData *AdditionalIssuanceData) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		json, _ := json.Marshal(additionalIssuanceData)
		log.Debug("CollectAdditionalIssuance", "blockNumber", blockNumber, "additionalIssuanceData", string(json))
		exeBlockData.AdditionalIssuanceData = additionalIssuanceData
	}
}

func CollectStakingFrozenItem(blockNumber uint64, nodeId NodeID, nodeAddress NodeAddress, frozenEpochNo uint64, recovery bool) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectStakingFrozenItem", "blockNumber", blockNumber, "nodeId", Bytes2Hex(nodeId[:]), "nodeAddress", nodeAddress.Hex(), "frozenEpochNo", frozenEpochNo, "recovery", recovery)
		exeBlockData.StakingFrozenItemList = append(exeBlockData.StakingFrozenItemList, &StakingFrozenItem{NodeID: nodeId, NodeAddress: Address(nodeAddress), FrozenEpochNo: frozenEpochNo, Recovery: recovery})
	}
}

func CollectRestrictingReleaseItem(blockNumber uint64, destAddress Address, releaseAmount *big.Int, lackingAmount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectRestrictingReleaseItem", "blockNumber", blockNumber, "destAddress", destAddress, "releaseAmount", releaseAmount, "lackingAmount", lackingAmount)
		exeBlockData.RestrictingReleaseItemList = append(exeBlockData.RestrictingReleaseItemList, &RestrictingReleaseItem{DestAddress: destAddress, ReleaseAmount: releaseAmount, LackingAmount: lackingAmount})
	}
}

func CollectBlockRewardData(blockNumber uint64, blockRewardAmount *big.Int, delegatorReward bool) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectBlockRewardData", "blockNumber", blockNumber, "blockRewardAmount", blockRewardAmount, "delegatorReward", delegatorReward)
		if exeBlockData.RewardData == nil {
			exeBlockData.RewardData = new(RewardData)
		}
		exeBlockData.RewardData.BlockRewardAmount = blockRewardAmount
		exeBlockData.RewardData.DelegatorReward = delegatorReward
	}
}

func CollectStakingRewardData(blockNumber uint64, stakingRewardAmount *big.Int, candidateInfoList []*CandidateInfo) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectStakingRewardData", "blockNumber", blockNumber, "stakingRewardAmount", stakingRewardAmount)
		for _, candidateInfo := range candidateInfoList {
			log.Debug("nodeID:" + Bytes2Hex(candidateInfo.NodeID[:]))
		}
		if exeBlockData.RewardData == nil {
			exeBlockData.RewardData = new(RewardData)
		}
		exeBlockData.RewardData.StakingRewardAmount = stakingRewardAmount
		exeBlockData.RewardData.CandidateInfoList = candidateInfoList
	}
}

func CollectDuplicatedSignSlashingSetting(blockNumber uint64, penaltyRatioByValidStakings, rewardRatioByPenalties uint32) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectDuplicatedSignSlashingSetting", "blockNumber", blockNumber, "penaltyRatioByValidStakings", penaltyRatioByValidStakings, "rewardRatioByPenalties", rewardRatioByPenalties)
		if exeBlockData.DuplicatedSignSlashingSetting == nil {
			//在同一个区块中，只要设置一次即可
			exeBlockData.DuplicatedSignSlashingSetting = &DuplicatedSignSlashingSetting{PenaltyRatioByValidStakings: penaltyRatioByValidStakings, RewardRatioByPenalties: rewardRatioByPenalties}
		}
	}
}

func CollectStakingSetting(blockNumber uint64, operatingThreshold *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectStakingSetting", "blockNumber", blockNumber, "operatingThreshold", operatingThreshold)
		if exeBlockData.StakingSetting == nil {
			exeBlockData.StakingSetting = &StakingSetting{OperatingThreshold: operatingThreshold}
		}
	}
}

func CollectZeroSlashingItem(blockNumber uint64, nodeId NodeID, slashingAmount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectZeroSlashingItem", "blockNumber", blockNumber, "nodeId", Bytes2Hex(nodeId[:]), "slashingAmount", slashingAmount)
		exeBlockData.ZeroSlashingItemList = append(exeBlockData.ZeroSlashingItemList, &ZeroSlashingItem{NodeID: nodeId, SlashingAmount: slashingAmount})
	}
}

/*func CollectZeroSlashingItem(blockNumber uint64, zeroSlashingItemList []*ZeroSlashingItem) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		json, _ := json.Marshal(zeroSlashingItemList)
		log.Debug("CollectZeroSlashingItem", "blockNumber", blockNumber, "zeroSlashingItemList", string(json))
		exeBlockData.ZeroSlashingItemList = zeroSlashingItemList
	}
}*/

func CollectEmbedTransferTx(blockNumber uint64, txHash Hash, from, to Address, amount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectEmbedTransferTx", "blockNumber", blockNumber, "txHash", txHash.Hex(), "from", from.Bech32(), "to", to.Bech32(), "amount", amount)
		amt := new(big.Int).Set(amount)
		exeBlockData.EmbedTransferTxList = append(exeBlockData.EmbedTransferTxList, &EmbedTransferTx{TxHash: txHash, From: from, To: to, Amount: amt})
	}
}

func CollectEmbedContractTx(blockNumber uint64, txHash Hash, from, contractAddress Address, input []byte) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectEmbedContractTx", "blockNumber", blockNumber, "txHash", txHash.Hex(), "contractAddress", from.Bech32(), "input", Bytes2Hex(input))
		exeBlockData.EmbedContractTxList = append(exeBlockData.EmbedContractTxList, &EmbedContractTx{TxHash: txHash, From: from, ContractAddress: contractAddress, Input: Bytes2Hex(input)})
	}
}

//撤消委托时，才需要收集委托奖励总金额
func CollectWithdrawDelegation(blockNumber uint64, txHash Hash, delegateAddress Address, nodeId NodeID, delegationRewardAmount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectWithdrawDelegation", "blockNumber", blockNumber, "txHash", txHash.Hex(), "delegateAddress", delegateAddress.Bech32(), "nodeId", Bytes2Hex(nodeId[:]), "delegationRewardAmount", delegationRewardAmount)
		amt := new(big.Int).Set(delegationRewardAmount)
		exeBlockData.WithdrawDelegationList = append(exeBlockData.WithdrawDelegationList, &WithdrawDelegation{TxHash: txHash, DelegateAddress: delegateAddress, NodeID: nodeId, RewardAmount: amt})
	}
}

func CollectActiveVersion(blockNumber uint64, newVersion uint32) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectActiveVersion", "blockNumber", blockNumber, "newVersion", newVersion)
		exeBlockData.ActiveVersion = FormatVersion(newVersion)
	}
}

func CollectAutoStakingTx(blockNumber uint64, txHash Hash, restrictingAmount *big.Int, balanceAmount *big.Int) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectAutoStakingTx", "blockNumber", blockNumber, "txHash", txHash.Hex(), "restrictingAmount", restrictingAmount.String(), "balanceAmount", balanceAmount.String())
		exeBlockData.AutoStakingMap[txHash] = &AutoStakingTx{RestrictingAmount: restrictingAmount, BalanceAmount: balanceAmount}
	}
}

func CollectEpochElection(blockNumber uint64, nodeIdList []NodeID) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectEpochElection", "blockNumber", blockNumber, "nodeIdList", nodeIdList)
		exeBlockData.EpochElection = nodeIdList
	}
}
func CollectConsensusElection(blockNumber uint64, nodeIdList []NodeID) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectConsensusElection", "blockNumber", blockNumber, "nodeIdList", nodeIdList)
		exeBlockData.ConsensusElection = nodeIdList
	}
}
