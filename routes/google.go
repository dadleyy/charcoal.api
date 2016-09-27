package routes

import "os"
import "fmt"
import "errors"
import "net/url"
import "golang.org/x/oauth2"
import "golang.org/x/oauth2/google"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/context"

const ERR_BAD_RUNTIME = "BAD_RUNTIME"
const ERR_BAD_AUTH_CODE = "BAD_AUTH_CODE"
const GOOGLE_AUTH_ENDPOINT = "https://accounts.google.com/o/oauth2/v2/auth"
const GOOGLE_TOKEN_ENDPOINT = "https://www.googleapis.com/oauth2/v4/token"

func GoogleOauthRedirect(ectx echo.Context) error {
	clientid := os.Getenv("GOOGLE_CLIENT_ID")
	redir := os.Getenv("GOOGLE_REDIRECT_URL")
	fin, err := url.Parse(GOOGLE_AUTH_ENDPOINT)

	if err != nil {
		return err
	}

	queries := make(url.Values)

	queries.Set("response_type", "code")
	queries.Set("redirect_uri", redir)
	queries.Set("client_id", clientid)
	queries.Set("scope", "https://www.googleapis.com/auth/plus.login email")
	queries.Set("access_type", "offline")

	fin.RawQuery = queries.Encode()
	ectx.Logger().Infof("generating full url: %s", fmt.Sprintf("%s", fin.String()))

	ectx.Redirect(301, fin.String())

	return nil
}

func GoogleOauthReceiveCode(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Miritos)

	if ok != true {
		return errors.New(ERR_BAD_RUNTIME)
	}

	code := runtime.QueryParam("code")

	if len(code) == 0 {
		return runtime.ErrorOut(errors.New(ERR_BAD_AUTH_CODE))
	}

	redir := os.Getenv("GOOGLE_REDIRECT_URL")
	clientid := os.Getenv("GOOGLE_CLIENT_ID")
	secret := os.Getenv("GOOGLE_CLIENT_SECRET")

	client := &oauth2.Config{
		RedirectURL: redir,
		ClientID: clientid,
		ClientSecret: secret,
		Scopes: []string{"https://www.googleapis.com/auth/plus.login", "email"},
		Endpoint: google.Endpoint,
	}

	token, err := client.Exchange(oauth2.NoContext, code)

	if err != nil {
		return runtime.ErrorOut(err)
	}

	runtime.Logger().Infof("google auth callback received token[%s]", token.AccessToken)

	return nil
}
