package crypto

type PrivateKey struct {
}

type PublicKey struct {
}

type Signature struct {
}

// GenerateKey generates a private key.
func GenerateKey() (*PrivateKey, error) {
	return &PrivateKey{}, nil
}

// Recover recover a private key from hex string.
func (pk *PrivateKey) Recover(s string) error {
	return nil
}

// Public returns the public key corresponding to private key.
func (pk *PrivateKey) Public() *PublicKey {
	return nil
}

// Sign signs digest with private key.
func (pk *PrivateKey) Sign(digest []byte) (*Signature, error) {
	return nil, nil
}

// String returns a hex string representation of the private key.
func (pk *PrivateKey) String() string {
	return ""
}

// Bytes returns a byte slice representation of the private key.
func (pk *PrivateKey) Bytes() []byte {
	return []byte{}
}

// Add add a private key to aggregate.
func (pk *PrivateKey) Add(rhs *PrivateKey) {
}

// Recover generates a public key from hex string.
func (pk *PublicKey) Recover(s string) error {
	return nil
}

// String returns a hex string representation of the public key.
func (pk *PublicKey) String() string {
	return ""
}

// Bytes returns a byte slice representation of the public key.
func (pk *PublicKey) Bytes() []byte {
	return []byte{}
}

// Add add a public key to aggregate.
func (pk *PublicKey) Add(rhs *PublicKey) {
}

// Equal returns a boolean reporting whether pk and rhs are the same.
func (pk *PublicKey) Equal(rhs *PublicKey) bool {
	return true
}

// Recover recover a signature from hex string.
func (sig *Signature) Recover(s string) error {
	return nil
}

// Verify verifies the signature using the public key. Its return
// value records whether the signature is valid.
func (sig *Signature) Verify(pub *PublicKey, m string) bool {
	return true
}

// Add add the signature to aggregate.
func (sig *Signature) Add(rhs *Signature) {
}

// String returns a hex string representation of the signature.
func (sig *Signature) String() string {
	return ""
}

// Bytes returns a byte slice representation of the signature.
func (sig *Signature) Bytes() []byte {
	return []byte{}
}
