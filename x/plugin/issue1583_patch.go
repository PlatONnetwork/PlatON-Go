package plugin

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/x/reward"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

//给没有领取委托奖励的账户平账 , https://github.com/PlatONnetwork/PlatON-Go/issues/1583
func NewFixIssue1583Plugin() *FixIssue1583Plugin {
	fix := new(FixIssue1583Plugin)
	return fix
}

type FixIssue1583Plugin struct{}

func (a *FixIssue1583Plugin) fix(blockHash common.Hash, chainID *big.Int, state xcom.StateDB) error {
	if chainID.Cmp(params.AlayaChainConfig.ChainID) != 0 {
		return nil
	}
	accounts, err := newIssue1583Accounts()
	if err != nil {
		return err
	}
	for _, account := range accounts {
		receiveReward := account.RewardPer.CalDelegateReward(account.delegationAmount)
		if err := rm.ReturnDelegateReward(account.addr, receiveReward, state); err != nil {
			log.Error("fix issue 1583,return delegate reward fail", "account", account.addr, "err", err)
			return common.InternalError
		}
	}
	return nil
}

type issue1583Accounts struct {
	addr             common.Address
	delegationAmount *big.Int
	RewardPer        reward.DelegateRewardPer
}

func newIssue1583Accounts() ([]issue1583Accounts, error) {
	type delegationInfo struct {
		account          string
		delegationAmount string
	}

	//node f2ec2830850  in Epoch216
	node1DelegationInfo := []delegationInfo{
		{"atp12trrqnpqkj2kn03cwz5ae4v7lfshvqetqrkeez", "1000000000000000000"},
		{"atp143ml5dd3qz3wykmg4p3vnp9eqlugd9sxmgpsux", "2000000000000000000"},
		{"atp1687slxxcghuhxgv3uy6v2epftn8nhn9jss2sd6", "9361750080000000000"},
		{"atp1dtmhmexryrg7h8tzsufrg4d9y48sne7rr5ezsj", "1150000000000000000"},
		{"atp1y4arxmjpy5grkp9attefax07z56wcuq2n5937u", "2806452320000000000"},
		{"atp1z3w63q0q6rnrqw55s4hhu7ewfep0vcqaxt508d", "1612074080000000000"},
	}
	node1DelegationAmount, _ := new(big.Int).SetString("3185100027705555555556", 10)
	node1DelegationReward, _ := new(big.Int).SetString("8599922061705855512", 10)
	node1RewardPer := reward.DelegateRewardPer{
		Delegate: node1DelegationAmount,
		Reward:   node1DelegationReward,
	}

	accounts := make([]issue1583Accounts, 0)
	for _, c := range node1DelegationInfo {
		addr, err := common.Bech32ToAddress(c.account)
		if err != nil {
			return nil, err
		}
		amount, _ := new(big.Int).SetString(c.delegationAmount, 10)
		accounts = append(accounts, issue1583Accounts{
			addr:             addr,
			delegationAmount: amount,
			RewardPer:        node1RewardPer,
		})
	}

	//fff1010bbf176 in epoch475
	node2DelegationInfos := []delegationInfo{
		{"atp1rek8y8nz07tdp4v469xmuyymar8qe0pnejknyw", "18729736290000000000"},
		{"atp1szy8d7094kl0q82la9l2pfz3hjy99zh730flw4", "89940584120000000000"},
		{"atp15wekgs8q07rs24dmnwmdgd0qhqwvqlus6dh5g5", "62385695880000000000"},
		{"atp1cmfc2nea3am2znutaunaze2q8t6ttyrcu7mjqh", "16252405540000000000"},
	}

	node2DelegationAmount, _ := new(big.Int).SetString("10986475785670000000000", 10)
	node2DelegationReward, _ := new(big.Int).SetString("8790854215894931178", 10)
	node2RewardPer := reward.DelegateRewardPer{
		Delegate: node2DelegationAmount,
		Reward:   node2DelegationReward,
	}

	for _, c := range node2DelegationInfos {
		addr, err := common.Bech32ToAddress(c.account)
		if err != nil {
			return nil, err
		}
		amount, _ := new(big.Int).SetString(c.delegationAmount, 10)
		accounts = append(accounts, issue1583Accounts{
			addr:             addr,
			delegationAmount: amount,
			RewardPer:        node2RewardPer,
		})
	}

	return accounts, nil
}
