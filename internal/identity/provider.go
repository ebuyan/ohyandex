package identity

import (
	"encoding/base64"
	"errors"
	"github.com/ebuyan/ohyandex/pkg/openhab"
	"github.com/ebuyan/ohyandex/pkg/session"
	"github.com/ebuyan/ohyandex/pkg/view"
	"net/http"
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

func (p Provider) checkUsernameAndPassword(_ http.ResponseWriter, r *http.Request) (credentials string, err error) {
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
	err = errors.New("invalid username or password")
	return
}
