package consensus

type EvidenceType int32

type Evidence interface {
	//Verify(ecdsa.PublicKey) error
	//Equal(Evidence) bool
	//return lowest number
	BlockNumber() uint64
	Epoch() uint64
	ViewNumber() uint64
	Hash() []byte
	//Address() common.Address
	//Validate() error
	//Type() EvidenceType
}

type Evidences []Evidence

type EvidencePool interface {
	//Deserialization of evidence
	UnmarshalEvidence([]byte) (Evidence, error)
	//Get current evidences
	Evidences() []Evidence
	//Clear all evidences
	Clear(viewNumber uint64)
	Close()
}
