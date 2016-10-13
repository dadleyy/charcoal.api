package net

type ResponseBucket struct {
	errors []error
	results []Result
	meta map[string]interface{}
}
