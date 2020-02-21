package arrival

import (
	"net/http"
)

type (
	Arrival interface {
		Arrive(*http.Request) ([]byte, error)
	}
)
