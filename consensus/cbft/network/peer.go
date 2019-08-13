package network

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"math/big"
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
	errClosed                  = errors.New("peer set is closed")
	errAlreadyRegistered       = errors.New("peer is already registered")
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
)

// Peer represents a node in the network.
type peer struct {
	*p2p.Peer
	id      string
	rw      p2p.MsgReadWriter
	version int           // Protocol version negotiated
	term    chan struct{} // Termination channel to stop the broadcaster

	// Node status information
	highestQCBn *big.Int
	qcLock      sync.RWMutex
	lockedBn    *big.Int
	lLock       sync.RWMutex
	commitBn    *big.Int
	cLock       sync.RWMutex

	// Record the message received by the peer node.
	// If the threshold is exceeded, the queue tail
	// record is popped up and then added.
	knownMessageHash mapset.Set

	PingList *list.List
}

// NewPeer creates a new peer.
func NewPeer(pv int, p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
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
		PingList:         list.New(),
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
	p.pingLoop()
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
			if p.PingList.Len() > 5 {
				front := p.PingList.Front()
				p.PingList.Remove(front)
			}
			p.PingList.PushBack(pingTime)

			log.Trace("Send a ping message", "peerID", p.ID(), "pingTimeNano", pingTime, "PingList.Len", p.PingList.Len())
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

// PeerInfo represents the node information of the CBFT protocol.
type PeerInfo struct {
	ProtocolVersion int    `json:"protocol_version"`
	HighestQCBn     uint64 `json:"highest_qc_bn"`
	LockedBn        uint64 `json:"locked_bn"`
	CommitBn        uint64 `json:"commit_bn"`
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

// Peer retrieves the registered peer with the given id.
func (ps *PeerSet) Get(id string) (*peer, error) {
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

// PeersWithConsensus retrieves a list of peers that exist with the PeerSet based
// on the incoming consensus node ID array.
func (ps *PeerSet) PeersWithConsensus(consensusNodes []discover.NodeID) []*peer {
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

// PeersWithoutConsensus retrieves a list of peer that does not contain consensus nodes.
func (ps *PeerSet) PeersWithoutConsensus(consensusNodes []discover.NodeID) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	consensusNodeMap := make(map[string]string)
	for _, nodeID := range consensusNodes {
		nodeID := nodeID.TerminalString()
		consensusNodeMap[nodeID] = nodeID
	}

	list := make([]*peer, 0, len(ps.peers))
	for nodeId, peer := range ps.peers {
		if _, ok := consensusNodeMap[nodeId]; !ok {
			list = append(list, peer)
		}
	}

	return list
}

// Peers retrieves a list of peer from the PeerSet.
func (ps *PeerSet) Peers() []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		list = append(list, p)
	}
	return list
}

// Returns a list of nodes that are larger than the qcNumber of the highest qc block.
func (ps *PeerSet) PeersWithHighestQCBn(qcNumber uint64) []*peer {
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

// Returns a list of nodes that are larger than the lockedNumber of the highest locked block.
func (ps *PeerSet) PeersWithHighestLockedBn(lockedNumber uint64) []*peer {
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

// Returns a list of nodes that are larger than the commitNumber of the highest locked block.
func (ps *PeerSet) PeersWithHighestCommitBn(commitNumber uint64) []*peer {
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

func (ps *PeerSet) printPeers() {
	// Output in 2 seconds
	outTimer := time.NewTicker(time.Second * 5)
	for {
		if ps.closed {
			break
		}
		select {
		case <-outTimer.C:
			peers := ps.Peers()
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
