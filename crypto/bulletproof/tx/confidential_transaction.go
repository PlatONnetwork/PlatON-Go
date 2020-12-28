package tx

/*
#cgo CFLAGS: -I${SRCDIR}/confidential_transaction/confidentialtx_ffi/include
#cgo LDFLAGS: -L${SRCDIR}/confidential_transaction/target/release
#cgo linux LDFLAGS: -lconfidentialtx -lm -lpthread -ldl
#cgo darwin LDFLAGS: -lconfidentialtx

#include "confidentialtx.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

type TxType uint8
type TxLog struct {
	/// Type of Tx
	TxType TxType

	/// Input notes of the transaction
	Inputs []Note

	/// Output notes of the transaction
	Outputs []Note
}
type Note struct {
	EphemeralPk []byte

	SpendingPk []byte
	/// Confidential token.
	Token []byte
}

func VerifyConfidentialTx(proof []byte) ([]byte, error) {
	var err C.ConfidentialTxError
	res := C.confidential_tx_verify((*C.uchar)(unsafe.Pointer(&proof[0])), C.int(len(proof)), &err)
	if err.code != 0 {
		return nil, errors.New(C.GoString(err.message))
	}
	defer func() {
		C.confidential_tx_destroy_string(err.message)
		C.confidential_tx_destroy_bytebuffer(res)
	}()

	return C.GoBytes(unsafe.Pointer(res.data), C.int(res.len)), nil
}

func CreateConfidentialTx(proof []byte) ([]byte, error) {
	var err C.ConfidentialTxError
	res := C.create_confidential_tx((*C.uchar)(unsafe.Pointer(&proof[0])), C.int(len(proof)), &err)
	if err.code != 0 {
		return nil, errors.New(C.GoString(err.message))
	}
	defer func() {
		C.confidential_tx_destroy_string(err.message)
		C.confidential_tx_destroy_bytebuffer(res)
	}()

	return C.GoBytes(unsafe.Pointer(res.data), C.int(res.len)), nil
}
