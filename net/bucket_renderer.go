package net

import "net/http"

type BucketRenderer interface {
	Render(*ResponseBucket, http.ResponseWriter) error
}
