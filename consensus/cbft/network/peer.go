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

package network

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	mapset "github.com/deckarep/golang-set"
)

var (
	errClosed                  = errors.New("peer set is closed")               // errClosed represents the node is closed error message description.
	errAlreadyRegistered       = errors.New("peer is already registered")       // errAlreadyRegistered represents a node registered error.
	errNotRegistered           = errors.New("peer is not registered")           // errNotRegistered represents that the node is not registered.
	errInvalidHandshakeMessage = errors.New("invalid handshake message params") // The parameters passed in the node handshake are not correct.
	errForkBlock               = errors.New("forked block")                     // Means that when the block heights are equal and the block hashes are not equal.
)

const (
	// The maximum number of queues for message packets
	// that are communicated by peers.
	maxKnownMessageHash = 20000

	// Protocol handshake timeout period, handshake failure after timeout.
	handshakeTimeout = 5 * time.Second

	// Heartbeat detection interval (unit: second).
	pingInterval = 15 * time.Second

	// maxQueueSize is maximum threshold for the queue of messages waiting to sent.
	maxQueueSize = 128
)

// Peer represents a node in the network.
type peer struct {
	*p2p.Peer                   // Network layer p2p reference.
	id        string            // Peer id identifier
	rw        p2p.MsgReadWriter //
	version   int               // Protocol version negotiated
	term      chan struct{}     // Termination channel to stop the broadcaster

	// Node status information
	highestQCBn *big.Int     // The highest QC height of the node.
	qcLock      sync.RWMutex //
	lockedBn    *big.Int     // The highest Lock height of the node.
	lLock       sync.RWMutex //
	commitBn    *big.Int     // The highest Commit height of the node.
	cLock       sync.RWMutex //

	// Record the message received by the peer node.
	// If the threshold is exceeded, the queue tail
	// record is popped up and then added.
	knownMessageHash mapset.Set

	pingList *list.List
	listLock sync.RWMutex

	// Message sending queue, the queue stores
	// messages to be sent to the peer.
	sendQueue chan *types.MsgPackage
}

// newPeer creates a new peer.
func newPeer(pv int, p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
	return &peer{
		Peer:             p,
		rw:               rw,
		id:               p.ID().TerminalString(),
		term:             make(chan struct{}),
		version:          pv,
		highestQCBn:      new(big.Int),
		lockedBn:         new(big.Int),
		commitBn:         new(big.Int),
		knownMessageHash: mapset.NewSet(),
		pingList:         list.New(),
		sendQueue:        make(chan *types.MsgPackage, maxQueueSize),
	}
}

// Return the id of peer.
func (p *peer) PeerID() string {
	return p.id
}

// Return p2p.MsgReadWriter from peer.
func (p *peer) ReadWriter() p2p.MsgReadWriter {
	return p.rw
}

// ListLen returns the number of elements of list l.
// The complexity is O(1).
func (p *peer) ListLen() int {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	return p.pingList.Len()
}

// ListFront returns the first element of list l
// or nil if the list is empty.
func (p *peer) ListFront() *list.Element {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	return p.pingList.Front()
}

// ListRemove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (p *peer) ListRemove(e *list.Element) interface{} {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	return p.pingList.Remove(e)
}

// ListPushFront inserts a new element e with value v at the
// front of list l and returns e.
func (p *peer) ListPushFront(v interface{}) *list.Element {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	return p.pingList.PushFront(v)
}

// ListPushBack inserts a new element e with value v at the
// back of list l and returns e.
func (p *peer) ListPushBack(v interface{}) *list.Element {
	p.listLock.Lock()
	defer p.listLock.Unlock()
	return p.pingList.PushBack(v)
}

