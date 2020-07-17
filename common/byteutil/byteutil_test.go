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
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/x/restricting"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func TestBytesToString(t *testing.T) {
	msg := "test string"
	data, err := rlp.EncodeToBytes(msg)
	assert.Nil(t, err)
	dmsg := BytesToString(data)
	assert.Equal(t, msg, dmsg)
}

func TestBytesToUint8(t *testing.T) {
	var msg uint8
	msg = 255
	data, err := rlp.EncodeToBytes(msg)
	assert.Nil(t, err)
	dmsg := BytesToUint8(data)
	assert.Equal(t, msg, dmsg)
}

func TestBytesToUint16(t *testing.T) {
	var msg uint16
	msg = 65535
	data, err := rlp.EncodeToBytes(msg)
	assert.Nil(t, err)
	dmsg := BytesToUint16(data)
	assert.Equal(t, msg, dmsg)
}

func TestBytesToUint32(t *testing.T) {
	var msg uint32
	msg = 4294967295
	data, err := rlp.EncodeToBytes(msg)
	assert.Nil(t, err)
	dmsg := BytesToUint32(data)
	assert.Equal(t, msg, dmsg)
}

func TestBytesToUint64(t *testing.T) {
	var msg uint64
	msg = 18446744073709551612
	data, err := rlp.EncodeToBytes(msg)
	assert.Nil(t, err)
	dmsg := BytesToUint64(data)
	assert.Equal(t, msg, dmsg)
}

func TestBytesToBigInt(t *testing.T) {
	msg := new(big.Int).SetUint64(18446744073709551612)
	data, err := rlp.EncodeToBytes(msg)
	assert.Nil(t, err)
	dmsg := BytesToBigInt(data)
	assert.Equal(t, msg, dmsg)
}

func TestBytesToBigIntArr(t *testing.T) {
	big1 := new(big.Int).SetUint64(18446744073709551611)
	big2 := new(big.Int).SetUint64(18446744073709551612)
	big3 := new(big.Int).SetUint64(18446744073709551613)
	bigArr := [3]*big.Int{big1, big2, big3}
	data, err := rlp.EncodeToBytes(bigArr)
	assert.Nil(t, err)
	dbigArr := BytesToBigIntArr(data)
	assert.Equal(t, len(bigArr), len(dbigArr))
	for i := 0; i < 3; i++ {
		assert.Equal(t, bigArr[i], dbigArr[i])
	}
}

func TestBytesToNodeId(t *testing.T) {
	ecdsaKey, _ := crypto.GenerateKey()
	nodeID := discover.PubkeyID(&ecdsaKey.PublicKey)
	data, err := rlp.EncodeToBytes(nodeID)
	assert.Nil(t, err)
	dnodeID := BytesToNodeId(data)
	assert.Equal(t, nodeID, dnodeID)
	assert.NotNil(t, PrintNodeID(dnodeID))
}

func TestBytesToNodeIdArr(t *testing.T) {
	nodeIdArr := make([]discover.NodeID, 0, 3)
	for i := 0; i < 3; i++ {
		ecdsaKey, _ := crypto.GenerateKey()
		nodeID := discover.PubkeyID(&ecdsaKey.PublicKey)
		nodeIdArr = append(nodeIdArr, nodeID)
	}
	data, err := rlp.EncodeToBytes(nodeIdArr)
	assert.Nil(t, err)
	dnodeIdArr := BytesToNodeIdArr(data)
	assert.Equal(t, len(nodeIdArr), len(dnodeIdArr))
	for i := 0; i < 3; i++ {
		assert.Equal(t, nodeIdArr[i], dnodeIdArr[i])
	}
}

func TestBytesToHash(t *testing.T) {
	hash := common.BytesToHash(utils.Rand32Bytes(32))
	data, err := rlp.EncodeToBytes(hash)
	assert.Nil(t, err)
	dhash := BytesToHash(data)
	assert.Equal(t, hash, dhash)
}

func TestBytesToHashArr(t *testing.T) {
	hashArr := make([]common.Hash, 0, 3)
	for i := 0; i < 3; i++ {
		hash := common.BytesToHash(utils.Rand32Bytes(32))
		hashArr = append(hashArr, hash)
	}
	data, err := rlp.EncodeToBytes(hashArr)
	assert.Nil(t, err)
	dhashArr := BytesToHashArr(data)
	assert.Equal(t, len(hashArr), len(dhashArr))
	for i := 0; i < 3; i++ {
		assert.Equal(t, hashArr[i], dhashArr[i])
	}
}

func TestBytesToAddress(t *testing.T) {
	address := common.BytesToAddress(utils.Rand32Bytes(20))
	data, err := rlp.EncodeToBytes(address)
	assert.Nil(t, err)
	daddress := BytesToAddress(data)
	assert.Equal(t, address, daddress)
}

func TestBytesToAddressArr(t *testing.T) {
	addressArr := make([]common.Address, 0, 3)
	for i := 0; i < 3; i++ {
		address := common.BytesToAddress(utils.Rand32Bytes(20))
		addressArr = append(addressArr, address)
	}
	data, err := rlp.EncodeToBytes(addressArr)
	assert.Nil(t, err)
	daddressArr := BytesToAddressArr(data)
	assert.Equal(t, len(addressArr), len(daddressArr))
	for i := 0; i < 3; i++ {
		assert.Equal(t, addressArr[i], daddressArr[i])
	}
}

func TestBytesToVersionSign(t *testing.T) {
	versionSign := common.BytesToVersionSign(utils.Rand32Bytes(65))
	data, err := rlp.EncodeToBytes(versionSign)
	assert.Nil(t, err)
	dversionSign := BytesToVersionSign(data)
	assert.Equal(t, versionSign, dversionSign)
}

func TestBytesToVersionSignArr(t *testing.T) {
	versionSignArr := make([]common.VersionSign, 0, 3)
	for i := 0; i < 3; i++ {
		versionSign := common.BytesToVersionSign(utils.Rand32Bytes(65))
		versionSignArr = append(versionSignArr, versionSign)
	}
	data, err := rlp.EncodeToBytes(versionSignArr)
	assert.Nil(t, err)
	dversionSignArr := BytesToVersionSignArr(data)
	assert.Equal(t, len(versionSignArr), len(dversionSignArr))
	for i := 0; i < 3; i++ {
		assert.Equal(t, versionSignArr[i], dversionSignArr[i])
	}
}

func TestBytesToRestrictingPlanArr(t *testing.T) {
	pArr := make([]restricting.RestrictingPlan, 0, 3)
	for i := 0; i < 3; i++ {
		p := restricting.RestrictingPlan{
			Epoch:  1,
			Amount: new(big.Int).SetInt64(10),
		}
		pArr = append(pArr, p)
	}
	data, err := rlp.EncodeToBytes(pArr)
	assert.Nil(t, err)
	dpArr := BytesToRestrictingPlanArr(data)
	assert.Equal(t, len(pArr), len(dpArr))
	for i := 0; i < 3; i++ {
		assert.Equal(t, pArr[i].Epoch, dpArr[i].Epoch)
		assert.Equal(t, pArr[i].Amount.Int64(), dpArr[i].Amount.Int64())
	}
}
