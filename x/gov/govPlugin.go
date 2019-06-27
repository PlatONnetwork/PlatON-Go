package gov
//
//import (
//	"github.com/PlatONnetwork/PlatON-Go/core/types"
//	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
//	"sync"
//)
//
//var pluginOnce sync.Once
//var govPlugin *GovPlugin
//
//type GovPlugin struct {
//	gov *Gov
//}
//
//// Instance a global GovPlugin
//func GovPluginInstance(gov *Gov) *GovPlugin {
//	pluginOnce.Do(func() {
//		govPlugin = &GovPlugin{gov: gov}
//	})
//	return govPlugin
//}
//
//func (govPlugin *GovPlugin) BeginBlock(header *types.Header, state xcom.StateDB) (bool, error) {
//	return gov.BeginBlock(header.Hash(), state)
//}
//
//func (govPlugin *GovPlugin) EndBlock(header *types.Header, state xcom.StateDB) (bool, error) {
//	return gov.EndBlock(header.Hash(), state)
//}
//
//func (govPlugin *GovPlugin) Confirmed(block *types.Block) error {
//	return nil
//}