// Handshake passes each other's status data and verifies the protocol version,
// the successful handshake can successfully establish a connection by peer.
func (p *peer) Handshake(outStatus *protocols.CbftStatusData) (*protocols.CbftStatusData, error) {
	if nil == outStatus {
		return nil, errInvalidHandshakeMessage
	}
	errc := make(chan error, 2)
	var inStatus protocols.CbftStatusData
	// Asynchronously send status information of the local node.
	go func() {
		errc <- p2p.Send(p.rw, protocols.CBFTStatusMsg, outStatus)
	}()
	// Asynchronously waiting to receive status data sent by the peer.
	go func() {
		errc <- p.readStatus(&inStatus)
	}()
	timeout := time.NewTicker(handshakeTimeout)
	defer timeout.Stop()
	for i := 0; i < 2; i++ {
		select {
		case err := <-errc:
			if err != nil {
				return nil, err
			}
		case <-timeout.C:
			return nil, p2p.DiscReadTimeout
		}
	}
	// If the height of the peer block is less than local,
	// determine whether it belongs to the fork block.
	if inStatus.QCBn.Uint64() == outStatus.QCBn.Uint64() && inStatus.QCBlock != outStatus.QCBlock {
		log.Error("Unmatched block on the QC", "localNumber", outStatus.QCBn.Uint64(),
			"localHash", outStatus.QCBlock.TerminalString(),
			"remoteNumber", inStatus.QCBn.Uint64(),
			"remoteHash", inStatus.QCBlock.TerminalString())
		return nil, errForkBlock
	}
	if inStatus.LockBn.Uint64() == outStatus.LockBn.Uint64() && inStatus.LockBlock != outStatus.LockBlock {
		log.Error("Unmatched block on the locked", "localNumber", outStatus.LockBn.Uint64(),
			"localHash", outStatus.LockBlock.TerminalString(),
			"remoteNumber", inStatus.LockBn.Uint64(),
			"remoteHash", inStatus.LockBlock.TerminalString())
		return nil, errForkBlock
	}
	if inStatus.CmtBn.Uint64() == outStatus.CmtBn.Uint64() && inStatus.CmtBlock != outStatus.CmtBlock {
		log.Error("Unmatched block on the commit", "localNumber", outStatus.CmtBn.Uint64(),
			"localHash", outStatus.CmtBlock.TerminalString(),
			"remoteNumber", inStatus.CmtBn.Uint64(),
			"remoteHash", inStatus.CmtBlock.TerminalString())
		return nil, errForkBlock
	}

	// 1、If the QCBlock from another peer is less than the current node,
	// determine if the local node contains a block height and a hash that matches it.
	// qcBn/lockedBn/commitBn.
	p.highestQCBn, p.lockedBn, p.commitBn = inStatus.QCBn, inStatus.LockBn, inStatus.CmtBn
	log.Debug("Handshake success and done", "remoteQCBn", p.QCBn(),
		"remoteLockedBn", p.LockedBn(), "remoteCommitBn", p.CommitBn())

	return &inStatus, nil
}

// readStatus receive status data from another.
func (p *peer) readStatus(status *protocols.CbftStatusData) error {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Code != protocols.CBFTStatusMsg {
		return types.ErrResp(types.ErrNoStatusMsg, "first msg has code %x (!= %x)", msg.Code, protocols.CBFTStatusMsg)
	}
	if msg.Size > protocols.CbftProtocolMaxMsgSize {
		return types.ErrResp(types.ErrMsgTooLarge, "%v > %v", msg.Size, protocols.CbftProtocolMaxMsgSize)
	}
	if err := msg.Decode(&status); err != nil {
		return types.ErrResp(types.ErrDecode, "msg %v: %v", msg, err)
	}
	if int(status.ProtocolVersion) != p.version {
		return types.ErrResp(types.ErrCbftProtocolVersionMismatch, "%d (!= %d)", status.ProtocolVersion, p.version)
	}
	return nil
}

// MarkMessageHash is used to record the hash value of each message from the peer node.
// If the queue is full, remove the bottom element and add a new one.
func (p *peer) MarkMessageHash(hash common.Hash) {
	for p.knownMessageHash.Cardinality() >= maxKnownMessageHash {
		p.knownMessageHash.Pop()
	}
	p.knownMessageHash.Add(hash)
}

// ContainsMessageHash determines if the specified message hash is included.
func (p *peer) ContainsMessageHash(hash common.Hash) bool {
	return p.knownMessageHash.Contains(hash)
}

// RemoveMessageHash remove the msg from knownMessageHash.
func (p *peer) RemoveMessageHash(hash common.Hash) {
	p.knownMessageHash.Remove(hash)
}

// Close terminates the running state of the peer.
func (p *peer) Close() {
	close(p.term)
}

// SetHighest saves the highest QC block.
func (p *peer) SetQcBn(qcBn *big.Int) {
	if qcBn != nil {
		p.qcLock.Lock()
		defer p.qcLock.Unlock()
		log.Trace("Set QCBn", "peerID", p.id, "oldQCBn", p.highestQCBn.Uint64(), "newQCBn", qcBn.Uint64())
		p.highestQCBn.Set(qcBn)
	}
}

