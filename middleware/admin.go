package middleware

import "fmt"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/defs"
import "github.com/dadleyy/charcoal.api/services"

func RequireAdmin(handler net.HandlerFunc) net.HandlerFunc {
	check := func(runtime *net.RequestRuntime) *net.ResponseBucket {
		uman := services.UserManager{runtime.DB, runtime.Logger}

		if uman.IsAdmin(&runtime.User) != true || runtime.Client.System != true {
			return runtime.SendErrors(fmt.Errorf(defs.ErrUnauthorizedAdmin))
		}

		runtime.Debugf("user checks out as admin, continuing...")
		return handler(runtime)
	}

	return check
}
