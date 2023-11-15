// Copyright 2020 The go-ethereum Authors
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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"math/big"
	"math/rand"

	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"
)

// handleGetBlockHeaders66 is the eth/66 version of handleGetBlockHeaders
func handleGetBlockHeaders66(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the complex header query
	var query GetBlockHeadersPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := ServiceGetBlockHeadersQuery(backend.Chain(), query.GetBlockHeadersPacket, peer)
	return peer.ReplyBlockHeadersRLP(query.RequestId, response)
}

// ServiceGetBlockHeadersQuery assembles the response to a header query. It is
// exposed to allow external packages to test protocol behavior.
func ServiceGetBlockHeadersQuery(chain *core.BlockChain, query *GetBlockHeadersPacket, peer *Peer) []rlp.RawValue {
	if query.Skip == 0 {
		// The fast path: when the request is for a contiguous segment of headers.
		return serviceContiguousBlockHeaderQuery(chain, query)
	} else {
		return serviceNonContiguousBlockHeaderQuery(chain, query, peer)
	}
}

func serviceNonContiguousBlockHeaderQuery(chain *core.BlockChain, query *GetBlockHeadersPacket, peer *Peer) []rlp.RawValue {
	hashMode := query.Origin.Hash != (common.Hash{})
	first := true
	maxNonCanonical := uint64(100)

	// Gather headers until the fetch or network limits is reached
	var (
		bytes   common.StorageSize
		headers []rlp.RawValue
		unknown bool
		lookups int
	)
	for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit &&
		len(headers) < maxHeadersServe && lookups < 2*maxHeadersServe {
		lookups++
		// Retrieve the next header satisfying the query
		var origin *types.Header
		if hashMode {
			if first {
				first = false
				origin = chain.GetHeaderByHash(query.Origin.Hash)
				if origin != nil {
					query.Origin.Number = origin.Number.Uint64()
				}
			} else {
				origin = chain.GetHeader(query.Origin.Hash, query.Origin.Number)
			}
		} else {
			origin = chain.GetHeaderByNumber(query.Origin.Number)
		}
		if origin == nil {
			break
		}
		if rlpData, err := rlp.EncodeToBytes(origin); err != nil {
			log.Crit("Unable to decode our own headers", "err", err)
		} else {
			headers = append(headers, rlp.RawValue(rlpData))
			bytes += common.StorageSize(len(rlpData))
		}
		// Advance to the next header of the query
		switch {
		case hashMode && query.Reverse:
			// Hash based traversal towards the genesis block
			ancestor := query.Skip + 1
			if ancestor == 0 {
				unknown = true
			} else {
				query.Origin.Hash, query.Origin.Number = chain.GetAncestor(query.Origin.Hash, query.Origin.Number, ancestor, &maxNonCanonical)
				unknown = (query.Origin.Hash == common.Hash{})
			}
		case hashMode && !query.Reverse:
			// Hash based traversal towards the leaf block
			var (
				current = origin.Number.Uint64()
				next    = current + query.Skip + 1
			)
			if next <= current {
				infos, _ := json.MarshalIndent(peer.Peer.Info(), "", "  ")
				peer.Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
				unknown = true
			} else {
				if header := chain.GetHeaderByNumber(next); header != nil {
					nextHash := header.Hash()
					expOldHash, _ := chain.GetAncestor(nextHash, next, query.Skip+1, &maxNonCanonical)
					if expOldHash == query.Origin.Hash {
						query.Origin.Hash, query.Origin.Number = nextHash, next
					} else {
						unknown = true
					}
				} else {
					unknown = true
				}
			}
		case query.Reverse:
			// Number based traversal towards the genesis block
			if query.Origin.Number >= query.Skip+1 {
				query.Origin.Number -= query.Skip + 1
			} else {
				unknown = true
			}

		case !query.Reverse:
			// Number based traversal towards the leaf block
			query.Origin.Number += query.Skip + 1
		}
	}
	return headers
}

