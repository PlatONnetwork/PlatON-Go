package plugin

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

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
	type candidate struct {
		Node   string
		Num    int
		Amount string
	}

	candidates := []candidate{
		{"89ca7ccb7fab8e4c8b1b24c747670757b9ef1b3b7631f64e6ea6b469c5936c501fcdcfa7fef2a77521072162c1fc0f8a1663899d31ebb1bc7d00678634ef746c", 522944, "25839533916550000000000"},
		{"ab74f5500dd35497ce09b2dc92a3da26ea371dd9f6d438559b6e19c8f1622ee630951b510cb370aca8267f9bb9a9108bc532ec48dd077474cb79a48122f2ab03", 507203, "30031800000000000000000"},
		{"a2340b4acd4f7b743d7e8785e8ff297490b0e333f25cfe31d17df006f7e554553c6dc502f7c9f7b8798ab3ccf74624065a6cc20603842b1015793c0b37de9b15", 2944652, " 22640000000000000000000"},
		{"fff1010bbf1762d13bf13828142c612a7d287f0f1367f8104a78f001145fd788fb44b87e9eac404bc2e880602450405850ff286658781dce130aee981394551d", 902037, "21794743754935111385799"},
		{"1fd9fd7d9c31dad117384c7cc2f223a9e76f7aa81e30f38f030f24212be3fa20ca1f067d878a8ae97deb89b81efbd0542b3880cbd428b4ffae494fcd2c31834b", 518839, "15890357356590015813760"},
		{"19b8fb478a8502a25e461270122ece3135b15dc0de118264495bae30d39af81fd9134ed95364e6a39c3eebfba57fbffa7961a5158d3dac0f0da0313dac7af024", 514281, "7017624238740000000000"},
		{"8bc8734315acf2af4c92a458f077f1f8c96f0530fb43510c11361f1d6469631423206ef76cd879ade849ee15fbcaeb042e3721168614b4fad4eecd60a6aa3e94", 618133, "11339285782500000000000"},
		{"f1efed4e853d00ff3f1be65fd497bc8f0a3d5f66b285069c9190653567e1838ab635b88940d3ce786747af549a1a5bf9b7173e9dc3a3aea9f10363613581a9e0", 509125, "12734262030390000000000"},
		{"1dbe057f33d9748e1d396d624f4c2554f67742f18247e6be6c615c56c70a6e18a6604dd887fd1e9ffdf9708486fb76b711cb5d8e66ccc69d2cee09428832aa98", 504281, "287542747579596068100"},
		{"e2053e04f95afa5c8378677de62212edb972e21b40421786c53de57141853ce870481a80b68a449903479751da114a30d27a568812e178b919bd1e9b2c82f92e", 508822, "2004160057699128540000"},
		{"94bdbf207f6390354debfc2b3ff30ed101bc52d339e8310b2e2a1dd235cb9d40d27d65013c04030e05bbaceeea635dfdbfbbb47c683d29c6751f2bb3159e6abd", 557309, "900000000000000000000"},
		{"ed552a64f708696ac53962b88e927181688c8bc260787c82e1c9c21a62da4ce59c31fc594e48249e89392ce2e6e2a0320d6688b38ad7884ff6fe664faf4b12d9", 2405590, " 5823052137310000000000"},
	}

	nodes := make([]issue1654Candidate, 0)
	for _, c := range candidates {
		amount, _ := new(big.Int).SetString(c.Amount, 10)
		nodes = append(nodes, issue1654Candidate{
			nodeID:     discover.MustHexID(c.Node),
			stakingNum: uint64(c.Num),
			shouldSub:  amount,
		})
	}
	return nodes, nil
}
