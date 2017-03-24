package routes

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"

func FindGameMembershipHistory(runtime *net.RequestRuntime) error {
	results := make([]models.GameMembershipHistory, 0)
	blueprint := runtime.Blueprint(runtime.DB)

	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Debugf("invalid blueprint apply: %s", err.Error())
		return err
	}

	for _, r := range results {
		runtime.AddResult(r)
	}

	runtime.SetTotal(total)

	return nil
}
