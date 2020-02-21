package departure

import (
	"context"
)

type (
	Departure interface {
		Depart(context.Context, string) error
	}
)
