package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"

	"github.com/PlatONnetwork/PlatON-Go/rpc"

	"github.com/tealeg/xlsx"
)

const DefaultViewNumber = uint64(0)

type ViewCountMap map[uint64]uint64

type AnalystEntity struct {
	BeginNumber        uint64
	EndNumber          uint64
	ViewBlockRate      uint64
	ViewCountMap       ViewCountMap
	MissViewList       []uint64
	TotalProduceTime   uint64
	AverageProduceTime uint64
	TopArray           [][]uint64
	TxCount            uint64
	Tps                uint64
}

func (txg *TxGenAPI) CalRes(configPaths []string, output string, t int) error {
	x := make(BlockInfos, 0)
	sendTotal := uint64(0)
	for _, path := range configPaths {
		file, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			return fmt.Errorf("Failed to open config file:%v", err)
		}
		defer file.Close()
		var res TxGenResData
		if err := json.NewDecoder(file).Decode(&res); err != nil {
			return fmt.Errorf("invalid res file r:%v", err)
		}

		for _, ttf := range res.Ttf {
			x = append(x, ttf)
		}
		sendTotal += res.TotalTxSend
	}
	sort.Sort(x)
	endTime := common.MillisToTime(x[0].ProduceTime).Add(time.Second * time.Duration(t))
	txConut := 0
	timeUse := int64(0)
	analysts := make([][3]int64, 0)
	total := 0

	for _, ttf := range x {
		total += ttf.TxLength
		if !common.MillisToTime(ttf.ProduceTime).Before(endTime) {
			analysts = append(analysts, [3]int64{endTime.Unix(), time.Duration(int64(float64(timeUse) / float64(txConut))).Milliseconds(), int64(txConut) / int64(t)})
			endTime = endTime.Add(time.Second * time.Duration(t))
			txConut = 0
			timeUse = 0
		}
		txConut += ttf.TxLength
		timeUse += ttf.TimeUse
	}

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("ttf statistics")
	if err != nil {
		return err
	}

	// add title
	row := sheet.AddRow()
	cell_1 := row.AddCell()
	cell_1.Value = "time"
	cell_2 := row.AddCell()
	cell_2.Value = "ttf"
	cell_3 := row.AddCell()
	cell_3.Value = "tps"
	cell_4 := row.AddCell()
	cell_4.Value = "totalReceive"
	cell_5 := row.AddCell()
	cell_5.Value = "totalSend"

	//add data
	for i, d := range analysts {
		row := sheet.AddRow()
		time := row.AddCell()
		time.Value = strconv.FormatInt(d[0], 10)
		ttf := row.AddCell()
		ttf.Value = strconv.FormatInt(d[1], 10)
		tps := row.AddCell()
		tps.Value = strconv.FormatInt(d[2], 10)
		if i == 0 {
			totalReceive := row.AddCell()
			totalReceive.Value = strconv.FormatInt(int64(total), 10)
			totalSend := row.AddCell()
			totalSend.Value = strconv.FormatInt(int64(sendTotal), 10)
		}
	}
	err = xlsxFile.Save(output)
	if err != nil {
		return err
	}
	return nil
}

type Tps struct {
	BlockProduceTime time.Time
	TxLength         int
}

