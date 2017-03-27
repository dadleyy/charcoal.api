package net

import "io"
import "fmt"
import "time"
import "net/http"
import "github.com/jinzhu/gorm"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/filestore"

type ServerRuntime struct {
	*log.Logger
	*gorm.DB

	Streams map[string](chan activity.Message)
	Mux     *Multiplexer
}

// request
//
// Given http.Request and UrlParam references, this function will return the request context
// that will ultimately be sent down the handlerfunc chain matched by the multiplexer.
func (server *ServerRuntime) Request(request *http.Request, params *UrlParams) *RequestRuntime {
	fs := filestore.S3FileStore{}

	runtime := RequestRuntime{
		Request:   request,
		UrlParams: params,
		FileSaver: fs,
		Logger:    server.Logger,
		DB:        server.DB,

		streams: server.Streams,
	}

	return &runtime
}

// ServeHTTP
//
// Used by the http.Server instance to handle requests. always renders json
func (server *ServerRuntime) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler, params, found := server.Mux.Find(request.Method, request.URL.Path)

	if response.Header().Get("status") != "" {
		server.Debugf("has status: %v", response.Header().Get("status"))
		return
	}

	// build the request runtime
	result := &ResponseBucket{Errors: []error{fmt.Errorf("not-found")}}

	if found == true {
		runtime := server.Request(request, &params)
		defer runtime.Close()
		result = handler(runtime)
	}

	if result == nil {
		result = &ResponseBucket{}
	}

	if result.Meta == nil {
		result.Meta = make(map[string]interface{})
	}

	result.Set("time", time.Now())

	if len(result.Redirect) >= 1 {
		outh := response.Header()
		outh.Set("Location", result.Redirect)
		response.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if len(result.Proxy) >= 1 {
		resp, err := http.Get(result.Proxy)

		if err != nil {
			server.Logger.Debugf("unable to download file: %s", err.Error())
			fmt.Fprintf(response, "not found")
			return
		}

		outh := response.Header()

		outh.Set("Content-Length", resp.Header.Get("Content-Length"))
		outh.Set("Content-Type", resp.Header.Get("Content-Type"))

		response.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		io.Copy(response, resp.Body)
		return
	}

	var renderer BucketRenderer

	switch request.Header.Get("accepts") {
	default:
		renderer = JsonRenderer{}
	}

	renderer.Render(result, response)
}
