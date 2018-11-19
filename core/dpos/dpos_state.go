package depos

import (
	"fmt"
	"sync"
	"strings"
	"encoding/json"
	"Platon-go/common"
	"Platon-go/rlp"
	"Platon-go/log"
	"Platon-go/params"
	"Platon-go/core/state"
	"Platon-go/p2p/discover"
	"math/big"
	"Platon-go/core"
	"errors"
	//"github.com/go-errors/errors"
)


const(
	// 即时入围竞选人
	ImmediatePrefix 	= "id"
	ImmediateListPrefix = "iL"
	// 见证人
	WitnessPrefix 		= "wn"
	WitnessListPrefix 	= "wL"
	// 需要退款的
	DefeatPrefix 		= "df"
	DefeatListPrefix 	= "dL"
)

var (
	// 即时入围竞选人
	ImmediateBtyePrefix 	= []byte("id")
	ImmediateListBtyePrefix = []byte("iL")
	// 见证人
	WitnessBtyePrefix 		= []byte("wn")
	WitnessListBtyePrefix 	= []byte("wL")
	// 需要退款的
	DefeatBtyePrefix 		= []byte("df")
	DefeatListBtyePrefix 	= []byte("dL")
	CandidateAddr 			= common.HexToAddress("0x1....10")



	CandidateEncodeErr    		= errors.New("Candidate encoding err")
	CandidateDecodeErr 			= errors.New("Cnadidate decoding err")
	WithdrawPriceErr 			= errors.New("Withdraw Price err")
	CandidateEmptyErr 			= errors.New("Candidate is empty")
	ContractBalanceNotEnoughErr = errors.New("Contract's balance is not enough")
	CandidateOwnerErr 			= errors.New("CandidateOwner Addr is illegal")
)

type CandidatePool struct {
	// 当前入选者数目
	count 					uint64
	// 最大允许入选人数目
	maxCount				uint64
	// 最大允许见证人数目
	maxChair				uint64
	// 允许退款的块间隔
	RefundBlockNumber 		uint64

	// 上一轮选出的见证人集
	originCandidates  		map[discover.NodeID]*Candidate
	// 即时的入选人集
	immediateCandates 		map[discover.NodeID]*Candidate
	// 质押失败的竞选人集 (退款用)
	defeatCandidates 		map[discover.NodeID][]*Candidate
	blockChain     			*core.BlockChain
	cacheState 				*state.StateDB

	// cache
	candidateCacheArr 		[]*Candidate
	lock 					*sync.RWMutex
}

