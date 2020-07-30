package eth

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/mroth/weightedrand"

	"github.com/PlatONnetwork/PlatON-Go/event"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

func NewTxGenAPI(eth *Ethereum) *TxGenAPI {
	return &TxGenAPI{eth: eth}
}

type TxGenAPI struct {
	eth         *Ethereum
	txGenExitCh chan struct{}
	start       bool
	txfeed      event.Subscription
}

//begin make tx ,Broadcast transactions directly through p2p, without entering the transaction pool
// normalTx, evmTx, wasmTx，The proportion of normal transactions and contract transactions sent out
// such as 1:1:1,this should send 1 normal transactions,1 evm transaction, 1 wasm transaction
// txPer,How many transactions are sent at once
// sendingAmount,Send amount of normal transaction
// txTime,How many ms to send a batch of transactions
// accountPath,Account configuration address
// start, end ,Start and end account
func (txg *TxGenAPI) Start(normalTx, evmTx, wasmTx uint, txPer, txTime uint, sendingAmount uint64, accountPath string, start, end uint) error {
	if txg.start {
		return errors.New("the tx maker is working")
	}

	//make sure when the txGen is start ,the node will not receive txs from other node,
	//so this node can keep in sync with other nodes
	atomic.StoreUint32(&txg.eth.protocolManager.acceptRemoteTxs, 1)

	txch := make(chan types.Transactions, 20)
	txg.txfeed = txg.eth.blockchain.SubscribeBlockTxsEvent(txch)
	txg.txGenExitCh = make(chan struct{})

	if err := txg.makeTransaction(normalTx, evmTx, wasmTx, txPer, txTime, sendingAmount, accountPath, start, end, txch); err != nil {
		return err
	}
	txg.start = true
	return nil
}

func (txg *TxGenAPI) makeTransaction(tx, evm, wasm uint, txPer, txTime uint, sendingAmount uint64, accountPath string, start, end uint, txch chan types.Transactions) error {
	state, err := txg.eth.blockchain.State()
	if err != nil {
		return err
	}
	txm, err := NewTxMakeManger(tx, evm, wasm, sendingAmount, txg.eth.txPool.Nonce, state.GetCodeSize, accountPath, start, end)
	if err != nil {
		state.ClearReference()
		return err
	}
	state.ClearReference()

	singine := types.NewEIP155Signer(new(big.Int).SetInt64(txg.eth.chainConfig.ChainID.Int64()))

	txsCh := make(chan []*types.Transaction, 1)

	go func() {
		for {
			select {
			case res := <-txch:
				txm.blockProduceTime = time.Now()
				if len(res) > 0 {
					for _, receipt := range res {
						if account, ok := txm.accounts[receipt.FromAddr(singine)]; ok {
							account.ReceiptsNonce = receipt.Nonce()
						}
					}
				}
			case txs := <-txsCh:
				txg.eth.txPool.AddRemotes(txs)
			//	txg.eth.protocolManager.txsCh <- core.NewTxsEvent{txs}
			case <-txg.txGenExitCh:
				log.Debug("MakeTransaction get receipt nonce  exit")
				return
			}
		}
	}()

	log.Info("begin to MakeTransaction")
	gasPrice := new(big.Int).SetInt64(50000000000)

	go func() {
		shouldmake := time.NewTicker(time.Millisecond * time.Duration(txTime))
		shouldReport := time.NewTicker(time.Second * 10)
		for {
			if time.Since(txm.blockProduceTime) >= 10*time.Second {
				log.Debug("MakeTx should sleep", "time", time.Since(txm.blockProduceTime))
				time.Sleep(time.Second * 5)
				continue
			}
			select {
			case <-shouldmake.C:
				now := time.Now()
				txs := make([]*types.Transaction, 0, txPer)
				toAdd := txm.pickTxReceive()
				for _, account := range txm.accounts {
					if account.Nonce >= account.ReceiptsNonce+10 {
						if time.Since(account.SendTime) >= time.Second*20 {
							log.Debug("wait account 20s", "account", account.Address, "nonce", account.Nonce, "receiptnonce", account.ReceiptsNonce, "wait time", time.Since(account.SendTime))
							account.Nonce = account.ReceiptsNonce + 1
							delete(txm.sleepAccounts, account.Address)
						} else {
							if _, ok := txm.sleepAccounts[account.Address]; !ok {
								txm.sleepAccounts[account.Address] = struct{}{}
							}
							continue
						}
					} else {
						delete(txm.sleepAccounts, account.Address)
					}

					txContractInputData, txReceive, gasLimit, amount := txm.generateTxParams(toAdd)

					tx := types.NewTransaction(account.Nonce, txReceive, amount, gasLimit, gasPrice, txContractInputData)
					newTx, err := types.SignTx(tx, singine, account.Priv)
					if err != nil {
						log.Crit(fmt.Errorf("sign error,%s", err.Error()).Error())
					}
					txs = append(txs, newTx)
					txm.sendDone(account)
					if len(txs) >= int(txPer) {
						break
					}
				}
				if len(txs) != 0 {
					log.Debug("make Transaction time use", "use", time.Since(now), "txs", len(txs))
					txsCh <- txs
				}
			case <-shouldReport.C:
				sleepAccountsLength := len(txm.sleepAccounts)
				log.Debug("MakeTx info", "sleepAccount", sleepAccountsLength, "perTx", txPer)
			case <-txg.txGenExitCh:
				shouldmake.Stop()
				shouldReport.Stop()
				log.Debug("MakeTransaction exit")
				return
			}
		}
	}()
	return nil
}

