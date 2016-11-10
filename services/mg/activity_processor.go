package mg

import "github.com/sizethree/miritos.api/activity"

type ActivityProcessor interface {
	Process(*Message) ([]activity.Message, error)
}
