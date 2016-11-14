package mg

import "fmt"
import "net/http"
import "encoding/json"

type Client struct {
	ApiKey string
}

func (client Client) Retreive(url string) (Message, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return Message{}, err
	}

	get := &http.Client{}
	req.SetBasicAuth("api", client.ApiKey)

	response, err := get.Do(req)

	if err != nil {
		return Message{}, err
	}

	defer response.Body.Close()

	deocode := json.NewDecoder(response.Body)

	result := Message{ContentMap: make(ContentIdMap)}

	if err := deocode.Decode(&result); err != nil {
		return Message{}, err
	}

	if len(result.ContentMap) == 0 {
		return result, fmt.Errorf("NO_CONTENTS")
	}

	return result, nil
}
