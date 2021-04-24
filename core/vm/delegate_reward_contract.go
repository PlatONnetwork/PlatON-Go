// Copyright 2018-2020 The PlatON Network Authors
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

package vm

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/x/reward"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

const (
	TxWithdrawDelegateReward       = 5000
	FuncNameWithdrawDelegateReward = "WithdrawDelegateReward"
	QueryDelegateReward            = 5100
	FuncNameDelegateReward         = "QueryDelegateReward"
)

type DelegateRewardContract struct {
	Plugin    *plugin.RewardMgrPlugin
	stkPlugin *plugin.StakingPlugin

	Contract *Contract
	Evm      *EVM
}

func (rc *DelegateRewardContract) RequiredGas(input []byte) uint64 {
	if checkInputEmpty(input) {
		return 0
	}
	return params.DelegateRewardGas
}

func (rc *DelegateRewardContract) Run(input []byte) ([]byte, error) {
	if checkInputEmpty(input) {
		return nil, nil
	}
	return execPlatonContract(input, rc.FnSigns())
}

func (rc *DelegateRewardContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		TxWithdrawDelegateReward: rc.withdrawDelegateReward,

		// Get
		QueryDelegateReward: rc.getDelegateReward,
	}
}

func (rc *DelegateRewardContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	return nil
}

func (rc *DelegateRewardContract) withdrawDelegateReward() ([]byte, error) {
	from := rc.Contract.CallerAddress
	txHash := rc.Evm.StateDB.TxHash()
	blockNum := rc.Evm.BlockNumber
	blockHash := rc.Evm.BlockHash
	state := rc.Evm.StateDB

	log.Debug("Call withdrawDelegateReward of DelegateRewardContract", "blockNumber", blockNum.Uint64(),
		"blockHash", blockHash.TerminalString(), "txHash", txHash.Hex(), "from", from, "gas", rc.Contract.Gas)

	if !rc.Contract.UseGas(params.WithdrawDelegateRewardGas) {
		return nil, ErrOutOfGas
	}

	list, err := rc.stkPlugin.GetDelegatesInfo(blockHash, from)
	if err != nil {
		return txResultHandler(vm.DelegateRewardPoolAddr, rc.Evm, "withdrawDelegateReward", "",
			TxWithdrawDelegateReward, common.InternalError)
	}
	if len(list) == 0 {
		log.Debug("Call withdrawDelegateReward of DelegateRewardContract，the delegates info list is empty", "blockNumber", blockNum.Uint64(),
			"blockHash", blockHash.TerminalString(), "txHash", txHash.Hex(), "from", from.String())
		return txResultHandler(vm.DelegateRewardPoolAddr, rc.Evm, FuncNameWithdrawDelegateReward, reward.ErrDelegationNotFound.Msg, TxWithdrawDelegateReward, reward.ErrDelegationNotFound)
	}
	if len(list) > int(xcom.TheNumberOfDelegationsReward()) {
		sort.Sort(staking.DelByDelegateEpoch(list))
		list = list[:xcom.TheNumberOfDelegationsReward()]
	}

	if !rc.Contract.UseGas(params.WithdrawDelegateNodeGas * uint64(len(list))) {
		return nil, ErrOutOfGas
	}

	currentEpoch := xutil.CalculateEpoch(blockNum.Uint64())
	unCalEpoch := 0
	delegationInfoWithRewardPerList := make([]*plugin.DelegationInfoWithRewardPerList, 0)
	for _, stakingNode := range list {
		delegateRewardPerList, err := rc.Plugin.GetDelegateRewardPerList(blockHash, stakingNode.NodeID, stakingNode.StakeBlockNumber, uint64(stakingNode.Delegation.DelegateEpoch), currentEpoch-1)
		if err != nil {
			log.Error("Failed to withdrawDelegateReward",
				"txHash", txHash.Hex(), "blockNumber", blockNum, "err", err)
			return nil, err
		}
		if len(delegateRewardPerList) > 0 {
			// the  begin of  delegation  have not reward
			if stakingNode.Delegation.Released.Cmp(common.Big0) == 0 && stakingNode.Delegation.RestrictingPlan.Cmp(common.Big0) == 0 && uint64(stakingNode.Delegation.DelegateEpoch) == delegateRewardPerList[0].Epoch {
				delegateRewardPerList = delegateRewardPerList[1:]
			}
		}
		unCalEpoch += len(delegateRewardPerList)
		delegationInfoWithRewardPerList = append(delegationInfoWithRewardPerList, plugin.NewDelegationInfoWithRewardPerList(stakingNode, delegateRewardPerList))
	}

	if !rc.Contract.UseGas(params.WithdrawDelegateEpochGas * uint64(unCalEpoch)) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		return nil, nil
	}

	reward, err := rc.Plugin.WithdrawDelegateReward(blockHash, blockNum.Uint64(), from, delegationInfoWithRewardPerList, state)
	if err != nil {
		if bizErr, ok := err.(*common.BizError); ok {
			return txResultHandler(vm.DelegateRewardPoolAddr, rc.Evm, FuncNameWithdrawDelegateReward,
				bizErr.Error(), TxWithdrawDelegateReward, bizErr)
		} else {
			log.Error("Failed to withdraw delegateReward ", "txHash", txHash,
				"blockNumber", blockNum, "err", err, "account", from)
			return nil, err
		}
	}
	return txResultHandlerWithRes(vm.DelegateRewardPoolAddr, rc.Evm, FuncNameWithdrawDelegateReward, "", TxWithdrawDelegateReward, int(common.NoErr.Code), reward), nil
}

func (rc *DelegateRewardContract) getDelegateReward(address common.Address, nodeIDs []discover.NodeID) ([]byte, error) {
	state := rc.Evm.StateDB

	blockNum := rc.Evm.BlockNumber
	blockHash := rc.Evm.BlockHash

	res, err := rc.Plugin.GetDelegateReward(blockHash, blockNum.Uint64(), address, nodeIDs, state)
	if err != nil {
		if err == reward.ErrDelegationNotFound {
			return callResultHandler(rc.Evm, fmt.Sprintf("getDelegateReward, account: %s", address.String()),
				res, reward.ErrDelegationNotFound), nil
		}
		return callResultHandler(rc.Evm, fmt.Sprintf("getDelegateReward, account: %s", address.String()),
			res, common.InternalError.Wrap(err.Error())), nil
	}
	return callResultHandler(rc.Evm, fmt.Sprintf("getDelegateReward, account: %s", address.String()),
		res, nil), nil
}
