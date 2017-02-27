package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyToken(t *testing.T) {
	token := "SECRET"

	tests := map[string]struct {
		validator  TokenValidator
		statusCode int
	}{
		"InvalidToken": {
			validator: ValidatorFunc(func(t string) bool {
				if t == token {
					return false
				}
				return true
			}),
			statusCode: http.StatusUnauthorized,
		},

		"ValidToken": {
			validator: ValidatorFunc(func(t string) bool {
				if t == token {
					return true
				}
				return false

			}),
			statusCode: http.StatusOK,
		},
	}

	for testName, test := range tests {
		req := httptest.NewRequest("POST", "/", nil)
		req.SetBasicAuth(token, "")
		resp := okOnSuccess(test.validator, req)
		if resp.StatusCode != test.statusCode {
			t.Errorf("%s failed: Status Codes did not match.\n  Expected: %d, Got: %d", testName, test.statusCode, resp.StatusCode)
		}

	}

}

func VerifyTokenEmptyHeader(t *testing.T) {

	tests := map[string]struct {
		validator  TokenValidator
		statusCode int
	}{
		"InvalidToken": {
			validator: ValidatorFunc(func(string) bool {
				return false
			}),
			statusCode: http.StatusUnauthorized,
		},

		"ValidToken": {
			validator: ValidatorFunc(func(token string) bool {
				return true
			}),
			statusCode: http.StatusOK,
		},
	}

	for testName, test := range tests {
		resp := okOnSuccess(test.validator, httptest.NewRequest("POST", "/", nil))
		if resp.StatusCode != test.statusCode {
			t.Errorf("%s failed: Status Codes did not match.\n  Expected: %d, Got: %d", testName, test.statusCode, resp.StatusCode)
		}

	}

}

func okOnSuccess(validator TokenValidator, req *http.Request) *http.Response {
	var (
		basicHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}

		handler = VerifyTokens(validator, basicHandler)
		w       = httptest.NewRecorder()
	)

	handler.ServeHTTP(w, req)
	return w.Result()
}
