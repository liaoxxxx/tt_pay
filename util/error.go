package util

import (
	"fmt"
)

const (
	ErrorPattern = `{"code": "%s", "msg": "%s", "sub_code": "%s", "sub_msg": "%s", "detail": "%s"}`
)

// Error为请求失败错误
// 当出现此Error时，意味着网络连接建立成功，但请求失败
// 典型的情况是请求参数设置错误
type Error struct {
	Code    string
	Msg     string
	SubCode string
	SubMsg  string
	Detail  string
}

func (e *Error) Error() string {
	if e == nil {
		return "e is nil"
	}
	return fmt.Sprintf(ErrorPattern, e.Code, e.Msg, e.SubCode, e.SubMsg, e.Detail)
}

// Wrap用来包装error，以提供trace信息帮助使用者debug
// 该方法借鉴了github.com/pkg/errors 包
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	e := &WithMessage{
		cause: err,
		msg:   message,
	}
	return e
}

type WithMessage struct {
	cause error
	msg   string
}

func (w *WithMessage) Error() string { return w.msg + ": " + w.cause.Error() }

// Cause函数用来提取原生error
func (w *WithMessage) Cause() error { return w.cause }
