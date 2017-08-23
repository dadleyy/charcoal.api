package net

import "net/http"
import "encoding/json"

type JsonRenderer struct {
}

type commonResult struct {
	Status string                 `json:"status"`
	Meta   map[string]interface{} `json:"meta"`
}

func (renderer JsonRenderer) Render(bucket *ResponseBucket, response http.ResponseWriter) error {
	headers := response.Header()
	headers["Content-Type"] = []string{"application/json"}

	writer := json.NewEncoder(response)

	if errors := bucket.Errors; len(errors) >= 1 {
		elist := make([]string, len(errors))

		for i, e := range errors {
			elist[i] = e.Error()
		}

		result := struct {
			commonResult
			Errors []string `json:"errors"`
		}{commonResult{"ERRORED", bucket.Meta}, elist}

		response.WriteHeader(http.StatusBadRequest)
		if err := writer.Encode(result); err != nil {
			return err
		}

		return nil
	}

	response.WriteHeader(http.StatusOK)
	result := struct {
		commonResult
		Results interface{} `json:"results"`
	}{commonResult{"SUCCESS", bucket.Meta}, bucket.Results}

	if err := writer.Encode(result); err != nil {
		return err
	}

	return nil
}
