// +build windows

package bls

/*
#cgo CFLAGS: -I./src/
#include <stdio.h>
#include <bls/bls.h>
*/
import "C"
import "fmt"

// Init --
// call this function before calling all the other operations
// this function is not thread safe
func Init(curve int) error {
	return nil
}

// ID --
type ID struct {
	v Fr
}

// getPointer --
func (id *ID) getPointer() (p *C.blsId) {
	fmt.Println("getPointer.")
	return nil
}

// GetLittleEndian --
func (id *ID) GetLittleEndian() []byte {
	return nil
}

// SetLittleEndian --
func (id *ID) SetLittleEndian(buf []byte) error {
	return nil
}

// GetHexString --
func (id *ID) GetHexString() string {
	return ""
}

// GetDecString --
func (id *ID) GetDecString() string {
	return ""
}

// SetHexString --
func (id *ID) SetHexString(s string) error {
	return nil
}

// SetDecString --
func (id *ID) SetDecString(s string) error {
	return nil
}

// IsEqual --
func (id *ID) IsEqual(rhs *ID) bool {
	return false
}

// SecretKey --
type SecretKey struct {
	v Fr
}

func LoadBLS(file string) (*SecretKey, error) {
	return nil, nil
}

// getPointer --
func (sec *SecretKey) getPointer() (p *C.blsSecretKey) {
	return nil
}

// GetLittleEndian --
func (sec *SecretKey) GetLittleEndian() []byte {
	return nil
}

// SetLittleEndian --
func (sec *SecretKey) SetLittleEndian(buf []byte) error {
	return nil
}

// GetHexString --
func (sec *SecretKey) GetHexString() string {
	return ""
}

// GetDecString --
func (sec *SecretKey) GetDecString() string {
	return ""
}

// SetHexString --
func (sec *SecretKey) SetHexString(s string) error {
	return nil
}

// SetDecString --
func (sec *SecretKey) SetDecString(s string) error {
	return nil
}

// IsEqual --
func (sec *SecretKey) IsEqual(rhs *SecretKey) bool {
	return false
}

// SetByCSPRNG --
func (sec *SecretKey) SetByCSPRNG() {

}

// Add --
func (sec *SecretKey) Add(rhs *SecretKey) {
}

// Mul --
func (sec *SecretKey) Mul(rhs *SecretKey) {
}

// Sub --
func (sec *SecretKey) Sub(rhs *SecretKey) {
}

// GetMasterSecretKey --
func (sec *SecretKey) GetMasterSecretKey(k int) (msk []SecretKey) {
	return nil
}

// GetMasterPublicKey --
func GetMasterPublicKey(msk []SecretKey) (mpk []PublicKey) {
	return nil
}

// Set --
func (sec *SecretKey) Set(msk []SecretKey, id *ID) error {
	return nil
}

// Recover --
func (sec *SecretKey) Recover(secVec []SecretKey, idVec []ID) error {
	return nil
}

// GetPop --
func (sec *SecretKey) GetPop() (sign *Sign) {
	return nil
}

// PublicKey --
type PublicKey struct {
	v G2
}

func (pub *PublicKey) getQ() (p *C.blsPublicKey) {
	return nil
}

// getPointer --
func (pub *PublicKey) getPointer() (p *C.blsPublicKey) {
	return nil
}

// Serialize --
func (pub *PublicKey) Serialize() []byte {
	return nil
}

// Deserialize --
func (pub *PublicKey) Deserialize(buf []byte) error {
	return nil
}

// GetHexString --
func (pub *PublicKey) GetHexString() string {
	return ""
}

// SetHexString --
func (pub *PublicKey) SetHexString(s string) error {
	return nil
}

// IsEqual --
func (pub *PublicKey) IsEqual(rhs *PublicKey) bool {
	return false
}

func (pub *PublicKey) MarshalText() ([]byte, error) {
	return nil, nil
}

func (pub *PublicKey) UnmarshalText(text []byte) error {
	return nil
}

// Add --
func (pub *PublicKey) Add(rhs *PublicKey) {
}

// Mul --
func (pub *PublicKey) Mul(rhs *SecretKey) {
}

// Set --
func (pub *PublicKey) Set(mpk []PublicKey, id *ID) error {
	return nil
}

// Recover --
func (pub *PublicKey) Recover(pubVec []PublicKey, idVec []ID) error {
	return nil
}

// Sign  --
type Sign struct {
	v G1
}

// getPointer --
func (sign *Sign) getPointer() (p *C.blsSignature) {
	return nil
}

// Serialize --
func (sign *Sign) Serialize() []byte {
	return nil
}

// Deserialize --
func (sign *Sign) Deserialize(buf []byte) error {
	return nil
}

// GetHexString --
func (sign *Sign) GetHexString() string {
	return ""
}

// SetHexString --
func (sign *Sign) SetHexString(s string) error {
	return nil
}

// IsEqual --
func (sign *Sign) IsEqual(rhs *Sign) bool {
	return false
}

// GetPublicKey --
func (sec *SecretKey) GetPublicKey() (pub *PublicKey) {
	return nil
}

// Sign -- Constant Time version
func (sec *SecretKey) Sign(m string) (sign *Sign) {
	return nil
}

// Add --
func (sign *Sign) Add(rhs *Sign) {

}

// Recover --
func (sign *Sign) Recover(signVec []Sign, idVec []ID) error {
	return nil
}

// Verify --
func (sign *Sign) Verify(pub *PublicKey, m string) bool {
	return false
}

// VerifyPop --
func (sign *Sign) VerifyPop(pub *PublicKey) bool {
	return false
}

// DHKeyExchange --
func DHKeyExchange(sec *SecretKey, pub *PublicKey) (out PublicKey) {
	return *pub
}

//add@20190716
//get G2
func GetGeneratorOfG2() (pub *PublicKey) {
	return nil
}

// PubBatchAdd --
func PubkeyBatchAdd(pkVec []PublicKey) (pub PublicKey) {
	return PublicKey{}
}

// SecBatchAdd --
func SeckeyBatchAdd(secVec []SecretKey) (sec SecretKey) {
	return SecretKey{}
}

func AggregateSign(sigVec []Sign) (sig Sign) {
	return Sign{}
}

func GTBatchMul(eVec []GT) (e GT) {
	return GT{}
}

func GTBatchAdd(eVec []GT) (e GT) {
	return GT{}
}

func MsgsToHashToG1(mVec []string) ([]Sign, error) {
	return nil, nil
}

func BatchVerifySameMsg(curve int, msg string, pkVec []PublicKey, sign Sign) error {
	return nil
}

func BatchVerifyDistinctMsg(curve int, pkVec []PublicKey, msgVec []Sign, sign Sign) error {
	return nil
}

func SynthSameMsg(curve int, pkVec []PublicKey, mVec []string) ([]PublicKey, []string, []string, error) {
	return nil, nil, nil, nil
}

// IsValid --
func G2IsValid(rhs *PublicKey) bool {
	return false
}

func Schnorr_test(curve int, r, c SecretKey, G, V, P PublicKey) error {
	return nil
}

type Proof struct {
	C, R SecretKey
}

func SchnorrNIZKProve(curve int, sec SecretKey) (*Proof, error) {
	return nil, nil
}

func SchnorrNIZKVerify(curve int, proof Proof, P PublicKey) error {
	return nil
}
