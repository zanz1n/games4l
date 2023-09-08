package errors

type StatusCodeErr interface {
	error
	Status() int
	SetMsg(s string)
}

type statusCodeErr struct {
	msg  string
	code int
}

func (se *statusCodeErr) Error() string {
	return se.msg
}

func (se *statusCodeErr) Status() int {
	return se.code
}

func (se *statusCodeErr) SetMsg(s string) {
	se.msg = s
}

func NewStatusCodeErr(msg string, code int) StatusCodeErr {
	return &statusCodeErr{
		msg:  msg,
		code: code,
	}
}