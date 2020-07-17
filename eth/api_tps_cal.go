package eth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/rpc"

	"github.com/tealeg/xlsx"
)

const DefaultViewNumber = uint64(0)

type ViewCountMap map[uint64]uint64

type AnalystEntity struct {
	beginNumber        uint64
	endNumber          uint64
	viewBlockRate      uint64
	viewCountMap       ViewCountMap
	missViewList       []uint64
	totalProduceTime   uint64
	averageProduceTime uint64
	topArray           []uint64
	txCount            uint64
	tps                uint64
}

func (txg *TxGenAPI) GetTps(ctx context.Context, beginBn, endBn uint64, interval uint64, resultPath string) error {
	if beginBn >= endBn || endBn < interval || endBn%interval != 0 || beginBn%interval != 1 {
		return errors.New(fmt.Sprintf("Invalid parameter, beginBn: %d, endBn: %d, interval: %d \n", beginBn, endBn, interval))
	}

	// cal current block hight
	currentNumber, _ := txg.eth.APIBackend.HeaderByNumber(ctx, rpc.LatestBlockNumber) // latest header should always be available

	if currentNumber.Number.Uint64() < beginBn+interval-1 {
		return errors.New(fmt.Sprintf("The current block number is too low to require statistics, beginBn: %d, endBn: %d, interval: %d, currentNumber: %d \n", beginBn, endBn, interval, currentNumber))
	}
	analystData := make([]*AnalystEntity, 0)
	// current round
	round := (endBn - beginBn + 1) / interval
	for i := uint64(0); i < round; i++ {
		beginNumber := beginBn
		endNumber := beginNumber + interval - 1

		for {
			currentNumber, _ = txg.eth.APIBackend.HeaderByNumber(ctx, rpc.LatestBlockNumber) // latest header should always be available
			if endNumber <= currentNumber.Number.Uint64() {
				break
			} else {
				//log2.Printf("Remote number is low, please wait...remote: %d, beginNumber: %d, endNumber: %d \n", currentNumber, beginNumber, endNumber)
				time.Sleep(5000 * time.Millisecond)
			}
		}

		// cal block time
		totalProduceTime, averageProduceTime, topArray, txCount, tps := AnalystProduceTime(beginNumber, endNumber, txg.eth.APIBackend)
		// cal view
		_, viewCountMap, missViewList, viewBlockRate, err := AnalystView(beginNumber, endNumber, txg.eth.Engine())
		if err != nil {
			return err
		}

		beginBn = beginNumber + interval

		// export excel
		entity := &AnalystEntity{
			beginNumber:        beginNumber,
			endNumber:          endNumber,
			viewBlockRate:      viewBlockRate,
			viewCountMap:       viewCountMap,
			missViewList:       missViewList,
			totalProduceTime:   totalProduceTime,
			averageProduceTime: averageProduceTime,
			topArray:           topArray,
			txCount:            txCount,
			tps:                tps,
		}
		analystData = append(analystData, entity)
	}
	return saveExcel(analystData, resultPath)
}

/*
	output parameter
		diffTimestamp 				current epoch  produce block use time(ms)
		diffTimestamp / diffNumber	Average block time（ms）
		topArray					The top 10 time-consuming blocks
		txCount						Total transactions
		tps							tps
*/
func AnalystProduceTime(beginNumber uint64, endNumber uint64, backend *EthAPIBackend) (uint64, uint64, []uint64, uint64, uint64) {
	beginHeader, _ := backend.HeaderByNumber(context.Background(), rpc.BlockNumber(beginNumber))
	endHeader, _ := backend.HeaderByNumber(context.Background(), rpc.BlockNumber(endNumber))

	preTimestamp := beginHeader.Time.Uint64()
	topArray := make([]uint64, 0, 250)
	for i := beginNumber + 1; i <= endNumber; i++ {
		header, _ := backend.HeaderByNumber(context.Background(), rpc.BlockNumber(int64(i)))
		diff := header.Time.Uint64() - preTimestamp
		topArray = append(topArray, diff)
		preTimestamp = header.Time.Uint64()
	}

	diffTimestamp := endHeader.Time.Uint64() - beginHeader.Time.Uint64()
	diffNumber := endHeader.Number.Uint64() - beginHeader.Number.Uint64() + 1

	// To transactions
	txCount := uint64(0)
	bh, _ := backend.HeaderByNumber(context.Background(), rpc.BlockNumber(int64(beginNumber)))
	eh, _ := backend.HeaderByNumber(context.Background(), rpc.BlockNumber(int64(endNumber)))
	for i := beginNumber; i <= endNumber; i++ {
		h, _ := backend.BlockByNumber(context.Background(), rpc.BlockNumber(int64(i)))
		c := hexutil.Uint(len(h.Transactions()))
		txCount = txCount + uint64(c)
	}
	tps := (txCount * 1000) / (eh.Time.Uint64() - bh.Time.Uint64())
	return diffTimestamp, diffTimestamp / diffNumber, topArray, txCount, tps
}

