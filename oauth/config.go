package oauth

import (
	"os"

	"github.com/go-oauth2/oauth2/models"
)

type OAuthConfig struct{ models.Client }

func NewOAuthConfig() OAuthConfig {
	return OAuthConfig{
		models.Client{
			ID:     os.Getenv("CLIENT_ID"),
			Secret: os.Getenv("CLIENT_SECRET"),
		},
	}
}
