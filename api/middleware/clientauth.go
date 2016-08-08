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

	runtime, _ := context.Get("runtime").(api.Runtime)
	var token models.ClientToken
	var client models.Client
	// var user models.User

	if e := runtime.DB.Where("token = ?", parts[0]).First(&token).Error; e != nil {
		glog.Infof("(client auth) unable to find token by %s, moving on\n", parts[0])
		context.Next()
		return
	}

	if e := runtime.DB.Where("client_id = ?", parts[1]).First(&client).Error; e != nil {
		glog.Infof("(client auth) unable to find client by %s, moving on\n", parts[1])
		context.Next()
		return
	}

	if client.ID != token.Client {
		glog.Infof("(client auth) mismatch: client was %d but token#%d points at %d\n", client.ID, token.ID, token.Client)
		context.Next()
		return
	}

	glog.Infof("found token with client %d and user %d\n", token.Client, token.User)

	context.Next()
}
