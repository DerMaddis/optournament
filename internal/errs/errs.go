package errs

type CustomError struct {
	Err        error
	StatusCode int
}

func (ce CustomError) Error() string {
	return ce.Err.Error()
}
