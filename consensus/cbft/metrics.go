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

package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/metrics"
)

var (
	blockMinedGauage       = metrics.NewRegisteredGauge("cbft/gauage/block/mined", nil)
	viewChangedTimer       = metrics.NewRegisteredTimer("cbft/timer/view/changed", nil)
	blockQCCollectedGauage = metrics.NewRegisteredGauge("cbft/gauage/block/qc_collected", nil)
	blockExecutedGauage    = metrics.NewRegisteredGauge("cbft/gauage/block/executed", nil)

	blockProduceMeter          = metrics.NewRegisteredMeter("cbft/meter/block/produce", nil)
	blockCheckFailureMeter     = metrics.NewRegisteredMeter("cbft/meter/block/check_failure", nil)
	signatureCheckFailureMeter = metrics.NewRegisteredMeter("cbft/meter/signature/check_failure", nil)
	blockConfirmedMeter        = metrics.NewRegisteredMeter("cbft/meter/block/confirmed", nil)

	masterCounter    = metrics.NewRegisteredCounter("cbft/counter/view/count", nil)
	consensusCounter = metrics.NewRegisteredCounter("cbft/counter/consensus/count", nil)
	minedCounter     = metrics.NewRegisteredCounter("cbft/counter/mined/count", nil)

	viewNumberGauage          = metrics.NewRegisteredGauge("cbft/gauage/view/number", nil)
	epochNumberGauage         = metrics.NewRegisteredGauge("cbft/gauage/epoch/number", nil)
	proposerIndexGauage       = metrics.NewRegisteredGauge("cbft/gauage/proposer/index", nil)
	validatorCountGauage      = metrics.NewRegisteredGauge("cbft/gauage/validator/count", nil)
	blockNumberGauage         = metrics.NewRegisteredGauge("cbft/gauage/block/number", nil)
	highestQCNumberGauage     = metrics.NewRegisteredGauge("cbft/gauage/block/qc/number", nil)
	highestLockedNumberGauage = metrics.NewRegisteredGauge("cbft/gauage/block/locked/number", nil)
	highestCommitNumberGauage = metrics.NewRegisteredGauge("cbft/gauage/block/commit/number", nil)
)
