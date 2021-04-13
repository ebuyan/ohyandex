package server

import (
	"net/http"
	"os"

	"ohyandex/identity"
	"ohyandex/oauth"
	"ohyandex/provider"

	"github.com/gorilla/mux"
)

type Http struct {
	host             string
	handler          Handler
	oAuthServer      *oauth.OAuthServer
	serviceProvider  provider.ServiceProvider
	identityProvider identity.IdentityProvider
}

func NewHttp() Http {
	oauthServer := oauth.NewOAuthServer()
	return Http{
		host:             os.Getenv("HTTP_HOST"),
		handler:          NewHandler(oauthServer),
		oAuthServer:      oauthServer,
		serviceProvider:  provider.NewServiceProvider(),
		identityProvider: identity.NewIdentityProvider(),
	}
}

func (h *Http) Start() {
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	r := mux.NewRouter()

	r.HandleFunc("/login", h.identityProvider.RenderLogin).Methods(http.MethodGet)
	r.HandleFunc("/login", h.identityProvider.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth", h.identityProvider.Auth).Methods(http.MethodGet)

	r.HandleFunc("/oauth/authorize", h.handler.HandleError(h.oAuthServer.Authorize)).Methods(http.MethodGet)
	r.HandleFunc("/oauth/token", h.handler.HandleError(h.oAuthServer.Token)).Methods(http.MethodPost)
	r.HandleFunc("/oauth/refresh", h.handler.HandleError(h.oAuthServer.Token)).Methods(http.MethodPost)

	r.HandleFunc("/provider", h.handler.Ping).Methods(http.MethodGet)
	r.HandleFunc("/provider/v1.0", h.handler.Ping).Methods(http.MethodGet)

	r.HandleFunc("/provider/v1.0/user/unlink", h.handler.HandleError(h.oAuthServer.DeleteToken)).Methods(http.MethodPost)

	r.HandleFunc("/provider/v1.0/user/devices", h.handler.BearerAuth(h.serviceProvider.Devices)).Methods(http.MethodGet)
	r.HandleFunc("/provider/v1.0/user/devices/query", h.handler.BearerAuth(h.serviceProvider.DevicesState)).Methods(http.MethodPost)
	r.HandleFunc("/provider/v1.0/user/devices/action", h.handler.BearerAuth(h.serviceProvider.ControlDevices)).Methods(http.MethodPost)

	http.Handle("/", r)
	http.ListenAndServe(h.host, nil)
}
