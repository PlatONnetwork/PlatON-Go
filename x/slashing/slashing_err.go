package slashing

import "github.com/PlatONnetwork/PlatON-Go/common"

var (
	ErrDuplicateSignVerify = common.NewBizError(303000, "duplicate signature verification failed")
	ErrSlashingExist       = common.NewBizError(303001, "punishment has been implemented")
	ErrBlockNumberTooHigh  = common.NewBizError(303002, "blockNumber too high")
	ErrIntervalTooLong     = common.NewBizError(303003, "evidence interval is too long")
	ErrGetCandidate        = common.NewBizError(303004, "failed to get certifier information")
	ErrAddrMismatch        = common.NewBizError(303005, "address does not match")
	ErrNodeIdMismatch      = common.NewBizError(303006, "nodeId does not match")
	ErrBlsPubKeyMismatch   = common.NewBizError(303007, "blsPubKey does not match")
	ErrSlashingFail        = common.NewBizError(303008, "slashing node fail")
	ErrNotValidator        = common.NewBizError(303009, "This node is not a validator")
	ErrSameAddr            = common.NewBizError(303010, "Can't report yourself")
)
