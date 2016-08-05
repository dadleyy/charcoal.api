package middleware

import "time"

import "github.com/golang/glog"
import "github.com/kataras/iris"

type MetaData struct {
	Total int
	Page int
	Time time.Time
}

type BucketResult interface {
	ToJson() map[string]interface{}
}

type Bucket struct {
	Meta MetaData
	Results []BucketResult
	Errors []error
}

type bucketAlias Bucket

type bucketJson struct {
	*bucketAlias
	Errors []string
	Results []interface{}
	Status string
}

func (b *Bucket) Render(ctx *iris.Context) {
	if ctx.IsStopped() {
		glog.Errorf("skipping jsonapi render - context was stopped\n")
		return
	}

	b.Meta.Time = time.Now()

	json := bucketJson{
		Status: "OK",
		bucketAlias: (*bucketAlias)(b),
	}

	status := iris.StatusOK

	if len(b.Errors) >= 1 {
		json.Status = "FAILED"

		for _, e := range b.Errors {
			json.Errors = append(json.Errors, e.Error())
		}

		status = iris.StatusBadRequest
	}

	if len(b.Results) >= 1 {
		for _, r := range b.Results {
			json.Results = append(json.Results, r.ToJson())
		}
	}

	ctx.JSON(status, json)
}

func JsonAPI(ctx *iris.Context) {
	bucket := Bucket{}
	ctx.Set("jsonapi", &bucket)
	defer bucket.Render(ctx)
	ctx.Next()
}
