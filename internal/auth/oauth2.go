package auth

import (
	"os"
	"strconv"

	"github.com/gorilla/sessions"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

func NewOAuth() error {
	//Get environment variables
	key := os.Getenv("GOTH_AUTH_SECRET")
	MaxAgeDays, err := strconv.Atoi(os.Getenv("OAUTH_AGE_DAYS"))
	if err != nil {
		return err
	}
	MaxAgeSeconds := MaxAgeDays * 3600 * 24
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	callbackURL := os.Getenv("GITHUB_CALLBACK_URL")
	IsProd := os.Getenv("IS_PROD") == "true"

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAgeSeconds)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store
	goth.UseProviders(
		github.New(clientId, clientSecret, callbackURL),
	)

	return nil

}
