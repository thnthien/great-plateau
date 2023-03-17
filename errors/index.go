package errors

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

const gray, resetColor = "\x1b[90m", "\x1b[0m"

func ValueOf(f zapcore.Field) any {
	switch {
	case f.Integer != 0:
		return f.Integer
	case f.String != "":
		return f.String
	}
	return f.Interface
}

func marshalLogLineField(v interface{}) []byte {
	if xerr, ok := v.(*CError); ok {
		data, _ := xerr.MarshalJSON()
		return data
	}
	data, err := json.Marshal(v)
	if err != nil {
		data, _ = json.Marshal(err)
	}
	return data
}

type IError interface {
	error
	IStack

	Format(st fmt.State, verb rune)
	Log(msg string, fields ...zapcore.Field) IError
}

type IStack interface {
	StackTrace() errors.StackTrace
}
