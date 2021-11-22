package errors

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/tkeel-io/kit/log"
)

type Error interface {
	error
	WithMetadata(map[string]string) Error
	WithMessage(string) Error
}

const (
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

var _ Error = &TError{}

type errKey string

var errs = map[errKey]*TError{}

// Register 注册错误信息
func Register(egoError *TError) {
	errs[errKey(egoError.Reason)] = egoError
}

// Error Error信息
func (x *TError) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", x.Code, x.Reason, x.Message, x.Metadata)
}

// Is 判断是否为根因错误
func (x *TError) Is(err error) bool {
	egoErr, flag := err.(*TError)
	if !flag {
		return false
	}
	return x.Reason == egoErr.Reason
}

// GRPCStatus returns the Status represented by se.
func (x *TError) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(x.Code), x.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   x.Reason,
			Metadata: x.Metadata,
		})
	return s
}

// WithMetadata with an MD formed by the mapping of key, value.
func (x *TError) WithMetadata(md map[string]string) Error {
	err := proto.Clone(x).(*TError)
	err.Metadata = md
	return err
}

// WithMessage set message to current TError
func (x *TError) WithMessage(msg string) Error {
	err := proto.Clone(x).(*TError)
	err.Message = msg
	return err
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *TError {
	return &TError{
		Code:    int32(code),
		Message: message,
		Reason:  reason,
	}
}

// ToHTTPStatusCode Get equivalent HTTP status code from x.Code
func (x *TError) ToHTTPStatusCode() int {
	return GRPCToHTTPStatusCode(codes.Code(x.Code))
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *TError {
	if err == nil {
		return nil
	}
	if se := new(TError); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				e, ok := errs[errKey(d.Reason)]
				if ok {
					return e.WithMessage(gs.Message()).WithMetadata(d.Metadata).(*TError)
				}
				return New(
					int(gs.Code()),
					d.Reason,
					gs.Message(),
				).WithMetadata(d.Metadata).(*TError)
			}
		}
	}
	return New(int(codes.Unknown), UnknownReason, err.Error())
}

// PrintErrLog ...
func PrintErrLog(msg string, err error) {
	switch e := err.(type) {
	case *TError:
		log.Error(e.GetMessage(), zap.Any("meta", e.GetMetadata()))
	default:
		log.Error(msg, zap.Error(e))
	}
}
