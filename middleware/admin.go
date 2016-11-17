package middleware

import "fmt"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/services"

func RequireAdmin(handler net.HandlerFunc) net.HandlerFunc {
	check := func(runtime *net.RequestRuntime) error {
		uman := services.UserManager{runtime.Database()}

		if uman.IsAdmin(&runtime.User) != true {
			return runtime.AddError(fmt.Errorf("BAD_PERMISSIONS"))
		}

		runtime.Debugf("user checks out as admin, continuing...")
		return handler(runtime)
	}

	return check
}
