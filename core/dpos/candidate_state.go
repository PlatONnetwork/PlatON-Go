package depos

import (
	"fmt"
	"sync"
	"encoding/json"
	"Platon-go/common"
	"Platon-go/rlp"
	"Platon-go/log"
	"Platon-go/params"
	"Platon-go/core/state"
	"Platon-go/p2p/discover"
	"math/big"
	"errors"
	"Platon-go/core/vm"
	"Platon-go/core/types"
	"net"
	"strconv"
)


const(
	// 即时入围竞选人
	ImmediatePrefix 			= "id"
	ImmediateListPrefix 		= "iL"
	// 见证人
	WitnessPrefix 				= "wn"
	WitnessListPrefix 			= "wL"
	// 下一轮见证人
	NextWitnessPrefix 			= "Nwn"
	NextWitnessListPrefix		= "NwL"
	// 需要退款的
	DefeatPrefix 				= "df"
	DefeatListPrefix 			= "dL"
)

var (
	// 即时入围竞选人
	ImmediateBtyePrefix 		= []byte(ImmediatePrefix)
	ImmediateListBtyePrefix 	= []byte(ImmediateListPrefix)
	// 见证人
	WitnessBtyePrefix 			= []byte(WitnessPrefix)
	WitnessListBtyePrefix 		= []byte(WitnessListPrefix)
	// 下一轮见证人
	NextWitnessBtyePrefix 		= []byte(NextWitnessPrefix)
	NextWitnessListBytePrefix 	= []byte(NextWitnessListPrefix)
	// 需要退款的
	DefeatBtyePrefix 			= []byte(DefeatPrefix)
	DefeatListBtyePrefix 		= []byte(DefeatListPrefix)



	CandidateEncodeErr    		= errors.New("Candidate encoding err")
	CandidateDecodeErr 			= errors.New("Cnadidate decoding err")
	WithdrawPriceErr 			= errors.New("Withdraw Price err")
	CandidateEmptyErr 			= errors.New("Candidate is empty")
	ContractBalanceNotEnoughErr = errors.New("Contract's balance is not enough")
	CandidateOwnerErr 			= errors.New("CandidateOwner Addr is illegal")
)

type CandidatePool struct {
	// 当前入围者数目
	//count 					uint64
	// 最大允许入围人数目
	maxCount				uint64
	// 最大允许见证人数目
	maxChair				uint64
	// 允许退款的块间隔
	RefundBlockNumber 		uint64

	// 本轮选出的见证人集
	originCandidates  		map[discover.NodeID]*types.Candidate
	// 下一轮见证人集
	nextOriginCandidates  	map[discover.NodeID]*types.Candidate
	// 即时的入选人集
	immediateCandates 		map[discover.NodeID]*types.Candidate
	// 质押失败的竞选人集 (退款用)
	defeatCandidates 		map[discover.NodeID][]*types.Candidate

	// cache
	candidateCacheArr 		[]*types.Candidate
	//candiateIds 			[]discover.NodeID
	lock 					*sync.RWMutex
}

var candidatePool *CandidatePool

// 初始化全局候选池对象
func NewCandidatePool(state vm.StateDB, configs *params.DposConfig, isgenesis bool) (*CandidatePool, error) {
//func NewCandidatePool(blockChain *core.BlockChain, configs *params.DposConfig) (*CandidatePool, error) {

	// 创世块的时候需要, 把配置的信息加载到stateDB
	if isgenesis {
		if err := loadConfig(configs, state); nil != err {
			log.Error("Failed to load config on NewCandidatePool", "err", err)
			return nil, err
		}
	}
	//var idArr []discover.NodeID
	//if ids, err := getImmediateIdsByState(state); nil != err {
	//	log.Error("Failed to decode immediateIds on NewCandidatePool", "err", err)
	//	return nil, err
	//}else {
	//	idArr = ids
	//}


	//var originMap, immediateMap map[discover.NodeID]*Candidate
	//var  defeatMap map[discover.NodeID][]*Candidate
	//// 非创世块，需要从db加载
	//if isgenesis {
	//	tr := state.StorageTrie(common.CandidateAddr)
	//	// 构造上一轮见证人的map
	//	var err error
	//	originMap, immediateMap, defeatMap, err = initCandidatesByTrie(tr)
	//	if nil != err {
	//		return nil, err
	//	}
	//}else {
	//	var err error
	//	originMap, immediateMap, err = buildByConfig(configs, state)
	//	if nil != err {
	//		return nil, err
	//	}
	//	defeatMap =  make(map[discover.NodeID][]*Candidate)
	//}

	candidatePool =  &CandidatePool{
		//count: 					uint64(len(idArr)),
		maxCount:				configs.MaxCount,
		maxChair:				configs.MaxChair,
		RefundBlockNumber: 		configs.RefundBlockNumber,
		originCandidates: 		make(map[discover.NodeID]*types.Candidate, 0),
		immediateCandates: 		make(map[discover.NodeID]*types.Candidate, 0),
		defeatCandidates: 		make(map[discover.NodeID][]*types.Candidate, 0),
		lock: 					&sync.RWMutex{},
	}
	return candidatePool, nil
}

