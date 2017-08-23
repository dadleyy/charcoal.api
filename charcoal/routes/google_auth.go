package routes

import "os"
import "fmt"
import "net/url"

import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/defs"
import "github.com/dadleyy/charcoal.api/charcoal/models"
import "github.com/dadleyy/charcoal.api/charcoal/services"

// GoogleOauthRedirect sends the user to the google auth page.
func GoogleOauthRedirect(runtime *net.RequestRuntime) *net.ResponseBucket {
	clientid := os.Getenv("GOOGLE_CLIENT_ID")
	redir := os.Getenv("GOOGLE_REDIRECT_URL")
	fin, err := url.Parse(services.EndpointGoogleAuth)

	if err != nil {
		runtime.Errorf("[google redir] trouble parsing google auth endpoint: %s", services.EndpointGoogleAuth)
		return runtime.ServerError()
	}

	query := runtime.URL.Query()

	requester := query.Get("client_id")

	if len(requester) == 0 {
		return runtime.LogicError("invalid-client")
	}

	var client models.Client

	if err := runtime.Where("client_id = ?", requester).First(&client).Error; err != nil {
		runtime.Errorf("[google redir] invalid client id used in google auth: %s", clientid)
		return runtime.LogicError(defs.ErrGoogleNoClientAssociated)
	}

	if len(client.RedirectURI) == 0 {
		runtime.Errorf("[google redir] client %d (%s) is missing a redirect uri", client.ID, client.Name)
		return runtime.LogicError(defs.ErrGoogleMissingClientRedirect)
	}

	queries := make(url.Values)

	queries.Set("response_type", "code")
	queries.Set("redirect_uri", redir)
	queries.Set("client_id", clientid)
	queries.Set("scope", "https://www.googleapis.com/auth/plus.login email")
	queries.Set("access_type", "offline")
	queries.Set("prompt", "consent")

	// set the state that gets sent to google (which will get sent back to us) to the client id
	// proved to us that represents the client opening this dialog.
	queries.Set("state", requester)

	fin.RawQuery = queries.Encode()

	return runtime.Redirect(fin.String())
}

func GoogleOauthReceiveCode(runtime *net.RequestRuntime) *net.ResponseBucket {
	query := runtime.URL.Query()

	// extract the code sent from google and the "state" which is the client id originally sent
	// during the outh prompt so that we know who to add a token to.
	code := query.Get("code")
	state := query.Get("state")

	if len(state) == 0 {
		runtime.Errorf("[google receive code] unable to find state sent back from google")
		return runtime.LogicError(defs.ErrGoogleInvalidGoogleResponse)
	}

	var client models.Client

	if err := runtime.Where("client_id = ?", state).First(&client).Error; err != nil {
		runtime.Errorf("[google receive code] invalid client id used in google auth: %s", state)
		return runtime.LogicError(defs.ErrGoogleNoClientAssociated)
	}

	if len(client.RedirectURI) == 0 {
		runtime.Errorf("[google receive code] unable to find auth code sent from google")
		return runtime.LogicError(defs.ErrGoogleBadAuthCode)
	}

	if len(code) == 0 {
		runtime.Errorf("[google receive code] unable to find auth code sent from google")
		return runtime.Redirect(fmt.Sprintf("%s?error=bad_code", client.RedirectURI))
	}

	authman := services.GoogleAuthentication{runtime.DB, runtime.Logger}

	result, err := authman.Process(&client, code)

	if err != nil {
		runtime.Errorf("[google receive code] unable to authenticate client /w code: %s", err.Error())
		return runtime.Redirect(fmt.Sprintf("%s?error=bad_code", client.RedirectURI))
	}

	fin, err := url.Parse(result.RedirectURI())

	if err != nil {
		runtime.Warnf("[google receive code] invalid redirect uri for client[%d]: %s", client.ID, err.Error())
		return runtime.LogicError("invalid-redirect-uri")
	}

	queries := make(url.Values)
	queries.Set("token", result.Token())
	fin.RawQuery = queries.Encode()

	return runtime.Redirect(fin.String())
}
