package net

import "fmt"
import "strings"
import "net/url"
import "net/http"
import "mime/multipart"

type HttpBody struct {
	url.Values
	Files    map[string]*multipart.FileHeader
	jsonBody []byte
}

func parseMultiPartBody(request *http.Request, limit int64) (*HttpBody, error) {
	if err := request.ParseMultipartForm(limit); err != nil {
		return nil, err
	}

	files := make(map[string]*multipart.FileHeader)
	result := &HttpBody{Values: url.Values{}, Files: files}

	for key, vals := range request.MultipartForm.Value {
		for _, val := range vals {
			result.Add(key, val)
		}
	}

	for key, files := range request.MultipartForm.File {
		if len(files) != 0 {
			result.AddFile(key, files[0])
		}
	}

	return result, nil
}

func parseUrlEncodedBody(request *http.Request, limit int64) (*HttpBody, error) {
	if err := request.ParseForm(); err != nil {
		return nil, err
	}

	files := make(map[string]*multipart.FileHeader)
	result := &HttpBody{Values: url.Values{}, Files: files}

	for key, vals := range request.PostForm {
		for _, val := range vals {
			result.Add(key, val)
		}
	}

	return result, nil
}

func ParseBody(request *http.Request, limit int64) (*HttpBody, error) {
	formatting := ""

	for _, value := range []string{"Content-Type", "content-type"} {
		if header := request.Header.Get(value); header != "" {
			formatting = header
		}
	}

	if strings.Contains(formatting, "multipart/form-data") {
		return parseMultiPartBody(request, limit)
	}

	if strings.Contains(formatting, "form-urlencoded") {
		return parseUrlEncodedBody(request, limit)
	}

	return nil, fmt.Errorf("BAD_CONTENT_TYPE")
}

func (d *HttpBody) AddFile(key string, file *multipart.FileHeader) {
	d.Files[key] = file
}