func loadConfig(configs *params.DposConfig, state vm.StateDB) error {
	if len(configs.Candidates) != 0 {
		// 如果配置过多，只取前面的
		if len(configs.Candidates) > int(configs.MaxCount) {
			configs.Candidates = (configs.Candidates)[:configs.MaxCount]
		}

		// can cache
		witcans 		:=	make([]*types.Candidate, 0)
		immcans 		:=  make([]*types.Candidate, 0)

		//witnessMap := make(map[discover.NodeID]*types.Candidate, 0)
		//immediateMap := make(map[discover.NodeID]*types.Candidate, 0)

		for i, canConfig := range configs.Candidates {
			can := &types.Candidate{
				Deposit:			canConfig.Deposit,
				BlockNumber: 		canConfig.BlockNumber,
				TxIndex: 		 	canConfig.TxIndex,
				CandidateId: 	 	canConfig.CandidateId,
				Host: 			 	canConfig.Host,
				Port: 			 	canConfig.Port,
				Owner: 				canConfig.Owner,
				From: 				canConfig.From,
			}
			// 详情追加到树中
			if val, err := rlp.EncodeToBytes(can); nil == err {
				// 追加见证人信息
				if uint64(i) < configs.MaxChair {
					fmt.Println("设置进去WitnessId = ", can.CandidateId.String())
					//state.SetState(common.CandidateAddr, WitnessKey(can.CandidateId), val)
					setWitnessState(state, can.CandidateId, val)
					witcans = append(witcans, can)
				}
				fmt.Println("设置进去ImmediateId = ", can.CandidateId.String())
				// 追加入围人信息
				//state.SetState(common.CandidateAddr,  ImmediateKey(can.CandidateId), val)
				setImmediateState(state, can.CandidateId, val)
				immcans = append(immcans, can)
			}else {
				log.Error("Failed to encode candidate object on loadConfig", "key", string(WitnessKey(can.CandidateId)), "err", err)
				return err
			}
		}
		// 索引排序
		candidateSort(witcans)
		candidateSort(immcans)

		// id cache
		witnessIds 		:= 	make([]discover.NodeID, 0)
		immediateIds 	:= 	make([]discover.NodeID, 0)

		for _, can := range witcans {
			witnessIds = append(witnessIds, can.CandidateId)
		}
		for _, can := range immcans {
			immediateIds = append(immediateIds, can.CandidateId)
		}

		// 索引上树
		if arrVal, err := rlp.EncodeToBytes(witnessIds); nil == err {
			setWitnessIdsState(state, arrVal)
		}else {
			log.Error("Failed to encode witnessIds on loadConfig", "err", err)
			return err
		}

		if arrVal, err := rlp.EncodeToBytes(immediateIds); nil == err {
			setImmediateIdsState(state, arrVal)
		}else {
			log.Error("Failed to encode immediateIds on loadConfig", "err", err)
			return err
		}
	}
	return nil
}

