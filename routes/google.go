package routes

import "os"
import "errors"
import "net/url"
import "net/http"
import "encoding/base64"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/services"

const ERR_BAD_RUNTIME = "BAD_RUNTIME"
const ERR_BAD_AUTH_CODE = "BAD_AUTH_CODE"
const ERR_NO_ASSOCIATED_CLIENT_GOOGLE_AUTH = "NO_ASSOICATED_CLIENT"
const GOOGLE_AUTH_ENDPOINT = "https://accounts.google.com/o/oauth2/v2/auth"

func GoogleOauthRedirect(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Miritos)

	clientid := os.Getenv("GOOGLE_CLIENT_ID")
	redir := os.Getenv("GOOGLE_REDIRECT_URL")
	fin, err := url.Parse(GOOGLE_AUTH_ENDPOINT)

	if !ok || err != nil {
		return errors.New("ERR_BAD_RUNTIME")
	}

	state := ectx.QueryParam("client_id")

	if len(state) == 0 {
		return runtime.ErrorOut(errors.New("BAD_CLIENT_ID"))
	}

	queries := make(url.Values)

	queries.Set("response_type", "code")
	queries.Set("redirect_uri", redir)
	queries.Set("client_id", clientid)
	queries.Set("scope", "https://www.googleapis.com/auth/plus.login email")
	queries.Set("access_type", "offline")
	queries.Set("state", base64.StdEncoding.EncodeToString([]byte(state)))
	fin.RawQuery = queries.Encode()

	ectx.Redirect(http.StatusTemporaryRedirect, fin.String())
	return nil
}

func GoogleOauthReceiveCode(ectx echo.Context) error {
	runtime, ok := ectx.(*context.Miritos)

	if ok != true {
		return errors.New(ERR_BAD_RUNTIME)
	}

	// extract the code sent from google and the "state" which is the client id
	// originally sent during the outh prompt
	code := runtime.QueryParam("code")
	state := runtime.QueryParam("state")

	if len(code) == 0 {
		return runtime.ErrorOut(errors.New(ERR_BAD_AUTH_CODE))
	}

	if len(state) == 0 {
		return runtime.ErrorOut(errors.New(ERR_NO_ASSOCIATED_CLIENT_GOOGLE_AUTH))
	}

	// decode the client id sent along in the redirect
	referrer, err := base64.StdEncoding.DecodeString(state)

	if err != nil {
		return runtime.ErrorOut(err)
	}

	authman := services.GoogleAuthentication{runtime.DB}

	result, err := authman.Process(string(referrer), code)

	if err != nil {
		return runtime.ErrorOut(err)
	}

	fin, err := url.Parse(result.RedirectUri())

	if err != nil {
		return runtime.ErrorOut(err)
	}

	queries := make(url.Values)
	queries.Set("token", result.Token())
	fin.RawQuery = queries.Encode()

	runtime.Redirect(http.StatusTemporaryRedirect, fin.String())

	return nil
}
