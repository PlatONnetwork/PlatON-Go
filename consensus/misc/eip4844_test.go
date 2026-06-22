// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package misc

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/params"
)

func TestCalcExcessBlobGas(t *testing.T) {
	var tests = []struct {
		excess uint64
		blobs  uint64
		want   uint64
	}{
		{0, 0, 0},
		{0, 1, 0},
		{0, params.BlobTxTargetBlobGasPerBlock / params.BlobTxBlobGasPerBlob, 0},

		{0, (params.BlobTxTargetBlobGasPerBlock / params.BlobTxBlobGasPerBlob) + 1, params.BlobTxBlobGasPerBlob},
		{1, (params.BlobTxTargetBlobGasPerBlock / params.BlobTxBlobGasPerBlob) + 1, params.BlobTxBlobGasPerBlob + 1},
		{1, (params.BlobTxTargetBlobGasPerBlock / params.BlobTxBlobGasPerBlob) + 2, 2*params.BlobTxBlobGasPerBlob + 1},

		{params.BlobTxTargetBlobGasPerBlock, params.BlobTxTargetBlobGasPerBlock / params.BlobTxBlobGasPerBlob, params.BlobTxTargetBlobGasPerBlock},
		{params.BlobTxTargetBlobGasPerBlock, (params.BlobTxTargetBlobGasPerBlock / params.BlobTxBlobGasPerBlob) - 1, params.BlobTxBlobGasPerBlob},
		{params.BlobTxTargetBlobGasPerBlock, (params.BlobTxTargetBlobGasPerBlock / params.BlobTxBlobGasPerBlob) - 2, 0},
		{params.BlobTxBlobGasPerBlob - 1, (params.BlobTxTargetBlobGasPerBlock / params.BlobTxBlobGasPerBlob) - 1, 0},
	}
	for _, tt := range tests {
		result := CalcExcessBlobGas(tt.excess, tt.blobs*params.BlobTxBlobGasPerBlob)
		if result != tt.want {
			t.Errorf("excess blob gas mismatch: have %v, want %v", result, tt.want)
		}
	}
}

func TestCalcBlobFee(t *testing.T) {
	tests := []struct {
		excessBlobGas uint64
		blobfee       int64
	}{
		{0, 1},
		{1542707, 2},
		{10 * 1024 * 1024, 111},
	}
	for i, tt := range tests {
		have := CalcBlobFee(tt.excessBlobGas)
		if have.Int64() != tt.blobfee {
			t.Errorf("test %d: blobfee mismatch: have %v want %v", i, have, tt.blobfee)
		}
	}
}