func (c *CandidatePool) initDataByState (state vm.StateDB) error {

	// 加载 见证人信息
	var witnessIds []discover.NodeID
	c.originCandidates = make(map[discover.NodeID]*types.Candidate, 0)
	if ids, err := getWitnessIdsByState(state); nil != err {
		log.Error("Failed to decode witnessIds on initDataByState", "err", err)
		return err
	}else {
		witnessIds = ids
	}
	PrintObject("witnessIds = ", witnessIds)
	for _, witnessId := range witnessIds {
		fmt.Println("witnessId = ", witnessId.String())
		var can *types.Candidate
		if c, err := getWitnessByState(state, witnessId); nil != err {
			log.Error("Failed to decode Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		}else {
			can = c
		}
		c.originCandidates[witnessId] = can
	}
	// 加载 下一轮见证人
	var nextWitnessIds []discover.NodeID
	c.nextOriginCandidates = make(map[discover.NodeID]*types.Candidate, 0)
	if ids, err := getNextWitnessIdsByState(state); nil != err {
		log.Error("Failed to decode nextWitnessIds on initDataByState", "err", err)
		return err
	}else {
		nextWitnessIds = ids
	}
	PrintObject("nextWitnessIds = ", nextWitnessIds)
	for _, witnessId := range nextWitnessIds {
		fmt.Println("nextwitnessId = ", witnessId.String())
		var can *types.Candidate
		if c, err := getNextWitnessByState(state, witnessId); nil != err {
			log.Error("Failed to decode Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		}else {
			can = c
		}
		c.nextOriginCandidates[witnessId] = can
	}
	// 加载 入围者
	var immediateIds []discover.NodeID
	c.immediateCandates = make(map[discover.NodeID]*types.Candidate, 0)
	if ids, err := getImmediateIdsByState(state); nil != err {
		log.Error("Failed to decode immediateIds on initDataByState", "err", err)
		return err
	}else {
		immediateIds = ids
	}

	// cache
	canCache := make([]*types.Candidate, 0)

	PrintObject("immediateIds = ", immediateIds)
	for _, immediateId := range immediateIds {
		fmt.Println("immediateId = ", immediateId.String())
		var can *types.Candidate
		if c, err := getImmediateByState(state, immediateId); nil != err {
			log.Error("Failed to decode Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		}else {
			can = c
		}
		c.immediateCandates[immediateId] = can
		canCache = append(canCache, can)
	}
	c.candidateCacheArr = canCache

	// 加载 需要退款信息
	var defeatIds []discover.NodeID
	c.defeatCandidates = make(map[discover.NodeID][]*types.Candidate, 0)
	if ids, err := getDefeatIdsByState(state); nil != err {
		log.Error("Failed to decode defeatIds on initDataByState", "err", err)
		return err
	}else {
		defeatIds = ids
	}
	PrintObject("defeatIds = ", defeatIds)
	for _, defeatId := range defeatIds {
		fmt.Println("defeatId = ", defeatId.String())
		var canArr []*types.Candidate
		if arr, err := getDefeatsByState(state, defeatId); nil != err {
			log.Error("Failed to decode CandidateArr on initDataByState", "err", err)
			return CandidateDecodeErr
		}else {
			canArr = arr
		}
		c.defeatCandidates[defeatId] = canArr
	}
	return nil
}




// 候选人抵押
func(c *CandidatePool) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	PrintObject("can:", can)
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on SetCandidate err", err)
		return err
	}
	// 先追加到 缓存数组中然后做排序
	if len(c.immediateCandates) != 0 && len(c.candidateCacheArr) == 0 {
		for _, v := range c.immediateCandates {
			c.candidateCacheArr = append(c.candidateCacheArr, v)
		}
	}

	// cache map
	//cache := c.immediateCandates

	// 判断当前质押人是否新来的
	// 只有是新来的，才会追加数组中
	var needSort bool
	if _, ok := c.immediateCandates[can.CandidateId]; !ok {
		c.candidateCacheArr = append(c.candidateCacheArr, can)
		needSort = true
	}

	c.immediateCandates[can.CandidateId] = can
	PrintObject("immediateMap:", c.immediateCandates)
	// 排序
	if needSort {
		candidateSort(c.candidateCacheArr)
	}
	// 把多余入围者移入落榜名单
	if len(c.candidateCacheArr) > int(c.maxCount) {
		// 截取出当前落榜的候选人
		tmpArr := (c.candidateCacheArr)[c.maxCount:]
		// 保留入围候选人
		c.candidateCacheArr = (c.candidateCacheArr)[:c.maxCount]
		// 处理落选人
		for _, tmpCan := range tmpArr {

			//var isDEL bool
			//if _, ok := cache[tmpCan.CandidateId]; ok {
			//	isDEL = true
			//}
			// 删除trie中的 入围者信息
			if err := c.delImmediate(state, tmpCan.CandidateId); nil != err {
				return err
			}
			// 追加到落榜集
			if err := c.setDefeat(state, tmpCan.CandidateId, tmpCan); nil != err {
				return err
			}
		}
		// 更新落选人索引
		if err := c.setDefeatIndex(state); nil != err {
			return err
		}
	}

	// cache id
	sortIds := make([]discover.NodeID, 0)

	// 入围者上树
	for _, can := range c.candidateCacheArr {
		//var isADD bool
		//if _, ok := cache[can.CandidateId]; !ok {
		//	isADD = true
		//}
		c.setImmediate(state, can.CandidateId, can)
		sortIds = append(sortIds, can.CandidateId)
	}
	// 更新入围者索引
	c.setImmediateIndex(state, sortIds)
	return nil
}

// 获取入围候选人信息
func (c *CandidatePool) GetCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error) {
	return c.getCandidate(state, nodeId)
}

