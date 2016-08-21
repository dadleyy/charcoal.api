package middleware

import "time"
import "github.com/golang/glog"
import "github.com/dadleyy/iris"
import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"

type MetaData map[string]interface{}

type BucketResult interface {
	ToJson() map[string]interface{}
}

type ResponseBucket struct {
	Meta MetaData
	Results []BucketResult
	Errors []error
}

type bucketAlias ResponseBucket

type bucketJson struct {
	*bucketAlias
	Errors []string
	Results []interface{}
	Status string
}

// Runtime stuct
//
// Defines a context that is created for every request that is available to each handler.
type Runtime struct {
	Bucket *ResponseBucket
	DB db.Client
	User models.User
	Client models.Client
}

func (runt *Runtime) Error(e error) {
	runt.Bucket.Errors = append(runt.Bucket.Errors, e)
}

func (runt *Runtime) Result(result BucketResult) {
	runt.Bucket.Results = append(runt.Bucket.Results, result)
}

func (runt *Runtime) Meta(key string, val interface{}) {
	runt.Bucket.Meta[key] = val
}

// Runtime.Render
// 
// Given an iris request context, this function will render out the runtime's reponse bucket
// and appropriately set the status code as well as status property of the response data.
func (runt *Runtime) Render(context *iris.Context) {
	// close the database connection - we're done with it
	runt.DB.Close()

	// do not render if the context is finished
	if context.IsStopped() {
		glog.Errorf("skipping jsonapi render - context was stopped\n")
		return
	}

	// update the response bucket's current time
	runt.Bucket.Meta["time"] = time.Now()

	json := bucketJson{
		Status: "OK",
		bucketAlias: (*bucketAlias)(runt.Bucket),
	}

	// start of with an OK status
	status := iris.StatusOK

	// if at any point an error was added via the `Error` method, update the status
	// and add all of the errors as strings to the json struct we're building
	if len(runt.Bucket.Errors) >= 1 {
		json.Status = "FAILED"

		// convert each error to a string and add it to the json error
		for _, e := range runt.Bucket.Errors {
			json.Errors = append(json.Errors, e.Error())
		}

		// update the overall status code of this request
		status = iris.StatusBadRequest
	}

	// if at any point there was a result added via the `Result` method, add each of them 
	// to the json object's Results array, converting to json along the way
	if results := runt.Bucket.Results; len(results) >= 1 {
		for _, r := range results {
			json.Results = append(json.Results, r.ToJson())
		}
	}

	// finish by rendering
	context.JSON(status, json)

	// forcefully stop execution here
	context.StopExecution()
}


// Runtime
//
// Definess a middleware function that will inject an instance of `api.Runtime` into 
// each request so that the handler can access things like the orm by retreiving the 
// value from the context: 
//
// runtime, ok := context.Get("runtime").(*api.Runtime)
// 
// where api.Runtime is a struct that could look like:
//
// type Runtime struct {
//   DB *gorm.DB
// }
// 
// In this example, we're using the "github.com/jinzhu/gorm" package to handle db
// related communication and modeling.
func InjectRuntime(context *iris.Context) {
	var runtime Runtime

	// attempt to connect to the mysql database
	client, err := db.Get()

	// if there was an issue opening the connection, send a 500 error - nothing we 
	// can do here to solve that problem
	if err != nil {
		context.Panic()
		return
	}

	// initialize the runtime struct with the gorm client; all other properties will
	// be default initialized
	runtime = Runtime{Bucket: &ResponseBucket{Meta: make(MetaData)}, DB: client}

	// after the all middleware has completed, let the runtime do whatever it needs to
	// do in order to complete this request.
	defer runtime.Render(context)

	// inject our runtime into the user context for this request
	context.Set("runtime", &runtime)

	// move on now that it is exposed
	context.Next()
}
