package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"runtime"

	"github.com/PlatONnetwork/PlatON-Go/core/vm"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
)

var (
	jsonFlag      = flag.String("json", "", "Json file path for contract parameters")
	funcTypesFlag = flag.String("funcTypes", "", "A collection of contract interface types")
)

// createStaking
type Ppos_1000 struct {
	Typ            uint16
	BenefitAddress common.Address
	NodeId         discover.NodeID
	ExternalId     string
	NodeName       string
	Website        string
	Details        string
	Amount         *big.Int
	ProgramVersion uint32
}

// editorCandidate
type Ppos_1001 struct {
	BenefitAddress common.Address
	NodeId         discover.NodeID
	ExternalId     string
	NodeName       string
	Website        string
	Details        string
}

// increaseStaking
type Ppos_1002 struct {
	NodeId discover.NodeID
	Typ    uint16
	Amount *big.Int
}

// withdrewCandidate
type Ppos_1003 struct {
	NodeId discover.NodeID
}

// delegate
type Ppos_1004 struct {
	Typ    uint16
	NodeId discover.NodeID
	Amount *big.Int
}

// withdrewDelegate
type Ppos_1005 struct {
	StakingBlockNum uint64
	NodeId          discover.NodeID
	Amount          *big.Int
}

// getRelatedListByDelAddr
type Ppos_1103 struct {
	Addr common.Address
}

// getDelegateInfo
type Ppos_1104 struct {
	StakingBlockNum uint64
	DelAddr         common.Address
	NodeId          discover.NodeID
}

// getCandidateInfo
type Ppos_1105 struct {
	NodeId discover.NodeID
}

// submitText
type Ppos_2000 struct {
	Verifier        discover.NodeID
	PIPID           string
	EndVotingRounds uint64
}

// submitVersion
type Ppos_2001 struct {
	Verifier        discover.NodeID
	PIPID           string
	NewVersion      uint32
	EndVotingRounds uint64
}

// submitCancel
type Ppos_2005 struct {
	Verifier        discover.NodeID
	PIPID           string
	EndVotingRounds uint64
	TobeCanceled    common.Hash
}

// vote
type Ppos_2003 struct {
	Verifier       discover.NodeID
	ProposalID     common.Hash
	Option         uint8
	ProgramVersion uint32
	VersionSign    common.VersionSign
}

type Ppos_20031 struct {
	Verifier       discover.NodeID
	ProposalID     common.Hash
	Option         uint8
	ProgramVersion uint32
	VersionSign    common.VersionSign
}

//declareVersion
type Ppos_2004 struct {
	Verifier       discover.NodeID
	ProgramVersion uint32
	VersionSign    common.VersionSign
}

// getProposal
type Ppos_2100 struct {
	ProposalID common.Hash
}

// getTallyResult
type Ppos_2101 struct {
	ProposalID common.Hash
}

// listProposal
type Ppos_2102 struct {
}

// getActiveVersion
type Ppos_2103 struct {
}

// getProgramVersion
type Ppos_2104 struct {
}

// ReportDuplicateSign
type Ppos_3000 struct {
	Data string
}

// CheckDuplicateSign
type Ppos_3001 struct {
	Etype       uint32
	Addr        common.Address
	BlockNumber uint64
}

// CreateRestrictingPlan
type Ppos_4000 struct {
	Account common.Address
	Plans   []restricting.RestrictingPlan
}

// GetRestrictingInfo
type Ppos_4100 struct {
	Account common.Address
}

type decDataConfig struct {
	P1000  Ppos_1000
	P1001  Ppos_1001
	P1002  Ppos_1002
	P1003  Ppos_1003
	P1004  Ppos_1004
	P1005  Ppos_1005
	P1103  Ppos_1103
	P1104  Ppos_1104
	P1105  Ppos_1105
	P2000  Ppos_2000
	P2001  Ppos_2001
	P2005  Ppos_2005
	P2003  Ppos_2003
	P20031 []Ppos_20031
	P2004  Ppos_2004
	P2100  Ppos_2100
	P2101  Ppos_2101
	P2102  Ppos_2102
	P2103  Ppos_2103
	P2104  Ppos_2104
	P3000  Ppos_3000
	P3001  Ppos_3001
	P4000  Ppos_4000
	P4100  Ppos_4100
}

