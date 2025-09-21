package models

type ModelError string

const (
	ErrEmailAlreadyExists ModelError = "email already exists"
)

func (e ModelError) Error() string {
	return string(e)
}
