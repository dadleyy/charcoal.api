package net

import "fmt"
import "time"
import "net/http"
import "github.com/labstack/gommon/log"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/filestore"

type ServerRuntime struct {
	Log *log.Logger
	DatabaseConnection *db.Connection
	Queue chan activity.Message
	Mux *Multiplexer
}

// request
// 
// Given http.Request and UrlParam references, this function will return the request context 
// that will ultimately be sent down the handlerfunc chain matched by the multiplexer.
func (server *ServerRuntime) request(request *http.Request, params *UrlParams) RequestRuntime {
	errors := make([]error, 0)
	results := make([]Result, 0)
	meta := make(map[string]interface{})

	bucket := ResponseBucket{errors, results, meta}

	meta["time"] = time.Now()

	fs := filestore.S3FileStore{}

	runtime := RequestRuntime{
		Request: request,
		UrlParams: params,
		database: server.DatabaseConnection,
		queue: server.Queue,
		log: server.Log,
		bucket: bucket,
		store: fs,
	}

	return runtime
}

// ServeHTTP
//
// Used by the http.Server instance to handle requests. always renders json
func (server *ServerRuntime) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler, params, found := server.Mux.Find(request.Method, request.URL.Path)

	// not found
	if found == false {
		server.Log.Debugf("error matching route: %s", request.URL.Path)
		response.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(response, "not found")
		return
	}

	// build the request runtime
	runtime := server.request(request, &params)

	if err := handler(&runtime); err != nil {
		server.Log.Debugf("error handling route: %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "server error")
		return
	}

	var renderer BucketRenderer

	switch request.Header.Get("accepts") {
	default:
		renderer = JsonRenderer{&runtime.bucket}
	}

	renderer.Render(response)
}
