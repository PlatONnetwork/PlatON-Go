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
	ImmediateBtyePrefix 		= []byte("id")
	ImmediateListBtyePrefix 	= []byte("iL")
	// 见证人
	WitnessBtyePrefix 			= []byte("wn")
	WitnessListBtyePrefix 		= []byte("wL")
	// 下一轮见证人
	NextWitnessBtyePrefix 		= []byte("Nwn")
	NextWitnessListBytePrefix 	= []byte("NwL")
	// 需要退款的
	DefeatBtyePrefix 			= []byte("df")
	DefeatListBtyePrefix 		= []byte("dL")

	// 内置候选池合约账户
	CandidateAddr 				= common.HexToAddress("0x1....10")



	CandidateEncodeErr    		= errors.New("Candidate encoding err")
	CandidateDecodeErr 			= errors.New("Cnadidate decoding err")
	WithdrawPriceErr 			= errors.New("Withdraw Price err")
	CandidateEmptyErr 			= errors.New("Candidate is empty")
	ContractBalanceNotEnoughErr = errors.New("Contract's balance is not enough")
	CandidateOwnerErr 			= errors.New("CandidateOwner Addr is illegal")
)

type CandidatePool struct {
	// 当前入围者数目
	count 					uint64
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
	//blockChain     			*core.BlockChain
	//cacheState 				*state.StateDB

	// cache
	//candidateCacheArr 		[]*Candidate
	candiateIds 			[]discover.NodeID
	lock 					*sync.RWMutex
}

// 初始化全局候选池对象
func NewCandidatePool(state *state.StateDB, configs *params.DposConfig, isgenesis bool) (*CandidatePool, error) {
//func NewCandidatePool(blockChain *core.BlockChain, configs *params.DposConfig) (*CandidatePool, error) {

	// 创世块的时候需要, 把配置的信息加载到stateDB
	if isgenesis {
		if err := loadConfig(configs, state); nil != err {
			return nil, err
		}
	}
	var idArr []discover.NodeID
	if valByte := state.GetState(CandidateAddr, ImmediateListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &idArr); nil != err {
			log.Error("Failed to decode immediateIds", "err", err)
			return nil, err
		}
	}


	//var originMap, immediateMap map[discover.NodeID]*Candidate
	//var  defeatMap map[discover.NodeID][]*Candidate
	//// 非创世块，需要从db加载
	//if isgenesis {
	//	tr := state.StorageTrie(CandidateAddr)
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

	return &CandidatePool{
		count: 					uint64(len(idArr)),
		maxCount:				configs.MaxCount,
		maxChair:				configs.MaxChair,
		RefundBlockNumber: 		configs.RefundBlockNumber,
		originCandidates: 		make(map[discover.NodeID]*types.Candidate, 0),
		immediateCandates: 		make(map[discover.NodeID]*types.Candidate, 0),
		defeatCandidates: 		make(map[discover.NodeID][]*types.Candidate, 0),
		lock: 					&sync.RWMutex{},
	}, nil
}
//// 根据配置文件构建 dpos原始见证人
//func buildByConfig(configs *params.DposConfig, state *state.StateDB) (map[discover.NodeID]*Candidate, map[discover.NodeID]*Candidate, error){
//	originMap, immediateMap := make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0)
//	if len(configs.Candidates) != 0 {
//		// 如果配置过多，只取前面的
//		if len(configs.Candidates) > int(configs.MaxCount) {
//			configs.Candidates = (configs.Candidates)[:configs.MaxCount]
//		}
//		for i, canConfig := range configs.Candidates {
//			can := &Candidate{
//				Deposit:			canConfig.Deposit,
//				BlockNumber: 		canConfig.BlockNumber,
//				TxIndex: 		 	canConfig.TxIndex,
//				CandidateId: 	 	canConfig.CandidateId,
//				Host: 			 	canConfig.Host,
//				Port: 			 	canConfig.Port,
//				Owner: 				canConfig.Owner,
//				From: 				canConfig.From,
//			}
//
//			// 追加到树中
//			if val, err := rlp.EncodeToBytes(can); nil == err {
//				if uint64(i) < configs.MaxChair {
//					state.SetState(CandidateAddr, WitnessKey(can.CandidateId), val)
//					originMap[can.CandidateId] = can
//				}
//				state.SetState(CandidateAddr,  ImmediateKey(can.CandidateId), val)
//				immediateMap[can.CandidateId] = can
//			}else {
//				log.Error("Failed to encode candidate object", "key", string(WitnessKey(can.CandidateId)), "err", err)
//				return make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0), err
//			}
//		}
//	}
//	return originMap, immediateMap, nil
//}

