package cache

import (
	"errors"
	"fmt"
)

type ErrCacheMiss struct {
	Key string
}

func (e ErrCacheMiss) Error() string {
	return fmt.Sprintf("key %s does not exist", e.Key)
}

func IsErrCacheMiss(err error) bool {
	return errors.As(err, new(ErrCacheMiss))
}
