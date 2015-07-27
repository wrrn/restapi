package response

import (
	"encoding/json"
	"net/http"
)

// ServerError is just a convience function that allows us to write a
// status code of 500 and a message of "Server Error" to the response
func ServerError(w http.ResponseWriter) {
	http.Error(w, "Server Error", http.StatusInternalServerError)
}

// MethodNotAllowed is just a convience function that allows us to write a
// status code of 405 and a message of "" to the response
func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "", http.StatusMethodNotAllowed)
}

// WriteResponse will attempt to write the data to response. On failure it will
// write a status code of 500 and a message of "Server Error" to the response.
func Write(w http.ResponseWriter, code int, data []byte) {
	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		ServerError(w)
	}
}

// WriteJsonResponse will attempt to write the data to response in json format.
// On failure it will  write a status code of 500 and a message of "Server Error"
// to the response.
func WriteJson(w http.ResponseWriter, code int, data interface{}) {
	rawJson, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
	Write(w, code, rawJson)

}
