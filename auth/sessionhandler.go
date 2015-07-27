package auth

import "net/http"

type sessionsHandler struct {
	http.Handler
	Auth
}

func (s sessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := s.CheckSession(r)
	if err != nil {
		Forbidden(w)
		return
	}

	s.Handler.ServeHTTP(w, r)
}
