package oauth

import "github.com/go-oauth2/oauth2/store"

type OAuthStore struct{ *store.ClientStore }

func NewOAuthStore() OAuthStore {
	config := NewOAuthConfig()
	store := store.NewClientStore()
	store.Set(config.ID, &config)
	return OAuthStore{store}
}
