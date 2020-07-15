// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package byteutil

import (
	"encoding/hex"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
)

var Bytes2X_CMD = map[string]interface{}{
	"string":   BytesToString,
	"[8]byte":  BytesTo8Bytes,
	"[16]byte": BytesTo16Bytes,
	"[32]byte": BytesTo32Bytes,
	"[64]byte": BytesTo64Bytes,

	"uint8":  BytesToUint8,
	"uint16": BytesToUint16,
	"uint32": BytesToUint32,
	"uint64": BytesToUint64,

	"*big.Int":              BytesToBigInt,
	"[]*big.Int":            BytesToBigIntArr,
	"discover.NodeID":       BytesToNodeId,
	"[]discover.NodeID":     BytesToNodeIdArr,
	"common.Hash":           BytesToHash,
	"[]common.Hash":         BytesToHashArr,
	"common.Address":        BytesToAddress,
	"[]common.Address":      BytesToAddressArr,
	"common.VersionSign":    BytesToVersionSign,
	"[]common.VersionSign":  BytesToVersionSignArr,
	"bls.PublicKeyHex":      BytesToPublicKeyHex,
	"[]bls.PublicKeyHex":    BytesToPublicKeyHexArr,
	"bls.SchnorrProofHex":   BytesToSchnorrProofHex,
	"[]bls.SchnorrProofHex": BytesToSchnorrProofHexArr,

	"[]restricting.RestrictingPlan": BytesToRestrictingPlanArr,
}

func BytesToString(curByte []byte) string {
	//return string(curByte)
	var str string
	if err := rlp.DecodeBytes(curByte, &str); nil != err {
		panic("BytesToString:" + err.Error())
	}
	return str
}

func BytesTo8Bytes(curByte []byte) [8]byte {
	var arr [8]byte
	if err := rlp.DecodeBytes(curByte, &arr); nil != err {
		panic("BytesTo8Bytes:" + err.Error())
	}
	return arr
}

func BytesTo16Bytes(curByte []byte) [16]byte {
	var arr [16]byte
	if err := rlp.DecodeBytes(curByte, &arr); nil != err {
		panic("BytesTo16Bytes:" + err.Error())
	}
	return arr
}

func BytesTo32Bytes(curByte []byte) [32]byte {
	/*var arr [32]byte
	copy(arr[:], curByte)
	return arr*/
	var arr [32]byte
	if err := rlp.DecodeBytes(curByte, &arr); nil != err {
		panic("BytesTo32Bytes:" + err.Error())
	}
	return arr
}

func BytesTo64Bytes(curByte []byte) [64]byte {
	/*var arr [64]byte
	copy(arr[:], curByte)
	return arr*/
	var arr [64]byte
	if err := rlp.DecodeBytes(curByte, &arr); nil != err {
		panic("BytesTo64Bytes:" + err.Error())
	}
	return arr
}

func BytesToUint8(b []byte) uint8 {
	var x uint8
	if err := rlp.DecodeBytes(b, &x); nil != err {
		panic("BytesToUint8:" + err.Error())
	}
	return x
}

func BytesToUint16(b []byte) uint16 {
	/*b = append(make([]byte, 2-len(b)), b...)
	return binary.BigEndian.Uint16(b)*/
	var x uint16
	if err := rlp.DecodeBytes(b, &x); nil != err {
		panic("BytesToUint16:" + err.Error())
	}
	return x
}

func BytesToUint32(b []byte) uint32 {
	/*b = append(make([]byte, 4-len(b)), b...)
	return binary.BigEndian.Uint32(b)*/
	var x uint32
	if err := rlp.DecodeBytes(b, &x); nil != err {
		panic("BytesToUint32:" + err.Error())
	}
	return x
}

func BytesToUint64(b []byte) uint64 {
	/*b = append(make([]byte, 8-len(b)), b...)
	return binary.BigEndian.Uint64(b)*/
	var x uint64
	if err := rlp.DecodeBytes(b, &x); nil != err {
		panic("BytesToUint64:" + err.Error())
	}
	return x
}

func BytesToBigInt(curByte []byte) *big.Int {
	//return new(big.Int).SetBytes(curByte)
	var bigInt *big.Int
	if err := rlp.DecodeBytes(curByte, &bigInt); nil != err {
		panic("BytesToBigInt:" + err.Error())
	}
	return bigInt
}

func BytesToBigIntArr(curByte []byte) []*big.Int {
	var arr []*big.Int
	if err := rlp.DecodeBytes(curByte, &arr); nil != err {
		panic("BytesToBigIntArr:" + err.Error())
	}
	return arr
}

