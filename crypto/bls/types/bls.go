package types

type ISecretKey interface {
	Sign(m string) (sign ISign)
	SetByCSPRNG()
	GetPublicKey() (pub IPublicKey)
	MakeSchnorrNIZKP() (ISchnorrProof, error)
	Serialize() []byte
	Deserialize(buf []byte) error
	SetLittleEndian(buf []byte) error
	GetLittleEndian() []byte
}

type IPublicKey interface {
	Add(rhs IPublicKey)
	Serialize() []byte
	Bytes() []byte
	UnmarshalText(text []byte) error
	MarshalText() ([]byte, error)
}

type IPublicKeyHex interface {
	ParseBlsPubKey() (IPublicKey, error)
	Bytes() []byte
	UnmarshalText(text []byte) error
}
type ISchnorrProof interface {
	VerifySchnorrNIZK(pk IPublicKey) error
	UnmarshalText(text []byte) error
	MarshalText() ([]byte, error)
	Serialize() []byte
	Deserialize(buf []byte) error
}

type ISchnorrProofHex interface {
	MarshalText() ([]byte, error)
	UnmarshalText(text []byte) error
}

type ISign interface {
	Add(sign ISign)
	Verify(pub IPublicKey, m string) bool
	Serialize() []byte
	Deserialize(buf []byte) error
}
