package net

import "net/http"

type BucketRenderer interface {
	Render(http.ResponseWriter) error
}

