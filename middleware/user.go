package middleware

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/services"

const ErrBadBearerToken = "ERR_BAD_BEARER"

func InjectUser(handler net.HandlerFunc) net.HandlerFunc {
	inject := func(runtime *net.RequestRuntime) *net.ResponseBucket {
		client := runtime.Client

		if valid := client.ID >= 1; valid != true {
			runtime.Debugf("[inject user] client missing - inject user cannot continue")
			return handler(runtime)
		}

		headers := runtime.Header
		bearer := headers.Get(defs.ClientUserTokenHeader)

		if len(bearer) < 1 {
			runtime.Debugf("[inject user] no bearer token header found while injecting user info")
			return handler(runtime)
		}

		clientmgr := services.UserClientManager{runtime.DB}

		if err := clientmgr.Validate(bearer, &client); err != nil {
			runtime.Debugf("[inject user] bad bearer - client[%d]", client.ID)
			return handler(runtime)
		}

		var token models.ClientToken

		if err := runtime.Where("token = ?", bearer).First(&token).Error; err != nil {
			runtime.Debugf("[inject user] unable to find token from %s", bearer)
			return handler(runtime)
		}

		if err := runtime.First(&runtime.User, token.UserID).Error; err != nil {
			runtime.Debugf("[inject user] unable to find user from %s", bearer)
			return handler(runtime)
		}

		result := handler(runtime)

		if result != nil {
			result.Set("user", runtime.User.Public())
		}

		return result
	}

	return inject
}

func RequireUser(handler net.HandlerFunc) net.HandlerFunc {
	require := func(runtime *net.RequestRuntime) *net.ResponseBucket {
		if valid := runtime.User.ID >= 1; valid != true {
			return runtime.LogicError("invalid-user")
		}

		return handler(runtime)
	}

	return require
}