func serviceContiguousBlockHeaderQuery(chain *core.BlockChain, query *GetBlockHeadersPacket) []rlp.RawValue {
	count := query.Amount
	if count > maxHeadersServe {
		count = maxHeadersServe
	}
	if query.Origin.Hash == (common.Hash{}) {
		// Number mode, just return the canon chain segment. The backend
		// delivers in [N, N-1, N-2..] descending order, so we need to
		// accommodate for that.
		from := query.Origin.Number
		if !query.Reverse {
			from = from + count - 1
		}
		headers := chain.GetHeadersFrom(from, count)
		if !query.Reverse {
			for i, j := 0, len(headers)-1; i < j; i, j = i+1, j-1 {
				headers[i], headers[j] = headers[j], headers[i]
			}
		}
		return headers
	}
	// Hash mode.
	var (
		headers []rlp.RawValue
		hash    = query.Origin.Hash
		header  = chain.GetHeaderByHash(hash)
	)
	if header != nil {
		rlpData, _ := rlp.EncodeToBytes(header)
		headers = append(headers, rlpData)
	} else {
		// We don't even have the origin header
		return headers
	}
	num := header.Number.Uint64()
	if !query.Reverse {
		// Theoretically, we are tasked to deliver header by hash H, and onwards.
		// However, if H is not canon, we will be unable to deliver any descendants of
		// H.
		if canonHash := chain.GetCanonicalHash(num); canonHash != hash {
			// Not canon, we can't deliver descendants
			return headers
		}
		descendants := chain.GetHeadersFrom(num+count-1, count-1)
		for i, j := 0, len(descendants)-1; i < j; i, j = i+1, j-1 {
			descendants[i], descendants[j] = descendants[j], descendants[i]
		}
		headers = append(headers, descendants...)
		return headers
	}
	{ // Last mode: deliver ancestors of H
		for i := uint64(1); header != nil && i < count; i++ {
			header = chain.GetHeaderByHash(header.ParentHash)
			if header == nil {
				break
			}
			rlpData, _ := rlp.EncodeToBytes(header)
			headers = append(headers, rlpData)
		}
		return headers
	}
}

func handleGetBlockBodies66(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the block body retrieval message
	var query GetBlockBodiesPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := ServiceGetBlockBodiesQuery(backend.Chain(), query.GetBlockBodiesPacket)
	return peer.ReplyBlockBodiesRLP(query.RequestId, response)
}

// ServiceGetBlockBodiesQuery assembles the response to a body query. It is
// exposed to allow external packages to test protocol behavior.
func ServiceGetBlockBodiesQuery(chain *core.BlockChain, query GetBlockBodiesPacket) []rlp.RawValue {
	// Gather blocks until the fetch or network limits is reached
	var (
		bytes  int
		bodies []rlp.RawValue
	)
	for lookups, hash := range query {
		if bytes >= softResponseLimit || len(bodies) >= maxBodiesServe ||
			lookups >= 2*maxBodiesServe {
			break
		}
		if data := chain.GetBodyRLP(hash); len(data) != 0 {
			bodies = append(bodies, data)
			bytes += len(data)
		}
	}
	return bodies
}

func handleGetNodeData66(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the trie node data retrieval message
	var query GetNodeDataPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := ServiceGetNodeDataQuery(backend.Chain(), query.GetNodeDataPacket)
	return peer.ReplyNodeData(query.RequestId, response)
}

// ServiceGetNodeDataQuery assembles the response to a node data query. It is
// exposed to allow external packages to test protocol behavior.
func ServiceGetNodeDataQuery(chain *core.BlockChain, query GetNodeDataPacket) [][]byte {
	// Gather state data until the fetch or network limits is reached
	var (
		bytes int
		nodes [][]byte
	)
	for lookups, hash := range query {
		if bytes >= softResponseLimit || len(nodes) >= maxNodeDataServe ||
			lookups >= 2*maxNodeDataServe {
			break
		}
		// Retrieve the requested state entry
		entry, err := chain.TrieNode(hash)
		if len(entry) == 0 || err != nil {
			// Read the contract code with prefix only to save unnecessary lookups.
			entry, err = chain.ContractCodeWithPrefix(hash)
		}
		if err == nil && len(entry) > 0 {
			nodes = append(nodes, entry)
			bytes += len(entry)
		}
	}
	return nodes
}

