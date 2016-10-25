package routes

import "fmt"
import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"

func FindClients(runtime *net.RequestRuntime) error {
	blueprint := runtime.Blueprint()
	var clients []models.Client

	total, err := blueprint.Apply(&clients, runtime.Database())

	if err != nil {
		runtime.Debugf("unable to query clients: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, client := range clients {
		runtime.AddResult(client)
	}

	runtime.SetMeta("total", total)

	return nil
}
