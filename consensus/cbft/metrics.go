package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/metrics"
)

var (
	blockMinedTimer       = metrics.NewRegisteredTimer("cbft/timer/block/mined", nil)
	viewChangedTimer      = metrics.NewRegisteredTimer("cbft/timer/view/changed", nil)
	blockQCCollectedTimer = metrics.NewRegisteredTimer("cbft/timer/block/qc_collected", nil)
	blockExecutedTimer    = metrics.NewRegisteredTimer("cbft/timer/block/executed", nil)

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
