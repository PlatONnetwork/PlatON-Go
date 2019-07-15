package cbft

import "github.com/PlatONnetwork/PlatON-Go/common/consensus"

type EvidencePool interface {
	//Deserialization of evidence
	UnmarshalEvidence([]byte) (consensus.Evidence, error)
	//Get current evidences
	Evidences() []consensus.Evidence
	//Clear all evidences
	Clear()
	Close()
}
