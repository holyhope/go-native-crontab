package god

import (
	"fmt"
	"reflect"
)

type InvalidOptionError struct {
	Key   string
	Value interface{}
}

func (err *InvalidOptionError) Error() string {
	return fmt.Sprintf("invalid option %s: %v", err.Key, err.Value)
}

func (err *InvalidOptionError) Is(err2 error) bool {
	if err2, ok := err2.(*InvalidOptionError); ok {
		return err.Key == err2.Key && reflect.DeepEqual(err.Value, err2.Value)
	}

	return false
}

type MissingOptionError InvalidOptionError

func NewMissingOptionError(key string) error {
	return &MissingOptionError{
		Key: key,
	}
}

func (err *MissingOptionError) Error() string {
	return fmt.Sprintf("missing option %s", err.Key)
}
