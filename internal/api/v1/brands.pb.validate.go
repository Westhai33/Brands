// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: brands.proto

package v1

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

// Validate checks the field values on GetBrandRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetBrandRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetBrandRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetBrandRequestMultiError, or nil if none found.
func (m *GetBrandRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetBrandRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetBrandId() <= 0 {
		err := GetBrandRequestValidationError{
			field:  "BrandId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetBrandRequestMultiError(errors)
	}

	return nil
}

// GetBrandRequestMultiError is an error wrapping multiple validation errors
// returned by GetBrandRequest.ValidateAll() if the designated constraints
// aren't met.
type GetBrandRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetBrandRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetBrandRequestMultiError) AllErrors() []error { return m }

// GetBrandRequestValidationError is the validation error returned by
// GetBrandRequest.Validate if the designated constraints aren't met.
type GetBrandRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetBrandRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetBrandRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetBrandRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetBrandRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetBrandRequestValidationError) ErrorName() string { return "GetBrandRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetBrandRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetBrandRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetBrandRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetBrandRequestValidationError{}

// Validate checks the field values on GetBrandResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetBrandResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetBrandResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetBrandResponseMultiError, or nil if none found.
func (m *GetBrandResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetBrandResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for BrandId

	// no validation rules for Name

	// no validation rules for Link

	// no validation rules for Description

	// no validation rules for LogoUrl

	// no validation rules for CoverImageUrl

	// no validation rules for FoundedYear

	// no validation rules for OriginCountry

	// no validation rules for Popularity

	// no validation rules for IsPremium

	// no validation rules for IsUpcoming

	if len(errors) > 0 {
		return GetBrandResponseMultiError(errors)
	}

	return nil
}

// GetBrandResponseMultiError is an error wrapping multiple validation errors
// returned by GetBrandResponse.ValidateAll() if the designated constraints
// aren't met.
type GetBrandResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetBrandResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetBrandResponseMultiError) AllErrors() []error { return m }

// GetBrandResponseValidationError is the validation error returned by
// GetBrandResponse.Validate if the designated constraints aren't met.
type GetBrandResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetBrandResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetBrandResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetBrandResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetBrandResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetBrandResponseValidationError) ErrorName() string { return "GetBrandResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetBrandResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetBrandResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetBrandResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetBrandResponseValidationError{}

// Validate checks the field values on GetAllBrandsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAllBrandsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAllBrandsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAllBrandsRequestMultiError, or nil if none found.
func (m *GetAllBrandsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAllBrandsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Filter

	// no validation rules for Sort

	if len(errors) > 0 {
		return GetAllBrandsRequestMultiError(errors)
	}

	return nil
}

// GetAllBrandsRequestMultiError is an error wrapping multiple validation
// errors returned by GetAllBrandsRequest.ValidateAll() if the designated
// constraints aren't met.
type GetAllBrandsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAllBrandsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAllBrandsRequestMultiError) AllErrors() []error { return m }

// GetAllBrandsRequestValidationError is the validation error returned by
// GetAllBrandsRequest.Validate if the designated constraints aren't met.
type GetAllBrandsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAllBrandsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAllBrandsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAllBrandsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAllBrandsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAllBrandsRequestValidationError) ErrorName() string {
	return "GetAllBrandsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetAllBrandsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAllBrandsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAllBrandsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAllBrandsRequestValidationError{}

// Validate checks the field values on GetAllBrandsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAllBrandsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAllBrandsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAllBrandsResponseMultiError, or nil if none found.
func (m *GetAllBrandsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAllBrandsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetBrands() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetAllBrandsResponseValidationError{
						field:  fmt.Sprintf("Brands[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetAllBrandsResponseValidationError{
						field:  fmt.Sprintf("Brands[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetAllBrandsResponseValidationError{
					field:  fmt.Sprintf("Brands[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetAllBrandsResponseMultiError(errors)
	}

	return nil
}

// GetAllBrandsResponseMultiError is an error wrapping multiple validation
// errors returned by GetAllBrandsResponse.ValidateAll() if the designated
// constraints aren't met.
type GetAllBrandsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAllBrandsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAllBrandsResponseMultiError) AllErrors() []error { return m }

// GetAllBrandsResponseValidationError is the validation error returned by
// GetAllBrandsResponse.Validate if the designated constraints aren't met.
type GetAllBrandsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAllBrandsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAllBrandsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAllBrandsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAllBrandsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAllBrandsResponseValidationError) ErrorName() string {
	return "GetAllBrandsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetAllBrandsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAllBrandsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAllBrandsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAllBrandsResponseValidationError{}

// Validate checks the field values on Brand with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Brand) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Brand with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in BrandMultiError, or nil if none found.
func (m *Brand) ValidateAll() error {
	return m.validate(true)
}

func (m *Brand) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for BrandId

	// no validation rules for Name

	// no validation rules for Link

	// no validation rules for Description

	// no validation rules for LogoUrl

	// no validation rules for CoverImageUrl

	// no validation rules for FoundedYear

	// no validation rules for OriginCountry

	// no validation rules for Popularity

	// no validation rules for IsPremium

	// no validation rules for IsUpcoming

	if len(errors) > 0 {
		return BrandMultiError(errors)
	}

	return nil
}

// BrandMultiError is an error wrapping multiple validation errors returned by
// Brand.ValidateAll() if the designated constraints aren't met.
type BrandMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BrandMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BrandMultiError) AllErrors() []error { return m }

// BrandValidationError is the validation error returned by Brand.Validate if
// the designated constraints aren't met.
type BrandValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BrandValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BrandValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BrandValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BrandValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BrandValidationError) ErrorName() string { return "BrandValidationError" }

