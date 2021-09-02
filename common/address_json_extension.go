package common

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type addressCodec struct{}

func (codec *addressCodec) IsEmpty(ptr unsafe.Pointer) bool {
	addr := *((*Address)(ptr))
	return addr == ZeroAddr
}
func (codec *addressCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	addr := *((*Address)(ptr))
	stream.WriteString(addr.Hex())
}

type addressPtrCodec struct {
	valType reflect2.Type
}

func (codec *addressPtrCodec) IsEmpty(ptr unsafe.Pointer) bool {
	addr := codec.valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(addr) {
		return true
	}
	return false
}
func (codec *addressPtrCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	obj := codec.valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(obj) {
		stream.WriteNil()
		return
	}
	addr := *((**Address)(ptr))
	stream.WriteString(addr.Hex())
}

type AddressExtension struct {
	jsoniter.EncoderExtension
}

var (
	addressPtrType = reflect2.TypeOfPtr((*Address)(nil))
	addressType    = reflect2.TypeOfPtr((*Address)(nil)).Elem()
)

// CreateEncoder get encoder from map
func (extension *AddressExtension) CreateEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	if typ == addressType {
		return &addressCodec{}
	} else if typ == addressPtrType {
		return &addressPtrCodec{typ}
	}
	return nil
}

func (extension *AddressExtension) CreateMapKeyEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	if typ == addressType {
		return &addressCodec{}
	} else if typ == addressPtrType {
		return &addressPtrCodec{typ}
	}
	return nil
}
