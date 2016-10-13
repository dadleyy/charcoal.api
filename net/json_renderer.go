package net

import "net/http"
import "encoding/json"

type JsonRenderer struct {
	bucket *ResponseBucket
}

type commonResult struct {
	Status string `json:"status"`
	Meta map[string]interface{} `json:"meta"`
}

func (renderer JsonRenderer) Render(response http.ResponseWriter) error {
	headers := response.Header()
	headers["Content-Type"] = []string{"application/json"}

	writer := json.NewEncoder(response)

	if errors := renderer.bucket.errors; len(errors) >= 1 {
		elist := make([]string, len(errors))

		for i, e := range errors {
			elist[i] = e.Error()
		}

		result := struct {
			commonResult
			Errors []string `json:"errors"`
		}{commonResult{"ERRORED", renderer.bucket.meta}, elist}

		response.WriteHeader(http.StatusBadRequest)
		if err := writer.Encode(result); err != nil {
			return err
		}

		return nil
	}

	response.WriteHeader(http.StatusOK)
	result := struct {
		commonResult
		Results []Result `json:"results"`
	}{commonResult{"SUCCESS", renderer.bucket.meta}, renderer.bucket.results}

	if err := writer.Encode(result); err != nil {
		return err
	}

	return nil
}

