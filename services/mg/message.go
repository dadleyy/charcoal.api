package mg

import "fmt"
import "strings"
import "golang.org/x/net/html"

type Message struct {
	From       string       `json:"From"`
	Subject    string       `json:"subject"`
	Body       string       `json:"body-html"`
	ContentMap ContentIdMap `json:"content-id-map"`
}

func (message *Message) Images() []ContentItem {
	reader := strings.NewReader(message.Body)
	tokenizer := html.NewTokenizer(reader)
	result := make([]ContentItem, 0)

	for {
		chunk := tokenizer.Next()

		if chunk == html.ErrorToken {
			break
		}

		if chunk != html.StartTagToken {
			continue
		}

		nameb, attrs := tokenizer.TagName()

		if string(nameb) != "img" || attrs != true {
			continue
		}

		key, val, attrs := tokenizer.TagAttr()

		for ; attrs != false; key, val, attrs = tokenizer.TagAttr() {
			if string(key) != "src" || strings.HasPrefix(string(val), "cid:") != true {
				continue
			}

			cid := strings.TrimPrefix(string(val), "cid:")
			item, ok := message.ContentMap[fmt.Sprintf("<%s>", cid)]

			if !ok {
				continue
			}

			result = append(result, item)
		}
	}

	return result
}
