package middleware

import "fmt"
import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/errors"
import "github.com/dadleyy/charcoal.api/services"

func RequireAdmin(handler net.HandlerFunc) net.HandlerFunc {
	check := func(runtime *net.RequestRuntime) error {
		uman := services.UserManager{runtime.DB}

		if uman.IsAdmin(&runtime.User) != true || runtime.Client.System != true {
			return runtime.AddError(fmt.Errorf(errors.ErrUnauthorizedAdmin))
		}

		runtime.Debugf("user checks out as admin, continuing...")
		return handler(runtime)
	}

	return check
}
