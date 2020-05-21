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

package types

import (
	"errors"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// EncodeExtra encode cbft version and `QuorumCert` as extra data.
func EncodeExtra(cbftVersion byte, qc *QuorumCert) ([]byte, error) {
	extra := []byte{cbftVersion}
	bxBytes, err := rlp.EncodeToBytes(qc)
	if err != nil {
		return nil, err
	}
	extra = append(extra, bxBytes...)
	return extra, nil
}

// DecodeExtra decode extra data as cbft version and `QuorumCert`.
func DecodeExtra(extra []byte) (byte, *QuorumCert, error) {
	if len(extra) == 0 {
		return 0, nil, errors.New("empty extra")
	}
	version := extra[0]
	var qc QuorumCert
	err := rlp.DecodeBytes(extra[1:], &qc)
	if err != nil {
		return 0, nil, err
	}
	return version, &qc, nil
}
