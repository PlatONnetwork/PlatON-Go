// create by platon
package byteutil

import (
	"bytes"
	"encoding/binary"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
	"strings"
)

var Command = map[string]interface{}{
	"string":            BytesToString,
	"[]uint8":           OriginBytes,
	"[64]uint8":         BytesTo64Bytes,
	"[32]uint8":         BytesTo32Bytes,
	"int":               BytesToInt,
	"*big.Int":          BytesToBigInt,
	"uint32":            binary.BigEndian.Uint32,
	"uint64":            binary.BigEndian.Uint64,
	"int32":             common.BytesToInt32,
	"int64":             common.BytesToInt64,
	"float32":           common.BytesToFloat32,
	"float64":           common.BytesToFloat64,
	"discover.NodeID":   BytesToNodeId,
	"[]discover.NodeID": ArrBytesToNodeId,
	"common.Hash":       BytesToHash,
	"common.Address":    BytesToAddress,
}

func BytesToAddress(curByte []byte) common.Address {
	str := BytesToString(curByte)
	return common.HexToAddress(str)
}

func BytesToNodeId(curByte []byte) discover.NodeID {
	str := BytesToString(curByte)
	nodeId, _ := discover.HexID(str)
	return nodeId
}

func ArrBytesToNodeId(curByte []byte) []discover.NodeID {
	str := BytesToString(curByte)
	strArr := strings.Split(str, ":")
	var ANodeID []discover.NodeID
	for i := 0; i < len(strArr); i++ {
		nodeId, _ := discover.HexID(strArr[i])
		ANodeID = append(ANodeID, nodeId)
	}
	return ANodeID
}

func BytesToHash(curByte []byte) common.Hash {
	str := BytesToString(curByte)
	return common.HexToHash(str)
}

func BytesTo32Bytes(curByte []byte) [32]byte {
	var arr [32]byte
	copy(arr[:], curByte)
	return arr
}

func BytesTo64Bytes(curByte []byte) [64]byte {
	var arr [64]byte
	copy(arr[:], curByte)
	return arr
}

func OriginBytes(curByte []byte) []byte {
	return curByte
}

func BytesToBigInt(curByte []byte) *big.Int {
	return new(big.Int).SetBytes(curByte)
}

func BytesToInt(curByte []byte) int {
	bytesBuffer := bytes.NewBuffer(curByte)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	b := int(x)
	return b
}

func BytesToString(curByte []byte) string {
	return string(curByte)
}
