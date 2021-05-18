package server

import (
	"net/http"
	"os"

	"ohyandex/internal/device"
	"ohyandex/internal/identity"
	"ohyandex/internal/oauth"

	"github.com/gorilla/mux"
)

type Http struct {
	host             string
	handler          Handler
	oAuthServer      *oauth.OAuthServer
	deviceProvider   device.Provider
	identityProvider identity.Provider
}

func NewHttp() Http {
	return Http{
		host:             os.Getenv("HTTP_HOST"),
		handler:          NewHandler(oauth.NewOAuthServer()),
		deviceProvider:   device.NewProvider(),
		identityProvider: identity.NewProvider(),
	}
}

func (h *Http) Start() {
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	r := mux.NewRouter()

	r.HandleFunc("/login", h.identityProvider.RenderLogin).Methods(http.MethodGet)
	r.HandleFunc("/login", h.identityProvider.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth", h.identityProvider.Auth).Methods(http.MethodGet)

	r.HandleFunc("/oauth/authorize", h.handler.HandleError(h.handler.Authorize)).Methods(http.MethodGet)
	r.HandleFunc("/oauth/token", h.handler.HandleError(h.handler.Token)).Methods(http.MethodPost)
	r.HandleFunc("/oauth/refresh", h.handler.HandleError(h.handler.Token)).Methods(http.MethodPost)

	r.HandleFunc("/provider", h.handler.Ping).Methods(http.MethodGet)
	r.HandleFunc("/provider/v1.0", h.handler.Ping).Methods(http.MethodGet)

	r.HandleFunc("/provider/v1.0/user/unlink", h.handler.HandleError(h.handler.DeleteToken)).Methods(http.MethodPost)

	r.HandleFunc("/provider/v1.0/user/devices", h.handler.BearerAuth(h.deviceProvider.Devices)).Methods(http.MethodGet)
	r.HandleFunc("/provider/v1.0/user/devices/query", h.handler.BearerAuth(h.deviceProvider.DevicesState)).Methods(http.MethodPost)
	r.HandleFunc("/provider/v1.0/user/devices/action", h.handler.BearerAuth(h.deviceProvider.ControlDevices)).Methods(http.MethodPost)

	http.Handle("/", r)
	http.ListenAndServe(h.host, nil)
}
