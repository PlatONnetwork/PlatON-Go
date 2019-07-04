package vm_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	"testing"
)






func create_staking () (){

}


func TestRLP_encode (t *testing.T) {

	var params [][]byte
	params = make([][]byte, 0)

	fnType, err := rlp.EncodeToBytes(uint16(1100))
	if nil != err {
		fmt.Println("fnType err", err)
	}else {
		var num uint16
		rlp.DecodeBytes(fnType, &num)
		fmt.Println("num is ", num)
	}
	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("rlp stakingContract encode rlp data fail")
	} else {
		fmt.Println("rlp stakingContract data rlp: ", hexutil.Encode(buf.Bytes()))
	}
}


func TestStakingContract_createStaking(t *testing.T) {
	stakingContract := vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	  newEvm(),
	}

	//var govPlugin *plugin.GovPlugin

	plugin.GovPluginInstance()

	sndb := snapshotdb.Instance()
	sndb.NewBlock(big.NewInt(1), common.ZeroHash, blockHash)


	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benifitAddress, _ := rlp.EncodeToBytes(addrArr[1])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[0])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes("PlatON, China")
	website, _ := rlp.EncodeToBytes("https://www.platon.network")
	details, _ := rlp.EncodeToBytes("platon super node")
	StakeThreshold, _ := new(big.Int).SetString("1000000000000000000000000", 10)
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	processVersion, _ := rlp.EncodeToBytes(uint32(456))




	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benifitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, processVersion)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("createStaking encode rlp data fail")
	} else {
		fmt.Println("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}


	res, err := stakingContract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}
}



func TestStakingContract_editorCandidate(t *testing.T) {
	stakingContract := vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	  newEvm(),
	}

	//var govPlugin *plugin.GovPlugin

	plugin.GovPluginInstance()

	sndb := snapshotdb.Instance()
	sndb.NewBlock(big.NewInt(1), common.ZeroHash, blockHash)


	// create
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benifitAddress, _ := rlp.EncodeToBytes(addrArr[1])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[0])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes("PlatON, China")
	website, _ := rlp.EncodeToBytes("https://www.platon.network")
	details, _ := rlp.EncodeToBytes("platon super node")
	StakeThreshold, _ := new(big.Int).SetString("1000000000000000000000000", 10)
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	processVersion, _ := rlp.EncodeToBytes(uint32(456))




	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benifitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, processVersion)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("createStaking encode rlp data fail")
	} else {
		fmt.Println("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}


	res, err := stakingContract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}

	// edit





}


func TestStakingContract_getCandidateInfo (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()
	stakingContract := vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	  newEvm(),
	}

	fmt.Println("sender:", sender.Hex())
	//var govPlugin *plugin.GovPlugin

	plugin.GovPluginInstance()
	plugin.GetRestrictingInstance()

	sndb := snapshotdb.Instance()
	if err := sndb.NewBlock(blockNumer, common.ZeroHash, blockHash); nil != err {
		fmt.Println("newBlock err", err)
	}


	// create
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(1))
	benifitAddress, _ := rlp.EncodeToBytes(addrArr[1])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[0])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes("PlatON, China")
	website, _ := rlp.EncodeToBytes("https://www.platon.network")
	details, _ := rlp.EncodeToBytes("platon super node")
	StakeThreshold, _ := new(big.Int).SetString("1000000000000000000000000", 10)
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	processVersion, _ := rlp.EncodeToBytes(uint32(456))




	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benifitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, processVersion)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("createStaking encode rlp data fail")
	} else {
		fmt.Println("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}


	res, err := stakingContract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}

	sndb.Commit(blockHash)
	//sndb.Compaction()


	// get candidate Info
	params = make([][]byte, 0)

	fnType, _ = rlp.EncodeToBytes(uint16(1105))
	nodeId, _ = rlp.EncodeToBytes(nodeIdArr[0])

	params = append(params, fnType)
	params = append(params, nodeId)

	buf = new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("createStaking encode rlp data fail")
	} else {
		fmt.Println("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err = stakingContract.Run(buf.Bytes())
	if nil != err {
		t.Error("getCandidate err", err)
	}else {

		var r xcom.Result
		err = rlp.DecodeBytes(res, &r)
		if nil != err {
			fmt.Println(err)
		}

		if r.Status {
			dbyte, err := rlp.EncodeToBytes(r.Data)

			if nil != err {
				t.Error("rlp encode r.Data failed", err)
			}else {

				var c staking.Candidate

				if err = rlp.DecodeBytes(dbyte, &c); nil!= err {
					t.Error("decode failed", err)
				}else {

					rbyte, _ := json.Marshal(c)

					t.Log("", string(rbyte))

				}
			}
		}else {
			t.Error("getCandidate failed", r.ErrMsg)
		}
	}
}





func TestStakingContract_cleanSnapshotDB(t *testing.T) {
	sndb := snapshotdb.Instance()
	sndb.Clear()
}