// 候选人退出
func (c *CandidatePool) WithdrawCandidate (state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on WithdrawCandidate err", err)
		return err
	}

	if price.Cmp(new(big.Int).SetUint64(0)) <= 0 {
		log.Error("withdraw failed price invalid, price", price.String())
		return WithdrawPriceErr
	}
	can, ok := c.immediateCandates[nodeId]
	if !ok {
		log.Error("withdraw failed current Candidate is empty")
		return CandidateEmptyErr
	}
	if nil == can {
		log.Error("withdraw failed current Candidate is empty")
		return CandidateEmptyErr
	}
	// 判断退押 金额
	if can.Deposit.Cmp(price) < 0 {
		log.Error("withdraw failed refund price must less or equal deposit", "key", nodeId.String())
		return WithdrawPriceErr
	}else if can.Deposit.Cmp(price) == 0 { // 全额退出质押
		// 删除入围者信息
		if err := c.delImmediate(state, nodeId); nil != err {
			return err
		}
		// 追加到落选
		if err := c.setDefeat(state, nodeId, can); nil != err {
			return err
		}
		// 更新落榜索引
		if err := c.setDefeatIndex(state); nil != err {
			return err
		}
	}else { // 只退了一部分, 需要重新对入围者排序
		// 剩下部分
		canNew := &types.Candidate{
			Deposit:		new(big.Int).Sub(can.Deposit, price),
			BlockNumber: 	can.BlockNumber,
			TxIndex: 		can.TxIndex,
			CandidateId: 	can.CandidateId,
			Host: 			can.Host,
			Port: 			can.Port,
			Owner: 			can.Owner,
			From: 			can.From,
		}
		c.immediateCandates[nodeId] = canNew

		// 更新剩余部分
		if err := c.setImmediate(state, nodeId, canNew); nil != err {
			return err
		}
		// 退款部分新建退款信息
		canDefeat := &types.Candidate{
			Deposit: 		price,
			BlockNumber: 	blockNumber,
			TxIndex: 		can.TxIndex,
			CandidateId: 	can.CandidateId,
			Host: 			can.Host,
			Port: 			can.Port,
			Owner: 			can.Owner,
			From: 			can.From,
		}
		// 退出部分
		if err := c.setDefeat(state, nodeId, canDefeat); nil != err {
			return err
		}
		//更新退款索引
		if err := c.setDefeatIndex(state); nil != err {
			return err
		}
	}
	return nil
}

// 获取实时所有入围候选人
func (c *CandidatePool) GetChosens (state vm.StateDB) []*types.Candidate {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on WithdrawCandidate err", err)
		return nil
	}
	immediateIds, err := c.getImmediateIndex(state)
	if nil != err {
			log.Error("Failed to getImmediateIndex err", err)
		return nil
	}
	arr := make([]*types.Candidate, 0)
	for _, id := range immediateIds {
		arr = append(arr, c.immediateCandates[id])
	}
	return arr
}

// 获取所有见证人
func (c *CandidatePool) GetChairpersons (state vm.StateDB) []*types.Candidate {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on GetChairpersons err", err)
		return nil
	}
	witnessIds, err := c.getWitnessIndex(state)
	if nil != err {
		log.Error("Failed to getWitnessIndex on GetChairpersonserr", err)
		return nil
	}
	arr := make([]*types.Candidate, 0)
	for _, id := range witnessIds {
		arr = append(arr, c.originCandidates[id])
	}
	return arr
}

// 获取退款信息
func (c *CandidatePool) GetDefeat(state vm.StateDB, nodeId discover.NodeID) ([]*types.Candidate, error){
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on GetDefeat err", err)
		return nil, err
	}

	defeat, ok := c.defeatCandidates[nodeId]
	if !ok {
		log.Error("Candidate is empty")
		return nil, nil
	}
	return defeat, nil
}

// 判断是否落榜
func (c *CandidatePool) IsDefeat (state vm.StateDB, nodeId discover.NodeID) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on IsDefeat err", err)
		return false, err
	}

	if _, ok := c.immediateCandates[nodeId]; !ok {
		log.Error("Candidate is empty")
		return false, nil
	}
	return true, nil
}

// 根据nodeId查询 质押信息中的 受益者地址
func (c *CandidatePool) GetOwner (state vm.StateDB, nodeId discover.NodeID) common.Address {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on GetOwner err", err)
		return common.Address{}
	}
	// 先查见证人
	var addr common.Address
	or_can, or_ok := c.originCandidates[nodeId]
	im_can, im_ok := c.immediateCandates[nodeId]
	canArr, de_ok := c.defeatCandidates[nodeId]
	if or_ok {
		addr = or_can.Owner
		return addr
	}
	if im_ok {
		addr = im_can.Owner
		return addr
	}
	if de_ok {
		if len(canArr) != 0 {
			addr = canArr[0].Owner
			return addr
		}
	}
	return common.Address{}
}

