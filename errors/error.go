package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type ErrorWithStack struct {
	Stack errors.StackTrace
	Err   error
}

// Error returns error message of ErrorWithStack
func (e *ErrorWithStack) Error() string {
	return e.Err.Error()
}

// StackTrace implements IStack to be compatible with sentry
func (e *ErrorWithStack) StackTrace() errors.StackTrace {
	return e.Stack
}

type LogLine struct {
	Level   string
	File    string
	Line    int
	Message string
	Fields  []zapcore.Field
}

// MarshalJSON returns log current with json formatted
func (l LogLine) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 512)
	return l.MarshalTo(b), nil
}

// MarshalTo changes log current to json and assign to input byte array
func (l LogLine) MarshalTo(b []byte) []byte {
	b = append(b, '{')

	b = append(b, marshalLogLineField("@"+l.Level)...)
	b = append(b, ':')
	b = append(b, marshalLogLineField(l.Message)...)
	b = append(b, ',')

	b = append(b, `"@file":`...)
	b = append(b, marshalLogLineField(l.File+":"+strconv.Itoa(l.Line))...)

	for _, field := range l.Fields {
		b = append(b, ',')
		b = append(b, marshalLogLineField(field.Key)...)
		b = append(b, ':')

		if field.Integer != 0 {
			b = append(b, strconv.Itoa(int(field.Integer))...)
		} else if field.String != "" {
			b = append(b, marshalLogLineField(field.String)...)
		} else {
			b = append(b, marshalLogLineField(field.Interface)...)
		}
	}
	b = append(b, '}')
	return b
}

type CError struct {
	Code            Code
	Err             error
	Message         string
	OriginalMessage string
	Trace           bool
	Stack           errors.StackTrace
	Logs            []LogLine
}

func (e *CError) Log(msg string, fields ...zapcore.Field) IError {
	_, file, line, _ := runtime.Caller(1)
	e.Logs = append(e.Logs, LogLine{
		Level:   "error",
		File:    file,
		Line:    line,
		Fields:  fields,
		Message: msg,
	})
	return e
}

func (e *CError) Error() string {
	return e.Message
}

func (e *CError) Cause() error {
	err := e.Err
	if err == nil {
		err = e
	}
	return &ErrorWithStack{
		Err:   err,
		Stack: e.Stack,
	}
}

// StackTrace returns Stack of APIError
func (e *CError) StackTrace() errors.StackTrace {
	return e.Stack
}

// Format parse APIError to string with suitable format
// then write it to provided writer
func (e *CError) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('#') || st.Flag('+'):
			_, _ = fmt.Fprintf(st, "\ncode=%v message=%v", e.Code, e.Message)
			if e.OriginalMessage != "" {
				_, _ = fmt.Fprintf(st, " original=%s", e.OriginalMessage)
			}
			if e.Err != nil {
				_, _ = fmt.Fprintf(st, " cause=%+v", e.Err)
			}
			for _, log := range e.Logs {
				_, _ = fmt.Fprint(st, "\n\t", log.Line, " ", gray, log.File, ":", strconv.Itoa(log.Line), resetColor)
				for k, v := range log.Fields {
					_, _ = fmt.Fprint(st, " ", k, "=", ValueOf(v))
				}
			}
			fallthrough
		case st.Flag('+'):
			_, _ = fmt.Fprintf(st, "%+v", e.StackTrace())
		default:
			_, _ = io.WriteString(st, e.Error())
		}
	case 's':
		_, _ = io.WriteString(st, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(st, "%q", e.Error())
	}
}

// MarshalJSON jsonize APIError to bytes
func (e *CError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}

	b := make([]byte, 0, 2048)

	b = append(b, '{')
	b = append(b, `"code":`...)
	b = append(b, strconv.FormatInt(int64(e.Code), 10)...)

	if e.Err != nil {
		b = append(b, ',')
		b = append(b, `"err":`...)
		b = append(b, marshal(e.Err.Error())...)
	}

	b = append(b, ',')
	b = append(b, `"msg":`...)
	b = append(b, marshal(e.Message)...)

	if e.OriginalMessage != "" {
		b = append(b, ',')
		b = append(b, `"orig":`...)
		b = append(b, marshal(e.OriginalMessage)...)
	}

	b = append(b, ',')
	b = append(b, `"logs":`...)
	b = append(b, '[')
	for i, line := range e.Logs {
		if i > 0 {
			b = append(b, ',')
		}
		b = line.MarshalTo(b)
	}
	b = append(b, ']')

	if e.Trace {
		b = append(b, ',')
		b = append(b, `"stack":`...)
		b = append(b, marshal(fmt.Sprintf("%+v", e.Stack))...)
	}

	b = append(b, '}')
	return b, nil
}

func marshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		data, _ = json.Marshal(err)
	}
	return data
}
