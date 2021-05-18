package server

import (
	"net/http"

	"ohyandex/internal/oauth"
	"ohyandex/pkg/logger"
)

type Handler struct {
	*oauth.OAuthServer
}

func NewHandler(oAuthServer *oauth.OAuthServer) Handler {
	return Handler{oAuthServer}
}

func (h Handler) Ping(w http.ResponseWriter, r *http.Request) {}

func (h Handler) BearerAuth(f func(w http.ResponseWriter, r *http.Request, userId string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := h.Auth(r)
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

func (h Handler) HandleError(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			logger.Error(r, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
