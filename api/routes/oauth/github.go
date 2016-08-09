package oauth

import "os"
import "fmt"
import "bytes"
import "errors"
import "net/http"
import "encoding/json"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"

const token_url string = "https://github.com/login/oauth/access_token"

type OauthResponse map[string]string

func Github(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)
	code := context.URLParam("code")

	if !ok {
		glog.Error("bad runtime")
		context.Panic()
		context.StopExecution()
		return
	}

	client_id, secret := os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET")

	if len(code) < 2 {
		runtime.Error(errors.New("invalid user id"))
		context.Next()
		return
	}

	body := map[string]string {
		"code": code,
		"client_id": client_id,
		"client_secret": secret,
		"scope": "user:email user read:org",
	}

	glog.Infof("github authorization attempt w/ code: %s\n", code)

	jsondata, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", token_url, bytes.NewBuffer(jsondata))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	if err != nil {
		runtime.Error(errors.New("problem connecting to github"))
		context.Next()
		return
	}

	client := &http.Client{}

	glog.Info(string(jsondata))

	response, err := client.Do(request)

	if err != nil {
		apierr := errors.New(fmt.Sprintf("bad github api response: %s", err.Error()))
		runtime.Error(apierr)
		context.Next()
		return
	}

	if response.Status != "200 OK" {
		apierr := errors.New(fmt.Sprintf("bad github api response: %s"))
		runtime.Error(apierr)
		context.Next()
		return
	}

	defer response.Body.Close()

	var rdata OauthResponse

	err = json.NewDecoder(response.Body).Decode(&rdata)

	if err != nil {
		runtime.Error(errors.New("invalid json data"))
		context.Next()
		return
	}

	token := rdata["access_token"]

	glog.Infof("received token: %s\n", token)

	runtime.DB.Begin().Commit()
}
