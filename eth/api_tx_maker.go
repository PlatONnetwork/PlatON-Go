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

	"github.com/PlatONnetwork/PlatON-Go/event"

	"github.com/PlatONnetwork/PlatON-Go/core"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
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

func (txg *TxGenAPI) Start(txPer, txTime int, accountPath string, start, end int) error {
	if txg.start {
		return errors.New("the tx maker is working")
	}

	//make sure when the txGen is start ,the node will not receive txs from other node,
	//so this node can keep in sync with other nodes
	atomic.StoreUint32(&txg.eth.protocolManager.acceptRemoteTxs, 1)

	txch := make(chan types.Transactions, 20)
	txg.txfeed = txg.eth.blockchain.SubscribeBlockTxsEvent(txch)
	txg.txGenExitCh = make(chan struct{})

	if err := txg.makeTransaction(txPer, txTime, accountPath, start, end, txch); err != nil {
		return err
	}
	txg.start = true
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

type PriAccount struct {
	Priv    *ecdsa.PrivateKey
	Nonce   uint64
	Address common.Address

	ReceiptsNonce uint64
	SendTime      time.Time
}

type TxMakeManger struct {
	accounts      map[common.Address]*PriAccount
	sleepAccounts map[common.Address]struct{}
	toPool        []common.Address

	accountSize      uint
	BlockProduceTime time.Time
}

func NewTxMakeManger(pendingState *state.ManagedState, accountPath string, start, end int) (*TxMakeManger, error) {
	file, err := os.Open(accountPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read genesis file:%v", err)
	}
	defer file.Close()

	var priKey []PriKeyJson
	if err := json.NewDecoder(file).Decode(&priKey); err != nil {
		return nil, fmt.Errorf("invalid genesis file chain id:%v", err)
	}

	t := new(TxMakeManger)
	t.accounts = make(map[common.Address]*PriAccount)
	t.toPool = make([]common.Address, 0)

	for i := start; i <= end; i++ {
		privateKey, err := crypto.HexToECDSA(priKey[i].Pri)
		if err != nil {
			return nil, fmt.Errorf("NewTxMakeManger HexToECDSA fail:%v", err)
		}
		address, err := common.Bech32ToAddress(priKey[i].Add)
		if err != nil {
			return nil, fmt.Errorf("NewTxMakeManger Bech32ToAddress fail:%v", err)
		}
		nonce := pendingState.GetNonce(address)
		t.accounts[address] = &PriAccount{privateKey, nonce, address, nonce, time.Now()}
		t.toPool = append(t.toPool, address)
	}
	t.sleepAccounts = make(map[common.Address]struct{})
	return t, nil
}

type PriKeyJson struct {
	Pri string `json:"private_key"`
	Add string `json:"address"`
}

func (txg *TxGenAPI) makeTransaction(txPer, txTime int, accountPath string, start, end int, txch chan types.Transactions) error {

	txm, err := NewTxMakeManger(txg.eth.txPool.State(), accountPath, start, end)
	if err != nil {
		return err
	}
	singine := types.NewEIP155Signer(new(big.Int).SetInt64(txg.eth.chainConfig.ChainID.Int64()))

	txsCh := make(chan []*types.Transaction, 1)

	go func() {
		for {
			select {
			case res := <-txch:
				txm.BlockProduceTime = time.Now()
				if len(res) > 0 {
					for _, receipt := range res {
						if account, ok := txm.accounts[receipt.FromAddr(singine)]; ok {
							account.ReceiptsNonce = receipt.Nonce()
						}
					}
				}
			case txs := <-txsCh:
				txg.eth.protocolManager.txsCh <- core.NewTxsEvent{txs}
			case <-txg.txGenExitCh:
				log.Debug("MakeTransaction get receipt nonce  exit")
				return
			}
		}
	}()

	log.Info("begin to MakeTransaction")
	gasPrice := new(big.Int).SetInt64(50000000000)
	amount := new(big.Int).SetInt64(1)

	go func() {
		shouldmake := time.NewTicker(time.Millisecond * time.Duration(txTime))
		shouldReport := time.NewTicker(time.Second * 10)
		length := len(txm.toPool)

		for {
			if time.Since(txm.BlockProduceTime) >= 10*time.Second {
				log.Debug("MakeTx should sleep", "time", time.Since(txm.BlockProduceTime))
				time.Sleep(time.Second * 5)
				continue
			}
			select {
			case <-shouldmake.C:
				now := time.Now()
				txs := make([]*types.Transaction, 0)
				toAdd := txm.toPool[rand.Intn(length)]
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
					tx := types.NewTransaction(account.Nonce, toAdd, amount, 30000, gasPrice, nil)
					newTx, err := types.SignTx(tx, singine, account.Priv)
					if err != nil {
						log.Crit(fmt.Errorf("sign error,%s", err.Error()).Error())
					}
					txs = append(txs, newTx)
					account.Nonce++
					account.SendTime = time.Now()
					if len(txs) >= txPer {
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