// 初始化全局候选池对象
//func NewCandidatePool(state *state.StateDB, configs *params.DposConfig, isgenesis bool) (*CandidatePool, error) {
func NewCandidatePool(blockChain *core.BlockChain, configs *params.DposConfig) (*CandidatePool, error) {

	currentBlock := blockChain.CurrentBlock()
	var state *state.StateDB
	if statedb, err := blockChain.State(); nil != err {
		log.Error("reference statedb failed", err)
	}else {
		state = statedb
	}

	var originMap, immediateMap map[discover.NodeID]*Candidate
	var  defeatMap map[discover.NodeID][]*Candidate
	// 非创世块，需要从db加载
	if blockChain.Genesis().NumberU64() != currentBlock.NumberU64() {
		tr := state.StorageTrie(CandidateAddr)
		// 构造上一轮见证人的map
		var err error
		originMap, immediateMap, defeatMap, err = initCandidatesByTrie(tr)
		if nil != err {
			return nil, err
		}
	}else {
		//fmt.Printf("config：%+v \n", configs)
		var err error
		originMap, immediateMap, err = buildByConfig(configs, state)
		if nil != err {
			return nil, err
		}
		defeatMap =  make(map[discover.NodeID][]*Candidate)
	}

	return &CandidatePool{
		count: 					uint64(len(immediateMap)),
		maxCount:				configs.MaxCount,
		maxChair:				configs.MaxChair,
		RefundBlockNumber: 		configs.RefundBlockNumber,
		originCandidates: 		originMap,
		immediateCandates: 		immediateMap,
		defeatCandidates: 		defeatMap,
		blockChain: 			blockChain,
		cacheState: 			state,
		candidateCacheArr: 		make([]*Candidate, 0),
		lock: 					&sync.RWMutex{},
	}, nil
}
// 根据配置文件构建 dpos原始见证人
func buildByConfig(configs *params.DposConfig, state *state.StateDB) (map[discover.NodeID]*Candidate, map[discover.NodeID]*Candidate, error){
	originMap, immediateMap := make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0)
	if len(configs.Candidates) != 0 {
		// 如果配置过多，只取前面的
		if len(configs.Candidates) > int(configs.MaxCount) {
			configs.Candidates = (configs.Candidates)[:configs.MaxCount]
		}
		for i, canConfig := range configs.Candidates {
			can := &Candidate{
				Deposit:			canConfig.Deposit,
				BlockNumber: 		canConfig.BlockNumber,
				TxIndex: 		 	canConfig.TxIndex,
				CandidateId: 	 	canConfig.CandidateId,
				Host: 			 	canConfig.Host,
				Port: 			 	canConfig.Port,
				Owner: 				canConfig.Owner,
				From: 				canConfig.From,
			}

			// 追加到树中
			if val, err := rlp.EncodeToBytes(can); nil == err {
				if uint64(i) < configs.MaxChair {
					state.SetState(CandidateAddr, WitnessKey(can.CandidateId), val)
					originMap[can.CandidateId] = can
				}
				state.SetState(CandidateAddr,  ImmediateKey(can.CandidateId), val)
				immediateMap[can.CandidateId] = can
			}else {
				log.Error("Failed to encode candidate object", "key", string(WitnessKey(can.CandidateId)), "err", err)
				return make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0), err
			}
		}
	}
	return originMap, immediateMap, nil
}

func (c *CandidatePool) CommitTrie (deleteEmptyObjects bool) (root common.Hash, err error) {
	return c.cacheState.Commit(deleteEmptyObjects)
}

func (c *CandidatePool)IteratorTrie(s string){
	iteratorTrie(s, c.cacheState.StorageTrie(CandidateAddr))
}

func iteratorTrie(s string, tr state.Trie){
	it := tr.NodeIterator(nil)
	for it.Next(true) {
		if it.Leaf() {
			var a Candidate
			rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &a)
			fmt.Println(s, string(tr.GetKey(it.LeafKey())), "== ", &a)
		}
	}
}

func printObject(s string, obj interface{}){
	objs, _ := json.Marshal(obj)
	fmt.Println(s, string(objs), "\n")
}

func (c *CandidatePool)buildCandidatesByTrie(prefix string) map[common.Hash]*Candidate {
	tr := c.cacheState.StorageTrie(CandidateAddr)
	it := tr.NodeIterator(nil)
	candidates := make(map[common.Hash]*Candidate, 0)
	for it.Next(true) {
		if it.Leaf() {
			trieKey := tr.GetKey(it.LeafKey())
			cleanKey := trieKey[len([]byte(CandidateAddr.String())):]

			// 判断前缀
			if strings.HasPrefix(string(cleanKey), prefix){
				key := common.BytesToHash(cleanKey[len([]byte(prefix)):])
				var candidate Candidate
				rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &candidate)
				//rlp.DecodeBytes(it.LeafBlob(), candidate)
				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), candidate)
				candidates[key] = &candidate
			}
		}
	}
	return candidates
}

func (c *CandidatePool)buildCandidateArrByTrie(prefix string) map[common.Hash][]*Candidate {
	tr := c.cacheState.StorageTrie(CandidateAddr)
	it := tr.NodeIterator(nil)
	candidates := make(map[common.Hash][]*Candidate, 0)
	for it.Next(true) {
		if it.Leaf() {
			trieKey := tr.GetKey(it.LeafKey())
			cleanKey := trieKey[len([]byte(CandidateAddr.String())):]

			// 判断前缀
			if strings.HasPrefix(string(cleanKey), prefix){
				key := common.BytesToHash(cleanKey[len([]byte(prefix)):])
				var arr []*Candidate
				rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &arr)
				//rlp.DecodeBytes(it.LeafBlob(), candidate)
				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), arr)
				candidates[key] = arr
			}
		}
	}
	return candidates
}