// Get the highest block height signature collected by the current node.
func (p *peer) QCBn() uint64 {
	p.qcLock.RLock()
	defer p.qcLock.RUnlock()
	return p.highestQCBn.Uint64()
}

// SetLockedBn saves the highest locked block.
func (p *peer) SetLockedBn(lockedBn *big.Int) {
	if lockedBn != nil {
		p.lLock.Lock()
		defer p.lLock.Unlock()
		log.Debug("Set lockedBn", "peerID", p.id, "oldLockedBn", p.lockedBn.Uint64(), "newLockedBn", lockedBn.Uint64())
		p.lockedBn.Set(lockedBn)
	}
}

// Get the highest block height locked by the current node.
func (p *peer) LockedBn() uint64 {
	p.lLock.RLock()
	defer p.lLock.RUnlock()
	return p.lockedBn.Uint64()
}

// SetLockedBn saves the highest commit block.
func (p *peer) SetCommitdBn(commitBn *big.Int) {
	if commitBn != nil {
		p.cLock.Lock()
		defer p.cLock.Unlock()
		log.Debug("Set commitBn", "peerID", p.id, "oldCommitBn", p.commitBn.Uint64(), "newCommitBn", commitBn.Uint64())
		p.commitBn.Set(commitBn)
	}
}

// Get the highest block height submitted by the current node.
func (p *peer) CommitBn() uint64 {
	p.cLock.RLock()
	defer p.cLock.RUnlock()
	return p.commitBn.Uint64()
}

// Start the loop that the peer uses to maintain its
// own functions.
func (p *peer) Run() {
	go p.pingLoop()
	go p.sendLoop()
}

// Send send the message
func (p *peer) Send(msg *types.MsgPackage) {
	select {
	case p.sendQueue <- msg:
	default:
		log.Debug("Send message fail, message queue blocking", "peer", p.PeerID(), "type", reflect.TypeOf(msg.Message()), "msgHash", msg.Message().MsgHash(), "msg", msg.Message().String())
	}
}

// The loop of heartbeat detection is mainly responsible for
// confirming the connection of the connection.
func (p *peer) pingLoop() {
	ping := time.NewTimer(pingInterval)
	defer ping.Stop()
	for {
		select {
		case <-ping.C:
			// Send a ping message directly and the response message
			// is processed at the CBFT layer.
			pingTime := strconv.FormatInt(time.Now().UnixNano(), 10)
			if p.ListLen() > 5 {
				front := p.ListFront()
				p.ListRemove(front)
			}
			p.ListPushBack(pingTime)

			log.Trace("Send a ping message", "peerID", p.ID(), "pingTimeNano", pingTime, "pingList.Len", p.pingList.Len())
			if err := p2p.SendItems(p.rw, protocols.PingMsg, pingTime); err != nil {
				log.Error("Send ping message failed", "err", err)
				return
			}
			ping.Reset(pingInterval)
		case <-p.term:
			log.Trace("Ping loop term", "peerID", p.ID().TerminalString())
			return
		}
	}
}

// sendLoop the loop reads data from message queue and sends it.
func (p *peer) sendLoop() {
	for {
		select {
		case msg := <-p.sendQueue:
			msgType := protocols.MessageType(msg.Message())
			if err := p2p.Send(p.rw, msgType, msg.Message()); err != nil {
				log.Error("Send message fail", "peer", p.PeerID(), "msg", msg.Message().String(), "err", err)
			} else {
				if msg.Mode() == types.FullMode {
					p.MarkMessageHash(msg.Message().MsgHash())
				}
			}
		case <-p.term:
			return
		}
	}
}

// PeerInfo represents the node information of the CBFT protocol.
type PeerInfo struct {
	ProtocolVersion int    `json:"protocolVersion"`
	HighestQCBn     uint64 `json:"highestQCBn"`
	LockedBn        uint64 `json:"lockedBn"`
	CommitBn        uint64 `json:"commitBn"`
}

// Info output status information of the current peer.
func (p *peer) Info() *PeerInfo {
	pv, qc, locked, commit := p.version, p.highestQCBn.Uint64(), p.lockedBn.Uint64(), p.commitBn.Uint64()
	return &PeerInfo{
		ProtocolVersion: pv,
		HighestQCBn:     qc,
		LockedBn:        locked,
		CommitBn:        commit,
	}
}

// PeerSet represents the collection of active peers currently participating
// in the Cbft protocol.
type PeerSet struct {
	peers  map[string]*peer
	lock   sync.RWMutex
	closed bool
}

