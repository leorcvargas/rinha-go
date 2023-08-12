package people

import "errors"

var (
	ErrNicknameTaken  = errors.New("nickname taken")
	ErrPersonNotFound = errors.New("person not found")
)
