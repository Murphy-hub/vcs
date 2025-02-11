/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package resterr

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	SystemError     ErrorCode = "system-error"
	Unauthorized    ErrorCode = "unauthorized"
	InvalidValue    ErrorCode = "invalid-value"
	AlreadyExist    ErrorCode = "already-exist"
	DoesntExist     ErrorCode = "doesnt-exist"
	ConditionNotMet ErrorCode = "condition-not-met"
)

func (c ErrorCode) Name() string {
	return string(c)
}

type CustomError struct {
	Code            ErrorCode
	IncorrectValue  string
	FailedOperation string
	Component       string
	Err             error
}

func NewSystemError(component, failedOperation string, err error) *CustomError {
	return &CustomError{
		Code:            SystemError,
		FailedOperation: failedOperation,
		Component:       component,
		Err:             err,
	}
}

func NewValidationError(code ErrorCode, incorrectValue string, err error) *CustomError {
	return &CustomError{
		Code:           code,
		IncorrectValue: incorrectValue,
		Err:            err,
	}
}

func NewUnauthorizedError(err error) *CustomError {
	return &CustomError{
		Code: Unauthorized,
		Err:  err,
	}
}

func (e *CustomError) Error() string {
	if e.Code == SystemError {
		return fmt.Sprintf("%s[%s, %s]: %v", SystemError, e.Component, e.FailedOperation, e.Err)
	}
	if e.Code == Unauthorized {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}

	return fmt.Sprintf("%s[%s]: %v", e.Code, e.IncorrectValue, e.Err)
}

func (e *CustomError) HTTPCodeMsg() (int, interface{}) {
	var code int

	switch e.Code {
	case SystemError:
		return http.StatusInternalServerError, map[string]interface{}{
			"code":      SystemError.Name(),
			"component": e.Component,
			"operation": e.FailedOperation,
			"message":   e.Err.Error(),
		}
	case Unauthorized:
		return http.StatusUnauthorized, map[string]interface{}{
			"code":    Unauthorized.Name(),
			"message": e.Err.Error(),
		}
	case AlreadyExist:
		code = http.StatusConflict

	case DoesntExist:
		code = http.StatusNotFound

	case ConditionNotMet:
		code = http.StatusPreconditionFailed

	case InvalidValue:
		fallthrough

	default:
		code = http.StatusBadRequest
	}

	return code, map[string]interface{}{
		"code":           e.Code.Name(),
		"incorrectValue": e.IncorrectValue,
		"message":        e.Err.Error(),
	}
}