// 一键提款
func (c *CandidatePool) RefundBalance (state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on RefundBalance err", err)
		return err
	}

	var canArr []*types.Candidate
	if defeatArr, ok := c.defeatCandidates[nodeId]; ok {
		canArr = defeatArr
	}else {
		log.Error("Failed to refundbalance cnadidate is empty")
		return CandidateDecodeErr
	}
	// cache
	// 用来做校验用，即正常情况应该每个nodeId的质押退款信息中的收益者(owner)应该为同一个
	var addr common.Address
	// 累计需要一次性退款的金额
	var amount uint64
	// 中转需要删除的退款信息
	delCanArr := make([]*types.Candidate, 0)

	//contractBalance := state.GetBalance(common.CandidateAddr)
	//currentNum := new(big.Int).SetUint64(blockNumber)

	// 遍历该nodeId下的所有 退款信息
	for index, can := range canArr {
		sub := new(big.Int).Sub(blockNumber, can.BlockNumber)
		fmt.Println("当前块高:", blockNumber.String(), "质押块高:", can.BlockNumber.String(), "相差:", sub.String())
		fmt.Println("当前nodeId:", can.CandidateId.String())
		if sub.Cmp(new(big.Int).SetUint64(c.RefundBlockNumber)) >= 0 { // 允许退款
			delCanArr = append(delCanArr, can)
			canArr = append(canArr[:index], canArr[index+1:]...)
			// 累加一次性退款金额
			amount += can.Deposit.Uint64()
		}else {
			fmt.Println("块高不匹配，不给予退款...")
			continue
		}

		if addr == common.ZeroAddr {
			addr = can.Owner
		}else {
			if addr != can.Owner {
				log.Info("Failed to refundbalance 发现抵押节点nodeId下有不同受益者地址", "nodeId", nodeId.String(), "addr1", addr.String(), "addr2", can.Owner)
				if len(canArr) != 0 {
					canArr = append(delCanArr, canArr...)
				}else {
					canArr = delCanArr
				}
				c.defeatCandidates[nodeId] = canArr
				fmt.Println("Failed to refundbalance 发现抵押节点nodeId下有不同受益者地址", "nodeId", nodeId.String(), "addr1", addr.String(), "addr2", can.Owner)
				return CandidateOwnerErr
			}
		}

		//if (contractBalance.Cmp(new(big.Int).SetUint64(amount))) < 0 {
		//	log.Error("Failed to refundbalance constract account insufficient balance ", state.GetBalance(common.CandidateAddr).String(), "amount", amount)
		//	if len(arr) != 0 {
		//		arr = append(delCanArr, arr...)
		//	}else {
		//		arr = delCanArr
		//	}
		//	c.defeatCandidates[nodeId] = arr
		//	return ContractBalanceNotEnoughErr
		//}
	}

	// 统一更新树
	if len(canArr) == 0 {
		//state.SetState(common.CandidateAddr, DefeatKey(nodeId), []byte{})
		//delete(c.defeatCandidates, nodeId)
		if err := c.delDefeat(state, nodeId); nil != err {
			log.Error("RefundBalance failed to delDefeat err", err)
			return err
		}
	}else {
		// 如果没被退完, 更新剩下的
		if arrVal, err := rlp.EncodeToBytes(canArr); nil != err {
			log.Error("Failed to encode candidate object on RefundBalance", "key", nodeId.String(), "err", err)
			canArr = append(delCanArr, canArr...)
			c.defeatCandidates[nodeId] = canArr
			return CandidateDecodeErr
		}else {
			// 更新退款详情
			setDefeatState(state, nodeId, arrVal)
			// 把没退完的设置回 map
			c.defeatCandidates[nodeId] = canArr
		}
	}
	// 扣减合约余额
	state.SubBalance(common.CandidateAddr, new(big.Int).SetUint64(amount))
	// 增加收益账户余额
	state.AddBalance(addr, new(big.Int).SetUint64(amount))
	return nil
}



// 揭榜见证人
func (c *CandidatePool) Election(state *state.StateDB) bool{
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on Election err", err)
		return false
	}
	//immediateIds, err := c.getImmediateIndex(state)
	//if nil != err {
	//	log.Error("Failed to getImmediateIndex on Election err", err)
	//	return false
	//}

	// 对当前所有入围者排序
	candidateSort(c.candidateCacheArr)

	// cache ids
	immediateIds := make([]discover.NodeID, 0)
	for _, can := range c.candidateCacheArr {
		immediateIds = append(immediateIds, can.CandidateId)
	}

	// 缓存前面一定数量的见证人
	var nextWitIds []discover.NodeID
	// 如果入选人数不超过见证人数
	if len(immediateIds) <= int(c.maxChair) {
		nextWitIds = make([]discover.NodeID, len(immediateIds))
		copy(nextWitIds, immediateIds)

	} else {
		// 入选人数超过了见证人数，提取前N名
		nextWitIds = make([]discover.NodeID, c.maxChair)
		copy(nextWitIds, immediateIds)
	}

	// cache map
	nextWits := make(map[discover.NodeID]*types.Candidate, 0)


	// copy 见证人信息
	copyCandidateMapByIds(nextWits, c.immediateCandates, nextWitIds)

	// 清空所有旧的
	for nodeId, _ := range c.nextOriginCandidates {
		if err := c.delNextWitness(state, nodeId); nil != err {
			log.Error("failed to delNextWitness on election err", err)
			return false
		}
	}
	// 设置新的 can
	for nodeId, can := range nextWits {
		if err := c.setNextWitness(state, nodeId, can); nil != err {
			log.Error("failed to setNextWitness on election err", err)
			return false
		}
	}
	// 更新索引
	if err := c.setNextWitnessIndex(state, nextWitIds); nil != err {
		log.Error("failed to setNextWitnessIndex on election err", err)
		return false
	}
	// 将新的揭榜后的见证人信息，置入下轮见证人集
	c.nextOriginCandidates = nextWits
	return len(c.nextOriginCandidates) > 0
}

