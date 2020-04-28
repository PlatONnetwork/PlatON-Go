package plugin

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
)

// Provides an API interface to obtain data related to the economic model
type PublicPPOSAPI struct {
	snapshotDB snapshotdb.DB
}

func NewPublicPPOSAPI() *PublicPPOSAPI {
	return &PublicPPOSAPI{snapshotdb.Instance()}
}

// Get node list of zero-out blocks
func (p *PublicPPOSAPI) GetWaitSlashingNodeList() string {
	list, err := slash.getWaitSlashingNodeList(0, common.ZeroHash)
	if nil != err || len(list) == 0 {
		return ""
	}
	return fmt.Sprintf("%+v", list)
}
