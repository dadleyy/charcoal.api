package services

import "os"
import "net/url"
import "net/http"
import "encoding/json"
import "golang.org/x/oauth2"
import "golang.org/x/oauth2/google"

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

const GOOGLE_INFO_ENDPOINT = "https://www.googleapis.com/oauth2/v2/userinfo"

type GoogleUserInfo struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
}

type GoogleAuthentication struct {
	*context.Database
}

type GoogleAuthenticationResult struct {
	Client models.Client
	User models.User
	GoogleAccount models.GoogleAccount
	ClientToken models.ClientToken
}

func (result *GoogleAuthenticationResult) RedirectUri() string {
	return result.Client.RedirectUri
}

func (result *GoogleAuthenticationResult) Token() string {
	return result.ClientToken.Token
}

func (manager *GoogleAuthentication) Process(referrer, code string) (GoogleAuthenticationResult, error) {
	var result GoogleAuthenticationResult

	if err := manager.Where("client_id = ?", referrer).First(&result.Client).Error; err != nil {
		return result, err
	}

	auth := &oauth2.Config{
		RedirectURL: os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/plus.login", "email"},
		Endpoint: google.Endpoint,
	}

	token, err := auth.Exchange(oauth2.NoContext, code)

	if err != nil {
		return result, err
	}

	lookup, err := url.Parse(GOOGLE_INFO_ENDPOINT)

	if err != nil {
		return result, err
	}

	queries := make(url.Values)
	queries.Set("access_token", token.AccessToken)
	lookup.RawQuery = queries.Encode()

	response, err := http.Get(lookup.String())

	if err != nil {
		return result, err
	}

	var info GoogleUserInfo
	err = json.NewDecoder(response.Body).Decode(&info)

	if err != nil {
		return result, err
	}

	response.Body.Close()

	result.User = models.User{Email: info.Email, Name: info.Name}

	usrmgr := UserManager{manager.Database}

	if err := usrmgr.FindOrCreate(&result.User); err != nil {
		return result, err
	}

	clientmgr := UserClientManager{manager.Database}

	result.ClientToken, err = clientmgr.AssociateClient(&result.User, &result.Client)

	if err != nil {
		return result, err
	}

	result.GoogleAccount = models.GoogleAccount{
		GoogleID: info.ID,
		User: result.User.ID,
		AccessToken: token.AccessToken,
	}

	if err := manager.FirstOrCreate(&result.GoogleAccount, models.GoogleAccount{GoogleID: info.ID}).Error; err != nil {
		return GoogleAuthenticationResult{}, err
	}

	return result, nil
}