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

func (s Session) Flush() {
	s.Store.Flush()
}

func (s Session) GetQueryParams() (url.Values, bool) {
	if v, ok := s.Store.Get("ReturnUri"); ok {
		return v.(url.Values), true
	}
	return nil, false
}

func (s Session) SetQueryParams(params url.Values) {
	s.Store.Set("ReturnUri", params)
	s.Store.Save()
}

func (s Session) GetId() (string, bool) {
	if v, ok := s.Store.Get("SessionId"); ok {
		return v.(string), true
	}
	return "", false
}

func (s Session) SetId(id string) {
	s.Store.Set("SessionId", id)
	s.Store.Save()
}

func (s Session) IsLoggedIn() bool {
	_, ok := s.Store.Get("SessionId")
	return ok
}