func (txg *TxGenAPI) DeployContracts(prikey string, configPath string) error {
	return handelTxGenConfig(configPath, func(txgenInput *TxGenInput) error {
		pri, err := crypto.HexToECDSA(prikey)
		if err != nil {
			return err
		}
		currentState, err := txg.eth.blockchain.State()
		if err != nil {
			return err
		}
		defer currentState.ClearReference()
		account := crypto.PubkeyToAddress(pri.PublicKey)
		nonce := currentState.GetNonce(account)
		singine := types.NewEIP155Signer(new(big.Int).SetInt64(txg.eth.chainConfig.ChainID.Int64()))
		gasPrice := txg.eth.txPool.GasPrice()

		for _, input := range [][]*TxGenInputContractConfig{txgenInput.Wasm, txgenInput.Evm} {
			for _, config := range input {
				tx := types.NewContractCreation(nonce, nil, config.DeployGasLimit, gasPrice, common.Hex2Bytes(config.ContractsCode))
				newTx, err := types.SignTx(tx, singine, pri)
				if err != nil {
					return err
				}
				if err := txg.eth.TxPool().AddRemote(newTx); err != nil {
					return err
				}
				config.DeployContractTxHash = newTx.Hash().String()
				nonce++
			}
		}
		return nil
	})
}

func (txg *TxGenAPI) UpdateConfig(configPath string) error {
	return handelTxGenConfig(configPath, func(txgenInput *TxGenInput) error {
		for _, input := range [][]*TxGenInputContractConfig{txgenInput.Wasm, txgenInput.Evm} {
			for _, config := range input {
				hash := common.HexToHash(config.DeployContractTxHash)
				tx, blockHash, _, index := rawdb.ReadTransaction(txg.eth.ChainDb(), hash)
				if tx == nil {
					return fmt.Errorf("the tx not find yet,tx:%s", hash.String())

				}
				receipts := txg.eth.blockchain.GetReceiptsByHash(blockHash)
				if len(receipts) <= int(index) {
					return fmt.Errorf("the tx receipts not find yet,tx:%s", hash.String())
				}
				receipt := receipts[index]
				config.ContractsAddress = receipt.ContractAddress.String()
			}

		}
		return nil
	})
}

func handelTxGenConfig(configPath string, handle func(*TxGenInput) error) error {
	file, err := os.OpenFile(configPath, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("Failed to read genesis file:%v", err)
	}
	defer file.Close()
	var txgenInput TxGenInput
	if err := json.NewDecoder(file).Decode(&txgenInput); err != nil {
		return fmt.Errorf("invalid genesis file chain id:%v", err)
	}

	if err := handle(&txgenInput); err != nil {
		return err
	}

	output, err := json.MarshalIndent(txgenInput, "", "    ")
	if err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	if _, err := file.Write(output); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	return nil
}

func (txg *TxGenAPI) Stop() error {
	if !txg.start {
		return errors.New("the tx maker has been closed")
	}
	close(txg.txGenExitCh)
	txg.start = false
	txg.txfeed.Unsubscribe()
	atomic.StoreUint32(&txg.eth.protocolManager.acceptRemoteTxs, 0)

	return nil
}

type TxGenInput struct {
	Wasm []*TxGenInputContractConfig `json:"wasm"`
	Evm  []*TxGenInputContractConfig `json:"evm"`
	Tx   []*TxGenInputAccountConfig  `json:"tx"`
}

type TxGenInputAccountConfig struct {
	Pri string `json:"private_key"`
	Add string `json:"address"`
}

type TxGenInputContractConfig struct {
	//CreateContracts
	DeployContractTxHash string `json:"deploy_contract_tx_hash"`
	ContractsCode        string `json:"contracts_code"`
	DeployGasLimit       uint64 `json:"deploy_gas_limit"`

	ContractsAddress string `json:"contracts_address"`

	//CallContracts
	CallWeights  uint   `json:"call_weights"`
	CallGasLimit uint64 `json:"call_gas_limit"`
	CallInput    string `json:"call_input"`
}

type txGenSendAccount struct {
	Priv    *ecdsa.PrivateKey
	Nonce   uint64
	Address common.Address

	ReceiptsNonce uint64
	SendTime      time.Time
}

type txGenTxReceiver struct {
	ContractsAddress common.Address
	Data             []byte
	Weights          uint
	GasLimit         uint64
}

