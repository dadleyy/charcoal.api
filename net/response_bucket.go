package net

type Result interface {
}

type ResponseBucket struct {
	errors   []error
	results  []Result
	meta     map[string]interface{}
	proxy    string
	redirect string
}
