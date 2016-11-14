package mg

import "github.com/sizethree/miritos.api/activity"

type ProcessedItem struct {
	Error   error
	Message activity.Message
	Item    ContentItem
}

type ActivityProcessor interface {
	Process(*Message, chan ProcessedItem)
}
