package downloader

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"math/big"
	"math/rand"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

// Test chain parameters.
var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress = crypto.PubkeyToAddress(testKey.PublicKey)
	testDB      = rawdb.NewMemoryDatabase()
	testGenesis = core.GenesisBlockForTesting(testDB, testAddress, big.NewInt(1000000000))
	//contract Store {
	//    uint256[] arr = [4,9,5,6,8,7,1,2,3,10];
	//
	//    constructor() public {
	//        for (uint256 i = 0; i < 10; i++) {
	//            arr[i] = i;
	//        }
	//    }
	//}
	evmContract = hexutil.MustDecode("0x6080604052604051806101400160405280600460ff168152602001600960ff168152602001600560ff168152602001600660ff168152602001600860ff168152602001600760ff168152602001600160ff168152602001600260ff168152602001600360ff168152602001600a60ff16815250600090600a6100829291906100ce565b5034801561008f57600080fd5b5060008090505b600a8110156100c85780600082815481106100ad57fe5b90600052602060002001819055508080600101915050610096565b50610145565b82805482825590600052602060002090810192821561010f579160200282015b8281111561010e578251829060ff169055916020019190600101906100ee565b5b50905061011c9190610120565b5090565b61014291905b8082111561013e576000816000905550600101610126565b5090565b90565b603e806101536000396000f3fe6080604052600080fdfea265627a7a72315820d78ddbad6fa6c2fe65bc00b227816a5b7a4f17e0b43fc14f4122adebad43e3be64736f6c63430005100032")
)

const (
	snapshotDBBaseNum = 300
	blockSyncItems    = 1024
)

// The common prefix of all test chains:
var testChainBase = newTestChain(blockSyncItems, testGenesis)

// Different forks on top of the base chain:
//var testChainForkLightA, testChainForkLightB, testChainForkHeavy *testChain

func init() {
	//	var forkLen = int(1000000 + 50)
	//	var wg sync.WaitGroup
	//	wg.Add(3)
	//	go func() { testChainForkLightA = testChainBase.makeFork(forkLen, false, 1); wg.Done() }()
	//	go func() { testChainForkLightB = testChainBase.makeFork(forkLen, false, 2); wg.Done() }()
	//	go func() { testChainForkHeavy = testChainBase.makeFork(forkLen, true, 3); wg.Done() }()
	//	wg.Wait()
}

type testChain struct {
	genesis  *types.Block
	chain    []common.Hash
	headerm  map[common.Hash]*types.Header
	blockm   map[common.Hash]*types.Block
	receiptm map[common.Hash][]*types.Receipt
	pposData [][2][]byte
	baseNum  int
}

// newTestChain creates a blockchain of the given length.
func newTestChain(length int, genesis *types.Block) *testChain {
	tc := new(testChain).copy(length)
	tc.genesis = genesis
	tc.chain = append(tc.chain, genesis.Hash())
	tc.headerm[tc.genesis.Hash()] = tc.genesis.Header()
	tc.blockm[tc.genesis.Hash()] = tc.genesis
	tc.generate(length-1, 0, genesis, false)
	return tc
}

// makeFork creates a fork on top of the test chain.
func (tc *testChain) makeFork(length int, heavy bool, seed byte) *testChain {
	fork := tc.copy(tc.len() + length)
	fork.generate(length, seed, tc.headBlock(), heavy)
	return fork
}

// shorten creates a copy of the chain with the given length. It panics if the
// length is longer than the number of available blocks.
func (tc *testChain) shorten(length int) *testChain {
	if length > tc.len() {
		panic(fmt.Errorf("can't shorten test chain to %d blocks, it's only %d blocks long", length, tc.len()))
	}
	return tc.copy(length)
}

func (tc *testChain) copy(newlen int) *testChain {
	cpy := &testChain{
		genesis:  tc.genesis,
		headerm:  make(map[common.Hash]*types.Header, newlen),
		blockm:   make(map[common.Hash]*types.Block, newlen),
		receiptm: make(map[common.Hash][]*types.Receipt, newlen),
		pposData: make([][2][]byte, 0),
	}
	for i := 0; i < len(tc.chain) && i < newlen; i++ {
		hash := tc.chain[i]
		cpy.chain = append(cpy.chain, tc.chain[i])
		cpy.blockm[hash] = tc.blockm[hash]
		cpy.headerm[hash] = tc.headerm[hash]
		cpy.receiptm[hash] = tc.receiptm[hash]
	}
	if len(tc.pposData) > 0 {
		cpy.pposData = tc.pposData[0:newlen]
	}
	if newlen < tc.baseNum {
		cpy.baseNum = newlen - 1
	} else {
		cpy.baseNum = tc.baseNum
	}
	return cpy
}

