package utilities

type Error struct {
	StatusCode int
	Err        error
}

func (r *Error) Error() string {
	return r.Err.Error()
}

func ErrorRequest(err error, code int) error {
	return &Error{
		StatusCode: code,
		Err:        err,
	}
}

func ParseError(r error) *Error {
	errInfo, ok := r.(*Error)
	if ok {
		return errInfo
	}
	return nil
}
