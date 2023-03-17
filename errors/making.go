package errors

import (
	"context"

	"github.com/pkg/errors"
)

func ErrUnauthenticated(ctx context.Context, err error) IError {
	return newIError(ctx, Unauthorized, err)
}

func ErrPermissionDenied(ctx context.Context, err error) IError {
	return newIError(ctx, Forbidden, err)
}

func ErrInternal(ctx context.Context, err error) IError {
	return newIError(ctx, InternalServerError, err, "InternalError")
}

func ErrFailedPrecondition(ctx context.Context, err error) IError {
	return newIError(ctx, PreconditionFailed, err)
}

func ErrInvalidArgument(ctx context.Context, err error) IError {
	return newIError(ctx, BadRequest, err)
}

func newIError(ctx context.Context, code Code, err error, message ...string) IError {
	msg := err.Error()
	if len(message) > 0 {
		msg = message[0]
	}
	return ErrorTraceCtx(ctx, code, msg, err)
}

func Error(code Code, message string, errs ...error) *CError {
	return newError(false, code, message, errs...)
}

func ErrorTrace(code Code, message string, errs ...error) *CError {
	return newError(true, code, message, errs...)
}

func ErrorTraceCtx(ctx context.Context, code Code, message string, errs ...error) *CError {
	return newError(true, code, message, errs...)
}

func newError(trace bool, code Code, message string, errs ...error) *CError {
	if message == "" {
		message = code.String()
	}

	var err error
	if len(errs) > 0 {
		err = errs[0]
	}

	if cerr, ok := err.(*CError); ok && cerr != nil {
		if cerr.OriginalMessage == "" {
			cerr.OriginalMessage = cerr.Message
		}
		cerr.Code = code
		cerr.Message = message
		cerr.Trace = trace || cerr.Trace
		return cerr
	}

	oriMsg := ""
	if err != nil {
		oriMsg = err.Error()
	}

	return &CError{
		Code:            code,
		Err:             err,
		Message:         message,
		OriginalMessage: oriMsg,
		Stack:           errors.New("").(IStack).StackTrace(),
		Trace:           trace,
		Logs:            nil,
	}
}
