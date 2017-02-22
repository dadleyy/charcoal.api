package routes

import "os"
import "fmt"
import "net/url"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

const ErrBadAuthCode = "BAD_AUTH_CODE"
const ErrNoClientAssociated = "NO_ASSOICATED_CLIENT"
const ErrMissingClientRedirect = "NO_REDIRECT_URI"
const ErrMissingAuthEndpoint = "BAD_AUTH_ENDPOINT"
const ErrInvalidGoogleResponse = "BAD_GOOGLE_RESPONSE"

func GoogleOauthRedirect(runtime *net.RequestRuntime) error {
	clientid := os.Getenv("GOOGLE_CLIENT_ID")
	redir := os.Getenv("GOOGLE_REDIRECT_URL")
	fin, err := url.Parse(services.EndpointGoogleAuth)

	if err != nil {
		runtime.Errorf("trouble parsing google auth endpoint: %s", services.EndpointGoogleAuth)
		return runtime.AddError(fmt.Errorf("SERVER_ERROR"))
	}

	query := runtime.URL.Query()

	requester := query.Get("client_id")

	if len(requester) == 0 {
		return runtime.AddError(fmt.Errorf(ErrNoClientAssociated))
	}

	var client models.Client

	if err := runtime.Where("client_id = ?", requester).First(&client).Error; err != nil {
		runtime.Errorf("invalid client id used in google auth: %s", clientid)
		return runtime.AddError(fmt.Errorf(ErrNoClientAssociated))
	}

	if len(client.RedirectUri) == 0 {
		runtime.Errorf("client %d (%s) is missing a redirect uri", client.ID, client.Name)
		return runtime.AddError(fmt.Errorf(ErrMissingClientRedirect))
	}

	queries := make(url.Values)

	queries.Set("response_type", "code")
	queries.Set("redirect_uri", redir)
	queries.Set("client_id", clientid)
	queries.Set("scope", "https://www.googleapis.com/auth/plus.login email")
	queries.Set("access_type", "offline")

	// set the state that gets sent to google (which will get sent back to us) to the client id
	// proved to us that represents the client opening this dialog.
	queries.Set("state", requester)

	fin.RawQuery = queries.Encode()

	runtime.Redirect(fin.String())
	return nil
}

func GoogleOauthReceiveCode(runtime *net.RequestRuntime) error {
	query := runtime.URL.Query()

	// extract the code sent from google and the "state" which is the client id originally sent
	// during the outh prompt so that we know who to add a token to.
	code := query.Get("code")
	state := query.Get("state")

	if len(state) == 0 {
		runtime.Errorf("unable to find state sent back from google")
		return runtime.AddError(fmt.Errorf(ErrInvalidGoogleResponse))
	}

	var client models.Client

	if err := runtime.Where("client_id = ?", state).First(&client).Error; err != nil {
		runtime.Errorf("invalid client id used in google auth: %s", state)
		return runtime.AddError(fmt.Errorf(ErrNoClientAssociated))
	}

	if len(client.RedirectUri) == 0 {
		runtime.Errorf("unable to find auth code sent from google")
		return runtime.AddError(fmt.Errorf(ErrBadAuthCode))
	}

	if len(code) == 0 {
		runtime.Errorf("unable to find auth code sent from google")
		runtime.Redirect(fmt.Sprint("%s?error=bad_code", client.RedirectUri))
		return nil
	}

	authman := services.GoogleAuthentication{runtime.DB}

	result, err := authman.Process(&client, code)

	if err != nil {
		runtime.Errorf("unable to authenticate client /w code: %s", err.Error())
		runtime.Redirect(fmt.Sprint("%s?error=bad_code", client.RedirectUri))
		return nil
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
