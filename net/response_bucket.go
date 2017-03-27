package net

type Result interface {
}

type ResponseBucket struct {
	Errors   []error
	Results  interface{}
	Meta     map[string]interface{}
	Proxy    string
	Redirect string
}

func (b *ResponseBucket) Set(key string, value interface{}) {
	if b.Meta == nil {
		return
	}

	b.Meta[key] = value
}