func BytesToNodeId(curByte []byte) discover.NodeID {
	//str := BytesToString(curByte)
	//nodeId, _ := discover.HexID(str)
	//return nodeId
	var nodeId discover.NodeID
	if err := rlp.DecodeBytes(curByte, &nodeId); nil != err {
		panic("BytesToNodeId:" + err.Error())
	}
	return nodeId
}

func BytesToNodeIdArr(curByte []byte) []discover.NodeID {
	/*str := BytesToString(curByte)
	strArr := strings.Split(str, ":")
	var ANodeID []discover.NodeID
	for i := 0; i < len(strArr); i++ {
		nodeId, _ := discover.HexID(strArr[i])
		ANodeID = append(ANodeID, nodeId)
	}
	return ANodeID*/
	var nodeIdArr []discover.NodeID
	if err := rlp.DecodeBytes(curByte, &nodeIdArr); nil != err {
		panic("BytesToNodeIdArr:" + err.Error())
	}
	return nodeIdArr
}

func BytesToHash(curByte []byte) common.Hash {
	//str := BytesToString(curByte)
	//return common.HexToHash(str)
	var hash common.Hash
	if err := rlp.DecodeBytes(curByte, &hash); nil != err {
		panic("BytesToHash:" + err.Error())
	}
	return hash
}

func BytesToHashArr(curByte []byte) []common.Hash {
	/*str := BytesToString(curByte)
	strArr := strings.Split(str, ":")
	var AHash []common.Hash
	for i := 0; i < len(strArr); i++ {
		AHash = append(AHash, common.HexToHash(strArr[i]))
	}
	return AHash*/

	var hashArr []common.Hash
	if err := rlp.DecodeBytes(curByte, &hashArr); nil != err {
		panic("BytesToHashArr:" + err.Error())
	}
	return hashArr
}

func BytesToAddress(curByte []byte) common.Address {
	//str := BytesToString(curByte)
	//return common.HexToAddress(str)
	var addr common.Address
	if err := rlp.DecodeBytes(curByte, &addr); nil != err {
		panic("BytesToAddress:" + err.Error())
	}
	return addr
}

func BytesToAddressArr(curByte []byte) []common.Address {
	//str := BytesToString(curByte)
	//return common.HexToAddress(str)
	var addrArr []common.Address
	if err := rlp.DecodeBytes(curByte, &addrArr); nil != err {
		panic("BytesToAddressArr:" + err.Error())
	}
	return addrArr
}

func BytesToVersionSign(currByte []byte) common.VersionSign {
	var version common.VersionSign
	if err := rlp.DecodeBytes(currByte, &version); nil != err {
		panic("BytesToVersionSign:" + err.Error())
	}
	return version
}

func BytesToVersionSignArr(currByte []byte) []common.VersionSign {
	var arr []common.VersionSign
	if err := rlp.DecodeBytes(currByte, &arr); nil != err {
		panic("BytesToVersionSignArr:" + err.Error())
	}
	return arr
}

func BytesToPublicKeyHex(currByte []byte) bls.PublicKeyHex {
	var pub bls.PublicKeyHex
	if err := rlp.DecodeBytes(currByte, &pub); nil != err {
		panic("BytesToPublicKeyHex:" + err.Error())
	}
	return pub
}

func BytesToPublicKeyHexArr(currByte []byte) []bls.PublicKeyHex {
	var arr []bls.PublicKeyHex
	if err := rlp.DecodeBytes(currByte, &arr); nil != err {
		panic("BytesToPublicKeyHexArr:" + err.Error())
	}
	return arr
}

func BytesToSchnorrProofHex(currByte []byte) bls.SchnorrProofHex {
	var proof bls.SchnorrProofHex
	if err := rlp.DecodeBytes(currByte, &proof); nil != err {
		panic("BytesToSchnorrProofHex:" + err.Error())
	}
	return proof
}

func BytesToSchnorrProofHexArr(currByte []byte) []bls.SchnorrProofHex {
	var arr []bls.SchnorrProofHex
	if err := rlp.DecodeBytes(currByte, &arr); nil != err {
		panic("BytesToSchnorrProofHexArr:" + err.Error())
	}
	return arr
}

func BytesToRestrictingPlanArr(curByte []byte) []restricting.RestrictingPlan {
	var planArr []restricting.RestrictingPlan
	if err := rlp.DecodeBytes(curByte, &planArr); nil != err {
		panic("BytesToAddressArr:" + err.Error())
	}
	return planArr
}

func PrintNodeID(nodeID discover.NodeID) string {
	return hex.EncodeToString(nodeID.Bytes()[:8])
}

func Concat(s1 []byte, s2 ...byte) []byte {
	r := make([]byte, len(s1)+len(s2))
	copy(r, s1)
	copy(r[len(s1):], s2)
	return r
}
