package database

import (
	"errors"
	"fmt"
)

type ErrDataInvalid struct {
	Err string
}

func IsErrDataInvalid(err error) bool {
	return errors.As(err, new(ErrDataInvalid))
}

func (e ErrDataInvalid) Error() string {
	return fmt.Sprintf("invalid data to insert: %s", e.Err)
}

type ErrInternal struct {
	Err string
}

func (e ErrInternal) Error() string {
	return fmt.Sprintf("internal error while order insert: %s", e.Err)
}

type ErrOrderNotFound struct {
	Id string
}

func IsErrOrderNotFound(err error) bool {
	return errors.As(err, new(ErrOrderNotFound))
}

func (e ErrOrderNotFound) Error() string {
	return fmt.Sprintf("order with id: %s not found", e.Id)
}
