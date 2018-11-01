package cbft

import (
	"Platon-go/p2p/discover"
	"crypto/ecdsa"
)

type CbftConfig struct {
	Period           uint64  `json:"period"`           // Number of seconds between blocks to enforce
	Epoch            uint64  `json:"epoch"`            // Epoch length to reset votes and checkpoint
	MaxLatency       int64   `json:"maxLatency"`       //共识节点间最大网络延迟时间，单位：毫秒
	LegalCoefficient float64 `json:"legalCoefficient"` //检查块的合法性时的用到的时间系数
	Duration         int64   `json:"duration"`         //每个出块节点的出块时长，单位：秒
	//mock
	InitialNodes []discover.Node   `json:"initialNodes"`
	NodeID       discover.NodeID   `json:"nodeID,omitempty"`
	PrivateKey   *ecdsa.PrivateKey `json:"PrivateKey,omitempty"`
}
