package routes

import "os"
import "fmt"
import "net/url"
import "encoding/base64"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/services"

const ERR_BAD_RUNTIME = "BAD_RUNTIME"
const ERR_BAD_AUTH_CODE = "BAD_AUTH_CODE"
const ERR_NO_ASSOCIATED_CLIENT_GOOGLE_AUTH = "NO_ASSOICATED_CLIENT"
const GOOGLE_AUTH_ENDPOINT = "https://accounts.google.com/o/oauth2/v2/auth"

func GoogleOauthRedirect(runtime *net.RequestRuntime) error {
	clientid := os.Getenv("GOOGLE_CLIENT_ID")
	redir := os.Getenv("GOOGLE_REDIRECT_URL")
	fin, err := url.Parse(GOOGLE_AUTH_ENDPOINT)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_AUTH_CONFIG"))
	}

	query := runtime.URL.Query()

	state := query.Get("client_id")

	if len(state) == 0 {
		return runtime.AddError(fmt.Errorf("BAD_AUTH_CONFIG"))
	}

	var client models.Client

	if err := runtime.Database().Where("client_id = ?", state).First(&client).Error; err != nil {
		runtime.Errorf("invalid client id used in google auth: %s", clientid)
		return runtime.AddError(fmt.Errorf("BAD_CLIENT_ID"))
	}

	if len(client.RedirectUri) == 0 {
		runtime.Errorf("client %d (%s) is missing a redirect uri", client.ID, client.Name)
		return runtime.AddError(fmt.Errorf("MISSING_REDIRECT_URI"))
	}

	queries := make(url.Values)

	queries.Set("response_type", "code")
	queries.Set("redirect_uri", redir)
	queries.Set("client_id", clientid)
	queries.Set("scope", "https://www.googleapis.com/auth/plus.login email")
	queries.Set("access_type", "offline")
	queries.Set("state", base64.StdEncoding.EncodeToString([]byte(state)))
	fin.RawQuery = queries.Encode()

	runtime.Redirect(fin.String())
	return nil
}

func GoogleOauthReceiveCode(runtime *net.RequestRuntime) error {
	query := runtime.URL.Query()
	// extract the code sent from google and the "state" which is the client id
	// originally sent during the outh prompt
	code := query.Get("code")
	state := query.Get("state")

	if len(code) == 0 {
		runtime.Errorf("unable to find auth code sent from google")
		return runtime.AddError(fmt.Errorf(ERR_BAD_AUTH_CODE))
	}

	if len(state) == 0 {
		runtime.Errorf("unable to find state sent back from google")
		return runtime.AddError(fmt.Errorf(ERR_NO_ASSOCIATED_CLIENT_GOOGLE_AUTH))
	}

	// decode the client id sent along in the redirect
	referrer, err := base64.StdEncoding.DecodeString(state)

	if err != nil {
		return runtime.AddError(err)
	}

	authman := services.GoogleAuthentication{runtime.Database()}

	result, err := authman.Process(string(referrer), code)

	if err != nil {
		runtime.Errorf("unable to authenticate client /w code: %s", err.Error())
		return runtime.AddError(fmt.Errorf("ERR_BAD_CLIENT_CODE"))
	}

	fin, err := url.Parse(result.RedirectUri())

	if err != nil {
		return runtime.AddError(err)
	}

	queries := make(url.Values)
	queries.Set("token", result.Token())
	fin.RawQuery = queries.Encode()

	runtime.Redirect(fin.String())

	return nil
}