// 触发替换下轮见证人列表
func (c *CandidatePool) Switch(state *state.StateDB) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on Switch err", err)
		return false
	}
	// 清空所有旧的
	for nodeId,_ := range c.originCandidates {
		if err := c.delWitness(state, nodeId); nil != err {
			log.Error("Failed to delWitness on Switch err", err)
			return false
		}
	}

	// cache can
	caches := make([]*types.Candidate, 0)

	// 设置新的
	for nodeId, can := range c.nextOriginCandidates {
		if err := c.setWitness(state, nodeId, can); nil != err {
			log.Error("Failed to setWitness on Switch err", err)
			return false
		}
		caches = append(caches, can)
	}
	// 排序
	candidateSort(caches)

	// cache ids
	canIds := make([]discover.NodeID, 0)

	for _, can := range caches {
		canIds = append(canIds, can.CandidateId)
	}

	// 替换新索引
	if err := c.setWitnessindex(state, canIds); nil != err {
		log.Error("Failed to setWitnessindex on Switch err", err)
		return false
	}
	c.originCandidates = c.nextOriginCandidates
	c.nextOriginCandidates = make(map[discover.NodeID]*types.Candidate)
	return true
}

// 获取见证人节点列表
func (c *CandidatePool) GetWitness (state *state.StateDB) []*discover.Node {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on GetWitness err", err)
		return nil
	}
	witnessIds, err := c.getWitnessIndex(state)
	if nil != err {
		log.Error("Failed to getWitnessIndex on GetWitness err", err)
		return nil
	}
	arr := make([]*discover.Node, 0)
	for _, id := range witnessIds {
		can := c.originCandidates[id]
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build Node on GetWitness err", err, "nodeId", can.CandidateId.String())
			continue
		}else {
			arr = append(arr, node)
		}
	}
	return arr
}





func (c *CandidatePool) setImmediate(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate/*, isADD bool*/) error {
	c.immediateCandates[candidateId] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setImmediate", "key", candidateId.String(), "err", err)
		return CandidateEncodeErr
	}else {
		// 实时的入围候选人 input trie
		setImmediateState(state, candidateId, value)
		//if isADD {
		//	c.count ++
		//}
	}
	return nil
}

func (c *CandidatePool) getImmediateIndex (state vm.StateDB) ([]discover.NodeID, error) {
	return getImmediateIdsByState(state)
}

// 删除自动更新索引
func (c *CandidatePool) delImmediate (state vm.StateDB, candidateId discover.NodeID/*, isDEL bool*/) error {

	// trie 中删掉实时信息
	setImmediateState(state, candidateId, []byte{})
	// map 中删掉
	delete(c.immediateCandates, candidateId)
	//if isDEL {
	//	// map 中删掉
	//	delete(c.immediateCandates, candidateId)
	//	c.count --
	//}
	// 删除索引中的对应id
	var canIds []discover.NodeID
	if ids, err := getImmediateIdsByState(state); nil != err {
		log.Error("Failed to decode ImmediateIds err", err)
		return err
	}else {
		canIds = ids
	}

	var flag bool
	for i, id := range canIds {
		if id == candidateId {
			flag = true
			canIds = append(canIds[:i], canIds[i+1:]...)
		}
	}
	if flag {
		if val, err := rlp.EncodeToBytes(canIds); nil != err {
			log.Error("Failed to encode ImmediateIds err", err)
			return err
		}else {
			setImmediateIdsState(state, val)
		}
	}
	return nil
}

func (c *CandidatePool) setImmediateIndex (state vm.StateDB, nodeIds []discover.NodeID) error {
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode ImmediateIds err", err)
		return err
	}else {
		setImmediateIdsState(state, val)
	}
	return nil
}

// 设置退款信息
func (c *CandidatePool) setDefeat(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate) error {

	var defeatArr []*types.Candidate
	// 追加退款信息
	if defeatArrTmp, ok := c.defeatCandidates[can.CandidateId]; ok {
		defeatArrTmp = append(defeatArrTmp, can)
		c.defeatCandidates[can.CandidateId] = defeatArrTmp
		defeatArr = defeatArrTmp
	}else {
		defeatArrTmp = make([]*types.Candidate, 0)
		defeatArrTmp = append(defeatArr, can)
		c.defeatCandidates[can.CandidateId] = defeatArrTmp
		defeatArr = defeatArrTmp
	}
	// trie 中添加 退款信息
	if value ,err := rlp.EncodeToBytes(&defeatArr); nil != err {
		log.Error("Failed to encode candidate object on setDefeat", "key", candidateId.String(), "err", err)
		return CandidateEncodeErr
	}else {
		setDefeatState(state, candidateId, value)
	}
	return nil
}

func (c *CandidatePool) delDefeat(state vm.StateDB, nodeId discover.NodeID) error {
	delete(c.defeatCandidates, nodeId)
	setDefeatState(state, nodeId, []byte{})

	// 删除索引中的对应id
	var canIds []discover.NodeID
	if ids, err := getDefeatIdsByState(state); nil != err {
		log.Error("Failed to decode DefeatIds err", err)
		return err
	}else {
		canIds = ids
	}

	var flag bool
	for i, id := range canIds {
		if id == nodeId {
			flag = true
			canIds = append(canIds[:i], canIds[i+1:]...)
		}
	}
	if flag {
		if val, err := rlp.EncodeToBytes(canIds); nil != err {
			log.Error("Failed to encode ImmediateIds err", err)
			return err
		}else {
			setDefeatIdsState(state, val)
		}
	}
	return nil
}