func parseConfigJson(configPath string, v *decDataConfig) error {
	if configPath == "" {
		panic(fmt.Errorf("parse config file error"))
	}

	file, err := os.Open(configPath)
	if err != nil {
		utils.Fatalf("Failed to read config file: %v", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(v)
	if err != nil {
		panic(fmt.Errorf("parse config to json error,%s", err.Error()))
	}

	return nil
}

func getRlpData(funcType uint16, cfg *decDataConfig) string {
	rlpData := ""

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(funcType)
	params = append(params, fnType)

	switch funcType {
	case 1000:
		{
			typ, _ := rlp.EncodeToBytes(cfg.P1000.Typ)
			benefitAddress, _ := rlp.EncodeToBytes(cfg.P1000.BenefitAddress.Bytes())
			nodeId, _ := rlp.EncodeToBytes(cfg.P1000.NodeId)
			externalId, _ := rlp.EncodeToBytes(cfg.P1000.ExternalId)
			nodeName, _ := rlp.EncodeToBytes(cfg.P1000.NodeName)
			website, _ := rlp.EncodeToBytes(cfg.P1000.Website)
			details, _ := rlp.EncodeToBytes(cfg.P1000.Details)
			amount, _ := rlp.EncodeToBytes(cfg.P1000.Amount)
			programVersion, _ := rlp.EncodeToBytes(cfg.P1000.ProgramVersion)
			params = append(params, typ)
			params = append(params, benefitAddress)
			params = append(params, nodeId)
			params = append(params, externalId)
			params = append(params, nodeName)
			params = append(params, website)
			params = append(params, details)
			params = append(params, amount)
			params = append(params, programVersion)
		}
	case 1001:
		{
			benefitAddress, _ := rlp.EncodeToBytes(cfg.P1001.BenefitAddress.Bytes())
			nodeId, _ := rlp.EncodeToBytes(cfg.P1001.NodeId)
			externalId, _ := rlp.EncodeToBytes(cfg.P1001.ExternalId)
			nodeName, _ := rlp.EncodeToBytes(cfg.P1001.NodeName)
			website, _ := rlp.EncodeToBytes(cfg.P1001.Website)
			details, _ := rlp.EncodeToBytes(cfg.P1001.Details)
			params = append(params, benefitAddress)
			params = append(params, nodeId)
			params = append(params, externalId)
			params = append(params, nodeName)
			params = append(params, website)
			params = append(params, details)
		}
	case 1002:
		{
			nodeId, _ := rlp.EncodeToBytes(cfg.P1002.NodeId)
			typ, _ := rlp.EncodeToBytes(cfg.P1002.Typ)
			amount, _ := rlp.EncodeToBytes(cfg.P1002.Amount)
			params = append(params, nodeId)
			params = append(params, typ)
			params = append(params, amount)
		}
	case 1003:
		{
			nodeId, _ := rlp.EncodeToBytes(cfg.P1003.NodeId)
			params = append(params, nodeId)
		}
	case 1004:
		{
			typ, _ := rlp.EncodeToBytes(cfg.P1004.Typ)
			nodeId, _ := rlp.EncodeToBytes(cfg.P1004.NodeId)
			amount, _ := rlp.EncodeToBytes(cfg.P1004.Amount)
			params = append(params, typ)
			params = append(params, nodeId)
			params = append(params, amount)
		}
	case 1005:
		{
			stakingBlockNum, _ := rlp.EncodeToBytes(cfg.P1005.StakingBlockNum)
			nodeId, _ := rlp.EncodeToBytes(cfg.P1005.NodeId)
			amount, _ := rlp.EncodeToBytes(cfg.P1005.Amount)
			params = append(params, stakingBlockNum)
			params = append(params, nodeId)
			params = append(params, amount)
		}
	case 1100:
	case 1101:
	case 1102:
	case 1103:
		{
			addr, _ := rlp.EncodeToBytes(cfg.P1103.Addr.Bytes())
			params = append(params, addr)
		}
	case 1104:
		{
			stakingBlockNum, _ := rlp.EncodeToBytes(cfg.P1104.StakingBlockNum)
			delAddr, _ := rlp.EncodeToBytes(cfg.P1104.DelAddr.Bytes())
			nodeId, _ := rlp.EncodeToBytes(cfg.P1104.NodeId)

			params = append(params, stakingBlockNum)
			params = append(params, delAddr)
			params = append(params, nodeId)
		}
	case 1105:
		{
			nodeId, _ := rlp.EncodeToBytes(cfg.P1105.NodeId)
			params = append(params, nodeId)
		}
	case 2000:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2000.Verifier)
			pipID, _ := rlp.EncodeToBytes(cfg.P2000.PIPID)
			endVotingRounds, _ := rlp.EncodeToBytes(cfg.P2000.EndVotingRounds)
			params = append(params, verifier)
			params = append(params, pipID)
			params = append(params, endVotingRounds)
		}
	case 2001:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2001.Verifier)
			pipID, _ := rlp.EncodeToBytes(cfg.P2000.PIPID)
			newVersion, _ := rlp.EncodeToBytes(cfg.P2001.NewVersion)
			endVotingRounds, _ := rlp.EncodeToBytes(cfg.P2000.EndVotingRounds)
			params = append(params, verifier)
			params = append(params, pipID)
			params = append(params, newVersion)
			params = append(params, endVotingRounds)
		}
	case 2005:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2005.Verifier)
			pipID, _ := rlp.EncodeToBytes(cfg.P2000.PIPID)
			endVotingRounds, _ := rlp.EncodeToBytes(cfg.P2005.EndVotingRounds)
			tobeCanceled, _ := rlp.EncodeToBytes(cfg.P2005.TobeCanceled)
			params = append(params, verifier)
			params = append(params, pipID)
			params = append(params, endVotingRounds)
			params = append(params, tobeCanceled)
		}
	case 2003:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2003.Verifier)
			proposalID, _ := rlp.EncodeToBytes(cfg.P2003.ProposalID.Bytes())
			op, _ := rlp.EncodeToBytes(cfg.P2003.Option)
			programVersion, _ := rlp.EncodeToBytes(cfg.P2003.ProgramVersion)
			versionSign, _ := rlp.EncodeToBytes(cfg.P2003.VersionSign)
			params = append(params, verifier)
			params = append(params, proposalID)
			params = append(params, op)
			params = append(params, programVersion)
			params = append(params, versionSign)
		}
	case 20031:
		{
			for i := 0; i < len(cfg.P20031); i++ {
				fnType, _ := rlp.EncodeToBytes(2003)
				verifier, _ := rlp.EncodeToBytes(cfg.P20031[i].Verifier)
				proposalID, _ := rlp.EncodeToBytes(cfg.P20031[i].ProposalID.Bytes())
				op, _ := rlp.EncodeToBytes(cfg.P20031[i].Option)
				programVersion, _ := rlp.EncodeToBytes(cfg.P20031[i].ProgramVersion)
				versionSign, _ := rlp.EncodeToBytes(cfg.P20031[i].VersionSign)
				params = make([][]byte, 0)
				params = append(params, fnType)
				params = append(params, verifier)
				params = append(params, proposalID)
				params = append(params, op)
				params = append(params, programVersion)
				params = append(params, versionSign)

				buf := new(bytes.Buffer)
				err := rlp.Encode(buf, params)
				if err != nil {
					panic(fmt.Errorf("%d encode rlp data fail: %v", funcType, err))
				} else {
					rlpData = hexutil.Encode(buf.Bytes())
					fmt.Printf("RLP= %s\n", rlpData)
				}
			}
			return ""
		}

	case 2004:
		{
			verifier, _ := rlp.EncodeToBytes(cfg.P2004.Verifier)
			programVersion, _ := rlp.EncodeToBytes(cfg.P2004.ProgramVersion)
			versionSign, _ := rlp.EncodeToBytes(cfg.P2004.VersionSign)
			params = append(params, verifier)
			params = append(params, programVersion)
			params = append(params, versionSign)
		}
	case 2100:
		{
			proposalID, _ := rlp.EncodeToBytes(cfg.P2100.ProposalID.Bytes())
			params = append(params, proposalID)
		}
	case 2101:
		{
			proposalID, _ := rlp.EncodeToBytes(cfg.P2101.ProposalID.Bytes())
			params = append(params, proposalID)
		}
	case 2102:
	case 2103:
	case 2104:
	case 3000:
		{
			data, _ := rlp.EncodeToBytes(cfg.P3000.Data)
			params = append(params, data)
		}
	case 3001:
		{
			etype, _ := rlp.EncodeToBytes(cfg.P3001.Etype)
			addr, _ := rlp.EncodeToBytes(cfg.P3001.Addr.Bytes())
			blockNumber, _ := rlp.EncodeToBytes(cfg.P3001.BlockNumber)
			params = append(params, etype)
			params = append(params, addr)
			params = append(params, blockNumber)
		}
	case 4000:
		{
			account, _ := rlp.EncodeToBytes(cfg.P4000.Account.Bytes())
			plans, _ := rlp.EncodeToBytes(cfg.P4000.Plans)
			params = append(params, account)
			params = append(params, plans)
		}
	case 4100:
		{
			account, _ := rlp.EncodeToBytes(cfg.P4100.Account.Bytes())
			params = append(params, account)
		}
	default:
		{
			//	panic(fmt.Errorf("funcType:%d is unknown!!!!", funcType))
			fmt.Printf("funcType:%d is unknown!!!!\n", funcType)
			return ""
		}
	}

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		panic(fmt.Errorf("%d encode rlp data fail: %v", funcType, err))
	} else {
		rlpData = hexutil.Encode(buf.Bytes())
		fmt.Printf("funcType:%d rlp data = %s\n", funcType, rlpData)
	}

	return rlpData
}

