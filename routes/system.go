package routes

import "github.com/sizethree/miritos.api/net"

func System(runtime *net.RequestRuntime) error {
	runtime.AddResult("OK")
	return nil
}
