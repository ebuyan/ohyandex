package identity

import (
	"encoding/base64"
	"errors"
	"net/http"
	"ohyandex/pkg/openhab"
	"ohyandex/pkg/session"
	"ohyandex/pkg/view"
)

type Provider struct {
	openhab.Client
}

func NewProvider() Provider {
	return Provider{openhab.NewClient()}
}

func (p Provider) RenderLogin(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "login", nil)
}

func (p Provider) Login(w http.ResponseWriter, r *http.Request) {
	session := session.GetSession(w, r)
	credentials, err := p.checkUsernameAndPassword(w, r)
	if err != nil {
		view.Render(w, "login", err)
		return
	}
	session.SetId(credentials)
	w.Header().Set("Location", "/auth")
	w.WriteHeader(http.StatusFound)
}

func (p Provider) Auth(w http.ResponseWriter, r *http.Request) {
	session := session.GetSession(w, r)

	if session.IsLoggedIn() {
		view.Render(w, "auth", nil)
		return
	}

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusFound)
}

func (p Provider) checkUsernameAndPassword(w http.ResponseWriter, r *http.Request) (credentials string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	credentials = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	ok, err := p.Ping(credentials)
	if ok {
		return
	}
	if err != nil {
		return
	}
	err = errors.New("Invalid username or password.")
	return
}
