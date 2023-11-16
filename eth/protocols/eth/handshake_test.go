// Copyright 2014 The go-ethereum Authors
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

package eth

import (
	"errors"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
)

// Tests that handshake failures are detected and reported correctly.
func TestHandshake66(t *testing.T) { testHandshake(t, ETH66) }

func testHandshake(t *testing.T, protocol uint) {
	t.Parallel()

	// Create a test backend only to have some valid genesis chain
	backend := newTestBackend(3)
	defer backend.close()

	var (
		genesis = backend.chain.Genesis()
		head    = backend.chain.CurrentBlock()
	)
	tests := []struct {
		code uint64
		data interface{}
		want error
	}{
		{
			code: TransactionsMsg, data: []interface{}{},
			want: errNoStatusMsg,
		},
		{
			code: StatusMsg, data: StatusPacket{10, 1, head.Number(), head.Hash(), genesis.Hash()},
			want: errProtocolVersionMismatch,
		},
		{
			code: StatusMsg, data: StatusPacket{uint32(protocol), 999, head.Number(), head.Hash(), genesis.Hash()},
			want: errNetworkIDMismatch,
		},
		{
			code: StatusMsg, data: StatusPacket{uint32(protocol), 1, head.Number(), head.Hash(), common.Hash{3}},
			want: errGenesisMismatch,
		},
	}
	for i, test := range tests {
		// Create the two peers to shake with each other
		app, net := p2p.MsgPipe()
		defer app.Close()
		defer net.Close()

		peer := NewPeer(protocol, p2p.NewPeer(enode.ID{}, "peer", nil), net, nil, nil)
		defer peer.Close()

		// Send the junk test with one peer, check the handshake failure
		go p2p.Send(app, test.code, test.data)

		err := peer.Handshake(1, head.Number(), head.Hash(), genesis.Hash(), nil)
		if err == nil {
			t.Errorf("test %d: protocol returned nil error, want %q", i, test.want)
		} else if !errors.Is(err, test.want) {
			t.Errorf("test %d: wrong error: got %q, want %q", i, err, test.want)
		}
	}
}