func initCandidatesByTrie (tr state.Trie) (map[discover.NodeID]*Candidate, map[discover.NodeID]*Candidate, map[discover.NodeID][]*Candidate, error){
	it := tr.NodeIterator(nil)
	// 见证人
	originCandidates  := make(map[discover.NodeID]*Candidate, 0)
	// 即时入围者
	immediateCandates := make(map[discover.NodeID]*Candidate, 0)
	// 需要退款信息
	defeatCandidates  := make(map[discover.NodeID][]*Candidate, 0)
	for it.Next(true) {
		if it.Leaf() {
			trieKey := tr.GetKey(it.LeafKey())
			cleanKey := trieKey[len([]byte(CandidateAddr.String())):]

			// 根据前缀获取 见证人信息
			if strings.HasPrefix(string(cleanKey), WitnessPrefix){
				key := discover.MustBytesID(cleanKey[len([]byte(WitnessPrefix)):])
				//key := common.BytesToHash(cleanKey[len([]byte(WitnessPrefix)):])
				var candidate Candidate
				if err := rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &candidate); nil != err {
					return make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID][]*Candidate, 0), err
				}
				//rlp.DecodeBytes(it.LeafBlob(), candidate)
				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), candidate)
				originCandidates[key] = &candidate
			}

			// 根据前缀获取 入围竞选人信息
			if strings.HasPrefix(string(cleanKey), ImmediatePrefix){
				key := discover.MustBytesID(cleanKey[len([]byte(ImmediatePrefix)):])
				//key := common.BytesToHash(cleanKey[len([]byte(ImmediatePrefix)):])
				var candidate Candidate
				if err := rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &candidate); nil != err {
					return make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID][]*Candidate, 0), err
				}
				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), candidate)
				immediateCandates[key] = &candidate
			}

			// 根据前缀获取 落榜需要退款信息
			if strings.HasPrefix(string(cleanKey), DefeatPrefix){
				key := discover.MustBytesID(cleanKey[len([]byte(DefeatPrefix)):])
				//key := common.BytesToHash(cleanKey[len([]byte(DefeatPrefix)):])
				var candidateArr []*Candidate
				if err := rlp.DecodeBytes(tr.GetKey(it.LeafBlob()), &candidateArr); nil != err {
					return make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID][]*Candidate, 0), err
				}
				fmt.Printf("遍历出来的k-v: %v == %+v \n", string(trieKey), candidateArr)
				defeatCandidates[key] = candidateArr
			}
		}
	}
	return originCandidates, immediateCandates, defeatCandidates, nil
}

// 候选人抵押 不带前缀的key
func(c *CandidatePool) SetCandidate(nodeId discover.NodeID, can *Candidate) error{
	c.lock.Lock()
	defer c.lock.Unlock()
	// 先追加到 缓存数组中然后做排序
	if len(c.immediateCandates) != 0 && len(c.candidateCacheArr) == 0 {
		for _, v := range c.immediateCandates {
			c.candidateCacheArr = append(c.candidateCacheArr, v)
		}
	}
	c.candidateCacheArr = append(c.candidateCacheArr, can)
	if err := c.setImmediate(nodeId, can); nil != err {
		return err
	}
	candidateSort(c.candidateCacheArr)
	// 把多余入围者移入落榜名单
	if len(c.candidateCacheArr) > int(c.maxCount) {
		// 截取出当前入围之外的候选人
		tmpArr := (c.candidateCacheArr)[c.maxCount:]
		// 保留入围候选人
		c.candidateCacheArr = (c.candidateCacheArr)[:c.maxCount]
		// 处理落选人
		for _, tmp := range tmpArr {
			c.delImmediate(tmp.CandidateId)
			// 追加到落榜集
			if err := c.setDefeat(tmp.CandidateId, tmp); nil != err {
				return err
			}
		}
	}
	return nil
}


