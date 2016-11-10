package routes

import "os"
import "fmt"
import "strings"
import "github.com/albrow/forms"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/services/mg"

func MailgunUploadHook(runtime *net.RequestRuntime) error {
	body, err := forms.Parse(runtime.Request)
	secret := os.Getenv("MAILGUN_WEBHOOK_SECRET")

	if err != nil {
		return runtime.AddError(err)
	}

	query := runtime.URL.Query()

	if v, ok := query["secret"]; !ok || len(v) != 1 || v[0] != secret {
		return runtime.AddError(fmt.Errorf("UNAUTHORIZED"))
	}

	validator := body.Validator()
	validator.Require("storage")
	validator.Require("message")

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		for _, m := range validator.Messages() {
			runtime.AddError(fmt.Errorf(m))
		}

		return nil
	}

	info := struct {
		Storage struct {
			Url string `json:"url"`
		}
	}{}

	if err := body.BindJSON(&info); err != nil {
		return err
	}

	runtime.Debugf("found storage url: %s", info.Storage.Url)

	mailgun := mg.Client{os.Getenv("MAILGUN_API_KEY")}

	message, err := mailgun.Retreive(info.Storage.Url)

	if err != nil {
		return runtime.AddError(err)
	}

	var processor mg.ActivityProcessor
	start := strings.Split(message.Subject, ":")[0]

	switch start {
	case "image", "photo":
		processor = &mg.ImageProcessor{runtime.Database(), runtime.Photos(), os.Getenv("MAILGUN_API_KEY")}
	default:
		return fmt.Errorf("INVALID_SUBJECT_LINE")
	}

	activity, err := processor.Process(&message)

	if err != nil {
		return err
	}

	runtime.Debugf("received message \"%s\" from: %s", message.Subject, message.From)

	for _, item := range activity {
		runtime.Publish(item)
	}

	return nil
}
