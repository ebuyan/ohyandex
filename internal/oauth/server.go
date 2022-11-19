package oauth

import (
	"github.com/ebuyan/ohyandex/pkg/session"
	"net/http"

	"github.com/go-oauth2/oauth2/server"
	"gopkg.in/oauth2.v3"
)

type Server struct{ *server.Server }

func NewOAuthServer() *Server {
	srv := server.NewServer(server.NewConfig(), NewOAuthManager())
	s := Server{srv}

	s.SetUserAuthorizationHandler(userAuthorizeHandler)
	s.SetClientInfoHandler(clientInfoHandler)

	return &s
}

func (s Server) Auth(r *http.Request) (ti oauth2.TokenInfo, err error) {
	ti, err = s.ValidationBearerToken(r)
	return
}

func (s Server) Authorize(w http.ResponseWriter, r *http.Request) (err error) {
	session := session.GetSession(w, r)
	if v, ok := session.GetQueryParams(); ok {
		r.Form = v
	}
	err = s.HandleAuthorizeRequest(w, r)
	return
}

func (s Server) Token(w http.ResponseWriter, r *http.Request) (err error) {
	err = s.HandleTokenRequest(w, r)
	return
}

func (s Server) DeleteToken(w http.ResponseWriter, r *http.Request) (err error) {
	token, err := s.Auth(r)
	if err != nil {
		return
	}
	s.Manager.RemoveAccessToken(token.GetAccess())
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
		err = r.ParseForm()
	}
	clientID = r.Form.Get("client_id")
	clientSecret = r.Form.Get("client_secret")
	return
}
