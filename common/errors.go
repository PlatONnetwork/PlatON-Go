// Copyright 2021 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

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
	Code uint32 `json:"code"`
	Msg  string `json:"message"`
}

func (e *BizError) Error() string {
	return e.Msg
}

func (e *BizError) ErrorData() interface{} {
	return e
}

// ErrorCode returns the JSON error code for a revertal.
func (e *BizError) ErrorCode() int {
	return 4
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