// Error satisfies the builtin error interface
func (e BrandValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBrand.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BrandValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BrandValidationError{}

// Validate checks the field values on GetModelRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetModelRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetModelRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetModelRequestMultiError, or nil if none found.
func (m *GetModelRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetModelRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ModelId

	if len(errors) > 0 {
		return GetModelRequestMultiError(errors)
	}

	return nil
}

// GetModelRequestMultiError is an error wrapping multiple validation errors
// returned by GetModelRequest.ValidateAll() if the designated constraints
// aren't met.
type GetModelRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetModelRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetModelRequestMultiError) AllErrors() []error { return m }

// GetModelRequestValidationError is the validation error returned by
// GetModelRequest.Validate if the designated constraints aren't met.
type GetModelRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetModelRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetModelRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetModelRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetModelRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetModelRequestValidationError) ErrorName() string { return "GetModelRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetModelRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetModelRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetModelRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetModelRequestValidationError{}

// Validate checks the field values on GetModelResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetModelResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetModelResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetModelResponseMultiError, or nil if none found.
func (m *GetModelResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetModelResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ModelId

	// no validation rules for BrandId

	// no validation rules for Name

	// no validation rules for ReleaseDate

	// no validation rules for IsUpcoming

	// no validation rules for IsLimited

	if len(errors) > 0 {
		return GetModelResponseMultiError(errors)
	}

	return nil
}

// GetModelResponseMultiError is an error wrapping multiple validation errors
// returned by GetModelResponse.ValidateAll() if the designated constraints
// aren't met.
type GetModelResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetModelResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetModelResponseMultiError) AllErrors() []error { return m }

// GetModelResponseValidationError is the validation error returned by
// GetModelResponse.Validate if the designated constraints aren't met.
type GetModelResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetModelResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetModelResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetModelResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetModelResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetModelResponseValidationError) ErrorName() string { return "GetModelResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetModelResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetModelResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetModelResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetModelResponseValidationError{}

// Validate checks the field values on GetAllModelsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAllModelsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAllModelsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAllModelsRequestMultiError, or nil if none found.
func (m *GetAllModelsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAllModelsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Filter

	// no validation rules for Sort

	if len(errors) > 0 {
		return GetAllModelsRequestMultiError(errors)
	}

	return nil
}

// GetAllModelsRequestMultiError is an error wrapping multiple validation
// errors returned by GetAllModelsRequest.ValidateAll() if the designated
// constraints aren't met.
type GetAllModelsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAllModelsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAllModelsRequestMultiError) AllErrors() []error { return m }

// GetAllModelsRequestValidationError is the validation error returned by
// GetAllModelsRequest.Validate if the designated constraints aren't met.
type GetAllModelsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAllModelsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAllModelsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAllModelsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAllModelsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAllModelsRequestValidationError) ErrorName() string {
	return "GetAllModelsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetAllModelsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAllModelsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAllModelsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAllModelsRequestValidationError{}

// Validate checks the field values on GetAllModelsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAllModelsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAllModelsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetAllModelsResponseMultiError, or nil if none found.
func (m *GetAllModelsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAllModelsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetModels() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetAllModelsResponseValidationError{
						field:  fmt.Sprintf("Models[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetAllModelsResponseValidationError{
						field:  fmt.Sprintf("Models[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetAllModelsResponseValidationError{
					field:  fmt.Sprintf("Models[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetAllModelsResponseMultiError(errors)
	}

	return nil
}

// GetAllModelsResponseMultiError is an error wrapping multiple validation
// errors returned by GetAllModelsResponse.ValidateAll() if the designated
// constraints aren't met.
type GetAllModelsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAllModelsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAllModelsResponseMultiError) AllErrors() []error { return m }

// GetAllModelsResponseValidationError is the validation error returned by
// GetAllModelsResponse.Validate if the designated constraints aren't met.
type GetAllModelsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAllModelsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAllModelsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAllModelsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAllModelsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAllModelsResponseValidationError) ErrorName() string {
	return "GetAllModelsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetAllModelsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAllModelsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAllModelsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAllModelsResponseValidationError{}

// Validate checks the field values on Model with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Model) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Model with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ModelMultiError, or nil if none found.
func (m *Model) ValidateAll() error {
	return m.validate(true)
}

func (m *Model) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ModelId

	// no validation rules for BrandId

	// no validation rules for Name

	// no validation rules for ReleaseDate

	// no validation rules for IsUpcoming

	// no validation rules for IsLimited

	if len(errors) > 0 {
		return ModelMultiError(errors)
	}

	return nil
}

// ModelMultiError is an error wrapping multiple validation errors returned by
// Model.ValidateAll() if the designated constraints aren't met.
type ModelMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ModelMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ModelMultiError) AllErrors() []error { return m }

// ModelValidationError is the validation error returned by Model.Validate if
// the designated constraints aren't met.
type ModelValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ModelValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ModelValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ModelValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ModelValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ModelValidationError) ErrorName() string { return "ModelValidationError" }

// Error satisfies the builtin error interface
func (e ModelValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sModel.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ModelValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ModelValidationError{}
