package middleware

import "errors"
import "strings"
import "encoding/base64"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/models"

func ClientAuthentication(context *iris.Context) {
	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok {
		glog.Error("bad runtime found while looking up auth header")
		context.Panic()
		context.StopExecution()
		return
	}

	header := context.RequestHeader("X-CLIENT-AUTH")
	clientid := context.RequestHeader("X-CLIENT-ID")

	if len(clientid) < 1 {
		glog.Infof("(client auth) unacceptable request - no client id found")
		runtime.Render(context)
		return
	}

	if e := runtime.DB.Where("client_id = ?", clientid).First(&runtime.Client).Error; e != nil {
		glog.Infof("(client auth) unacceptable request - bad client id")
		runtime.Error(errors.New("invalid client id"))
		runtime.Render(context)
		return
	}

	if len(header) < 1 {
		glog.Infof("no client header found, continuing\n");
		context.Next()
		return
	}

	data, err := base64.StdEncoding.DecodeString(header)

	if err != nil {
		glog.Errorf("unable to decode client auth header: %s", err.Error())
		context.Next()
		return
	}

	parts := strings.Split(string(data), ":")

	if len(parts) != 2 {
		glog.Errorf("bad token header (too many parts): %s", data)
		context.Next()
		return
	}

	var token models.ClientToken

	if e := runtime.DB.Where("token = ?", parts[0]).First(&token).Error; e != nil {
		glog.Infof("(client auth) unable to find token by %s, moving on\n", parts[0])
		context.Next()
		return
	}

	client := runtime.Client

	if client.ID != token.Client {
		glog.Infof("(client auth) mismatch: client was %d but token#%d points at %d\n", client.ID, token.ID, token.Client)
		context.Next()
		return
	}

	if e := runtime.DB.Where("ID = ?", token.User).First(&runtime.User).Error; e != nil {
		glog.Infof("(client auth) unable to find user %d\n", token.User)
		context.Next()
		return
	}

	glog.Infof("found user %s\n", runtime.User.Email)
	context.Next()
}
