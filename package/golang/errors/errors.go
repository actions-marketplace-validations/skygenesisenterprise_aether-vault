package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrCodeInvalidRequest   ErrorCode = "INVALID_REQUEST"
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeConflict         ErrorCode = "CONFLICT"
	ErrCodeInternal         ErrorCode = "INTERNAL_ERROR"
	ErrCodeUnavailable      ErrorCode = "UNAVAILABLE"
	ErrCodeTimeout          ErrorCode = "TIMEOUT"
	ErrCodeRateLimited      ErrorCode = "RATE_LIMITED"
	ErrCodeInvalidToken     ErrorCode = "INVALID_TOKEN"
	ErrCodeExpiredToken     ErrorCode = "EXPIRED_TOKEN"
	ErrCodeSecretNotFound   ErrorCode = "SECRET_NOT_FOUND"
	ErrCodeSecretExists     ErrorCode = "SECRET_EXISTS"
	ErrCodePolicyViolation  ErrorCode = "POLICY_VIOLATION"
	ErrCodeTOTPFailed       ErrorCode = "TOTP_FAILED"
	ErrCodeIdentityNotFound ErrorCode = "IDENTITY_NOT_FOUND"
	ErrCodeAuditFailed      ErrorCode = "AUDIT_FAILED"
)

type VaultError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Status  int       `json:"-"`
}

func (e *VaultError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *VaultError) Is(target error) bool {
	if t, ok := target.(*VaultError); ok {
		return e.Code == t.Code
	}
	return false
}

func NewError(code ErrorCode, message string) *VaultError {
	return &VaultError{
		Code:    code,
		Message: message,
		Status:  statusCodeForError(code),
	}
}

func NewErrorWithDetails(code ErrorCode, message, details string) *VaultError {
	return &VaultError{
		Code:    code,
		Message: message,
		Details: details,
		Status:  statusCodeForError(code),
	}
}

func WrapError(err error, code ErrorCode, message string) *VaultError {
	if err == nil {
		return NewError(code, message)
	}
	return &VaultError{
		Code:    code,
		Message: message,
		Details: err.Error(),
		Status:  statusCodeForError(code),
	}
}

func statusCodeForError(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidRequest, ErrCodeInvalidToken, ErrCodeSecretExists:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeExpiredToken:
		return http.StatusUnauthorized
	case ErrCodeForbidden, ErrCodePolicyViolation:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeSecretNotFound, ErrCodeIdentityNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeRateLimited:
		return http.StatusTooManyRequests
	case ErrCodeUnavailable:
		return http.StatusServiceUnavailable
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	case ErrCodeInternal, ErrCodeTOTPFailed, ErrCodeAuditFailed:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func ErrorCodeFromStatus(status int) ErrorCode {
	switch status {
	case http.StatusBadRequest:
		return ErrCodeInvalidRequest
	case http.StatusUnauthorized:
		return ErrCodeUnauthorized
	case http.StatusForbidden:
		return ErrCodeForbidden
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusConflict:
		return ErrCodeConflict
	case http.StatusTooManyRequests:
		return ErrCodeRateLimited
	case http.StatusServiceUnavailable:
		return ErrCodeUnavailable
	case http.StatusRequestTimeout:
		return ErrCodeTimeout
	case http.StatusInternalServerError:
		return ErrCodeInternal
	default:
		return ErrCodeInternal
	}
}

func IsVaultError(err error) bool {
	_, ok := err.(*VaultError)
	return ok
}

func GetErrorCode(err error) ErrorCode {
	if vaultErr, ok := err.(*VaultError); ok {
		return vaultErr.Code
	}
	return ErrCodeInternal
}
