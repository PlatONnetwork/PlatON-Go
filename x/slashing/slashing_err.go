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
