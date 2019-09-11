package common

import "fmt"

var (
	SuccessCode      = uint16(0)
	Success          = &BizError{Code: SuccessCode, Msg: "Success"}
	InternalError    = &BizError{Code: 1, Msg: "System error"}
	NotFound         = &BizError{Code: 2, Msg: "Object not found"}
	InvalidParameter = &BizError{Code: 3, Msg: "Invalid parameter"}
)

// business error, Gas will not be returned back to caller
type BizError struct {
	Code uint16
	Msg  string
	Err  error
}

func (e *BizError) Error() string {
	return e.Msg
}

func NewBizError(code uint16, text string) *BizError {
	return &BizError{Code: code, Msg: text}
}

func NewBizErrorf(code uint16, format string, a ...interface{}) *BizError {
	return NewBizError(code, fmt.Sprintf(format, a...))
}

func NewBizErrorw(code uint16, text string, err error) *BizError {
	return &BizError{Code: code, Msg: text, Err: err}
}

func (be *BizError) Wrap(text string) *BizError {

	return &BizError{Code: be.Code, Msg: be.Msg + " " + text, Err: be.Err}

}

func (be *BizError) Wrapf(format string, a ...interface{}) *BizError {

	return &BizError{Code: be.Code, Msg: be.Msg + " " + fmt.Sprintf(format, a), Err: be.Err}

}

func DecodeError(err error) (uint16, string) {
	if err == nil {
		return Success.Code, Success.Msg
	}
	switch typed := err.(type) {
	case *BizError:
		return typed.Code, typed.Msg
	default:
	}
	return InternalError.Code, err.Error()
}