func loadConfig(configs *params.DposConfig, state *state.StateDB) error {
	if len(configs.Candidates) != 0 {
		// 如果配置过多，只取前面的
		if len(configs.Candidates) > int(configs.MaxCount) {
			configs.Candidates = (configs.Candidates)[:configs.MaxCount]
		}

		// id cache
		witnessIds 	:= make([]discover.NodeID, 0)
		immediateIds := make([]discover.NodeID, 0)

		witnessMap := make(map[discover.NodeID]*types.Candidate, 0)
		immediateMap := make(map[discover.NodeID]*types.Candidate, 0)

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
					state.SetState(CandidateAddr, WitnessKey(can.CandidateId), val)
					witnessIds = append(witnessIds, can.CandidateId)
					witnessMap[can.CandidateId] = can
				}
				fmt.Println("设置进去ImmediateId = ", can.CandidateId.String())
				// 追加入围人信息
				state.SetState(CandidateAddr,  ImmediateKey(can.CandidateId), val)
				immediateIds = append(immediateIds, can.CandidateId)
				immediateMap[can.CandidateId] = can
			}else {
				log.Error("Failed to encode candidate object", "key", string(WitnessKey(can.CandidateId)), "err", err)
				return err
			}
		}
		// 索引排序
		candidateSort(witnessIds, witnessMap)
		candidateSort(immediateIds, immediateMap)
		// 索引上树
		if arrVal, err := rlp.EncodeToBytes(witnessIds); nil == err {
			state.SetState(CandidateAddr, WitnessListKey(), arrVal)
		}else {
			log.Error("Failed to encode immediateIds", "err", err)
			return err
		}

		if arrVal, err := rlp.EncodeToBytes(immediateIds); nil == err {
			state.SetState(CandidateAddr, ImmediateListKey(), arrVal)
		}else {
			log.Error("Failed to encode witnessIds", "err", err)
			return err
		}
	}
	return nil
}
//
//func (c *CandidatePool) CommitTrie (deleteEmptyObjects bool) (root common.Hash, err error) {
//	return c.cacheState.Commit(deleteEmptyObjects)
//}

//func (c *CandidatePool)IteratorTrie(s string){
//	iteratorTrie(s, c.cacheState.StorageTrie(CandidateAddr))
//}

func iteratorTrie(s string, tr state.Trie){
	it := tr.NodeIterator(nil)
	for it.Next(true) {
		if it.Leaf() {
			var a types.Candidate
			rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &a)
			fmt.Println(s, string(tr.GetKey(it.LeafKey())), "== ", &a)
		}
	}
}

func printObject(s string, obj interface{}){
	objs, _ := json.Marshal(obj)
	fmt.Println(s, string(objs), "\n")
}

