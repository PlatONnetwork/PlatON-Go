package common

import "fmt"

var (
	Success          = &BizError{Code: 0, Msg: "Success"}
	InternalError    = &BizError{Code: -1, Msg: "System error"}
	NotFound         = &BizError{Code: 1, Msg: "Object not found"}
	InvalidParameter = &BizError{Code: 2, Msg: "Invalid parameter"}
)

// business error, Gas will not be returned back to caller
type BizError struct {
	Code int
	Msg  string
	Err  error
}

func (e *BizError) Error() string {
	return e.Msg
}

func NewBizError(code int, text string) *BizError {
	return &BizError{Code: code, Msg: text}
}

func NewBizErrorf(code int, format string, a ...interface{}) *BizError {
	return NewBizError(code, fmt.Sprintf(format, a...))
}

func NewBizErrorw(code int, text string, err error) *BizError {
	return &BizError{Code: code, Msg: text, Err: err}
}

func (be *BizError) Wrap(text string) *BizError {

	return &BizError{Code: be.Code, Msg: be.Msg + " " + text, Err: be.Err}

}

func DecodeError(err error) (int, string) {
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
