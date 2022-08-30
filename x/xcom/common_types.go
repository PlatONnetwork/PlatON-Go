// Copyright 2021 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package xcom

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// StateDB is an Plugin database for full state querying.
type StateDB interface {
	CreateAccount(common.Address)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	GetCommittedState(common.Address, []byte) []byte
	//GetState(common.Address, common.Hash) common.Hash
	//SetState(common.Address, common.Hash, common.Hash)
	GetState(common.Address, []byte) []byte
	SetState(common.Address, []byte, []byte)

	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	ForEachStorage(common.Address, func([]byte, []byte) bool)

	//ppos add
	TxHash() common.Hash
	TxIdx() uint32

	IntermediateRoot(deleteEmptyObjects bool) common.Hash
}

type Result struct {
	Code uint32
	Ret  interface{}
}

func NewResult(err *common.BizError, data interface{}) []byte {
	var res *Result
	if err != nil && err != common.NoErr {
		res = &Result{err.Code, err.Msg}
	} else {
		res = &Result{common.NoErr.Code, data}
	}
	bs, _ := json.Marshal(res)
	return bs
}

// addLog let the result add to event.
// 参数datas可为空,里面的值不能为空
// Log.data字段编码规则:
// 如果datas为空,  rlp([code]),
// 如果datas不为空,rlp([code,rlp(data1),rlp(data2)...]),
func AddLogWithRes(state StateDB, blockNumber uint64, contractAddr common.Address, event, code string, datas ...interface{}) {
	buf := new(bytes.Buffer)
	logData := [][]byte{[]byte(code)}
	if len(datas) != 0 && datas[0] != nil {
		for _, res := range datas {
			resByte, err := rlp.EncodeToBytes(res)
			if err != nil {
				log.Error("Cannot RlpEncode the log datas", "datas", datas, "err", err, "event", event)
				panic("Cannot RlpEncode the log data")
			}
			logData = append(logData, resByte)
		}
	}
	if err := rlp.Encode(buf, logData); nil != err {
		log.Error("Cannot RlpEncode the log data", "data", code, "err", err, "event", event)
		panic("Cannot RlpEncode the log data")
	}

	state.AddLog(&types.Log{
		Address:     contractAddr,
		Topics:      nil, //[]common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data:        buf.Bytes(),
		BlockNumber: blockNumber,
	})
}
