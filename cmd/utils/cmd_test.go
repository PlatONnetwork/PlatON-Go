package utils

import (
	"testing"
)

func TestImportChain(t *testing.T) {
	/*
		rootDir := "/tmp/testchain" + strconv.Itoa(int(time.Now().UnixNano()))

		sourceDir := rootDir + "_1"

		fileName := "000001.log"

		cbft := cbft.NewFaker()

		sdb, err := ethdb.NewLDBDatabase(sourceDir, 0, 0)
		if err != nil {
			Fatalf("Open db from file, dir: %s, err: %v", sourceDir, err)
		}

		genesis := core.DefaultGenesisBlock()
		genesisBlock := genesis.MustCommit(sdb)

		chain1, err := core.NewBlockChain(sdb, nil, params.TestChainConfig, cbft, vm.Config{}, nil)
		if err != nil {
			Fatalf("Can't create source BlockChain: %v", err)
		}

		// inserChain
		// Full block-chain requested

		blocks, _ := core.GenerateChain(params.TestChainConfig, types.NewBlockWithHeader(genesisBlock.Header()), cbft, sdb, 100, func(i int, b *core.BlockGen) {
			b.SetCoinbase(common.Address{0: byte(12), 19: byte(i)})
		})

		_, err = chain1.InsertChain(blocks)

		if nil != err {
			Fatalf("Failed to InsertChain: %v", err)
		}

		//chainDb := ethdb.NewMemDatabase()

		chainDb, err := ethdb.NewLDBDatabase(rootDir+"_2", 0, 0)
		if err != nil {
			Fatalf("Open db from file, dir: %s, err: %v", sourceDir, err)
		}

		genesis.MustCommit(chainDb)

		chain2, err := core.NewBlockChain(chainDb, nil, params.TestChainConfig, cbft, vm.Config{}, nil)
		if err != nil {
			Fatalf("Can't create target BlockChain: %v", err)
		}

		err = ImportChain(chain2, sourceDir+"/"+fileName)
		fmt.Println("importchain: " + err.Error())*/
}