// NewPeerSet creates a new PeerSet to track the active participants.
func NewPeerSet() *PeerSet {
	ps := &PeerSet{
		peers: make(map[string]*peer),
	}
	// start a goroutine timing output A connection status information
	go ps.printPeers()
	return ps
}

// Register injects a new peer into the working set, or
// returns an error if the peer is already known. If a new peer it registered,
// its broadcast loop is also started.
func (ps *PeerSet) Register(p *peer) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	if ps.closed {
		return errClosed
	}
	if _, ok := ps.peers[p.id]; ok {
		return errAlreadyRegistered
	}
	ps.peers[p.id] = p
	return nil
}

// Unregister removes a remote peer from the active set, disabling any further
// actions to/from that particular entity.
func (ps *PeerSet) Unregister(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	p, ok := ps.peers[id]
	if !ok {
		return errNotRegistered
	}
	delete(ps.peers, id)
	p.Close()

	return nil
}

// Get retrieves the registered peer with the given id.
func (ps *PeerSet) get(id string) (*peer, error) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	p, ok := ps.peers[id]
	if !ok {
		return nil, errNotRegistered
	}

	return p, nil
}

// Len returns if the current number of peers in the set.
func (ps *PeerSet) Len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// peersWithConsensus retrieves a list of peers that exist with the PeerSet based
// on the incoming consensus node ID array.
func (ps *PeerSet) peersWithConsensus(consensusNodes []discover.NodeID) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(consensusNodes))
	for _, nodeID := range consensusNodes {
		nodeID := nodeID.TerminalString()
		if peer, ok := ps.peers[nodeID]; ok {
			list = append(list, peer)
		}
	}
	return list
}

// peersWithoutConsensus retrieves a list of peer that does not contain consensus nodes.
func (ps *PeerSet) peersWithoutConsensus(consensusNodes []discover.NodeID) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	consensusNodeMap := make(map[string]string)
	for _, nodeID := range consensusNodes {
		nodeID := nodeID.TerminalString()
		consensusNodeMap[nodeID] = nodeID
	}

	list := make([]*peer, 0, len(ps.peers))
	for nodeID, peer := range ps.peers {
		if _, ok := consensusNodeMap[nodeID]; !ok {
			list = append(list, peer)
		}
	}

	return list
}

// Peers retrieves a list of peer from the PeerSet.
func (ps *PeerSet) allPeers() []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		list = append(list, p)
	}
	return list
}

// peersWithHighestQCBn returns a list of nodes that are larger than the qcNumber of the highest qc block.
func (ps *PeerSet) peersWithHighestQCBn(qcNumber uint64) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.QCBn() > qcNumber {
			list = append(list, p)
		}
	}
	log.Trace("QCBnHighestPeer done", "count", len(list), "peers", formatPeers(list))
	return list
}

// peersWithHighestLockedBn returns a list of nodes that are larger than the lockedNumber
// of the highest locked block.
func (ps *PeerSet) peersWithHighestLockedBn(lockedNumber uint64) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.LockedBn() > lockedNumber {
			list = append(list, p)
		}
	}
	log.Trace("LockedBnHighestPeer done", "count", len(list), "peers", formatPeers(list))
	return list
}

// peersWithHighestCommitBn returns a list of nodes that are larger than the commitNumber
// of the highest locked block.
func (ps *PeerSet) peersWithHighestCommitBn(commitNumber uint64) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if p.CommitBn() > commitNumber {
			list = append(list, p)
		}
	}
	log.Trace("CommitBnHighestPeer done", "count", len(list), "peers", formatPeers(list))
	return list
}

// Close disconnects all peers. No new peers can be registered
// after Close has returned.
func (ps *PeerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}

// printPeers timing printout output list
// of neighbor nodes of the current node.
func (ps *PeerSet) printPeers() {
	// Output in 2 seconds
	outTimer := time.NewTicker(time.Second * 5)
	for {
		if ps.closed {
			break
		}
		select {
		case <-outTimer.C:
			peers := ps.allPeers()
			if peers != nil {
				neighborPeerGauage.Update(int64(len(peers)))
			}
			var bf bytes.Buffer
			for idx, peer := range peers {
				bf.WriteString(peer.id)
				if idx < len(peers)-1 {
					bf.WriteString(",")
				}
			}
			pInfo := bf.String()
			log.Debug(fmt.Sprintf("The neighbor node owned by the current peer is : {%v}, size: {%d}", pInfo, len(peers)))
		}
	}
}
