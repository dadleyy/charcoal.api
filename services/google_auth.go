package services

import "os"
import "fmt"
import "net/url"
import "net/http"
import "encoding/json"
import "golang.org/x/oauth2"
import "golang.org/x/oauth2/google"

import "github.com/dadleyy/charcoal.api/db"
import "github.com/dadleyy/charcoal.api/models"

const EndpointGoogleInfo = "https://www.googleapis.com/oauth2/v2/userinfo"
const EndpointGoogleAuth = "https://accounts.google.com/o/oauth2/v2/auth"

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GoogleAuthentication struct {
	*db.Connection
}

type GoogleAuthenticationResult struct {
	Client        models.Client
	User          models.User
	GoogleAccount models.GoogleAccount
	ClientToken   models.ClientToken
}

func (result *GoogleAuthenticationResult) RedirectUri() string {
	return result.Client.RedirectUri
}

func (result *GoogleAuthenticationResult) Token() string {
	return result.ClientToken.Token
}

func (manager *GoogleAuthentication) Process(client *models.Client, code string) (GoogleAuthenticationResult, error) {
	var result GoogleAuthenticationResult

	if client == nil {
		return GoogleAuthenticationResult{}, fmt.Errorf("BAD_CLIENT")
	}

	if err := manager.Where("client_id = ?", client.ClientID).First(&result.Client).Error; err != nil {
		return result, err
	}

	auth := &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/plus.login", "email"},
		Endpoint:     google.Endpoint,
	}

	token, err := auth.Exchange(oauth2.NoContext, code)

	if err != nil {
		return result, err
	}

	lookup, err := url.Parse(EndpointGoogleInfo)

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

	cursor := manager.Where("email = ?", info.Email)

	// check to see if we have an existing google account and if so, return the user associated
	// with it, as well as the token and client
	if err := cursor.First(&result.GoogleAccount).Error; err == nil && result.GoogleAccount.ID >= 1 {
		cursor := manager.Where("id = ?", result.GoogleAccount.User)

		if err := cursor.First(&result.User).Error; err != nil {
			return GoogleAuthenticationResult{}, fmt.Errorf("FOUND_DUPLICATE_NO_USER: %s", err.Error())
		}

		cursor = manager.Where("user = ? AND client = ?", result.User.ID, result.Client.ID)

		if err := cursor.First(&result.ClientToken).Error; err != nil {
			return GoogleAuthenticationResult{}, fmt.Errorf("FOUND_DUPLICATE_NO_CLIENT")
		}

		return result, nil
	}

	// at this point we know that there is no google account already associated with the email, so we should
	// attempt to create (or find) the user record and associate it with a fresh google account record.
	result.User = models.User{Email: &info.Email, Name: &info.Name}

	usrmgr := UserManager{manager.Connection}

	if usrmgr.ValidDomain(info.Email) != true {
		return GoogleAuthenticationResult{}, fmt.Errorf(ErrUnauthorizedDomain)
	}

	if err := usrmgr.FindOrCreate(&result.User); err != nil {
		return result, err
	}

	clientmgr := UserClientManager{manager.Connection}

	result.ClientToken, err = clientmgr.Associate(&result.User, &result.Client)

	if err != nil {
		return result, err
	}

	result.GoogleAccount = models.GoogleAccount{
		GoogleID:    info.ID,
		User:        result.User.ID,
		AccessToken: token.AccessToken,
		Email:       info.Email,
		Name:        info.Name,
	}

	if err := manager.FirstOrCreate(&result.GoogleAccount, models.GoogleAccount{GoogleID: info.ID}).Error; err != nil {
		return GoogleAuthenticationResult{}, err
	}

	return result, nil
}