func (txg *TxGenAPI) CalBlockTps(ctx context.Context, beginBn, endBn uint64, output string) error {

	res := make([]Tps, 0)

	for i := uint64(beginBn); i < endBn; i++ {
		block, err := txg.eth.APIBackend.BlockByNumber(ctx, rpc.BlockNumber(i))
		if err != nil {
			return err
		}
		res = append(res, Tps{common.MillisToTime(block.Header().Time.Int64()), block.Transactions().Len()})
	}
	t := 10
	endTime := res[0].BlockProduceTime.Add(time.Second * time.Duration(t))
	blockCount := 0
	beginTimestamp := int64(0)
	endTimestamp := int64(0)
	txConut := 0
	analysts := make([][3]int64, 0)
	for i := 0; i < len(res); i++ {
		if res[i].BlockProduceTime.Before(endTime) {
			blockCount += 1
			if blockCount == 1 {
				beginTimestamp = common.Millis(res[i].BlockProduceTime)
			}
			txConut += res[i].TxLength
		} else {
			endTimestamp = common.Millis(res[i-1].BlockProduceTime)
			fmt.Println("rowInfo", "endTime", endTime.Unix(), "beginTimestamp", beginTimestamp, "endTimestamp", endTimestamp, "txConut", txConut, "blockCount", blockCount)
			if blockCount == 1 && beginTimestamp == endTimestamp {
				analysts = append(analysts, [3]int64{endTime.Unix(), int64(txConut / t), 0})
			} else {
				analysts = append(analysts, [3]int64{endTime.Unix(), int64(txConut / t), (endTimestamp - beginTimestamp) / int64(blockCount-1)})
			}
			endTime = endTime.Add(time.Second * time.Duration(t))
			blockCount = 0
			beginTimestamp = int64(0)
			endTimestamp = int64(0)
			txConut = 0
			blockCount += 1
			beginTimestamp = common.Millis(res[i].BlockProduceTime)
			txConut += res[i].TxLength
		}
	}

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("tps statistics")
	if err != nil {
		return err
	}

	// add title
	row := sheet.AddRow()
	cell_1 := row.AddCell()
	cell_1.Value = "time"
	cell_2 := row.AddCell()
	cell_2.Value = "tps"
	cell_3 := row.AddCell()
	cell_3.Value = "avg interval"

	//add data
	for _, d := range analysts {
		row := sheet.AddRow()
		beginNumber := row.AddCell()
		beginNumber.Value = strconv.FormatInt(d[0], 10)
		tps := row.AddCell()
		tps.Value = strconv.FormatInt(d[1], 10)
		interval := row.AddCell()
		interval.Value = strconv.Itoa(int(d[2]))
	}
	err = xlsxFile.Save(output)
	if err != nil {
		return err
	}
	return nil
}

func (txg *TxGenAPI) CalBlockAnalyst(ctx context.Context, beginBn, endBn uint64, interval uint64, resultPath string) ([]*AnalystEntity, error) {
	if beginBn >= endBn || endBn < interval || endBn%interval != 0 || beginBn%interval != 1 {
		return nil, fmt.Errorf("Invalid parameter, beginBn: %d, endBn: %d, interval: %d \n", beginBn, endBn, interval)
	}

	// cal current block hight
	currentNumber, _ := txg.eth.APIBackend.HeaderByNumber(ctx, rpc.LatestBlockNumber) // latest header should always be available

	if currentNumber.Number.Uint64() < beginBn+interval-1 {
		return nil, fmt.Errorf("The current block number is too low to require statistics, beginBn: %d, endBn: %d, interval: %d, currentNumber: %d \n", beginBn, endBn, interval, currentNumber)
	}

	if endBn > currentNumber.Number.Uint64() {
		return nil, fmt.Errorf("the endBn is grearter than current block, beginBn: %d, endBn: %d, interval: %d, currentNumber: %d \n", beginBn, endBn, interval, currentNumber.Number)
	}

	analystData := make([]*AnalystEntity, 0)
	// current round
	round := (endBn - beginBn + 1) / interval
	for i := uint64(0); i < round; i++ {
		beginNumber := beginBn
		endNumber := beginNumber + interval - 1

		// cal block time and view
		totalProduceTime, averageProduceTime, topArray, txCount, tps, viewCountMap, missViewList, viewBlockRate, err := AnalystProduceTimeAndView(beginNumber, endNumber, txg.eth.APIBackend)
		if err != nil {
			return nil, err
		}

		beginBn = beginNumber + interval

		// export excel
		entity := &AnalystEntity{
			BeginNumber:        beginNumber,
			EndNumber:          endNumber,
			ViewBlockRate:      viewBlockRate,
			ViewCountMap:       viewCountMap,
			MissViewList:       missViewList,
			TotalProduceTime:   totalProduceTime,
			AverageProduceTime: averageProduceTime,
			TopArray:           topArray,
			TxCount:            txCount,
			Tps:                tps,
		}
		analystData = append(analystData, entity)
	}
	if resultPath != "" {
		if err := saveExcel(analystData, resultPath); err != nil {
			return nil, err
		}
	}
	return analystData, nil
}

