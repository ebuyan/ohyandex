package oauth

import (
	"time"

	"github.com/go-oauth2/oauth2/manage"
	"github.com/go-oauth2/oauth2/store"
)

type OAuthManager struct{ *manage.Manager }

func NewOAuthManager() OAuthManager {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewFileTokenStore("db"))
	manager.MapClientStorage(NewOAuthStore())
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{AccessTokenExp: time.Hour * 24 * 365, RefreshTokenExp: time.Hour * 24 * 370, IsGenerateRefresh: true})
	return OAuthManager{manager}
}