//func (c *CandidatePool)buildCandidatesByTrie(prefix string) map[common.Hash]*Candidate {
//	tr := c.cacheState.StorageTrie(CandidateAddr)
//	it := tr.NodeIterator(nil)
//	candidates := make(map[common.Hash]*Candidate, 0)
//	for it.Next(true) {
//		if it.Leaf() {
//			trieKey := tr.GetKey(it.LeafKey())
//			cleanKey := trieKey[len([]byte(CandidateAddr.String())):]
//
//			// 判断前缀
//			if strings.HasPrefix(string(cleanKey), prefix){
//				key := common.BytesToHash(cleanKey[len([]byte(prefix)):])
//				var candidate Candidate
//				rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &candidate)
//				//rlp.DecodeBytes(it.LeafBlob(), candidate)
//				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), candidate)
//				candidates[key] = &candidate
//			}
//		}
//	}
//	return candidates
//}

//func (c *CandidatePool)buildCandidateArrByTrie(prefix string) map[common.Hash][]*Candidate {
//	tr := c.cacheState.StorageTrie(CandidateAddr)
//	it := tr.NodeIterator(nil)
//	candidates := make(map[common.Hash][]*Candidate, 0)
//	for it.Next(true) {
//		if it.Leaf() {
//			trieKey := tr.GetKey(it.LeafKey())
//			cleanKey := trieKey[len([]byte(CandidateAddr.String())):]
//
//			// 判断前缀
//			if strings.HasPrefix(string(cleanKey), prefix){
//				key := common.BytesToHash(cleanKey[len([]byte(prefix)):])
//				var arr []*Candidate
//				rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &arr)
//				//rlp.DecodeBytes(it.LeafBlob(), candidate)
//				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), arr)
//				candidates[key] = arr
//			}
//		}
//	}
//	return candidates
//}

//func initCandidatesByTrie (tr state.Trie) (map[discover.NodeID]*Candidate, map[discover.NodeID]*Candidate, map[discover.NodeID][]*Candidate, error){
//	it := tr.NodeIterator(nil)
//	// 见证人
//	originCandidates  := make(map[discover.NodeID]*Candidate, 0)
//	// 即时入围者
//	immediateCandates := make(map[discover.NodeID]*Candidate, 0)
//	// 需要退款信息
//	defeatCandidates  := make(map[discover.NodeID][]*Candidate, 0)
//	for it.Next(true) {
//		if it.Leaf() {
//			trieKey := tr.GetKey(it.LeafKey())
//			cleanKey := trieKey[len([]byte(CandidateAddr.String())):]
//
//			// 根据前缀获取 见证人信息
//			if strings.HasPrefix(string(cleanKey), WitnessPrefix){
//				key := discover.MustBytesID(cleanKey[len([]byte(WitnessPrefix)):])
//				//key := common.BytesToHash(cleanKey[len([]byte(WitnessPrefix)):])
//				var candidate Candidate
//				if err := rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &candidate); nil != err {
//					return make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID][]*Candidate, 0), err
//				}
//				//rlp.DecodeBytes(it.LeafBlob(), candidate)
//				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), candidate)
//				originCandidates[key] = &candidate
//			}
//
//			// 根据前缀获取 入围竞选人信息
//			if strings.HasPrefix(string(cleanKey), ImmediatePrefix){
//				key := discover.MustBytesID(cleanKey[len([]byte(ImmediatePrefix)):])
//				//key := common.BytesToHash(cleanKey[len([]byte(ImmediatePrefix)):])
//				var candidate Candidate
//				if err := rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &candidate); nil != err {
//					return make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID][]*Candidate, 0), err
//				}
//				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), candidate)
//				immediateCandates[key] = &candidate
//			}
//
//			// 根据前缀获取 落榜需要退款信息
//			if strings.HasPrefix(string(cleanKey), DefeatPrefix){
//				key := discover.MustBytesID(cleanKey[len([]byte(DefeatPrefix)):])
//				//key := common.BytesToHash(cleanKey[len([]byte(DefeatPrefix)):])
//				var candidateArr []*Candidate
//				if err := rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &candidateArr); nil != err {
//					return make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID][]*Candidate, 0), err
//				}
//				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), candidateArr)
//				defeatCandidates[key] = candidateArr
//			}
//		}
//	}
//	return originCandidates, immediateCandates, defeatCandidates, nil
//}

