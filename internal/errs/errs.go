package errs

import (
	"errors"
	"net/http"
)

type CustomError struct {
	Err        error
	StatusCode int
}

func (ce CustomError) Error() string {
	return ce.Err.Error()
}

var (
	ErrBadSongCount = CustomError{
		Err:        errors.New("bad song count"),
		StatusCode: http.StatusBadRequest,
	}
	ErrBadScore = CustomError{
		Err:        errors.New("bad score"),
		StatusCode: http.StatusBadRequest,
	}
)
