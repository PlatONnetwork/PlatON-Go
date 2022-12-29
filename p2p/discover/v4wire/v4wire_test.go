// Copyright 2019 The go-ethereum Authors
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

package v4wire

import (
	"bytes"
	"encoding/hex"
	"net"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// EIP-8 test vectors.
var testPackets = []struct {
	input      string
	wantPacket Packet
}{
	/*{
		input: "96e7a55f265b738379447058e63d7bad33b05b1e3e93d9537f05be9975090bd72dbe25ab1d57c3de65e0b089f6b884c2273d29b343211483557f7c46f3df0b2650779f3459fabbf34261f113b4f29f6d5f540ea7d3fe5e35b24264517f8529110001ea04cb847f000001820cfa8215a8d790000000000000000000000000000000018208ae820d058443b9a355",
		wantPacket: &Ping{
			Version:    4,
			From:       Endpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
			To:         Endpoint{net.ParseIP("::1"), 2222, 3333},
			Expiration: 1136239445,
		},
	},*/
	{
		input: "b045fd6e610bfe8e51393adb7aa058d60259a744cbcbceed6078b006857a7e881e32b9760ea0fd6c07fd9fbd9e73cda03c5ffa5d58c5b1cce8fc309ad2a702fc3f5a496f0e008bba521b77d5486151d391a76690d095b27adfb6a3e77d0193400101eb04cb847f000001820cfa8215a8d790000000000000000000000000000000018208ae820d058443b9a35503",
		wantPacket: &Ping{
			Version:    4,
			From:       Endpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
			To:         Endpoint{net.ParseIP("::1"), 2222, 3333},
			Expiration: 1136239445,
			ForkID:     []rlp.RawValue{{0x03}},
			//ENRSeq:     1,
		},
	},
	{
		input: "0cc5db6098d5d5a5d6af93731aee4629ebbdd967899e42e073dfee48f36162d3fee71339958ee7859a936d61e6e4e43f74f5dc119fffcd6b424df1929f55197b159aaef76f9bac9fed4f35677e85b049a618cdb62d5cdb70a3b238439c79bce30103f84eb840ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31387574077f301b421bc84df7266c44e9e6d569fc56be00812904767bf5ccd1fc7f8443b9a35582999983999999",
		wantPacket: &Findnode{
			Target:     hexPubkey("ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31387574077f301b421bc84df7266c44e9e6d569fc56be00812904767bf5ccd1fc7f"),
			Expiration: 1136239445,
			Rest:       []rlp.RawValue{{0x82, 0x99, 0x99}, {0x83, 0x99, 0x99, 0x99}},
		},
	},
	{
		input: "de3cd768c1b49fd8e7e4e0b50e28c371f4683181550d877a4ce6d32580ccd72c56b8409b9176e182b36ebc0715d6197a69b0eb806d6a7b7aa8615677891e15705c4cf3849f0ff477db229126dc4c0715e11f3ee9172659726dbb3eff8a64a1590004f9015bf90150f84d846321163782115c82115db8403155e1427f85f10a5c9a7755877748041af1bcd8d474ec065eb33df57a97babf54bfd2103575fa829115d224c523596b401065a97f74010610fce76382c0bf32f84984010203040101b840312c55512422cf9b8a4097e9a6ad79402e87a15ae909a4bfefa22398f03d20951933beea1e4dfa6f968212385e829f04c2d314fc2d4e255e0d3bc08792b069dbf8599020010db83c4d001500000000abcdef12820d05820d05b84038643200b172dcfef857492156971f0e6aa2c538d8b74010f8e140811d53b98c765dd2d96126051913f44582e8c199ad7c6d6819e9a56483f637feaac9448aacf8599020010db885a308d313198a2e037073488203e78203e8b8408dcab8618c3253b558d459da53bd8fa68935a719aff8b811197101a4b2b47dd2d47295286fc00cc081bb542d760717d1bdd6bec2c37cd72eca367d6dd3b9df738443b9a355010203",
		wantPacket: &Neighbors{
			Nodes: []Node{
				{
					ID:  hexPubkey("3155e1427f85f10a5c9a7755877748041af1bcd8d474ec065eb33df57a97babf54bfd2103575fa829115d224c523596b401065a97f74010610fce76382c0bf32"),
					IP:  net.ParseIP("99.33.22.55").To4(),
					UDP: 4444,
					TCP: 4445,
				},
				{
					ID:  hexPubkey("312c55512422cf9b8a4097e9a6ad79402e87a15ae909a4bfefa22398f03d20951933beea1e4dfa6f968212385e829f04c2d314fc2d4e255e0d3bc08792b069db"),
					IP:  net.ParseIP("1.2.3.4").To4(),
					UDP: 1,
					TCP: 1,
				},
				{
					ID:  hexPubkey("38643200b172dcfef857492156971f0e6aa2c538d8b74010f8e140811d53b98c765dd2d96126051913f44582e8c199ad7c6d6819e9a56483f637feaac9448aac"),
					IP:  net.ParseIP("2001:db8:3c4d:15::abcd:ef12"),
					UDP: 3333,
					TCP: 3333,
				},
				{
					ID:  hexPubkey("8dcab8618c3253b558d459da53bd8fa68935a719aff8b811197101a4b2b47dd2d47295286fc00cc081bb542d760717d1bdd6bec2c37cd72eca367d6dd3b9df73"),
					IP:  net.ParseIP("2001:db8:85a3:8d3:1319:8a2e:370:7348"),
					UDP: 999,
					TCP: 1000,
				},
			},
			Expiration: 1136239445,
			Rest:       []rlp.RawValue{{0x01}, {0x02}, {0x03}},
		},
	},
}

// This test checks that the decoder accepts packets according to EIP-8.
func TestForwardCompatibility(t *testing.T) {
	testkey, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	wantNodeKey := EncodePubkey(&testkey.PublicKey)

	for i, test := range testPackets {
		input, err := hex.DecodeString(test.input)
		if err != nil {
			t.Fatalf("invalid hex: %s", test.input)
		}
		packet, nodekey, _, err := Decode(input)
		if err != nil {
			t.Errorf("did not accept packet %s\n%v", test.input, err)
			continue
		}
		if !reflect.DeepEqual(packet, test.wantPacket) {
			t.Errorf("got %s\nwant %s,index:%d", spew.Sdump(packet), spew.Sdump(test.wantPacket), i)
		}
		if nodekey != wantNodeKey {
			t.Errorf("got id %v\nwant id %v", nodekey, wantNodeKey)
		}
	}
}

func TestPingEncode(t *testing.T) {
	packet, err := rlp.EncodeToBytes(&Ping{
		Version:    4,
		From:       Endpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
		To:         Endpoint{net.ParseIP("::1"), 2222, 3333},
		Expiration: 1136239445,
		ForkID:     []rlp.RawValue{{0x03}},
	})
	if err != nil {
		t.Error(err)
	}
	packet2, err2 := rlp.EncodeToBytes(&Ping{
		Version:    4,
		From:       Endpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
		To:         Endpoint{net.ParseIP("::1"), 2222, 3333},
		Expiration: 1136239445,
		ForkID:     []rlp.RawValue{{0x03}},
		ENRSeq:     1,
	})
	if err2 != nil {
		t.Error(err2)
	}
	packet3, err3 := rlp.EncodeToBytes(&PingV1{
		Version:    4,
		From:       Endpoint{net.ParseIP("127.0.0.1").To4(), 3322, 5544},
		To:         Endpoint{net.ParseIP("::1"), 2222, 3333},
		Expiration: 1136239445,
		Rest:       []rlp.RawValue{{0x03}},
	})
	if err3 != nil {
		t.Error(err3)
	}

	if !bytes.Equal(packet, packet2) && !bytes.Equal(packet3, packet2) {
		t.Error("should be same")
	}
}

func hexPubkey(h string) (ret Pubkey) {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	if len(b) != len(ret) {
		panic("invalid length")
	}
	copy(ret[:], b)
	return ret
}