func (c *CandidatePool) initDataByState (state vm.StateDB) error {

	// 加载 见证人信息
	var witnessIds []discover.NodeID
	c.originCandidates = make(map[discover.NodeID]*types.Candidate, 0)
	if valByte := state.GetState(CandidateAddr, WitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &witnessIds); nil != err {
			log.Error("Failed to decode witnessIds", "err", err)
			return err
		}
	}

	for _, witnessId := range witnessIds {
		fmt.Println("witnessId = ", witnessId.String())
		var can types.Candidate
		if err := rlp.DecodeBytes(state.GetState(CandidateAddr, WitnessKey(witnessId)), &can); nil != err {
			log.Error("Failed to decode Candidate", "err", err)
			return CandidateDecodeErr
		}
		c.originCandidates[witnessId] = &can
	}
	// 加载 下一轮见证人
	var nextWitnessIds []discover.NodeID
	c.nextOriginCandidates = make(map[discover.NodeID]*types.Candidate, 0)
	if valByte := state.GetState(CandidateAddr, NextWitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &nextWitnessIds); nil != err {
			log.Error("Failed to decode nextWitnessIds", "err", err)
			return err
		}
	}

	for _, witnessId := range nextWitnessIds {
		fmt.Println("nextwitnessId = ", witnessId.String())
		var can types.Candidate
		if err := rlp.DecodeBytes(state.GetState(CandidateAddr, NextWitnessKey(witnessId)), &can); nil != err {
			log.Error("Failed to decode Candidate", "err", err)
			return CandidateDecodeErr
		}
		c.nextOriginCandidates[witnessId] = &can
	}
	// 加载 入围者
	var immediateIds []discover.NodeID
	c.immediateCandates = make(map[discover.NodeID]*types.Candidate, 0)
	if valByte := state.GetState(CandidateAddr, ImmediateListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &immediateIds); nil != err {
			log.Error("Failed to decode immediateIds", "err", err)
			return err
		}
	}
	for _, immediateId := range immediateIds {
		fmt.Println("immediateId = ", immediateId.String())
		var can types.Candidate
		if err := rlp.DecodeBytes(state.GetState(CandidateAddr, ImmediateKey(immediateId)), &can); nil != err {
			log.Error("Failed to decode Candidate", "err", err)
			return CandidateDecodeErr
		}
		c.immediateCandates[immediateId] = &can
	}
	// 加载 需要退款信息
	var defeatIds []discover.NodeID
	c.defeatCandidates = make(map[discover.NodeID][]*types.Candidate, 0)
	if valBtye := state.GetState(CandidateAddr, DefeatListKey()); len(valBtye) != 0 {
		if err := rlp.DecodeBytes(valBtye, &defeatIds); nil != err {
			log.Error("Failed to decode defeatIds", "err", err)
			return err
		}
	}
	for _, defeatId := range defeatIds {
		fmt.Println("defeatId = ", defeatId.String())
		var canArr []*types.Candidate
		if err := rlp.DecodeBytes(state.GetState(CandidateAddr, DefeatKey(defeatId)), &canArr); nil != err {
			log.Error("Failed to decode CandidateArr", "err", err)
			return CandidateDecodeErr
		}
		c.defeatCandidates[defeatId] = canArr
	}
	return nil
}

