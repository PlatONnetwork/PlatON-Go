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

package node

import (
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
)

var (
	//priKey = crypto.HexMustToECDSA("8e1477549bea04b97ea15911e2e9b3041b7a9921f80bd6ddbe4c2b080473de22")
	priKey = crypto.HexMustToECDSA("8c56e4a0d8bb1f82b94231d535c499fdcbf6e06221acf669d5a964f5bb974903")

	nodeID = enode.HexID("fd070d294c4348d165c0bb89ed5079841ffc3b44a6ce391418fcc6603ea2d284")
)

func initChandlerHandler() {
	chandler = GetCryptoHandler()
	chandler.SetPrivateKey(priKey)
}

func TestCryptoHandler_IsSignedByNodeID(t *testing.T) {
	initChandlerHandler()
	version := uint32(1<<16 | 1<<8 | 0)
	sig := chandler.MustSign(version)
	if !chandler.IsSignedByNodeID(version, sig, nodeID) {
		t.Fatal("verify sign error")
	} else {
		t.Log("verify sign OK")
	}
}
