package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const fbVersion = "v2.12"

const (
	callbackURL = "http://localhost:8080/login/facebook/auth/callback"
)

var (
	fbClientID     string
	fbClientSecret string

	fbScope = []string{"email", "user_location", "user_friends", "user_birthday"}
)

const (
	fbAuthURL  = "https://www.facebook.com/" + fbVersion + "/dialog/oauth"
	fbTokenURL = "https://graph.facebook.com/" + fbVersion + "/oauth/access_token"
	fbMeURL    = "https://graph.facebook.com/" + fbVersion + "/me"
)

func init() {
	fbClientID = os.Getenv("TEST_CLIENT_ID")
	fbClientSecret = os.Getenv("TEST_SECRET")
}

func NewConfig() *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     fbClientID,
		ClientSecret: fbClientSecret,
		RedirectURL:  callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fbAuthURL,
			TokenURL: fbTokenURL,
		},
		Scopes: fbScope,
	}

	return c
}

func GetAuthCodeURL(state string) string {
	oc := NewConfig()
	url := oc.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return url
}

type Facebook interface {
	ExchangeCode(string) (*oauth2.Token, error)
	GetMe(*oauth2.Token, interface{}) error
}

type FacebookImpl struct {
}

func (f *FacebookImpl) ExchangeCode(code string) (*oauth2.Token, error) {
	oc := NewConfig()
	tok, err := oc.Exchange(context.Background(), code)
	if err != nil {
		err := errors.Wrap(err, "failed to exchange code.")
		return nil, err
	}
	if tok.Valid() == false {
		err = errors.New("invalid token.")
		return nil, err
	}
	return tok, nil
}

func (f *FacebookImpl) GetMe(tok *oauth2.Token, account interface{}) error {
	oc := NewConfig()
	client := oc.Client(context.Background(), tok)
	url := addAppSecretProofHMAC(fbMeURL, tok.AccessToken)

	resp, err := client.Get(url)
	if err != nil {
		err := errors.Wrap(err, "failed to send graph api request.")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		err = errors.New("facebook is unavailable")
		return err
	}

	if resp.StatusCode >= 400 {
		err = errors.New("facebook request is invalid")
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		err := errors.Wrap(err, "failed to decode json.")
		return err
	}

	return nil
}

func addAppSecretProofHMAC(url string, accessToken string) string {
	mac := hmac.New(sha256.New, []byte(fbClientSecret))
	mac.Write([]byte(accessToken))
	hash := hex.EncodeToString(mac.Sum(nil))

	url += "?appsecret_proof=" + hash
	return url
}

func validateFacebookCode(code string, state string, v interface{}) (err error) {
	if code == "" {
		err = errors.New("code should be set on query string")
		return
	}

	if state == "" {
		err = errors.New("state should be set on query string")
		return
	}

	if v == nil {
		err = errors.New("state hasn't be set")
		return
	}

	ss := v.(string)
	if state != ss {
		err = errors.New("state is invalid")
		return
	}

	return
}
