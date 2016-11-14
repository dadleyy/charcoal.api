package routes

import "os"
import "fmt"
import "strings"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/services/mg"

func MailgunUploadHook(runtime *net.RequestRuntime) error {
	secret := os.Getenv("MAILGUN_WEBHOOK_SECRET")
	query := runtime.URL.Query()

	if v, ok := query["secret"]; !ok || len(v) != 1 || v[0] != secret {
		return runtime.AddError(fmt.Errorf("UNAUTHORIZED"))
	}

	body, err := net.ParseBody(runtime.Request, 150000000)

	if err != nil {
		return err
	}

	runtime.Debugf("received message-id from mailgun: \"%s\"", body.Get("Message-Id"))
	return nil

	var message mg.Message
	message.ContentMap = make(mg.ContentIdMap)

	/*
		if err := body.BindJSON(&message); err != nil {
			return err
		}
	*/

	var processor mg.ActivityProcessor
	start := strings.Split(message.Subject, ":")[0]

	runtime.Debugf("received message, subject line: \"%s\"", message.Subject)

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
