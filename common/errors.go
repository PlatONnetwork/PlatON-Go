package common

import "fmt"

// business error, Gas will not be returned back to caller
type BizError struct {
	s string
}

func (e *BizError) Error() string {
	return e.s
}

func NewBizError(text string) error {
	return &BizError{text}
}

func BizErrorf(format string, a ...interface{}) error {
	return NewBizError(fmt.Sprintf(format, a...))
}

// system error, Gas will be returned back to caller
type SysError struct {
	s string
}

func (e *SysError) Error() string {
	return e.s
}

func NewSysError(text string) error {
	return &SysError{text}
}

func SysErrorf(format string, a ...interface{}) error {
	return NewSysError(fmt.Sprintf(format, a...))
}