// 更新退款信息索引
func (c *CandidatePool)setDefeatIndex(state vm.StateDB) error {
	newdefeatIds := make([]discover.NodeID, 0)
	for id, _ := range c.defeatCandidates {
		newdefeatIds = append(newdefeatIds, id)
	}
	if value ,err := rlp.EncodeToBytes(&newdefeatIds); nil != err {
		log.Error("Failed to encode candidate object on setDefeatIds err", err)
		return CandidateEncodeErr
	}else {
		setDefeatIdsState(state, value)
	}
	return nil
}

func (c *CandidatePool) setWitness (state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	c.originCandidates[nodeId] = can
	if val, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode Candidate on setWitness err", err)
		return err
	}else {
		setWitnessState(state, nodeId, val)
	}
	return nil
}

func (c *CandidatePool) setWitnessindex(state vm.StateDB, nodeIds []discover.NodeID) error {
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode WitnessIds err", err)
		return err
	}else {
		setWitnessIdsState(state, val)
	}
	return nil
}

func (c *CandidatePool) delWitness (state vm.StateDB, candidateId discover.NodeID) error {
	// map 中删掉
	delete(c.originCandidates, candidateId)
	// trie 中删掉实时信息
	setWitnessState(state, candidateId, []byte{})
	// 删除索引中的对应id
	var canIds []discover.NodeID
	if ids, err := getWitnessIdsByState(state); nil != err {
		log.Error("Failed to decode ImmediateIds err", err)
		return err
	}else {
		canIds = ids
	}

	var flag bool
	for i, id := range canIds {
		if id == candidateId {
			flag = true
			canIds = append(canIds[:i], canIds[i+1:]...)
		}
	}
	if flag {
		if arrVal, err := rlp.EncodeToBytes(canIds); nil != err {
			log.Error("Failed to encode ImmediateIds err", err)
			return err
		}else {
			setWitnessIdsState(state, arrVal)
		}
	}
	return nil
}

func (c *CandidatePool) getWitnessIndex (state vm.StateDB) ([]discover.NodeID, error) {
	return getWitnessIdsByState(state)
}

func (c *CandidatePool) setNextWitness(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	c.nextOriginCandidates[nodeId] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setImmediate", "key", nodeId.String(), "err", err)
		return CandidateEncodeErr
	}else {
		// 实时的入围候选人 input trie
		setNextWitnessState(state, nodeId, value)
	}
	return nil
}

func (c *CandidatePool) delNextWitness (state vm.StateDB, candidateId discover.NodeID) error {
	// map 中删掉
	delete(c.nextOriginCandidates, candidateId)
	// trie 中删掉实时信息
	setNextWitnessState(state, candidateId, []byte{})

	// 获取原有索引
	var canIds []discover.NodeID
	if ids, err := getNextWitnessIdsByState(state); nil != err {
		log.Error("Failed to decode ImmediateIds err", err)
		return err
	}else {
		canIds = ids
	}

	// 删除索引中的对应id
	var flag bool
	for i, id := range canIds {
		if id == candidateId {
			flag = true
			canIds = append(canIds[:i], canIds[i+1:]...)
		}
	}
	if flag {
		if arrVal, err := rlp.EncodeToBytes(canIds); nil != err {
			log.Error("Failed to encode ImmediateIds err", err)
			return err
		}else {
			setNextWitnessIdsState(state, arrVal)
		}
	}
	return nil
}

func (c *CandidatePool) setNextWitnessIndex (state vm.StateDB, nodeIds []discover.NodeID) error {
	if value ,err := rlp.EncodeToBytes(&nodeIds); nil != err {
		log.Error("Failed to encode candidate object on setDefeatIds err", err)
		return CandidateEncodeErr
	}else {
		setNextWitnessIdsState(state, value)
	}
	return nil
}

func (c *CandidatePool) getCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error){
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to initDataByState on getCandidate err", err)
		return nil, err
	}
	if candidatePtr, ok := c.immediateCandates[nodeId]; ok {
		return candidatePtr, nil
	}
	return nil, nil
}




func getWitnessIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var witnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidateAddr, WitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &witnessIds); nil != err {
			return nil, err
		}
	}
	return witnessIds, nil
}

func setWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, WitnessListKey(), arrVal)
}

func getWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, WitnessKey(id)), &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func setWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, WitnessKey(id), val)
}

func getNextWitnessIdsByState(state vm.StateDB) ([]discover.NodeID, error){
	var nextWitnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidateAddr, NextWitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &nextWitnessIds); nil != err {
			return nil, err
		}
	}
	return nextWitnessIds, nil
}

func setNextWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, NextWitnessListKey(), arrVal)
}

func getNextWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, NextWitnessKey(id)), &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func setNextWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, NextWitnessKey(id), val)
}

func getImmediateIdsByState(state vm.StateDB) ([]discover.NodeID, error){
	var immediateIds []discover.NodeID
	if valByte := state.GetState(common.CandidateAddr, ImmediateListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &immediateIds); nil != err {
			return nil, err
		}
	}
	return immediateIds, nil
}

func setImmediateIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, ImmediateListKey(), arrVal)
}

func getImmediateByState (state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, ImmediateKey(id)), &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func setImmediateState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, ImmediateKey(id), val)
}

func getDefeatIdsByState (state vm.StateDB) ([]discover.NodeID, error) {
	var defeatIds []discover.NodeID
	if valBtye := state.GetState(common.CandidateAddr, DefeatListKey()); len(valBtye) != 0 {
		if err := rlp.DecodeBytes(valBtye, &defeatIds); nil != err {
			return nil, err
		}
	}
	return defeatIds, nil
}

func setDefeatIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, DefeatListKey(), arrVal)
}

func getDefeatsByState(state vm.StateDB, id discover.NodeID) ([]*types.Candidate, error) {
	var canArr []*types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, DefeatKey(id)), &canArr); nil != err {
		return nil, err
	}
	return canArr, nil
}

func setDefeatState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, DefeatKey(id), val)
}

func copyCandidateMapByIds(target, source map[discover.NodeID]*types.Candidate, ids []discover.NodeID){
	for _, id := range ids {
		target[id] = source[id]
	}
}

func GetCandidatePtr () *CandidatePool {
	return candidatePool
}


func PrintObject(s string, obj interface{}){
	objs, _ := json.Marshal(obj)
	fmt.Println(s, string(objs), "\n")
}


func buildWitnessNode(can *types.Candidate) (*discover.Node, error) {
	if nil == can {
		return nil, CandidateEmptyErr
	}
	ip := net.ParseIP(can.Host)
	// uint16
	var port uint16
	if portInt, err := strconv.Atoi(can.Port); nil != err {
		return nil, err
	}else {
		port = uint16(portInt)
	}
	return discover.NewNode(can.CandidateId, ip, port, port), nil
}

func compare(c, can *types.Candidate) int {
	PrintObject("c", c)
	PrintObject("can", can)
	// 质押金大的放前面
	if c.Deposit.Cmp(can.Deposit) > 0 {
		return 1
	}else if c.Deposit.Cmp(can.Deposit) == 0 {
		// 块高小的放前面
		if c.BlockNumber.Cmp(can.BlockNumber) > 0 {
			return -1
		}else if c.BlockNumber.Cmp(can.BlockNumber) == 0 {
			// tx index 小的放前面
			if c.TxIndex > can.TxIndex {
				return -1
			}else if c.TxIndex == can.TxIndex {
				return 0
			}else {
				return 1
			}
		}else {
			return 1
		}
	}else {
		return -1
	}
}
// 候选人排序
func candidateSort(arr []*types.Candidate) {
	if len(arr) <= 1 {
		return
	}
	quickSort(arr, 0, len(arr) - 1)
}
func quickSort (arr []*types.Candidate, left, right int)  {
	if left < right {
		pivot := partition(arr, left, right)
		quickSort(arr, left, pivot - 1)
		quickSort(arr, pivot + 1, right)
	}
}
func partition(arr []*types.Candidate, left, right int) int {
	for left < right {
		for left < right && compare(arr[left], arr[right]) >= 0 {
			right --
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			left ++
		}
		for left < right && compare(arr[left], arr[right]) >= 0 {
			left ++
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			right --
		}
	}
	return left
}

func ImmediateKey(nodeId discover.NodeID) []byte {
	//key, _ := rlp.EncodeToBytes(nodeId)
	return immediateKey(nodeId.Bytes())
}
func immediateKey(key []byte) []byte {
	return append(ImmediateBtyePrefix, key...)
}

func WitnessKey(nodeId discover.NodeID) []byte {
	//key, _ := rlp.EncodeToBytes(nodeId)
	return witnessKey(nodeId.Bytes())
}
func witnessKey(key []byte) []byte {
	return append(WitnessBtyePrefix, key...)
}

func NextWitnessKey(nodeId discover.NodeID) []byte {
	return nextWitnessKey(nodeId.Bytes())
}
func nextWitnessKey(key []byte) []byte {
	return append(NextWitnessBtyePrefix, key...)
}

func DefeatKey(nodeId discover.NodeID) []byte {
	//key, _ := rlp.EncodeToBytes(nodeId)
	return defeatKey(nodeId.Bytes())
}
func defeatKey(key []byte) []byte {
	return append(DefeatBtyePrefix, key...)
}

func ImmediateListKey() []byte {
	return ImmediateListBtyePrefix
}

func WitnessListKey () []byte {
	return WitnessListBtyePrefix
}

func NextWitnessListKey () []byte {
	return NextWitnessListBytePrefix
}

func DefeatListKey () []byte {
	return DefeatListBtyePrefix
}