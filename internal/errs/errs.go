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
	ErrAlreadyInTournament = CustomError{
		Err:        errors.New("already in tournament"),
		StatusCode: http.StatusBadRequest,
	}
	ErrAlreadyStarted = CustomError{
		Err:        errors.New("already started"),
		StatusCode: http.StatusBadRequest,
	}
	ErrBadVote = CustomError{
		Err:        errors.New("bad vote"),
		StatusCode: http.StatusBadRequest,
	}
	ErrAlreadyVoted = CustomError{
		Err:        errors.New("already voted"),
		StatusCode: http.StatusBadRequest,
	}
)

type errNoPermission struct {
	CustomError
	Start       CustomError
	ForceSumbit CustomError
}

func (enp errNoPermission) Error() string {
	return enp.Err.Error()
}

var ErrNoPermission = errNoPermission{
	CustomError: CustomError{
		Err:        errors.New("no permission"),
		StatusCode: http.StatusForbidden,
	},
	Start: CustomError{
		Err:        errors.New("no permission to start"),
		StatusCode: http.StatusForbidden,
	},
	ForceSumbit: CustomError{
		Err:        errors.New("no permission to force submit"),
		StatusCode: http.StatusForbidden,
	},
}

type errNotFound struct {
	CustomError
	Tournament CustomError
	Invite     CustomError
}

func (enf errNotFound) Error() string {
	return enf.Err.Error()
}

var ErrNotFound = errNotFound{
	CustomError: CustomError{
		Err:        errors.New("not found"),
		StatusCode: http.StatusNotFound,
	},
	Tournament: CustomError{
		Err:        errors.New("tournament not found"),
		StatusCode: http.StatusNotFound,
	},
	Invite: CustomError{
		Err:        errors.New("invite not found"),
		StatusCode: http.StatusNotFound,
	},
}
