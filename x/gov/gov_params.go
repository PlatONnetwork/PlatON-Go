package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
)

var paramList = []*GovernParam{
	{
		ParamItem:     &ParamItem{"PPOS", "paramName1", "paramName1"},
		ParamValue:    &ParamValue{"", "10", 0},
		ParamVerifier: func(value string) bool { return true },
	},
	{
		ParamItem:     &ParamItem{"PPOS", "paramName2", "paramName2"},
		ParamValue:    &ParamValue{"", "100000", 0},
		ParamVerifier: func(value string) bool { return true },
	},
}

func InitGenesisGovernParam(snapDB snapshotdb.DB) error {
	var paramItemList []*ParamItem
	for _, param := range paramList {
		paramItemList = append(paramItemList, param.ParamItem)

		key := KeyParamValue(param.ParamItem.Module, param.ParamItem.Name)
		value := common.MustRlpEncode(param.ParamValue)
		if err := snapDB.PutBaseDB(key, value); err != nil {
			return err
		}
	}

	key := KeyParamItems()
	value := common.MustRlpEncode(paramItemList)
	if err := snapDB.PutBaseDB(key, value); err != nil {
		return err
	}
	return nil
}

func RegisterGovernParamVerifiers() {
	for _, param := range paramList {
		RegGovernParamVerifier(param.ParamItem.Module, param.ParamItem.Name, param.ParamVerifier)
	}
}
