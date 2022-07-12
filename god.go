package god

import "context"

var New func(context.Context, Options) (Unit, error)
