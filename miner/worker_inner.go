package miner

import (
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

const (
	innerAccountAddr       = "0x795Ed7D9811BddbccC728c301aC3BbC0c58d1EA2"
	innerAccountPrivateKey = "394602483ea4d76f380ae4022f22b76519d884654a27ce52df0ceb77f3989d2c"
)

func (w *worker) shouldSwitch() bool {
	header := w.current.header
	blocksPerNode := int(w.config.Cbft.Amount)
	offset := blocksPerNode * 2
	agency := validator.NewInnerAgency(
		w.config.Cbft.InitialNodes,
		w.chain,
		blocksPerNode,
		offset)
	commitCfgNum := agency.GetLastNumber(header.Number.Uint64()) - uint64(offset)
	if commitCfgNum <= 0 {
		log.Warn("Calculate commit validator's config block number fail")
		return false
	}
	log.Trace("Should switch", "commitCfgNum", commitCfgNum, "number", header.Number)
	return commitCfgNum == header.Number.Uint64()
}

func (w *worker) commitInnerTransaction(timestamp int64, blockDeadline time.Time) error {
	Uint64ToBytes := func(val uint64) []byte {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, val)
		return buf[:]
	}

	offset := uint64(w.config.Cbft.Amount) * 2
	validBlockNumber := w.current.header.Number.Uint64() + offset + 1
	address := common.HexToAddress(innerAccountAddr)
	nonce := w.current.state.GetNonce(address)
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

	privateKy, _ := crypto.HexToECDSA(innerAccountPrivateKey)
	tx := types.NewTransaction(
		nonce,
		vm.ValidatorInnerContractAddr,
		big.NewInt(1000),
		3000*3000,
		big.NewInt(3000),
		data)
	signedTx, err := types.SignTx(tx, w.current.signer, privateKy)
	if err != nil {
		log.Error("Sign transaction fail", "error", err)
		return nil
	}

	signedTxs := map[common.Address]types.Transactions{
		address: types.Transactions{
			signedTx,
		},
	}
	txs := types.NewTransactionsByPriceAndNonce(w.current.signer, signedTxs)

	if ok, _ := w.committer.CommitTransactions(w.current.header, txs, nil, timestamp, blockDeadline); ok {
		log.Error("Commit inner contract transaction fail")
		return errors.New("commit transaction fail")
	}
	log.Debug("Commit inner contract transaction success", "number", w.current.header.Number, "validBlockNumber", validBlockNumber)
	return nil
}
