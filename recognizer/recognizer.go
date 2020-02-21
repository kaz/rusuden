package recognizer

import (
	"context"
)

type (
	Recognizer interface {
		Recognize(context.Context, []byte) (string, error)
	}
)
