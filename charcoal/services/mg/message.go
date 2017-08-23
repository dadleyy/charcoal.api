package mg

import "strings"

type ContentItem struct {
	Url         string `json:"url"`
	ContentType string `json:"content-type"`
	Name        string `json:"name"`
}

type ContentIdMap map[string]ContentItem
type AttachmentList []ContentItem

type Message struct {
	From        string         `json:"From"`
	Subject     string         `json:"Subject"`
	Body        string         `json:"body-html"`
	ContentMap  ContentIdMap   `json:"content-id-map"`
	Attachments AttachmentList `json:"attachments"`
}

// Images
//
// returns an array of all the content items that are images.
func (message *Message) Images() []ContentItem {
	result := make([]ContentItem, 0)

	for _, item := range message.Attachments {
		if valid := strings.HasPrefix(item.ContentType, "image/"); valid != true {
			continue
		}

		result = append(result, item)
	}

	return result
}
