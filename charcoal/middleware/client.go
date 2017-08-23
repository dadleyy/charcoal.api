package middleware

import "strings"
import "encoding/base64"

import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/defs"

// InjectClient loads the Client value of the request based on the client auth header.
func InjectClient(handler net.HandlerFunc) net.HandlerFunc {
	inject := func(runtime *net.RequestRuntime) *net.ResponseBucket {
		headers := runtime.Header
		auth := headers.Get(defs.ClientTokenHeader)

		if len(auth) < 1 {
			return handler(runtime)
		}

		decoded, err := base64.StdEncoding.DecodeString(auth)

		if err != nil {
			runtime.Debugf("[client middleware] bad client auth header: %s", auth)
			return handler(runtime)
		}

		parts := strings.Split(string(decoded), ":")

		if len(parts) < 2 || len(parts[0]) < 1 || len(parts[1]) < 1 {
			return handler(runtime)
		}

		where := runtime.Where("client_id = ?", parts[0]).Where("client_secret = ?", parts[1])

		if e := where.First(&runtime.Client).Error; e != nil {
			runtime.Errorf("[client middleware] unable to find client: %s", e.Error())
			return handler(runtime)
		}

		result := handler(runtime)

		if result != nil {
			result.Set("client", runtime.Client)
		}

		return result
	}

	return inject
}

// RequireClient rejects the incoming request unless the client was loaded successfully.
func RequireClient(handler net.HandlerFunc) net.HandlerFunc {
	require := func(runtime *net.RequestRuntime) *net.ResponseBucket {
		if valid := runtime.Client.ID >= 1; valid == true {
			return handler(runtime)
		}

		return runtime.LogicError("invalid-client")
	}

	return require
}
