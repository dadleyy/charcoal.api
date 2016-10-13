package routes

import "github.com/sizethree/miritos.api/net"

func System(runtime *net.RequestRuntime) error {
	runtime.Result("OK")
	return nil
}