// generate creates a chain of n blocks starting at and including parent.
// the returned hash chain is ordered head->parent. In addition, every 22th block
// contains a transaction and every 5th an uncle to allow testing correct block
// reassembly.
func (tc *testChain) generate(n int, seed byte, parent *types.Block, heavy bool) {
	// start := time.Now()
	// defer func() { fmt.Printf("test chain generated in %v\n", time.Since(start)) }()

	blocks, receipts := core.GenerateChain(params.TestChainConfig, parent, &consensus.BftMock{}, testDB, n, func(i int, block *core.BlockGen) {
		block.SetCoinbase(common.Address{seed})
		// If a heavy chain is requested, delay blocks to raise difficulty
		if heavy {
			block.OffsetTime(-1)
		}
		// Include transactions to the miner to make blocks more interesting.
		if parent == tc.genesis && i%22 == 0 {
			signer := types.NewEIP155Signer(params.TestChainConfig.ChainID)
			// evm contract generate more storage, convenient for fast sync
			tx, err := types.SignTx(types.NewContractCreation(block.TxNonce(testAddress), big.NewInt(1000), params.TxGas*10, nil, evmContract), signer, testKey)
			if err != nil {
				panic(err)
			}
			block.AddTx(tx)
		}
	})

	// Convert the block-chain into a hash-chain and header/block maps
	for i, b := range blocks {
		hash := b.Hash()
		tc.chain = append(tc.chain, hash)
		tc.blockm[hash] = b
		tc.headerm[hash] = b.Header()
		tc.receiptm[hash] = receipts[i]
		tc.pposData = append(tc.pposData, [2][]byte{common.Int64ToBytes(rand.Int63()), common.Int64ToBytes(rand.Int63())})
	}
	tc.baseNum = snapshotDBBaseNum
}

// len returns the total number of blocks in the chain.
func (tc *testChain) len() int {
	return len(tc.chain)
}

// headBlock returns the head of the chain.
func (tc *testChain) headBlock() *types.Block {
	return tc.blockm[tc.chain[len(tc.chain)-1]]
}

// headersByHash returns headers in ascending order from the given hash.
func (tc *testChain) headersByHash(origin common.Hash, amount int, skip int) []*types.Header {
	num, _ := tc.hashToNumber(origin)
	return tc.headersByNumber(num, amount, skip)
}

// headersByNumber returns headers in ascending order from the given number.
func (tc *testChain) headersByNumber(origin uint64, amount int, skip int) []*types.Header {
	result := make([]*types.Header, 0, amount)
	for num := origin; num < uint64(len(tc.chain)) && len(result) < amount; num += uint64(skip) + 1 {
		if header, ok := tc.headerm[tc.chain[int(num)]]; ok {
			result = append(result, header)
		}
	}
	return result
}

// receipts returns the receipts of the given block hashes.
func (tc *testChain) receipts(hashes []common.Hash) [][]*types.Receipt {
	results := make([][]*types.Receipt, 0, len(hashes))
	for _, hash := range hashes {
		if receipt, ok := tc.receiptm[hash]; ok {
			results = append(results, receipt)
		}
	}
	return results
}

// bodies returns the block bodies of the given block hashes.
func (tc *testChain) bodies(hashes []common.Hash) ([][]*types.Transaction, [][]byte) {
	transactions := make([][]*types.Transaction, 0, len(hashes))
	ex := make([][]byte, len(hashes))
	for _, hash := range hashes {
		if block, ok := tc.blockm[hash]; ok {
			transactions = append(transactions, block.Transactions())
		}
	}
	return transactions, ex
}

func (tc *testChain) hashToNumber(target common.Hash) (uint64, bool) {
	for num, hash := range tc.chain {
		if hash == target {
			return uint64(num), true
		}
	}
	return 0, false
}