// 候选人抵押
func(c *CandidatePool) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		return err
	}
	// 先追加到 缓存数组中然后做排序
	if len(c.immediateCandates) != 0 && len(c.candiateIds) == 0 {
		for _, v := range c.immediateCandates {
			c.candiateIds = append(c.candiateIds, v.CandidateId)
		}
	}
	c.candiateIds = append(c.candiateIds, can.CandidateId)
	c.immediateCandates[can.CandidateId] = can

	// 排序
	candidateSort(c.candiateIds, c.immediateCandates)
	// 把多余入围者移入落榜名单
	if len(c.candiateIds) > int(c.maxCount) {
		// 截取出当前入围之外的候选人
		tmpArr := (c.candiateIds)[c.maxCount:]
		// 保留入围候选人
		c.candiateIds = (c.candiateIds)[:c.maxCount]
		// 处理落选人
		for _, tmpId := range tmpArr {
			tmp := c.immediateCandates[tmpId]
			// 删除trie中的 入围者信息
			if err := c.delImmediate(state, tmpId); nil != err {
				return err
			}
			// 追加到落榜集
			if err := c.setDefeat(state, tmpId, tmp); nil != err {
				return err
			}
			//delete(c.immediateCandates, tmpId)
		}
		// 更新落选人索引
		if err := c.setDefeatIndex(state); nil != err {
			return err
		}
	}
	// 入围者上树
	for _, id := range c.candiateIds {
		can := c.immediateCandates[id]
		if val, err := rlp.EncodeToBytes(can); nil != err {
			log.Error("Failed to encode Candidate err", err)
			return err
		}else {
			state.SetState(CandidateAddr, ImmediateKey(id), val)
		}
	}
	// 更新入围者索引
	if val, err := rlp.EncodeToBytes(c.candiateIds); nil != err {
		log.Error("Failed to encode ImmediateIds err", err)
		return err
	}else {
		state.SetState(CandidateAddr, ImmediateListKey(), val)
	}
	return nil
}


// 获取入围候选人信息
func (c *CandidatePool) GetCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error) {
	return c.getCandidate(state, nodeId)
}

// 候选人退出
func (c *CandidatePool) WithdrawCandidate (state vm.StateDB, nodeId discover.NodeID, price int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		return err
	}

	if price <= 0 {
		log.Error("withdraw failed price invalid, price", price)
		return WithdrawPriceErr
	}
	can, ok := c.immediateCandates[nodeId]
	if !ok {
		log.Error("current Candidate is empty")
		return CandidateEmptyErr
	}
	if nil == can {
		log.Error("current Candidate is empty")
		return CandidateEmptyErr
	}
	// 判断退押 金额
	if (can.Deposit.Cmp(new(big.Int).SetUint64(uint64(price)))) < 0 {
		log.Error("withdraw failed refund price must less or equal deposit", "key", nodeId.String())
		return WithdrawPriceErr
	}else if (can.Deposit.Cmp(new(big.Int).SetUint64(uint64(price)))) == 0 { // 全额退出质押
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
		for id, v := range c.immediateCandates {
			if id == nodeId {
				// 剩下部分
				canNew := &types.Candidate{
					Deposit:		new(big.Int).Sub(can.Deposit, new(big.Int).SetUint64(uint64(price))),
					BlockNumber: 	can.BlockNumber,
					TxIndex: 		can.TxIndex,
					CandidateId: 	v.CandidateId,
					Host: 			v.Host,
					Port: 			v.Port,
					Owner: 			can.Owner,
					From: 			can.From,
				}
				c.immediateCandates[id] = canNew

				// 更新剩余部分
				if err := c.setImmediate(state, nodeId, canNew); nil != err {
					return err
				}
				// 退款部分新建退款信息
				canDefeat := &types.Candidate{
					Deposit: 		new(big.Int).SetUint64(uint64(price)),
					BlockNumber: 	can.BlockNumber,
					TxIndex: 		can.TxIndex,
					CandidateId: 	v.CandidateId,
					Host: 			v.Host,
					Port: 			v.Port,
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
		}
	}
	return nil
}

// 获取实时所有入围候选人
func (c *CandidatePool) GetChosens (state vm.StateDB) []*types.Candidate {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
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
		return nil
	}
	witnessIds, err := c.getWitnessIndex(state)
	if nil != err {
		log.Error("Failed to getWitnessIndex err", err)
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
		return nil, err
	}

	defeat, ok := c.defeatCandidates[nodeId]
	if !ok {
		log.Error("Candidate is empty")
		return nil, CandidateDecodeErr
	}
	return defeat, nil
}

