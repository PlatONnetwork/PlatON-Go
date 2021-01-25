package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
)

//this is use fix validators staking shares error, https://github.com/PlatONnetwork/PlatON-Go/issues/1654
func NewFixIssue1654Plugin(sdb snapshotdb.DB) *FixIssue1654Plugin {
	fix := new(FixIssue1654Plugin)
	fix.sdb = sdb
	return fix
}

type FixIssue1654Plugin struct {
	sdb snapshotdb.DB
}

func (a *FixIssue1654Plugin) fix(blockHash common.Hash, chainID *big.Int, state xcom.StateDB) error {
	if chainID.Cmp(params.AlayaChainConfig.ChainID) != 0 {
		return nil
	}
	candidates, err := NewIssue1654Candidates()
	if err != nil {
		return err
	}
	for _, candidate := range candidates {
		canAddr, err := xutil.NodeId2Addr(candidate.nodeID)
		if nil != err {
			return err
		}
		can, err := stk.GetCandidateInfo(blockHash, canAddr)
		if snapshotdb.NonDbNotFoundErr(err) {
			return err
		}
		if can.IsNotEmpty() && can.StakingBlockNum == candidate.stakingNum {
			if can.Status.IsValid() {
				if err := stk.db.DelCanPowerStore(blockHash, can); nil != err {
					return err
				}
				can.SubShares(candidate.shouldSub)
				if err := stk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
					return err
				}
				if err := stk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable, gov.Gte0160VersionState(state)); nil != err {
					return err
				}
				log.Debug("fix issue1654,can is valid,update the can power", "nodeID", candidate.nodeID, "stakingNum", candidate.stakingNum, "sub", candidate.shouldSub, "newShare", can.Shares)
			} else {
				if can.Shares != nil {
					can.SubShares(candidate.shouldSub)
					if err := stk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable, gov.Gte0160VersionState(state)); nil != err {
						return err
					}
					log.Debug("fix issue1654,can is invalid", "nodeID", candidate.nodeID, "stakingNum", candidate.stakingNum, "sub", candidate.shouldSub, "newShare", can.Shares)
				}
			}
		}
	}
	return nil
}

type issue1654Candidate struct {
	nodeID     discover.NodeID
	stakingNum uint64
	shouldSub  *big.Int
}

func NewIssue1654Candidates() ([]issue1654Candidate, error) {
	//todo ,the candidates need set when alaya  upgrade to 0.15.0
	nodes := make([]issue1654Candidate, 0)
	return nodes, nil
}
