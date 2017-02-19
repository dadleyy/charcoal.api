package mg

import "github.com/dadleyy/charcoal.api/activity"

type ProcessedItem struct {
	Error   error
	Message activity.Message
	Item    ContentItem
}

type ActivityProcessor interface {
	Process(*Message, chan ProcessedItem)
}
