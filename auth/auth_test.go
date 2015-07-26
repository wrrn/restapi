package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func SetupDB() *sql.DB {
	db, err := sql.Open("postgres", "user=tenable password=insecure dbname=apitest")
	if err != nil {
		log.Fatal(err)
	}
	ResetDB(db)
	return db
}

func ResetDB(db *sql.DB) {
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM configuration")
	db.Exec("DELETE FROM sessions")
}

var auth = Auth{SetupDB()}

type Failure struct {
	Prefix   string
	Expected interface{}
	Actual   interface{}
}

func (f Failure) Error() string {
	str := f.Prefix
	if f.Expected != nil {
		str += fmt.Sprintf("\n Expected: %v", f.Expected)

	}

	if f.Actual != nil {
		str += fmt.Sprintf("\n Actual: %v", f.Actual)
	}

	return str

}

var handlerTests = map[string]struct {
	Setup func() error
	http.HandlerFunc
	*http.Request
	pass func(r *httptest.ResponseRecorder) error
}{

	"TestValidLogin": {
		HandlerFunc: auth.HandleLogin,
		Request:     generateLoginRequest(User{0, "john", "1234abc"}),
		pass: func(r *httptest.ResponseRecorder) error {
			if r.Code != http.StatusOK {
				return Failure{"", http.StatusOK, r.Code}
			}
			if r.Header().Get("Set-Cookie") == "" {
				return Failure{"Cookie Not Set", nil, nil}
			}
			body := r.Body.String()
			if body != "Authorized" {
				return Failure{"Body did not match", "Authorized", body}
			}

			return nil
		},
	},

	"TestInvalidLogin": {
		HandlerFunc: auth.HandleLogin,
		Request:     generateLoginRequest(User{0, "john", "1234a"}),
		pass: func(r *httptest.ResponseRecorder) error {
			if r.Code != http.StatusUnauthorized {
				return Failure{"", http.StatusUnauthorized, r.Code}
			}
			if r.Header().Get("Set-Cookie") != "" {
				return Failure{"Cookie Set", nil, nil}
			}
			body := strings.TrimSpace(r.Body.String())
			if body != "Unauthorized" {
				return Failure{"Body did not match", "Unauthorized", body}
			}

			return nil
		},
	},
}

func TestHandler(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	for testName, test := range handlerTests {
		auth.RegisterUser(User{0, "john", "1234abc"})
		r := httptest.NewRecorder()
		test.HandlerFunc(r, test.Request)
		if err := test.pass(r); err != nil {
			t.Error("Failed:", testName, err.Error())
		}
		ResetDB(auth.DB)
	}
}

func generateLoginRequest(user User) *http.Request {
	body, err := json.Marshal(user)
	if err != nil {
		return nil
	}
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	return req
}

func NewRequest(method, url string, reader io.Reader) (req *http.Request) {
	req, _ = http.NewRequest(method, url, reader)
	return req
}
