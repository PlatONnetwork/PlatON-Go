package consensus

type EvidenceType int32

type Evidence interface {
	//Verify(ecdsa.PublicKey) error
	Equal(Evidence) bool
	BlockNumber() uint64
	Epoch() uint64
	ViewNumber() uint64
	Hash() []byte
	//Address() common.Address
	Validate() error
	Type() EvidenceType
}

type Evidences []Evidence

func (e Evidences) Len() int {
	return len(e)
}

type EvidencePool interface {
	//Deserialization of evidence
	UnmarshalEvidence(data string) (Evidences, error)
	//Get current evidences
	Evidences() Evidences
	//Clear all evidences
	Clear(epoch uint64, viewNumber uint64)
	Close()
}