// 判断是否落榜
func (c *CandidatePool) IsDefeat (state vm.StateDB, nodeId discover.NodeID) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
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

// 揭榜见证人
func (c *CandidatePool) Election(state *state.StateDB) bool{
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to election err", err)
		return false
	}
	immediateIds, err := c.getImmediateIndex(state)
	if nil != err {
		log.Error("Failed to getImmediateIndex err", err)
		return false
	}
	// 对当前所有入围者排序
	candidateSort(immediateIds, c.immediateCandates)
	// 缓存前面一定数量的见证人
	//var nextWits map[discover.NodeID]*Candidate
	var nextWitIds []discover.NodeID
	// 如果入选人数不超过见证人数
	if len(immediateIds) <= int(c.maxChair) {
		nextWitIds = make([]discover.NodeID, len(immediateIds))
		copy(nextWitIds, immediateIds)

	}else {
		// 入选人数超过了见证人数，提取前N名
		nextWitIds = make([]discover.NodeID, c.maxChair)
		copy(nextWitIds, immediateIds)
	}

	nextWits := make(map[discover.NodeID]*types.Candidate, 0)
	// copy 见证人信息
	copyCandidateMapByIds(nextWits, c.immediateCandates, nextWitIds)

	// 删除旧的
	for nodeId, _ := range c.nextOriginCandidates {
		if err := c.delNextWitness(state, nodeId); nil != err {
			log.Error("failed to election err", err)
			return false
		}
	}
	// 设置新的
	for nodeId, can := range nextWits {
		if err := c.setNextWitness(state, nodeId, can); nil != err {
			log.Error("failed to election err", err)
			return false
		}
	}
	if err := c.setNextWitnessIndex(state, nextWitIds); nil != err {
		log.Error("failed to election err", err)
	}
	// 将新的揭榜后的见证人信息，置入下轮见证人集
	c.nextOriginCandidates = nextWits
	return len(c.nextOriginCandidates) > 0
}



// 一键提款
func (c *CandidatePool) RefundBalance (state vm.StateDB, nodeId discover.NodeID, blockNumber uint64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		log.Error("Failed to refundbalance err", err)
		return err
	}

	var arr []*types.Candidate
	if defeatArr, ok := c.defeatCandidates[nodeId]; ok {
		arr = defeatArr
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

	contractBalance := state.GetBalance(CandidateAddr)
	currentNum := new(big.Int).SetUint64(blockNumber)

	for index, can := range arr {
		sub := new(big.Int).Sub(currentNum, can.BlockNumber)
		if sub.Cmp(new(big.Int).SetUint64(c.RefundBlockNumber)) >= 0 { // 允许退款
			delCanArr = append(delCanArr, can)
			arr = append(arr[:index], arr[index+1:]...)
		}else {
			continue
		}

		if len(addr) == 0 {
			addr = can.Owner
		}else {
			if addr != can.Owner {
				log.Info("发现抵押节点nodeId下有不同受益者地址", "nodeId", nodeId.String(), "addr1", addr.String(), "addr2", can.Owner)
				if len(arr) != 0 {
					arr = append(delCanArr, arr...)
				}else {
					arr = delCanArr
				}
				c.defeatCandidates[nodeId] = arr
				return CandidateOwnerErr
			}
		}

		if (contractBalance.Cmp(new(big.Int).SetUint64(amount))) < 0 {
			log.Error("constract account insufficient balance ", state.GetBalance(CandidateAddr).String(), "amount", amount)
			if len(arr) != 0 {
				arr = append(delCanArr, arr...)
			}else {
				arr = delCanArr
			}
			c.defeatCandidates[nodeId] = arr
			return ContractBalanceNotEnoughErr
		}
	}

	// 统一更新树
	if len(arr) == 0 {
		state.SetState(CandidateAddr, DefeatKey(nodeId), []byte{})
		delete(c.defeatCandidates, nodeId)
	}else {
		if val, err := rlp.EncodeToBytes(arr); nil != err {
			log.Error("Failed to encode candidate object on RefundBalance", "key", nodeId.String(), "err", err)
			arr = append(delCanArr, arr...)
			c.defeatCandidates[nodeId] = arr
			return CandidateDecodeErr
		}else {
			state.SetState(CandidateAddr, DefeatKey(nodeId), val)
			c.defeatCandidates[nodeId] = arr
		}
	}
	// 扣减合约余额
	//sub := new(big.Int).Sub(state.GetBalance(CandidateAddr), new(big.Int).SetUint64(amount))
	//state.SetBalance(CandidateAddr, sub)
	state.SubBalance(CandidateAddr, new(big.Int).SetUint64(amount))
	// 增加收益账户余额
	//add := new(big.Int).Add(c.cacheState.GetBalance(addr), new(big.Int).SetUint64(amount))
	//c.cacheState.SetBalance(addr, add)
	state.AddBalance(addr, new(big.Int).SetUint64(amount))
	return nil
}

