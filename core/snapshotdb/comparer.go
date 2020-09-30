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

package snapshotdb

import "bytes"

type bytesComparer struct{}

func (bytesComparer) Compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

func (bytesComparer) Name() string {
	return "snapshotdb.BytewiseComparator"
}

func (bytesComparer) Separator(dst, a, b []byte) []byte {
	i, n := 0, len(a)
	if n > len(b) {
		n = len(b)
	}
	for ; i < n && a[i] == b[i]; i++ {
	}
	if i >= n {
		// Do not shorten if one string is a prefix of the other
	} else if c := a[i]; c < 0xff && c+1 < b[i] {
		dst = append(dst, a[:i+1]...)
		dst[len(dst)-1]++
		return dst
	}
	return nil
}

func (bytesComparer) Successor(dst, b []byte) []byte {
	for i, c := range b {
		if c != 0xff {
			dst = append(dst, b[:i+1]...)
			dst[len(dst)-1]++
			return dst
		}
	}
	return nil
}

// DefaultComparer are default implementation of the Comparer interface.
// It uses the natural ordering, consistent with bytes.Compare.
var DefaultComparer = bytesComparer{}
