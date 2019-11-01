package xcom

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func get(blockHash common.Hash, key []byte) ([]byte, error) {
	return snapshotdb.Instance().Get(blockHash, key)
}

func put(blockHash common.Hash, key []byte, value interface{}) error {
	bytes, err := rlp.EncodeToBytes(value)
	if err != nil {
		return err
	}
	return snapshotdb.Instance().Put(blockHash, key, bytes)
}

func del(blockHash common.Hash, key []byte) error {
	return snapshotdb.Instance().Del(blockHash, key)
}

func AddGovernParam(module, name, desc string, paramValue *ParamValue, blockHash common.Hash) error {
	itemList, err := ListGovernParamItem("", blockHash)
	if err != nil {
		return nil
	}
	itemList = append(itemList, &ParamItem{module, name, desc})
	if err := put(blockHash, keyPrefixParamItems, itemList); err != nil {
		return err
	}

	if err := put(blockHash, KeyParamValue(module, name), paramValue); err != nil {
		return err
	}
	return nil
}

func FindGovernParamValue(module, name string, blockHash common.Hash) (*ParamValue, error) {
	value, err := get(blockHash, KeyParamValue(module, name))
	if snapshotdb.NonDbNotFoundErr(err) {
		return nil, err
	}

	if len(value) > 0 {
		var paramValue ParamValue
		if err := rlp.DecodeBytes(value, &paramValue); err != nil {
			return nil, err
		} else {
			return &paramValue, nil
		}
	}
	return nil, nil
}

func UpdateGovernParamValue(module, name, newValue string, activeBlock uint64, blockHash common.Hash) error {
	value, err := get(blockHash, KeyParamValue(module, name))
	if snapshotdb.NonDbNotFoundErr(err) {
		return err
	}
	if len(value) > 0 {
		var paramValue ParamValue
		if err := rlp.DecodeBytes(value, &paramValue); err != nil {
			return err
		}
		paramValue.StaleValue = paramValue.Value
		paramValue.Value = newValue
		paramValue.ActiveBlock = activeBlock

		if err := put(blockHash, KeyParamValue(module, name), paramValue); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Not found the %s.%s Govern value", module, name)
}

func ListGovernParam(module string, blockHash common.Hash) ([]*GovernParam, error) {
	itemList, err := ListGovernParamItem(module, blockHash)
	if err != nil {
		return nil, err
	}
	var paraList []*GovernParam
	for _, item := range itemList {
		if value, err := FindGovernParamValue(item.Module, item.Name, blockHash); err != nil {
			return nil, err
		} else {
			param := &GovernParam{item, value, nil}
			paraList = append(paraList, param)
		}
	}
	return paraList, nil
}

func ListGovernParamItem(module string, blockHash common.Hash) ([]*ParamItem, error) {
	itemBytes, err := get(blockHash, KeyParamItems())
	if snapshotdb.NonDbNotFoundErr(err) {
		return nil, err
	}

	if len(itemBytes) > 0 {
		var itemList []*ParamItem
		if err := rlp.DecodeBytes(itemBytes, &itemList); err != nil {
			return nil, err
		}
		if len(module) == 0 {
			return itemList, nil
		} else {
			idx := 0
			for _, item := range itemList {
				if item.Module == module {
					itemList[idx] = item
					idx++
				}
			}
			return itemList[:idx], nil
		}
	}
	return nil, nil
}
