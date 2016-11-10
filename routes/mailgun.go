package routes

import "os"
import "fmt"
import "strings"
import "github.com/albrow/forms"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/services/mg"

func MailgunUploadHook(runtime *net.RequestRuntime) error {
	secret := os.Getenv("MAILGUN_WEBHOOK_SECRET")
	query := runtime.URL.Query()

	if v, ok := query["secret"]; !ok || len(v) != 1 || v[0] != secret {
		return runtime.AddError(fmt.Errorf("UNAUTHORIZED"))
	}

	body, err := forms.Parse(runtime.Request)
	if err != nil {
		return runtime.AddError(err)
	}

	var message mg.Message
	message.ContentMap = make(mg.ContentIdMap)

	if err := body.BindJSON(&message); err != nil {
		return err
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
