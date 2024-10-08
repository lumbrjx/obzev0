// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: common/proto/latency/latency.proto

package latency

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on TcpConfig with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *TcpConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TcpConfig with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TcpConfigMultiError, or nil
// if none found.
func (m *TcpConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *TcpConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetReqDelay() < 0 {
		err := TcpConfigValidationError{
			field:  "ReqDelay",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetResDelay() < 0 {
		err := TcpConfigValidationError{
			field:  "ResDelay",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetServer()) < 1 {
		err := TcpConfigValidationError{
			field:  "Server",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetClient()) < 1 {
		err := TcpConfigValidationError{
			field:  "Client",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return TcpConfigMultiError(errors)
	}

	return nil
}

// TcpConfigMultiError is an error wrapping multiple validation errors returned
// by TcpConfig.ValidateAll() if the designated constraints aren't met.
type TcpConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TcpConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TcpConfigMultiError) AllErrors() []error { return m }

// TcpConfigValidationError is the validation error returned by
// TcpConfig.Validate if the designated constraints aren't met.
type TcpConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TcpConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TcpConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TcpConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TcpConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TcpConfigValidationError) ErrorName() string { return "TcpConfigValidationError" }

// Error satisfies the builtin error interface
func (e TcpConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTcpConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TcpConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TcpConfigValidationError{}

// Validate checks the field values on RequestForTcp with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *RequestForTcp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RequestForTcp with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in RequestForTcpMultiError, or
// nil if none found.
func (m *RequestForTcp) ValidateAll() error {
	return m.validate(true)
}

func (m *RequestForTcp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetConfig() == nil {
		err := RequestForTcpValidationError{
			field:  "Config",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetConfig()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RequestForTcpValidationError{
					field:  "Config",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RequestForTcpValidationError{
					field:  "Config",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetConfig()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RequestForTcpValidationError{
				field:  "Config",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return RequestForTcpMultiError(errors)
	}

	return nil
}

// RequestForTcpMultiError is an error wrapping multiple validation errors
// returned by RequestForTcp.ValidateAll() if the designated constraints
// aren't met.
type RequestForTcpMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RequestForTcpMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RequestForTcpMultiError) AllErrors() []error { return m }

// RequestForTcpValidationError is the validation error returned by
// RequestForTcp.Validate if the designated constraints aren't met.
type RequestForTcpValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RequestForTcpValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RequestForTcpValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RequestForTcpValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RequestForTcpValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RequestForTcpValidationError) ErrorName() string { return "RequestForTcpValidationError" }

// Error satisfies the builtin error interface
func (e RequestForTcpValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRequestForTcp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RequestForTcpValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RequestForTcpValidationError{}

// Validate checks the field values on ResponseFromTcp with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ResponseFromTcp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ResponseFromTcp with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ResponseFromTcpMultiError, or nil if none found.
func (m *ResponseFromTcp) ValidateAll() error {
	return m.validate(true)
}

func (m *ResponseFromTcp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Message

	if len(errors) > 0 {
		return ResponseFromTcpMultiError(errors)
	}

	return nil
}

// ResponseFromTcpMultiError is an error wrapping multiple validation errors
// returned by ResponseFromTcp.ValidateAll() if the designated constraints
// aren't met.
type ResponseFromTcpMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ResponseFromTcpMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ResponseFromTcpMultiError) AllErrors() []error { return m }

// ResponseFromTcpValidationError is the validation error returned by
// ResponseFromTcp.Validate if the designated constraints aren't met.
type ResponseFromTcpValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ResponseFromTcpValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ResponseFromTcpValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ResponseFromTcpValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ResponseFromTcpValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ResponseFromTcpValidationError) ErrorName() string { return "ResponseFromTcpValidationError" }

// Error satisfies the builtin error interface
func (e ResponseFromTcpValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sResponseFromTcp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ResponseFromTcpValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ResponseFromTcpValidationError{}
