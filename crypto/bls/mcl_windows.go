// +build windows

package bls

/*
#cgo CFLAGS: -I./src/
#cgo CFLAGS:-DMCLBN_FP_UNIT_SIZE=6
#include <stdio.h>
#include <mcl/bn.h>
#include <bls/bls.h>
*/
import "C"

// CurveFp254BNb -- 254 bit curve
const CurveFp254BNb = C.mclBn_CurveFp254BNb

// CurveFp382_1 -- 382 bit curve 1
const CurveFp382_1 = C.mclBn_CurveFp382_1

// CurveFp382_2 -- 382 bit curve 2
const CurveFp382_2 = C.mclBn_CurveFp382_2

//const CurveSNARK1 = C.mclBn_CurveSNARK1

//const Bls12_CurveFp381 = C.mclBls12_CurveFp381

// GetMaxOpUnitSize --
func GetMaxOpUnitSize() int {
	return 0
}

// GetOpUnitSize --
// the length of Fr is GetOpUnitSize() * 8 bytes
func GetOpUnitSize() int {
	return 0
}

// GetCurveOrder --
// return the order of G1
func GetCurveOrder() string {
	return ""
}

// GetFieldOrder --
// return the characteristic of the field where a curve is defined
func GetFieldOrder() string {
	return ""
}

// Fr --
type Fr struct {
	v C.mclBnFr
}

// getPointer --
func (x *Fr) getPointer() (p *C.mclBnFr) {
	return nil
}

// Clear --
func (x *Fr) Clear() {
}

// SetInt64 --
func (x *Fr) SetInt64(v int64) {
}

// SetString --
func (x *Fr) SetString(s string, base int) error {
	return nil
}

// Deserialize --
func (x *Fr) Deserialize(buf []byte) error {
	return nil
}

// SetLittleEndian --
func (x *Fr) SetLittleEndian(buf []byte) error {
	return nil
}

// IsEqual --
func (x *Fr) IsEqual(rhs *Fr) bool {
	return false
}

// IsZero --
func (x *Fr) IsZero() bool {
	return false
}

// IsOne --
func (x *Fr) IsOne() bool {
	return false
}

// SetByCSPRNG --
func (x *Fr) SetByCSPRNG() {
}

// SetHashOf --
func (x *Fr) SetHashOf(buf []byte) bool {
	return false
}

// GetString --
func (x *Fr) GetString(base int) string {
	return ""
}

// Serialize --
func (x *Fr) Serialize() []byte {
	return nil
}

// FrNeg --
func FrNeg(out *Fr, x *Fr) {
}

// FrInv --
func FrInv(out *Fr, x *Fr) {
}

// FrAdd --
func FrAdd(out *Fr, x *Fr, y *Fr) {
}

// FrSub --
func FrSub(out *Fr, x *Fr, y *Fr) {
}

// FrMul --
func FrMul(out *Fr, x *Fr, y *Fr) {
}

// FrDiv --
func FrDiv(out *Fr, x *Fr, y *Fr) {
}

// G1 --
type G1 struct {
	v C.mclBnG1
}

// getPointer --
func (x *G1) getPointer() (p *C.mclBnG1) {
	return nil
}

// Clear --
func (x *G1) Clear() {
}

// SetString --
func (x *G1) SetString(s string, base int) error {
	return nil
}

// Deserialize --
func (x *G1) Deserialize(buf []byte) error {
	return nil
}

// IsEqual --
func (x *G1) IsEqual(rhs *G1) bool {
	return false
}

// IsZero --
func (x *G1) IsZero() bool {
	return false
}

// HashAndMapTo --
func (x *G1) HashAndMapTo(buf []byte) error {
	return nil
}

// GetString --
func (x *G1) GetString(base int) string {
	return ""
}

// Serialize --
func (x *G1) Serialize() []byte {
	return nil
}

// G1Neg --
func G1Neg(out *G1, x *G1) {
}

// G1Dbl --
func G1Dbl(out *G1, x *G1) {
}

// G1Add --
func G1Add(out *G1, x *G1, y *G1) {
}

// G1Sub --
func G1Sub(out *G1, x *G1, y *G1) {
}

// G1Mul --
func G1Mul(out *G1, x *G1, y *Fr) {
}

// G1MulCT -- constant time (depending on bit lengh of y)
func G1MulCT(out *G1, x *G1, y *Fr) {
}

// G2 --
type G2 struct {
	v C.mclBnG2
}

