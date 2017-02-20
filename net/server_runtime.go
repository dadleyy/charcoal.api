package net

import "io"
import "fmt"
import "time"
import "net/http"
import "github.com/labstack/gommon/log"

import "github.com/dadleyy/charcoal.api/db"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/filestore"

type RuntimeConfig struct {
	DB db.Config
}

type ServerRuntime struct {
	Logger *log.Logger
	Config RuntimeConfig
	Queue  chan activity.Message
	Mux    *Multiplexer
}

// request
//
// Given http.Request and UrlParam references, this function will return the request context
// that will ultimately be sent down the handlerfunc chain matched by the multiplexer.
func (server *ServerRuntime) Request(request *http.Request, params *UrlParams) (RequestRuntime, error) {
	errors := make([]error, 0)
	results := make([]Result, 0)
	meta := make(map[string]interface{})

	bucket := ResponseBucket{errors, results, meta, "", ""}

	meta["time"] = time.Now()

	fs := filestore.S3FileStore{}

	database := new(db.Connection)

	if err := db.Open(server.Config.DB, database); err != nil {
		return RequestRuntime{}, err
	}

	runtime := RequestRuntime{
		Request:   request,
		UrlParams: params,
		Logger:    server.Logger,
		queue:     server.Queue,
		bucket:    bucket,
		fs:        fs,
		database:  database,
	}

	return runtime, nil
}

// ServeHTTP
//
// Used by the http.Server instance to handle requests. always renders json
func (server *ServerRuntime) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler, params, found := server.Mux.Find(request.Method, request.URL.Path)

	// not found
	if found == false {
		server.Logger.Debugf("error matching route: %s", request.URL.Path)
		response.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(response, "not found")
		return
	}

	// build the request runtime
	runtime, err := server.Request(request, &params)

	// attempt to prepare a db connection for this request and error out if
	// something goes wrong along the way.
	if err != nil {
		server.Logger.Debugf("error matching route: %s", request.URL.Path)
		response.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(response, "not found")
		return
	}

	// once this function finishes we're done with the request.
	defer runtime.database.Close()

	var renderer BucketRenderer

	switch request.Header.Get("accepts") {
	default:
		renderer = JsonRenderer{&runtime.bucket}
	}

	if err := handler(&runtime); err != nil {
		server.Logger.Debugf("error handling route: %s", err.Error())
		renderer.Render(response)
		return
	}

	if len(runtime.bucket.redirect) >= 1 {
		outh := response.Header()
		outh.Set("Location", runtime.bucket.redirect)
		response.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if len(runtime.bucket.proxy) >= 1 {
		resp, err := http.Get(runtime.bucket.proxy)

		if err != nil {
			server.Logger.Debugf("unable to download file: %s", err.Error())
			fmt.Fprintf(response, "not found")
			return
		}

		outh := response.Header()

		outh.Set("Content-Length", resp.Header.Get("Content-Length"))
		outh.Set("Content-Type", resp.Header.Get("Content-Type"))

		response.WriteHeader(resp.StatusCode)
		server.Logger.Debugf("proxy-ing: \"%s\" | type[%s]", runtime.bucket.proxy, resp.Header.Get("Content-Type"))

		defer resp.Body.Close()

		io.Copy(response, resp.Body)
		return
	}

	renderer.Render(response)
}