// 获取入围候选人信息
func (c *CandidatePool) GetCandidate(nodeId discover.NodeID) (*Candidate, error) {
	return c.getCandidate(nodeId)
}

// 候选人退出
func (c *CandidatePool) WithdrawCandidate (nodeId discover.NodeID, price int) error {
	if price <= 0 {
		log.Error("withdraw failed price invalid, price", price)
		return WithdrawPriceErr
	}
	can, ok := c.immediateCandates[nodeId]
	if !ok {
		// 从树拿
		enc := c.cacheState.GetState(CandidateAddr, ImmediateKey(nodeId))
		var tmp Candidate
		if err := rlp.DecodeBytes(enc, &tmp); nil != err {
			log.Error("withdraw failed", "key", nodeId.String(), "err", err)
			return CandidateDecodeErr
		}
		can = &tmp
	}
	if nil == can {
		log.Error("current Candidate is empty")
		return CandidateEmptyErr
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	// 判断退押 金额
	if (can.Deposit.Cmp(new(big.Int).SetUint64(uint64(price)))) < 0 {
		log.Error("withdraw failed refund price must less or equal deposit", "key", nodeId.String())
		return WithdrawPriceErr
	}else if (can.Deposit.Cmp(new(big.Int).SetUint64(uint64(price)))) == 0 { // 全额退出质押
		c.delImmediate(nodeId)
		// 追加到落选
		if err := c.setDefeat(nodeId, can); nil != err {
			return err
		}
	}else { // 只退了一部分, 需要重新对入围者排序
		for i, v := range c.candidateCacheArr {
			if v.CandidateId == nodeId {
				// 剩下部分
				canNew := &Candidate{
					Deposit:		new(big.Int).Sub(can.Deposit, new(big.Int).SetUint64(uint64(price))),
					BlockNumber: 	can.BlockNumber,
					TxIndex: 		can.TxIndex,
					CandidateId: 	v.CandidateId,
					Host: 			v.Host,
					Port: 			v.Port,
					Owner: 			can.Owner,
					From: 			can.From,
				}
				c.candidateCacheArr[i] = canNew
				// 剩余部分
				if err := c.setImmediate(nodeId, canNew); nil != err {
					return err
				}
				canDefeat := &Candidate{
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
				if err := c.setDefeat(nodeId, canDefeat); nil != err {
					return err
				}
			}
		}

	}
	return nil
}

// 获取实时所有入围候选人
func (c *CandidatePool) GetChosens () []*Candidate {
	c.lock.Lock()
	defer c.lock.Unlock()
	arr := make([]*Candidate, 0)
	for _, v := range c.immediateCandates {
		arr = append(arr, v)
	}
	candidateSort(arr)
	return arr
}

// 获取所有见证人
func (c *CandidatePool) GetChairpersons () []*Candidate {
	arr := make([]*Candidate, 0)
	for _, v := range c.originCandidates {
		arr = append(arr, v)
	}
	candidateSort(arr)
	return arr
}


// 获取退款信息
func (c *CandidatePool) GetDefeat(nodeId discover.NodeID) ([]*Candidate, error){
	defeat, ok := c.defeatCandidates[nodeId]
	if !ok {
		enc := c.cacheState.GetState(CandidateAddr, DefeatKey(nodeId))
		var tmp []*Candidate
		if err := rlp.DecodeBytes(enc, &tmp); nil != err {
			log.Error("Failed to decode candidate object on GetDefeat", "key", nodeId.String(), "err", err)
			return nil, CandidateDecodeErr
		}
		defeat = tmp
		c.defeatCandidates[nodeId] = tmp
	}
	return defeat, nil
}

// 判断是否落榜
func (c *CandidatePool) IsDefeat (nodeId discover.NodeID) (bool, error) {

	if _, ok := c.immediateCandates[nodeId]; !ok {
		enc := c.cacheState.GetState(CandidateAddr, ImmediateKey(nodeId))
		var tmp Candidate
		if err := rlp.DecodeBytes(enc, &tmp); nil != err {
			log.Error("Failed to decode candidate object on IsDefeat", "key", nodeId.String(), "err", err)
			return false, CandidateDecodeErr
		}
		// 有点问题，不能直接这么写
		c.setImmediate(nodeId, &tmp)
	}else {
		return true, nil
	}
	return true, nil
}

// 根据nodeId查询 质押信息中的 受益者地址
func (c *CandidatePool) GetOwner (nodeId discover.NodeID) common.Address {
	c.lock.Lock()
	defer c.lock.Unlock()
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
func (c *CandidatePool) Election() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	// 先判断当前nodeId 是否为本轮见证人
	//if _, ok := c.originCandidates[nodeId]; !ok {
	//	log.Error("Current NODEID is not this round of witnesses")
	//	return false
	//}
	candidateSort(c.candidateCacheArr)
	// 缓存前面一定数量的见证人
	var arr []*Candidate
	// 如果入选人数不超过见证人数
	if len(c.candidateCacheArr) <= int(c.maxChair) {
		arr = make([]*Candidate, len(c.candidateCacheArr))
		copy(arr, c.candidateCacheArr)
	}else {
		// 入选人数超过了见证人数，提取前N名
		arr = make([]*Candidate, c.maxChair)
		copy(arr, c.candidateCacheArr)
	}
	tmpMap := make(map[discover.NodeID]*Candidate, 0)
	for _, v := range arr {
		tmpMap[v.CandidateId] = v
	}
	// 对比新的见证人集 和 上一轮的见证人集
	updateMap, delMap := breakUpMap(c.originCandidates, tmpMap)
	c.originCandidates = tmpMap
	// update trie
	if len(delMap) != 0 {
		for id, _ := range delMap {
			c.cacheState.SetState(CandidateAddr, WitnessKey(id), []byte{})
		}
	}
	for id, can := range updateMap {
		if val, err := rlp.EncodeToBytes(can); nil == err {
			c.cacheState.SetState(CandidateAddr, WitnessKey(id), val)
		}else {
			log.Error("Failed to encode candidate object on Election", "key", string(WitnessKey(id)), "err", err)
			continue
		}
	}
	return len(c.originCandidates) > 0
}


// 一键提款
func (c *CandidatePool) RefundBalance (nodeId discover.NodeID) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	var arr []*Candidate
	if defeatArr, ok := c.defeatCandidates[nodeId]; ok {
		arr = defeatArr
	}else {
		var canArr []*Candidate
		enc := c.cacheState.GetState(CandidateAddr, DefeatKey(nodeId))
		if err := rlp.DecodeBytes(enc, &canArr); nil != err {
			log.Error("Failed to decode candidate object on RefundBalance", "key", nodeId.String(), "err", err)
			return CandidateDecodeErr
		}
		arr = canArr
		c.defeatCandidates[nodeId] = canArr
	}
	// 遍历退款
	var addr common.Address  	// 用来做校验用，即正常情况应该每个nodeId的质押退款信息中的收益者(owner)应该为同一个
	var amount uint64 	 		// 累计需要一次性退款的金额
	delCanArr := make([]*Candidate, 0) // 中转需要删除的退款信息
	contractBalance := c.cacheState.GetBalance(CandidateAddr)
	currentBlock := c.blockChain.CurrentBlock()

	for index, can := range arr {
		sub := new(big.Int).Sub(currentBlock.Number(), can.BlockNumber)
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
			log.Error("constract account insufficient balance ", c.cacheState.GetBalance(CandidateAddr).String(), "amount", amount)
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
		c.cacheState.SetState(CandidateAddr, DefeatKey(nodeId), []byte{})
		delete(c.defeatCandidates, nodeId)
	}else {
		if val, err := rlp.EncodeToBytes(arr); nil != err {
			log.Error("Failed to encode candidate object on RefundBalance", "key", nodeId.String(), "err", err)
			arr = append(delCanArr, arr...)
			c.defeatCandidates[nodeId] = arr
			return CandidateDecodeErr
		}else {
			c.cacheState.SetState(CandidateAddr, DefeatKey(nodeId), val)
			c.defeatCandidates[nodeId] = arr
		}
	}
	// 扣减合约余额
	sub := new(big.Int).Sub(c.cacheState.GetBalance(CandidateAddr), new(big.Int).SetUint64(amount))
	c.cacheState.SetBalance(CandidateAddr, sub)
	// 增加收益账户余额
	add := new(big.Int).Add(c.cacheState.GetBalance(addr), new(big.Int).SetUint64(amount))
	c.cacheState.SetBalance(addr, add)
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




func (c *CandidatePool) setImmediate(key discover.NodeID, can *Candidate) error {
	c.immediateCandates[key] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setImmediate", "key", key.String(), "err", err)
		return CandidateEncodeErr
	}else {
		// 实时的入围候选人 input trie
		c.cacheState.SetState(CandidateAddr, ImmediateKey(key), value)
		c.count ++
	}
	return nil
}

func (c *CandidatePool) delImmediate (candidateId discover.NodeID) {
	// map 中删掉
	delete(c.immediateCandates, candidateId)
	// trie 中删掉实时信息
	c.cacheState.SetState(CandidateAddr, ImmediateKey(candidateId), []byte{})
	c.count --
}

// 设置退款信息
func (c *CandidatePool) setDefeat(candidateId discover.NodeID, can *Candidate) error {

	var defeatArr []*Candidate
	// 追加退款信息
	if defeatArrTmp, ok := c.defeatCandidates[can.CandidateId]; ok {
		defeatArrTmp = append(defeatArrTmp, can)
		c.defeatCandidates[can.CandidateId] = defeatArrTmp
		defeatArr = defeatArrTmp
	}else {
		defeatArrTmp = make([]*Candidate, 0)
		defeatArrTmp = append(defeatArr, can)
		c.defeatCandidates[can.CandidateId] = defeatArrTmp
		defeatArr = defeatArrTmp
	}

	// trie 中添加 退款信息
	if value ,err := rlp.EncodeToBytes(&defeatArr); nil != err {
		log.Error("Failed to encode candidate object on setDefeat", "key", candidateId.String(), "err", err)
		return CandidateEncodeErr
	}else {
		c.cacheState.SetState(CandidateAddr, DefeatKey(candidateId), value)
	}
	return nil
}




func (c *CandidatePool) getCandidate(nodeId discover.NodeID) (*Candidate, error){
	// 先去缓存map拿
	if candidatePtr, ok := c.immediateCandates[nodeId]; ok {
		return candidatePtr, nil
	}
	// 没有就去树上拿
	enc := c.cacheState.GetState(CandidateAddr, ImmediateKey(nodeId))
	var data Candidate
	if err := rlp.DecodeBytes(enc, &data); nil != err {
		log.Error("Failed to decode candidate object on getCandidate", "key", nodeId.String(), "err", err)
		return nil, CandidateDecodeErr
	}
	c.immediateCandates[nodeId] = &data
	return &data, nil
}


//func (c *CandidatePool) GetStorageTrie() Trie {
//
//}



func breakUpMap(origin, newData map[discover.NodeID]*Candidate) (map[discover.NodeID]*Candidate, map[discover.NodeID]struct{}){
	// 需要更新集 		需要删除集
	updateMap, delMap := make(map[discover.NodeID]*Candidate, 0), make(map[discover.NodeID]struct{}, 0)
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


func (c *Candidate) compare(can *Candidate) int {
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
func candidateSort(arr []*Candidate) {
	quickRealSort(arr, 0, len(arr) - 1)
}
func quickRealSort (arr []*Candidate, left, right int)  {
	if left < right {
		pivot := partition(arr, left, right)
		quickRealSort(arr, left, pivot - 1)
		quickRealSort(arr, pivot + 1, right)
	}
}
func partition(arr []*Candidate, left, right int) int {
	for left < right {
		for left < right && arr[left].compare(arr[right]) >= 0 {
			right --
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			left ++
		}
		for left < right && arr[left].compare(arr[right]) >= 0 {
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

func DefeatKey(nodeId discover.NodeID) []byte {
	//key, _ := rlp.EncodeToBytes(nodeId)
	return defeatKey(nodeId.Bytes())
}
func defeatKey(key []byte) []byte {
	return append(DefeatBtyePrefix, key...)
}