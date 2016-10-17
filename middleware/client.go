package middleware

import "fmt"
import "strings"
import "encoding/base64"

import "github.com/sizethree/miritos.api/net"

const ERR_BAD_RUNTIME = "BAD_RUNTIME"
const ERR_MISSING_CLIENT_ID = "MISSING_CLIENT_ID"
const ERR_BAD_CLIENT_ID = "BAD_CLIENT_ID"

func InjectClient(handler net.HandlerFunc) net.HandlerFunc {
	inject := func(runtime *net.RequestRuntime) error {
		headers := runtime.Header
		auth := headers.Get("X-CLIENT-AUTH")

		if len(auth) < 1 {
			return handler(runtime)
		}

		decoded, err := base64.StdEncoding.DecodeString(auth)

		if err != nil {
			runtime.Debugf("bad client auth header: %s", auth)

			return handler(runtime)
		}

		parts := strings.Split(string(decoded), ":")

		if len(parts) < 2 || len(parts[0]) < 1 || len(parts[1]) < 1 {
			return handler(runtime)
		}

		where := runtime.Database().Where("client_id = ?", parts[0]).Where("client_secret = ?", parts[1])

		if e := where.First(&runtime.Client).Error; e != nil {
			runtime.Errorf("unable to find client: %s", e.Error())
			return handler(runtime)
		}

		runtime.Debugf("injected client \"%d\" auth, continuing", runtime.Client.ID);
		return handler(runtime)
	}

	return inject
}

func RequireClient(handler net.HandlerFunc) net.HandlerFunc {
	require := func(runtime *net.RequestRuntime) error {
		if valid := runtime.Client.ID >= 1; valid == true {
			return handler(runtime)
		}

		runtime.Error(fmt.Errorf(ERR_BAD_CLIENT_ID))
		return nil
	}

	return InjectClient(require)
}
