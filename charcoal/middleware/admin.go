package middleware

import "fmt"
import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/defs"
import "github.com/dadleyy/charcoal.api/charcoal/services"

// RequireAdmin validates that the current user is associated with the admin user role.
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
