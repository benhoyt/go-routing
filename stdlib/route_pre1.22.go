//go:build !go1.22

package stdlib

import (
	"net/http"
)

var Serve http.Handler