// getPointer --
func (x *G2) getPointer() (p *C.mclBnG2) {
	return nil
}

// Clear --
func (x *G2) Clear() {
}

// SetString --
func (x *G2) SetString(s string, base int) error {
	return nil
}

// Deserialize --
func (x *G2) Deserialize(buf []byte) error {
	return nil
}

// IsEqual --
func (x *G2) IsEqual(rhs *G2) bool {
	return false
}

// IsZero --
func (x *G2) IsZero() bool {
	return false
}

// HashAndMapTo --
func (x *G2) HashAndMapTo(buf []byte) error {
	return nil
}

// GetString --
func (x *G2) GetString(base int) string {
	return ""
}

// Serialize --
func (x *G2) Serialize() []byte {
	return nil
}

// G2Neg --
func G2Neg(out *G2, x *G2) {
}

// G2Dbl --
func G2Dbl(out *G2, x *G2) {
}

// G2Add --
func G2Add(out *G2, x *G2, y *G2) {
}

// G2Sub --
func G2Sub(out *G2, x *G2, y *G2) {
}

// G2Mul --
func G2Mul(out *G2, x *G2, y *Fr) {
}

// GT --
type GT struct {
	v C.mclBnGT
}

// getPointer --
func (x *GT) getPointer() (p *C.mclBnGT) {
	return nil
}

// Clear --
func (x *GT) Clear() {
}

// SetInt64 --
func (x *GT) SetInt64(v int64) {
}

// SetString --
func (x *GT) SetString(s string, base int) error {
	return nil
}

// Deserialize --
func (x *GT) Deserialize(buf []byte) error {
	return nil
}

// IsEqual --
func (x *GT) IsEqual(rhs *GT) bool {
	return false
}

// IsZero --
func (x *GT) IsZero() bool {
	return false
}

// IsOne --
func (x *GT) IsOne() bool {
	return false
}

// GetString --
func (x *GT) GetString(base int) string {
	return ""
}

// Serialize --
func (x *GT) Serialize() []byte {
	return nil
}

// GTNeg --
func GTNeg(out *GT, x *GT) {
}

// GTInv --
func GTInv(out *GT, x *GT) {
}

// GTAdd --
func GTAdd(out *GT, x *GT, y *GT) {
}

// GTSub --
func GTSub(out *GT, x *GT, y *GT) {
}

// GTMul --
func GTMul(out *GT, x *GT, y *GT) {
}

// GTDiv --
func GTDiv(out *GT, x *GT, y *GT) {
}

// GTPow --
func GTPow(out *GT, x *GT, y *Fr) {
}

// Pairing --
func Pairing(out *GT, x *G1, y *G2) {
}

// FinalExp --
func FinalExp(out *GT, x *GT) {
}

// MillerLoop --
func MillerLoop(out *GT, x *G1, y *G2) {
}

// GetUint64NumToPrecompute --
func GetUint64NumToPrecompute() int {
	return 0
}

// PrecomputeG2 --
func PrecomputeG2(Qbuf []uint64, Q *G2) {
}

// PrecomputedMillerLoop --
func PrecomputedMillerLoop(out *GT, P *G1, Qbuf []uint64) {
}

// PrecomputedMillerLoop2 --
func PrecomputedMillerLoop2(out *GT, P1 *G1, Q1buf []uint64, P2 *G1, Q2buf []uint64) {
}

// FrEvaluatePolynomial -- y = c[0] + c[1] * x + c[2] * x^2 + ...
func FrEvaluatePolynomial(y *Fr, c []Fr, x *Fr) error {
	return nil
}

// G1EvaluatePolynomial -- y = c[0] + c[1] * x + c[2] * x^2 + ...
func G1EvaluatePolynomial(y *G1, c []G1, x *Fr) error {
	return nil
}

// G2EvaluatePolynomial -- y = c[0] + c[1] * x + c[2] * x^2 + ...
func G2EvaluatePolynomial(y *G2, c []G2, x *Fr) error {
	return nil
}

// FrLagrangeInterpolation --
func FrLagrangeInterpolation(out *Fr, xVec []Fr, yVec []Fr) error {
	return nil
}

// G1LagrangeInterpolation --
func G1LagrangeInterpolation(out *G1, xVec []Fr, yVec []G1) error {
	return nil
}

// G2LagrangeInterpolation --
func G2LagrangeInterpolation(out *G2, xVec []Fr, yVec []G2) error {
	return nil
}
