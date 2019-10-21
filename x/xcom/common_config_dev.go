// +build test

package xcom

import "github.com/PlatONnetwork/PlatON-Go/log"

func init() {
	log.Info("Init ppos common config", "network name", "DefaultTestNet", "network value", DefaultTestNet)
	GetEc(DefaultTestNet)
}
