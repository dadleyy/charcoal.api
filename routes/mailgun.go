package routes

import "os"
import "fmt"
import "strings"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/services/mg"

func MailgunUploadHook(runtime *net.RequestRuntime) error {
	secret := os.Getenv("MAILGUN_WEBHOOK_SECRET")
	key := os.Getenv("MAILGUN_API_KEY")
	query := runtime.URL.Query()

	if v, ok := query["secret"]; !ok || len(v) != 1 || v[0] != secret {
		runtime.Debugf("invalid attempt on mailgun webhook: %s", query["secret"])
		return runtime.AddError(fmt.Errorf("UNAUTHORIZED"))
	}

	body, err := net.ParseBody(runtime.Request, 150000000)

	if err != nil {
		return err
	}

	client := mg.Client{key}

	location := body.Get("message-url")

	if valid := len(location) >= 2; valid != true {
		return fmt.Errorf("BAD_MESSAGE_URL")
	}

	message, err := client.Retreive(location)

	var processor mg.ActivityProcessor
	start := strings.Split(message.Subject, ":")[0]

	runtime.Debugf("received mailgun message: subject[\"%s\"] sender[%s]", message.Subject, message.From)

	switch start {
	case "image", "photo":
		processor = &mg.ImageProcessor{runtime.DB, runtime.Photos(), key}
	default:
		return fmt.Errorf("INVALID_SUBJECT_LINE")
	}

	outputs := make(chan mg.ProcessedItem)
	go processor.Process(&message, outputs)

	for item := range outputs {
		if err := item.Error; err != nil {
			runtime.Debugf("failed processing item: %s", err.Error())
			continue
		}

		message := item.Message
		runtime.Debugf("successfully activity for object \"%s\", publishing...", message.Object.Url())
		runtime.Publish(message)
	}

	return nil
}
