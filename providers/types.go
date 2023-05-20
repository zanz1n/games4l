package providers

type StatusCodeErr interface {
	error
	Status() int
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

func NewStatusCodeErr(msg string, code int) StatusCodeErr {
	return &statusCodeErr{
		msg:  msg,
		code: code,
	}
}
