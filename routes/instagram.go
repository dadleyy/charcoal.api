package routes

import "os"
import "fmt"
import "net/url"
import "net/http"
import "io/ioutil"
import "encoding/json"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"

const INSTAGRAM_AUTH_ENDPOINT = "https://api.instagram.com/oauth/authorize"
const INSTAGRAM_TOKEN_ENDPOINT = "https://api.instagram.com/oauth/access_token"

func InstaOauthRedirect(runtime *net.RequestRuntime) error {
	clientid := os.Getenv("INSTAGRAM_CLIENT_ID")
	redir := os.Getenv("INSTAGRAM_REDIRECT_URL")
	fin, err := url.Parse(INSTAGRAM_AUTH_ENDPOINT)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_AUTH_CONFIG"))
	}

	query := runtime.URL.Query()

	requester := query.Get("client_id")

	if len(requester) == 0 {
		return runtime.AddError(fmt.Errorf("BAD_AUTH_CONFIG"))
	}

	var client models.Client

	if err := runtime.Database().Where("client_id = ?", requester).First(&client).Error; err != nil {
		runtime.Errorf("invalid client id used in insta auth: %s", clientid)
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

	// set the state that gets sent to instagram (which will get sent back to us) to the client id
	// proved to us that represents the client opening this dialog.
	queries.Set("state", requester)

	fin.RawQuery = queries.Encode()

	runtime.Debugf("sending user to %s", fin.String())

	runtime.Redirect(fin.String())
	return nil
}

func InstaOauthReceiveCode(runtime *net.RequestRuntime) error {
	query := runtime.URL.Query()

	// extract the code sent from instagram and the "state" which is the client id originally sent
	// during the outh prompt so that we know who to add a token to.
	code := query.Get("code")
	state := query.Get("state")

	if len(state) == 0 {
		runtime.Errorf("unable to find state sent back from instagram")
		return runtime.AddError(fmt.Errorf("BAD_CLIENT"))
	}

	var client models.Client

	if err := runtime.Database().Where("client_id = ?", state).First(&client).Error; err != nil {
		runtime.Errorf("invalid client id used in instagram auth: %s", state)
		return runtime.AddError(fmt.Errorf("BAD_CLIENT_ID"))
	}

	if len(client.RedirectUri) == 0 {
		runtime.Errorf("unable to find auth code sent from instagram")
		return runtime.AddError(fmt.Errorf(ERR_BAD_AUTH_CODE))
	}

	if len(code) == 0 {
		runtime.Errorf("unable to find auth code sent from instagram")
		runtime.Redirect(fmt.Sprint("%s?error=bad_code", client.RedirectUri))
		return nil
	}

	packet := url.Values{
		"client_secret": {os.Getenv("INSTAGRAM_CLIENT_SECRET")},
		"client_id":     {os.Getenv("INSTAGRAM_CLIENT_ID")},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {os.Getenv("INSTAGRAM_REDIRECT_URL")},
		"code":          {code},
	}

	result, err := http.PostForm(INSTAGRAM_TOKEN_ENDPOINT, packet)

	if err != nil {
		runtime.Debugf("bad instagram token resopnse: %s", err.Error())
		runtime.Redirect(fmt.Sprint("%s?error=bad_code", client.RedirectUri))
		return nil
	}

	defer result.Body.Close()

	tdata, err := ioutil.ReadAll(result.Body)

	if err != nil {
		runtime.Debugf("bad instagram token resopnse: %s", err.Error())
		runtime.Redirect(fmt.Sprint("%s?error=bad_code", client.RedirectUri))
		return nil
	}

	var tresult struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(tdata, &tresult); err != nil {
		runtime.Debugf("bad instagram token resopnse: %s", err.Error())
		runtime.Redirect(fmt.Sprint("%s?error=bad_code", client.RedirectUri))
		return nil
	}

	return nil
}
