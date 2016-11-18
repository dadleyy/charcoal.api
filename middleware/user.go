package middleware

import "fmt"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/services"

const ErrBadBearerToken = "ERR_BAD_BEARER"

func InjectUser(handler net.HandlerFunc) net.HandlerFunc {
	inject := func(runtime *net.RequestRuntime) error {
		client := runtime.Client

		if valid := client.ID >= 1; valid != true {
			runtime.Debugf("client missing - inject user cannot continue")
			return handler(runtime)
		}

		headers := runtime.Header
		bearer := headers.Get("X-CLIENT-BEARER-TOKEN")

		if len(bearer) < 1 {
			runtime.Debugf("no bearer token header found while injecting user info")
			return handler(runtime)
		}

		clientmgr := services.UserClientManager{runtime.Database()}

		if err := clientmgr.Validate(bearer, &client); err != nil {
			runtime.Debugf("unable to validate bearer \"%s\" for client \"%d", bearer, client.ID)
			return handler(runtime)
		}

		var token models.ClientToken

		if err := runtime.Database().Where("token = ?", bearer).First(&token).Error; err != nil {
			runtime.Debugf("unable to find token from %s", bearer)
			return handler(runtime)
		}

		if err := runtime.Database().First(&runtime.User, token.User).Error; err != nil {
			runtime.Debugf("unable to find user from %s", bearer)
			return handler(runtime)
		}

		runtime.SetMeta("user", runtime.User.Public())

		return handler(runtime)
	}

	return inject
}

func RequireUser(handler net.HandlerFunc) net.HandlerFunc {
	require := func(runtime *net.RequestRuntime) error {
		if valid := runtime.User.ID >= 1; valid != true {
			return runtime.AddError(fmt.Errorf(ErrBadBearerToken))
		}

		return handler(runtime)
	}

	return require
}
