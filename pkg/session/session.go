package session

import (
	"net/http"
	"net/url"

	"github.com/go-session/session"
)

type Session struct {
	session.Store
}

func GetSession(w http.ResponseWriter, r *http.Request) *Session {
	store, _ := session.Start(r.Context(), w, r)
	return &Session{store}
}

func (s Session) GetQueryParams() (url.Values, bool) {
	if v, ok := s.Get("ReturnUri"); ok {
		return v.(url.Values), true
	}
	return nil, false
}

func (s Session) SetQueryParams(params url.Values) {
	s.Set("ReturnUri", params)
	s.Save()
}

func (s Session) GetId() (string, bool) {
	if v, ok := s.Get("SessionId"); ok {
		return v.(string), true
	}
	return "", false
}

func (s Session) SetId(id string) {
	s.Set("SessionId", id)
	s.Save()
}

func (s Session) IsLoggedIn() bool {
	_, ok := s.Get("SessionId")
	return ok
}
