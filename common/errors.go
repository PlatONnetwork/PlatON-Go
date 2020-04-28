package common

import "fmt"

var (
	OkCode           = uint32(0)
	NoErr            = &BizError{Code: OkCode, Msg: "ok"}
	InternalError    = &BizError{Code: 1, Msg: "System error"}
	NotFound         = &BizError{Code: 2, Msg: "Object not found"}
	InvalidParameter = &BizError{Code: 3, Msg: "Invalid parameter"}
)

// business error, Gas will not be returned back to caller
type BizError struct {
	Code uint32
	Msg  string
}

func (e *BizError) Error() string {
	return e.Msg
}

func NewBizError(code uint32, text string) *BizError {
	return &BizError{Code: code, Msg: text}
}

func (be *BizError) Wrap(text string) *BizError {
	return &BizError{Code: be.Code, Msg: be.Msg + ":" + text}

}

func (be *BizError) Wrapf(format string, a ...interface{}) *BizError {
	return &BizError{Code: be.Code, Msg: be.Msg + ":" + fmt.Sprintf(format, a...)}
}

func (be *BizError) AppendMsg(msg string) {
	be.Msg = be.Msg + ":" + msg
}

func DecodeError(err error) (uint32, string) {
	if err == nil {
		return NoErr.Code, NoErr.Msg
	}
	switch typed := err.(type) {
	case *BizError:
		return typed.Code, typed.Msg
	default:
	}
	return InternalError.Code, err.Error()
}