func handleGetReceipts66(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the block receipts retrieval message
	var query GetReceiptsPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := ServiceGetReceiptsQuery(backend.Chain(), query.GetReceiptsPacket)
	return peer.ReplyReceiptsRLP(query.RequestId, response)
}

// ServiceGetReceiptsQuery assembles the response to a receipt query. It is
// exposed to allow external packages to test protocol behavior.
func ServiceGetReceiptsQuery(chain *core.BlockChain, query GetReceiptsPacket) []rlp.RawValue {
	// Gather state data until the fetch or network limits is reached
	var (
		bytes    int
		receipts []rlp.RawValue
	)
	for lookups, hash := range query {
		if bytes >= softResponseLimit || len(receipts) >= maxReceiptsServe ||
			lookups >= 2*maxReceiptsServe {
			break
		}
		// Retrieve the requested block's receipts
		results := chain.GetReceiptsByHash(hash)
		if results == nil {
			if header := chain.GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
				continue
			}
		}
		// If known, encode and queue for response packet
		if encoded, err := rlp.EncodeToBytes(results); err != nil {
			log.Error("Failed to encode receipt", "err", err)
		} else {
			receipts = append(receipts, encoded)
			bytes += len(encoded)
		}
	}
	return receipts
}

func handleNewBlockhashes(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of new block announcements just arrived
	ann := new(NewBlockHashesPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Mark the hashes as present at the remote node
	for _, block := range *ann {
		peer.markBlock(block.Hash)
	}
	// Deliver them all to the backend for queuing
	return backend.Handle(peer, ann)
}

func handleNewBlock(backend Backend, msg Decoder, peer *Peer) error {
	// Retrieve and decode the propagated block
	ann := new(NewBlockPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if err := ann.sanityCheck(); err != nil {
		return err
	}
	if hash := types.DeriveSha(ann.Block.Transactions(), trie.NewStackTrie(nil)); hash != ann.Block.TxHash() {
		log.Warn("Propagated block has invalid body", "have", hash, "exp", ann.Block.TxHash())
		return nil // TODO(karalabe): return error eventually, but wait a few releases
	}
	ann.Block.ReceivedAt = msg.Time()
	ann.Block.ReceivedFrom = peer

	// Mark the peer as owning the block
	peer.markBlock(ann.Block.Hash())

	return backend.Handle(peer, ann)
}

func handleBlockHeaders66(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of headers arrived to one of our previous requests
	res := new(BlockHeadersPacket66)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	metadata := func() interface{} {
		hashes := make([]common.Hash, len(res.BlockHeadersPacket))
		for i, header := range res.BlockHeadersPacket {
			hashes[i] = header.Hash()
		}
		return hashes
	}
	return peer.dispatchResponse(&Response{
		id:   res.RequestId,
		code: BlockHeadersMsg,
		Res:  &res.BlockHeadersPacket,
	}, metadata)
}

func handleBlockBodies66(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of block bodies arrived to one of our previous requests
	res := new(BlockBodiesPacket66)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	metadata := func() interface{} {
		var (
			txsHashes = make([]common.Hash, len(res.BlockBodiesPacket))
		)
		hasher := trie.NewStackTrie(nil)
		for i, body := range res.BlockBodiesPacket {
			txsHashes[i] = types.DeriveSha(types.Transactions(body.Transactions), hasher)
		}
		return [][]common.Hash{txsHashes}
	}
	return peer.dispatchResponse(&Response{
		id:   res.RequestId,
		code: BlockBodiesMsg,
		Res:  &res.BlockBodiesPacket,
	}, metadata)
}

func handleNodeData66(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of node state data arrived to one of our previous requests
	res := new(NodeDataPacket66)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	return peer.dispatchResponse(&Response{
		id:   res.RequestId,
		code: NodeDataMsg,
		Res:  &res.NodeDataPacket,
	}, nil) // No post-processing, we're not using this packet anymore
}

func handleReceipts66(backend Backend, msg Decoder, peer *Peer) error {
	// A batch of receipts arrived to one of our previous requests
	res := new(ReceiptsPacket66)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	metadata := func() interface{} {
		hasher := trie.NewStackTrie(nil)
		hashes := make([]common.Hash, len(res.ReceiptsPacket))
		for i, receipt := range res.ReceiptsPacket {
			hashes[i] = types.DeriveSha(types.Receipts(receipt), hasher)
		}
		return hashes
	}
	return peer.dispatchResponse(&Response{
		id:   res.RequestId,
		code: ReceiptsMsg,
		Res:  &res.ReceiptsPacket,
	}, metadata)
}

func handleNewPooledTransactionHashes(backend Backend, msg Decoder, peer *Peer) error {
	// New transaction announcement arrived, make sure we have
	// a valid and fresh chain to handle them
	if !backend.AcceptTxs() || !backend.AcceptRemoteTxs() {
		return nil
	}
	ann := new(NewPooledTransactionHashesPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Schedule all the unknown hashes for retrieval
	for _, hash := range *ann {
		peer.markTransaction(hash)
	}
	return backend.Handle(peer, ann)
}

func handleGetPooledTransactions66(backend Backend, msg Decoder, peer *Peer) error {
	// Decode the pooled transactions retrieval message
	var query GetPooledTransactionsPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	hashes, txs := answerGetPooledTransactions(backend, query.GetPooledTransactionsPacket, peer)
	return peer.ReplyPooledTransactionsRLP(query.RequestId, hashes, txs)
}

func answerGetPooledTransactions(backend Backend, query GetPooledTransactionsPacket, peer *Peer) ([]common.Hash, []rlp.RawValue) {
	// Gather transactions until the fetch or network limits is reached
	var (
		bytes  int
		hashes []common.Hash
		txs    []rlp.RawValue
	)
	for _, hash := range query {
		if bytes >= softResponseLimit {
			break
		}
		// Retrieve the requested transaction, skipping if unknown to us
		tx := backend.TxPool().Get(hash)
		if tx == nil {
			continue
		}
		// If known, encode and queue for response packet
		if encoded, err := rlp.EncodeToBytes(tx); err != nil {
			log.Error("Failed to encode transaction", "err", err)
		} else {
			hashes = append(hashes, hash)
			txs = append(txs, encoded)
			bytes += len(encoded)
		}
	}
	return hashes, txs
}

func handleTransactions(backend Backend, msg Decoder, peer *Peer) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if !backend.AcceptTxs() || !backend.AcceptRemoteTxs() {
		return nil
	}
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs TransactionsPacket
	if err := msg.Decode(&txs); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	for i, tx := range txs {
		// Validate and mark the remote transaction
		if tx == nil {
			return fmt.Errorf("%w: transaction %d is nil", errDecode, i)
		}
		peer.markTransaction(tx.Hash())
	}
	return backend.Handle(peer, &txs)
}

func handlePooledTransactions66(backend Backend, msg Decoder, peer *Peer) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if !backend.AcceptTxs() || !backend.AcceptRemoteTxs() {
		return nil
	}
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs PooledTransactionsPacket66
	if err := msg.Decode(&txs); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	for i, tx := range txs.PooledTransactionsPacket {
		// Validate and mark the remote transaction
		if tx == nil {
			return fmt.Errorf("%w: transaction %d is nil", errDecode, i)
		}
		peer.markTransaction(tx.Hash())
	}
	requestTracker.Fulfil(peer.id, peer.version, PooledTransactionsMsg, txs.RequestId)

	return backend.Handle(peer, &txs.PooledTransactionsPacket)
}

// handleGetPPOSStorageMsg handles PPOS Storage query, collect the requested Storage and reply
func handleGetPPOSStorageMsg(backend Backend, msg Decoder, peer *Peer) error {
	var query []interface{}
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if len(query) == 0 {
		return answerGetPPOSStorageMsgQuery(backend, peer)
	} else {
		if a, ok := query[0].([]byte); ok {
			return ServiceGetPPOSStorageMsgQueryV2(backend.Chain(), snapshotdb.Instance(), common.BytesToUint64(a), peer)
		}
	}
	log.Warn("handleGetPPOSStorageMsg seems not good,the query is wrong", "length", len(query), "val", query)
	return nil
}

func answerGetPPOSStorageMsgQuery(backend Backend, peer *Peer) error {
	f := func(num *big.Int, iter iterator.Iterator) error {
		var psInfo PposInfoPacket
		if num == nil {
			return errors.New("num should not be nil")
		}
		psInfo.Pivot = backend.Chain().GetHeaderByNumber(num.Uint64())
		psInfo.Latest = backend.Chain().CurrentHeader()
		if err := peer.SendPPOSInfo(psInfo); err != nil {
			peer.Log().Error("[GetPPOSStorageMsg]send last ppos meassage fail", "error", err)
			return err
		}
		var (
			byteSize int
			ps       PposStoragePacket
			count    int
		)
		ps.KVs = make([][2][]byte, 0)
		for iter.Next() {
			if bytes.Equal(iter.Key(), []byte(snapshotdb.CurrentHighestBlock)) || bytes.Equal(iter.Key(), []byte(snapshotdb.CurrentBaseNum)) || bytes.HasPrefix(iter.Key(), []byte(snapshotdb.WalKeyPrefix)) {
				continue
			}
			byteSize = byteSize + len(iter.Key()) + len(iter.Value())
			if count >= PPOSStorageKVSizeFetch || byteSize > softResponseLimit {
				if err := peer.SendPPOSStorage(ps); err != nil {
					peer.Log().Error("[GetPPOSStorageMsg]send ppos message fail", "error", err, "kvnum", ps.KVNum)
					return err
				}
				count = 0
				ps.KVs = make([][2][]byte, 0)
				byteSize = 0
			}
			k, v := make([]byte, len(iter.Key())), make([]byte, len(iter.Value()))
			copy(k, iter.Key())
			copy(v, iter.Value())
			ps.KVs = append(ps.KVs, [2][]byte{
				k, v,
			})
			ps.KVNum++
			count++
		}
		ps.Last = true
		if err := peer.SendPPOSStorage(ps); err != nil {
			peer.Log().Error("[GetPPOSStorageMsg]send last ppos message fail", "error", err)
			return err
		}
		return nil
	}
	go func() {
		if err := snapshotdb.Instance().WalkBaseDB(nil, f); err != nil {
			peer.Log().Error("[GetPPOSStorageMsg]send  ppos storage fail", "error", err)
		}
	}()
	return nil
}

func ServiceGetPPOSStorageMsgQueryV2(chain *core.BlockChain, sdb snapshotdb.DB, head uint64, peer *Peer) error {
	currentHeader := chain.CurrentHeader()
	if currentHeader.Number.Uint64() < head {
		return fmt.Errorf("faild to answerGetPPOSStorageMsgQueryV2 , the current header %v is less than the request header %v", currentHeader.Number.Uint64(), head)
	}

	f := func(baseBlock uint64, iter iterator.Iterator, blocks []rlp.RawValue) error {
		peer.Log().Debug("begin answerGetPPOSStorageMsgQueryV2", "blocks", len(blocks), "baseBlock", baseBlock)
		var (
			byteSize        int
			ps              PposStoragePacket
			count           int
			pposStorageRlps []rlp.RawValue
		)
		ps.KVs = make([][2][]byte, 0)
		for iter.Next() {
			if bytes.Equal(iter.Key(), []byte(snapshotdb.CurrentHighestBlock)) || bytes.Equal(iter.Key(), []byte(snapshotdb.CurrentBaseNum)) || bytes.HasPrefix(iter.Key(), []byte(snapshotdb.WalKeyPrefix)) {
				continue
			}
			byteSize = byteSize + len(iter.Key()) + len(iter.Value())
			if count >= PPOSStorageKVSizeFetch || byteSize > softResponseLimit {
				if err := peer.SendPPOSStorage(ps); err != nil {
					return fmt.Errorf("send ppos storage message fail,%v", err)
				}
				count = 0
				ps.KVs = make([][2][]byte, 0)
				byteSize = 0
			}
			k, v := make([]byte, len(iter.Key())), make([]byte, len(iter.Value()))
			copy(k, iter.Key())
			copy(v, iter.Value())
			ps.KVs = append(ps.KVs, [2][]byte{
				k, v,
			})
			ps.KVNum++
			count++
		}
		ps.Last = true
		if err := peer.SendPPOSStorage(ps); err != nil {
			return fmt.Errorf("send last ppos storage message fail,%v", err)
		}

		byteSize = 0
		if head == baseBlock {
			id := rand.Uint64()
			if err := peer.ReplyPPOSStorageV2(id, baseBlock, pposStorageRlps); err != nil {
				return fmt.Errorf("reply last ppos storage v2 message  fail,%v", err)
			}
		} else {
			for _, encoded := range blocks {
				if byteSize >= softResponseLimit {
					id := rand.Uint64()
					if err := peer.ReplyPPOSStorageV2(id, baseBlock, pposStorageRlps); err != nil {
						return fmt.Errorf("reply ppos storage v2 message  fail,%v", err)
					}
					pposStorageRlps = []rlp.RawValue{}
					byteSize = 0
				} else {
					pposStorageRlps = append(pposStorageRlps, encoded)
					byteSize += len(encoded)
				}
			}
			if len(pposStorageRlps) > 0 {
				id := rand.Uint64()
				if err := peer.ReplyPPOSStorageV2(id, baseBlock, pposStorageRlps); err != nil {
					return fmt.Errorf("reply last ppos storage v2 message  fail,%v", err)
				}
			}
		}
		peer.Log().Debug("end answerGetPPOSStorageMsgQueryV2", "blocks", len(blocks), "baseBlock", baseBlock)
		return nil
	}
	go func() {
		if err := sdb.WalkDB(head, f); err != nil {
			peer.Log().Error("answer GetPPOSStorageMsgQueryV2  fail", "error", err)
		}
	}()
	return nil
}

// handlePPosStorageMsg handles PPOS msg, collect the requested info and reply
func handlePPosStorageMsg(backend Backend, msg Decoder, peer *Peer) error {

	peer.Log().Debug("Received a broadcast message[PposStorageMsg]")
	var data PposStoragePacket
	if err := msg.Decode(&data); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Deliver all to the downloader
	return backend.Handle(peer, &data)
}

// handlePPosStorageMsg handles PPOS msg, collect the requested info and reply
func handlePPosStorageV2Msg(backend Backend, msg Decoder, peer *Peer) error {

	peer.Log().Debug("Received a broadcast message[PposStorageV2Msg]")
	var data PposStorageV2Packet
	if err := msg.Decode(&data); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Deliver all to the downloader
	return backend.Handle(peer, &data)
}

func handleGetOriginAndPivotMsg(backend Backend, msg Decoder, peer *Peer) error {

	peer.Log().Info("[GetOriginAndPivotMsg]Received a broadcast message")
	var query uint64
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	oHead := backend.Chain().GetHeaderByNumber(query)
	pivot, err := snapshotdb.Instance().BaseNum()
	if err != nil {
		peer.Log().Error("GetOriginAndPivot get snapshotdb baseNum fail", "err", err)
		return errors.New("GetOriginAndPivot get snapshotdb baseNum fail")
	}
	if pivot == nil {
		peer.Log().Error("[GetOriginAndPivot] pivot should not be nil")
		return errors.New("[GetOriginAndPivot] pivot should not be nil")
	}
	pHead := backend.Chain().GetHeaderByNumber(pivot.Uint64())

	data := make([]*types.Header, 0)
	data = append(data, oHead, pHead)
	if err := peer.SendOriginAndPivot(data); err != nil {
		peer.Log().Error("[GetOriginAndPivotMsg]send data meassage fail", "error", err)
		return err
	}
	return nil
}

func handleOriginAndPivotMsg(backend Backend, msg Decoder, peer *Peer) error {
	peer.Log().Debug("[OriginAndPivotMsg]Received a broadcast message")
	var data OriginAndPivotPacket
	if err := msg.Decode(&data); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Deliver all to the downloader
	return backend.Handle(peer, &data)
}

func handlePPOSInfoMsg(backend Backend, msg Decoder, peer *Peer) error {
	peer.Log().Debug("Received a broadcast message[PPOSInfoMsg]")
	var data PposInfoPacket
	if err := msg.Decode(&data); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Deliver all to the downloader
	return backend.Handle(peer, &data)
}