func Verify_tx_data(input []byte, command map[uint16]interface{}) (fn interface{}, FnParams []reflect.Value, err error) {

	defer func() {
		if er := recover(); nil != er {
			fn, FnParams, err = nil, nil, fmt.Errorf("parse tx data is panic: %s", er)
			log.Error("Failed to Verify PlatON inner contract tx data", "error", er)
		}
	}()

	var args [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &args); nil != err {
		log.Error("Failed to Verify PlatON inner contract tx data, Decode rlp input failed", "err", err)
		return nil, nil, err
	}

	//fmt.Println("the Function Type:", byteutil.BytesToUint16(args[0]))

	if fn, ok := command[byteutil.BytesToUint16(args[0])]; !ok {
		return nil, nil, err
	} else {

		funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		fmt.Println("The FuncName is", funcName)

		t := reflect.TypeOf(fn)
		fmt.Println("The FN", fn)
		fmt.Println("The TypeOf.Name", t.Name())
		fmt.Println("The TypeOf.String", t.String())

		v := reflect.ValueOf(fn)
		fmt.Println("The ValueOf", v.Type().Name())
		fmt.Println("The ValueOf.String", v.String())

		// the func params type list
		paramList := reflect.TypeOf(fn)
		// the func params len
		paramNum := paramList.NumIn()

		if paramNum != len(args)-1 {
			return nil, nil, errors.New("para num error")
		}
		params := make([]reflect.Value, paramNum)

		for i := 0; i < paramNum; i++ {
			//fmt.Println("byte:", args[i+1])

			targetType := paramList.In(i).String()
			inputByte := []reflect.Value{reflect.ValueOf(args[i+1])}
			params[i] = reflect.ValueOf(byteutil.Bytes2X_CMD[targetType]).Call(inputByte)[0]
			//fmt.Println("num", i+1, "type", targetType)
		}
		return fn, params, nil
	}
}

func main() {
	data := "0xe683820834a1a0373e89d01414ff4b02a638599b093c2a5cb7ae5a9c30c2653a451b320ec28ffe"
	bs, _ := hexutil.Decode(data)

	gc := &vm.GovContract{}
	if fn, _, err := Verify_tx_data(bs, gc.FnSigns()); err != nil {
		fmt.Print(err)
	} else {
		fmt.Print(fn)
	}

}

/*func main() {
	// Parse and ensure all needed inputs are specified
	flag.Parse()

	if *jsonFlag == "" && *funcTypesFlag == "" {
		fmt.Printf("json path is null or funcid is null\n")
		os.Exit(-1)
	}

	cfg := decDataConfig{}
	parseConfigJson(*jsonFlag, &cfg)

	funcTypes := strings.Split(*funcTypesFlag, "|")

	for _, fnType := range funcTypes {
		funcType, e := strconv.Atoi(fnType)
		if e != nil {
			fmt.Printf("funcType is error\n")
			os.Exit(-1)
		}
		getRlpData(uint16(funcType), &cfg)
	}

}
*/