type TxMakeManger struct {
	//from
	accounts map[common.Address]*txGenSendAccount

	amount *big.Int

	//to
	txReceiver   []common.Address
	evmReceiver  weightedrand.Chooser
	wsamReveiver weightedrand.Chooser

	sleepAccounts map[common.Address]struct{}

	blockProduceTime time.Time

	sendTx   uint
	sendEvm  uint
	sendWasm uint

	sendState uint
}

func (s *TxMakeManger) pickTxReceive() common.Address {
	return s.txReceiver[rand.Intn(len(s.txReceiver))]
}

func (s *TxMakeManger) generateTxParams(add common.Address) ([]byte, common.Address, uint64, *big.Int) {
	switch {
	case s.sendState < s.sendTx:
		return nil, add, 30000, s.amount
	case s.sendState < s.sendEvm:
		account := s.evmReceiver.Pick().(*txGenTxReceiver)
		return account.Data, account.ContractsAddress, account.GasLimit, nil
	case s.sendState < s.sendWasm:
		account := s.wsamReveiver.Pick().(*txGenTxReceiver)
		return account.Data, account.ContractsAddress, account.GasLimit, nil
	}
	log.Crit("generateTxParams fail,the sendState should not grate than the sendWasm", "state", s.sendState, "wasm", s.sendWasm)
	return nil, common.Address{}, 0, nil
}

func (s *TxMakeManger) sendDone(account *txGenSendAccount) {
	s.sendState++
	if s.sendState >= s.sendWasm {
		s.sendState = 0
	}
	account.Nonce++
	account.SendTime = time.Now()
}

func NewTxMakeManger(tx, evm, wasm uint, sendingAmount uint64, GetNonce func(addr common.Address) uint64, getCodeSize func(addr common.Address) int, accountPath string, start, end uint) (*TxMakeManger, error) {
	file, err := os.Open(accountPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read genesis file:%v", err)
	}
	defer file.Close()

	var txgenInput TxGenInput
	if err := json.NewDecoder(file).Decode(&txgenInput); err != nil {
		return nil, fmt.Errorf("invalid genesis file chain id:%v", err)
	}

	t := new(TxMakeManger)
	t.amount = new(big.Int).SetUint64(sendingAmount)
	t.accounts = make(map[common.Address]*txGenSendAccount)
	t.txReceiver = make([]common.Address, 0)

	for i := start; i <= end; i++ {
		privateKey, err := crypto.HexToECDSA(txgenInput.Tx[i].Pri)
		if err != nil {
			return nil, fmt.Errorf("NewTxMakeManger HexToECDSA fail:%v", err)
		}
		address, err := common.Bech32ToAddress(txgenInput.Tx[i].Add)
		if err != nil {
			return nil, fmt.Errorf("NewTxMakeManger Bech32ToAddress fail:%v", err)
		}
		nonce := GetNonce(address)
		t.accounts[address] = &txGenSendAccount{privateKey, nonce, address, nonce, time.Now()}
		t.txReceiver = append(t.txReceiver, address)
	}
	t.sleepAccounts = make(map[common.Address]struct{})

	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

	t.sendTx = tx
	t.sendEvm = tx + evm
	t.sendWasm = tx + evm + wasm
	if t.sendWasm == 0 {
		return nil, errors.New("new tx gen fail ,tx+evm+wasm size should not be zero")
	}
	if evm+wasm == 0 {
		return t, nil
	}

	if evm > 0 && len(txgenInput.Evm) == 0 {
		return nil, errors.New("new tx gen fail ,evm config not set")
	}
	if wasm > 0 && len(txgenInput.Wasm) == 0 {
		return nil, errors.New("new tx gen fail ,wasm config not set")
	}

	receiverChooser := make([]weightedrand.Choice, 0, len(txgenInput.Evm))
	for i, ContractConfigs := range [][]*TxGenInputContractConfig{txgenInput.Evm, txgenInput.Wasm} {
		for _, config := range ContractConfigs {
			if config.CallWeights != 0 {
				txReceiver := new(txGenTxReceiver)
				txReceiver.ContractsAddress = common.MustBech32ToAddress(config.ContractsAddress)
				if getCodeSize(txReceiver.ContractsAddress) <= 0 {
					return nil, fmt.Errorf("new tx gen fail the address don't have code,add:%s", txReceiver.ContractsAddress.String())
				}
				txReceiver.Data = common.Hex2Bytes(config.CallInput)
				txReceiver.Weights = config.CallWeights
				txReceiver.GasLimit = config.CallGasLimit
				receiverChooser = append(receiverChooser, weightedrand.NewChoice(txReceiver, txReceiver.Weights))
			}
		}
		if i == 0 {
			t.evmReceiver = weightedrand.NewChooser(receiverChooser...)
		} else {
			t.wsamReveiver = weightedrand.NewChooser(receiverChooser...)
		}
	}
	return t, nil
}
