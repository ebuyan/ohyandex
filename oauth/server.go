package oauth

import (
	"net/http"
	"ohyandex/session"

	"github.com/go-oauth2/oauth2/server"
	"gopkg.in/oauth2.v3"
)

type OAuthServer struct{ server *server.Server }

func NewOAuthServer() *OAuthServer {
	manager := NewOAuthManager()
	srv := server.NewServer(server.NewConfig(), manager.Manager)
	s := OAuthServer{server: srv}

	s.server.SetUserAuthorizationHandler(userAuthorizeHandler)
	s.server.SetClientInfoHandler(clientInfoHandler)

	return &s
}

func (s OAuthServer) Auth(r *http.Request) (ti oauth2.TokenInfo, err error) {
	ti, err = s.server.ValidationBearerToken(r)
	return
}

func (s OAuthServer) Authorize(w http.ResponseWriter, r *http.Request) (err error) {
	session := session.GetSession(w, r)
	if v, ok := session.GetQueryParams(); ok {
		r.Form = v
	}
	err = s.server.HandleAuthorizeRequest(w, r)
	return
}

func (s OAuthServer) Token(w http.ResponseWriter, r *http.Request) (err error) {
	err = s.server.HandleTokenRequest(w, r)
	return
}

func (s OAuthServer) DeleteToken(w http.ResponseWriter, r *http.Request) (err error) {
	token, err := s.Auth(r)
	if err != nil {
		return
	}
	s.server.Manager.RemoveAccessToken(token.GetAccess())
	session := session.GetSession(w, r)
	session.Flush()
	return
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (uid string, err error) {
	session := session.GetSession(w, r)
	if session.IsLoggedIn() {
		uid, _ := session.GetId()
		return uid, nil
	}
	if r.Form == nil {
		r.ParseForm()
	}
	session.SetQueryParams(r.Form)
	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusFound)
	return
}

func clientInfoHandler(r *http.Request) (clientID, clientSecret string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	clientID = r.Form.Get("client_id")
	clientSecret = r.Form.Get("client_secret")
	return
}
