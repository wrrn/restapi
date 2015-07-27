package request

import (
	"net/http"
	"strings"
)

// Is does a case insensitive check that the request's method matches
// the method argument
func Is(r *http.Request, method string) bool {
	return strings.ToUpper(r.Method) == strings.ToUpper(method)
}

// GetUrlVariables returns a list variable from a path delimited by the forward
// slash ("/") character
func GetURLVariables(path string) []string {
	path = strings.Trim(path, "/")
	variables := strings.Split(path, "/")
	return variables
}