// 触发替换下轮见证人列表
func (c *CandidatePool) Switch() bool {

	return false
}

// 根据块高重置 state
func (c *CandidatePool) ResetStateByBlockNumber (blockNumber uint64) bool {
	return false
}




func (c *CandidatePool) setImmediate(state vm.StateDB, key discover.NodeID, can *types.Candidate) error {
	c.immediateCandates[key] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setImmediate", "key", key.String(), "err", err)
		return CandidateEncodeErr
	}else {
		// 实时的入围候选人 input trie
		state.SetState(CandidateAddr, ImmediateKey(key), value)
		c.count ++
	}
	return nil
}

func (c *CandidatePool) getImmediateIndex (state vm.StateDB) ([]discover.NodeID, error) {
	var arr []discover.NodeID
	if err := rlp.DecodeBytes(state.GetState(CandidateAddr, ImmediateListKey()), &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

// 删除自动更新索引
func (c *CandidatePool) delImmediate (state vm.StateDB, candidateId discover.NodeID) error {
	// map 中删掉
	delete(c.immediateCandates, candidateId)
	// trie 中删掉实时信息
	state.SetState(CandidateAddr, ImmediateKey(candidateId), []byte{})
	c.count --
	// 删除索引中的对应id
	var ids []discover.NodeID
	if err := rlp.DecodeBytes(state.GetState(CandidateAddr, ImmediateListKey()), &ids); nil != err {
		log.Error("Failed to decode ImmediateIds err", err)
		return err
	}
	var flag bool
	for i, id := range ids {
		if id == candidateId {
			flag = true
			ids = append(ids[:i], ids[i+1:]...)
		}
	}
	if flag {
		if val, err := rlp.EncodeToBytes(ids); nil != err {
			log.Error("Failed to encode ImmediateIds err", err)
			return err
		}else {
			state.SetState(CandidateAddr, ImmediateListKey(), val)
		}
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
		state.SetState(CandidateAddr, DefeatKey(candidateId), value)
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
		state.SetState(CandidateAddr, DefeatListKey(), value)
	}
	return nil
}



func (c *CandidatePool) getWitnessIndex (state vm.StateDB) ([]discover.NodeID, error) {
	var arr []discover.NodeID
	if err := rlp.DecodeBytes(state.GetState(CandidateAddr, WitnessListKey()), &arr); nil != err {
		return nil, err
	}
	return arr, nil
}

func (c *CandidatePool) setNextWitness(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	c.nextOriginCandidates[nodeId] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setImmediate", "key", nodeId.String(), "err", err)
		return CandidateEncodeErr
	}else {
		// 实时的入围候选人 input trie
		state.SetState(CandidateAddr, NextWitnessKey(nodeId), value)
	}
	return nil
}

func (c *CandidatePool) delNextWitness (state vm.StateDB, candidateId discover.NodeID) error {
	// map 中删掉
	delete(c.nextOriginCandidates, candidateId)
	// trie 中删掉实时信息
	state.SetState(CandidateAddr, NextWitnessKey(candidateId), []byte{})
	// 删除索引中的对应id
	var ids []discover.NodeID
	if err := rlp.DecodeBytes(state.GetState(CandidateAddr, NextWitnessListKey()), &ids); nil != err {
		log.Error("Failed to decode ImmediateIds err", err)
		return err
	}
	var flag bool
	for i, id := range ids {
		if id == candidateId {
			flag = true
			ids = append(ids[:i], ids[i+1:]...)
		}
	}
	if flag {
		if val, err := rlp.EncodeToBytes(ids); nil != err {
			log.Error("Failed to encode ImmediateIds err", err)
			return err
		}else {
			state.SetState(CandidateAddr, NextWitnessListKey(), val)
		}
	}
	return nil
}

func (c *CandidatePool) setNextWitnessIndex (state vm.StateDB, nodeIds []discover.NodeID) error {
	if value ,err := rlp.EncodeToBytes(&nodeIds); nil != err {
		log.Error("Failed to encode candidate object on setDefeatIds err", err)
		return CandidateEncodeErr
	}else {
		state.SetState(CandidateAddr, DefeatListKey(), value)
	}
	return nil
}


func (c *CandidatePool) getCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error){
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state); nil != err {
		return nil, err
	}
	if candidatePtr, ok := c.immediateCandates[nodeId]; ok {
		return candidatePtr, nil
	}
	return nil, CandidateEmptyErr
}


