package god

import (
	"fmt"
	"reflect"
)

type InvalidOptionError struct {
	Key   OptionKey
	Value OptionValue
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
