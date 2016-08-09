package middleware

import "strings"
import "encoding/base64"
import "github.com/golang/glog"
import "github.com/kataras/iris"

import "github.com/sizethree/meritoss.api/api"
import "github.com/sizethree/meritoss.api/api/models"

func ClientAuthentication(context *iris.Context) {
	header := context.RequestHeader("X-CLIENT-AUTH")

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

	runtime, ok := context.Get("runtime").(*api.Runtime)

	if !ok {
		glog.Error("bad runtime found while looking up auth header")
		context.Panic()
		context.StopExecution()
		return
	}

	var token models.ClientToken

	if e := runtime.DB.Where("token = ?", parts[0]).First(&token).Error; e != nil {
		glog.Infof("(client auth) unable to find token by %s, moving on\n", parts[0])
		context.Next()
		return
	}

	if e := runtime.DB.Where("client_id = ?", parts[1]).First(&runtime.Client).Error; e != nil {
		glog.Infof("(client auth) unable to find client by %s, moving on\n", parts[1])
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
