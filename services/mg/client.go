package mg

import "fmt"
import "net/http"
import "encoding/json"

type ContentItem struct {
	Url         string `json:"url"`
	ContentType string `json:"content-type"`
	Name        string `json:"name"`
}

type ContentIdMap map[string]ContentItem

type Client struct {
	ApiKey string
}

func (client Client) Retreive(url string) (Message, error) {
	req, err := http.NewRequest("GET", url, nil)
	var result Message

	if err != nil {
		return result, err
	}

	get := &http.Client{}
	req.SetBasicAuth("api", client.ApiKey)

	response, err := get.Do(req)

	if err != nil {
		return result, err
	}

	defer response.Body.Close()

	deocode := json.NewDecoder(response.Body)

	result.ContentMap = make(ContentIdMap)

	if err := deocode.Decode(&result); err != nil {
		return result, err
	}

	if len(result.ContentMap) == 0 {
		return result, fmt.Errorf("NO_CONTENTS")
	}

	return result, nil
}
