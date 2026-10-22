package errs

type WrappedError struct {
	err error
	msg string
}

func NewWrappedError(m string) *WrappedError {
	return &WrappedError{
		msg: m,
	}
}

func (e *WrappedError) Wrap(r error) error {
	if e == nil {
		return &WrappedError{err: r}
	}

	if e.err == nil {
		return &WrappedError{err: r, msg: e.msg}
	}

	return &WrappedError{err: e, msg: e.msg}
}

func (e *WrappedError) Unwrap() error {
	return e.err
}

func (e *WrappedError) Error() string {
	switch {
	case e.msg != "" && e.err != nil:
		return e.msg + ": " + e.err.Error()
	case e.msg != "":
		return e.msg
	case e.err != nil:
		return e.err.Error()
	default:
		return "unknown error"
	}
}