/*
	output parameter
		epoch
		viewCountMap	each view produce blocks
		missViewList	missing view
		viewBlockRate   view produce block rate
*/
func AnalystView(beginNumber uint64, endNumber uint64, engine consensus.Engine) (uint64, ViewCountMap, []uint64, uint64, error) {
	beginQC := engine.GetPrepareQC(beginNumber)
	endQC := engine.GetPrepareQC(endNumber)
	if beginQC.Epoch != endQC.Epoch {
		return 0, nil, nil, 0, fmt.Errorf("Epoch is inconsistent")
	}
	epoch := beginQC.Epoch

	viewCountMap := make(ViewCountMap, 0)
	missViewList := make([]uint64, 0)
	// each view produce blocks
	for i := beginNumber; i <= endNumber; i++ {
		qc := engine.GetPrepareQC(i) // Get PrepareQC by blockNumber
		if count, ok := viewCountMap[qc.ViewNumber]; ok {
			viewCountMap[qc.ViewNumber] = count + 1
		} else {
			viewCountMap[qc.ViewNumber] = 1
		}
	}
	// missing view
	for i := DefaultViewNumber; i <= endQC.ViewNumber; i++ {
		if _, ok := viewCountMap[i]; !ok {
			missViewList = append(missViewList, i)
		}
	}
	// view produce block rate
	viewBlockRate := (endNumber - beginNumber + 1) * 100 / ((endQC.ViewNumber - DefaultViewNumber + 1) * 10)
	return epoch, viewCountMap, missViewList, viewBlockRate, nil
}

func saveExcel(data []*AnalystEntity, resultPath string) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Block statistics")
	if err != nil {
		log.Println(err.Error())
	}

	// add title
	row := sheet.AddRow()
	cell_1 := row.AddCell()
	cell_1.Value = "Start block"
	cell_2 := row.AddCell()
	cell_2.Value = "End block"
	cell_3 := row.AddCell()
	cell_3.Value = "view produce block rate"
	cell_4 := row.AddCell()
	cell_4.Value = "View actual produce number of blocks"
	cell_5 := row.AddCell()
	cell_5.Value = "Missing view"
	cell_6 := row.AddCell()
	cell_6.Value = "Produce block time (ms)"
	cell_7 := row.AddCell()
	cell_7.Value = "Average produce block time (ms)"
	cell_8 := row.AddCell()
	cell_8.Value = "Total transactions"
	cell_9 := row.AddCell()
	cell_9.Value = "TPS"
	cell_10 := row.AddCell()
	cell_10.Value = "Block interval"

	//add data
	for _, d := range data {
		row := sheet.AddRow()
		beginNumber := row.AddCell()
		beginNumber.Value = strconv.Itoa(int(d.beginNumber))
		endNumber := row.AddCell()
		endNumber.Value = strconv.Itoa(int(d.endNumber))
		viewBlockRate := row.AddCell()
		viewBlockRate.Value = strconv.Itoa(int(d.viewBlockRate))
		viewCountMap := row.AddCell()
		viewCountMap.Value = fmt.Sprintf("%v", d.viewCountMap)
		missViewList := row.AddCell()
		missViewList.Value = fmt.Sprintf("%v", d.missViewList)
		totalProduceTime := row.AddCell()
		totalProduceTime.Value = strconv.Itoa(int(d.totalProduceTime))
		averageProduceTime := row.AddCell()
		averageProduceTime.Value = strconv.Itoa(int(d.averageProduceTime))
		txCount := row.AddCell()
		txCount.Value = strconv.Itoa(int(d.txCount))
		tps := row.AddCell()
		tps.Value = strconv.Itoa(int(d.tps))
		topArray := row.AddCell()
		topArray.Value = fmt.Sprintf("%v", d.topArray)
	}
	err = file.Save(resultPath)
	if err != nil {
		return err
	}
	return nil
}
