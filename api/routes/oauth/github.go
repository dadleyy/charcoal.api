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
import "github.com/sizethree/meritoss.api/api/middleware"

const token_url string = "https://github.com/login/oauth/access_token"

type OauthResponse map[string]string

func Github(ctx *iris.Context) {
	bucket, _ := ctx.Get("jsonapi").(*middleware.Bucket)
	runtime, _ := ctx.Get("runtime").(api.Runtime)
	code := ctx.URLParam("code")

	client_id, secret := os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET")

	if len(code) < 2 {
		bucket.Errors = append(bucket.Errors, errors.New("invalid user id"))
		ctx.Next()
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
		bucket.Errors = append(bucket.Errors, errors.New("problem connecting to github"))
		ctx.Next()
		return
	}

	client := &http.Client{}

	glog.Info(string(jsondata))

	response, err := client.Do(request)

	if err != nil {
		apierr := errors.New(fmt.Sprintf("bad github api response: %s", err.Error()))
		bucket.Errors = append(bucket.Errors, apierr)
		ctx.Next()
		return
	}

	if response.Status != "200 OK" {
		apierr := errors.New(fmt.Sprintf("bad github api response: %s"))
		bucket.Errors = append(bucket.Errors, apierr)
		ctx.Next()
		return
	}

	defer response.Body.Close()

	var rdata OauthResponse

	err = json.NewDecoder(response.Body).Decode(&rdata)

	if err != nil {
		bucket.Errors = append(bucket.Errors, errors.New("fuck!"))
		ctx.Next()
		return
	}

	token := rdata["access_token"]

	glog.Infof("received token: %s\n", token)

	runtime.DB.Begin().Commit()
}