func breakUpMap(origin, newData map[discover.NodeID]*types.Candidate) (map[discover.NodeID]*types.Candidate, map[discover.NodeID]struct{}){
	// 需要更新集 		需要删除集
	updateMap, delMap := make(map[discover.NodeID]*types.Candidate, 0), make(map[discover.NodeID]struct{}, 0)
	for id, can := range origin {
		if _, ok := newData[id]; ok {
			updateMap[id] = can
		}else {
			delMap[id] = struct{}{}
		}
	}
	for id, can := range newData {
		if _, ok := origin[id]; !ok {
			updateMap[id] = can
		}
	}
	return updateMap, delMap
}


func copyCandidateMapByIds(target, source map[discover.NodeID]*types.Candidate, ids []discover.NodeID){
	for _, id := range ids {
		target[id] = source[id]
	}
}

func compare(c, can *types.Candidate) int {
	// 质押金大的放前面
	if c.Deposit.Cmp(can.Deposit) > 0 {
		return 1
	}else if c.Deposit == can.Deposit {
		// 块高小的放前面
		if i := c.BlockNumber.Cmp(can.BlockNumber); i == 1 {
			return -1
		}else if i == 0 {
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
func candidateSort(arr []discover.NodeID, candidates map[discover.NodeID]*types.Candidate) {
	quickRealSort(arr, candidates, 0, len(arr) - 1)
}
func quickRealSort (arr []discover.NodeID, candidates map[discover.NodeID]*types.Candidate, left, right int)  {
	if left < right {
		pivot := partition(arr, candidates, left, right)
		quickRealSort(arr, candidates, left, pivot - 1)
		quickRealSort(arr, candidates, pivot + 1, right)
	}
}
func partition(arr []discover.NodeID, candidates map[discover.NodeID]*types.Candidate, left, right int) int {
	for left < right {
		for left < right && compare(candidates[arr[left]], candidates[arr[right]]) >= 0 {
			right --
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			left ++
		}
		for left < right && compare(candidates[arr[left]], candidates[arr[right]]) >= 0 {
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