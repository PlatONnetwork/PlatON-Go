package miner

import (
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// innerAdminKey returns the locally configured key for the inner-mode switch
// transaction. The key is not part of the source tree; the address it must
// match is chainConfig.Cbft.InnerAdmin.
func (w *worker) innerAdminKey() (*ecdsa.PrivateKey, common.Address, error) {
	if w.chainConfig == nil || w.chainConfig.Cbft == nil || w.chainConfig.Cbft.InnerAdmin == (common.Address{}) {
		return nil, common.Address{}, errors.New("inner admin address is not configured")
	}
	if w.config == nil || len(w.config.InnerAdminKey) == 0 {
		return nil, common.Address{}, errors.New("inner admin key is not configured")
	}
	key, err := crypto.ToECDSA(w.config.InnerAdminKey)
	if err != nil {
		return nil, common.Address{}, err
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	if addr != w.chainConfig.Cbft.InnerAdmin {
		return nil, common.Address{}, errors.New("inner admin key does not match Cbft.InnerAdmin")
	}
	return key, addr, nil
}

func (w *worker) shouldSwitch(env *environment) bool {
	blocksPerNode := int(w.chainConfig.Cbft.Amount)
	offset := blocksPerNode * 2
	agency := validator.NewInnerAgency(
		w.chainConfig.Cbft.InitialNodes,
		w.chain,
		blocksPerNode,
		offset)
	commitCfgNum := agency.GetLastNumber(env.header.Number.Uint64()) - uint64(offset)
	if commitCfgNum <= 0 {
		log.Warn("Calculate commit validator's config block number fail")
		return false
	}
	log.Trace("Should switch", "commitCfgNum", commitCfgNum, "number", env.header.Number)
	return commitCfgNum == env.header.Number.Uint64()
}

func (w *worker) commitInnerTransaction(env *environment, timestamp int64, blockDeadline time.Time) error {
	Uint64ToBytes := func(val uint64) []byte {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, val)
		return buf[:]
	}

	offset := uint64(w.chainConfig.Cbft.Amount) * 2
	validBlockNumber := env.header.Number.Uint64() + offset + 1
	privateKy, address, err := w.innerAdminKey()
	if err != nil {
		log.Error("Inner admin is not available", "error", err)
		return err
	}
	nonce := env.state.GetNonce(address)
	param := [][]byte{
		common.Int64ToBytes(2003),
		[]byte("SwitchValidators"),
		Uint64ToBytes(validBlockNumber),
	}
	data, err := rlp.EncodeToBytes(param)
	if err != nil {
		log.Error("RLP encode fail", "error", err)
		return err
	}

	tx := types.NewTransaction(
		nonce,
		vm.ValidatorInnerContractAddr,
		big.NewInt(1000),
		3000*3000,
		big.NewInt(3000),
		data)
	signedTx, err := types.SignTx(tx, env.signer, privateKy)
	if err != nil {
		log.Error("Sign transaction fail", "error", err)
		return err
	}

	signedTxs := map[common.Address]types.Transactions{
		address: types.Transactions{
			signedTx,
		},
	}
	txs := types.NewTransactionsByPriceAndNonce(env.signer, signedTxs, nil)

	tempContractCache := make(map[common.Address]struct{})
	includedFrom := len(env.txs)
	if interrupted, _ := w.committer.CommitTransactions(env, txs, nil, timestamp, blockDeadline, tempContractCache); interrupted {
		log.Error("Commit inner contract transaction interrupted")
		return errors.New("commit transaction interrupted")
	}
	if !innerTxIncluded(env.txs[includedFrom:], signedTx.Hash()) {
		log.Error("Inner switch transaction was not included", "hash", signedTx.Hash(), "number", env.header.Number)
		return errors.New("inner switch transaction was not included")
	}
	log.Debug("Commit inner contract transaction success", "number", env.header.Number, "validBlockNumber", validBlockNumber)
	return nil
}

func innerTxIncluded(txs []*types.Transaction, hash common.Hash) bool {
	for _, tx := range txs {
		if tx.Hash() == hash {
			return true
		}
	}
	return false
}
