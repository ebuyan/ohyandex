package identity

import (
	"encoding/base64"
	"errors"
	"net/http"
	"ohyandex/provider"
	"ohyandex/session"
)

type IdentityProvider struct {
	view          View
	openhabClient provider.OpenhabClient
}

func NewIdentityProvider() IdentityProvider {
	return IdentityProvider{view: View{}, openhabClient: provider.NewOpenhabClient()}
}

func (p IdentityProvider) RenderLogin(w http.ResponseWriter, r *http.Request) {
	p.view.Render(w, "login", nil)
}

func (p IdentityProvider) Login(w http.ResponseWriter, r *http.Request) {
	session := session.GetSession(w, r)
	credentials, err := p.checkUsernameAndPassword(w, r)
	if err != nil {
		p.view.Render(w, "login", err)
		return
	}
	session.SetId(credentials)
	w.Header().Set("Location", "/auth")
	w.WriteHeader(http.StatusFound)
}

func (p IdentityProvider) Auth(w http.ResponseWriter, r *http.Request) {
	session := session.GetSession(w, r)

	if session.IsLoggedIn() {
		p.view.Render(w, "auth", nil)
		return
	}

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusFound)
}

func (p IdentityProvider) checkUsernameAndPassword(w http.ResponseWriter, r *http.Request) (credentials string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	credentials = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	ok, err := p.openhabClient.Ping(credentials)
	if ok {
		return
	}
	if err != nil {
		return
	}
	err = errors.New("Invalid username or password.")
	return
}
