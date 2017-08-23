package mg

import "github.com/dadleyy/charcoal.api/charcoal/bg"

type ProcessedItem struct {
	Error   error
	Message bg.Message
	Item    ContentItem
}

type ActivityProcessor interface {
	Process(*Message, chan ProcessedItem)
}