/*
	output parameter
		diffTimestamp 				current epoch  produce block use time(ms)
		diffTimestamp / diffNumber	Average block time（ms）
		TopArray					The top 10 time-consuming blocks
		TxCount						Total transactions
		Tps							Tps
		ViewCountMap	each view produce blocks
		MissViewList	missing view
		ViewBlockRate   view produce block rate
*/
func AnalystProduceTimeAndView(beginNumber uint64, endNumber uint64, backend *EthAPIBackend) (uint64, uint64, [][]uint64, uint64, uint64, ViewCountMap, []uint64, uint64, error) {
	ctx := context.Background()
	beginBlock, _ := backend.BlockByNumber(ctx, rpc.BlockNumber(beginNumber))
	endBlock, _ := backend.BlockByNumber(ctx, rpc.BlockNumber(endNumber))

	_, beginQC, err := ctypes.DecodeExtra(beginBlock.ExtraData())
	if err != nil {
		return 0, 0, nil, 0, 0, nil, nil, 0, fmt.Errorf("decodeExtra beginHeader Extra fail:%v", err)
	}

	_, endQC, err := ctypes.DecodeExtra(endBlock.ExtraData())
	if err != nil {
		return 0, 0, nil, 0, 0, nil, nil, 0, fmt.Errorf("decodeExtra endHeader Extra fail:%v", err)
	}

	if beginQC.Epoch != endQC.Epoch {
		return 0, 0, nil, 0, 0, nil, nil, 0, fmt.Errorf("Epoch is inconsistent")
	}

	viewCountMap := make(ViewCountMap, 0)
	missViewList := make([]uint64, 0)

	beginHeader := beginBlock.Header()
	endHeader := endBlock.Header()

	preTimestamp := beginHeader.Time.Uint64()
	topArray := make([][]uint64, 0, 250)

	viewCountMap[beginQC.ViewNumber] = 1

	// To transactions
	txCount := uint64(0)
	txCount += uint64(len(beginBlock.Transactions()))
	for i := beginNumber + 1; i <= endNumber; i++ {
		block, _ := backend.BlockByNumber(ctx, rpc.BlockNumber(int64(i)))
		header := block.Header()
		diff := header.Time.Uint64() - preTimestamp
		topArray = append(topArray, []uint64{diff, uint64(len(block.Transactions()))})
		preTimestamp = header.Time.Uint64()
		txCount = txCount + uint64(len(block.Transactions()))

		_, qc, err := ctypes.DecodeExtra(block.ExtraData())
		if err != nil {
			return 0, 0, nil, 0, 0, nil, nil, 0, fmt.Errorf("decode header Extra fail:%v", err)
		}
		if count, ok := viewCountMap[qc.ViewNumber]; ok {
			viewCountMap[qc.ViewNumber] = count + 1
		} else {
			viewCountMap[qc.ViewNumber] = 1
		}

	}

	diffTimestamp := endHeader.Time.Uint64() - beginHeader.Time.Uint64()
	diffNumber := endHeader.Number.Uint64() - beginHeader.Number.Uint64() + 1

	tps := (txCount * 1000) / (endHeader.Time.Uint64() - beginHeader.Time.Uint64())

	// missing view
	for i := DefaultViewNumber; i <= endQC.ViewNumber; i++ {
		if _, ok := viewCountMap[i]; !ok {
			missViewList = append(missViewList, i)
		}
	}

	// view produce block rate
	viewBlockRate := (endNumber - beginNumber + 1) * 100 / ((endQC.ViewNumber - DefaultViewNumber + 1) * 10)

	return diffTimestamp, diffTimestamp / diffNumber, topArray, txCount, tps, viewCountMap, missViewList, viewBlockRate, nil
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
		beginNumber.Value = strconv.Itoa(int(d.BeginNumber))
		endNumber := row.AddCell()
		endNumber.Value = strconv.Itoa(int(d.EndNumber))
		viewBlockRate := row.AddCell()
		viewBlockRate.Value = strconv.Itoa(int(d.ViewBlockRate))
		viewCountMap := row.AddCell()
		viewCountMap.Value = fmt.Sprintf("%v", d.ViewCountMap)
		missViewList := row.AddCell()
		missViewList.Value = fmt.Sprintf("%v", d.MissViewList)
		totalProduceTime := row.AddCell()
		totalProduceTime.Value = strconv.Itoa(int(d.TotalProduceTime))
		averageProduceTime := row.AddCell()
		averageProduceTime.Value = strconv.Itoa(int(d.AverageProduceTime))
		txCount := row.AddCell()
		txCount.Value = strconv.Itoa(int(d.TxCount))
		tps := row.AddCell()
		tps.Value = strconv.Itoa(int(d.Tps))
		topArray := row.AddCell()
		topArray.Value = fmt.Sprintf("%v", d.TopArray)
	}
	err = file.Save(resultPath)
	if err != nil {
		return err
	}
	return nil
}
