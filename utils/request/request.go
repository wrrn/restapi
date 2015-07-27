package request

import (
	"net/http"
	"strings"
)

func Is(r *http.Request, method string) bool {
	return strings.ToUpper(r.Method) == strings.ToUpper(method)
}
