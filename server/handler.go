package server

import (
	"net/http"

	"ohyandex/logger"
	"ohyandex/oauth"
)

type Handler struct {
	oAuthServer *oauth.OAuthServer
}

func NewHandler(oAuthServer *oauth.OAuthServer) Handler {
	return Handler{oAuthServer: oAuthServer}
}

func (h Handler) bearerAuth(f func(w http.ResponseWriter, r *http.Request, userId string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := h.oAuthServer.Auth(r)
		if err != nil {
			logger.Error(r, err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		err = f(w, r, token.GetUserID())
		if err != nil {
			logger.Error(r, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h Handler) handleError(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			logger.Error(r, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
