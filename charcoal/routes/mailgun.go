package routes

import "os"
import "strings"

import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/services/mg"

func MailgunUploadHook(runtime *net.RequestRuntime) *net.ResponseBucket {
	secret := os.Getenv("MAILGUN_WEBHOOK_SECRET")
	key := os.Getenv("MAILGUN_API_KEY")
	query := runtime.URL.Query()

	if v, ok := query["secret"]; !ok || len(v) != 1 || v[0] != secret {
		runtime.Warnf("[mailgun] invalid attempt on mailgun webhook: %s", query["secret"])
		return runtime.LogicError("unauthorized")
	}

	body, err := net.ParseBody(runtime.Request, 150000000)

	if err != nil {
		return runtime.LogicError("invalid-body")
	}

	client := mg.Client{key}

	location := body.Get("message-url")

	if valid := len(location) >= 2; valid != true {
		return runtime.LogicError("invalid-location")
	}

	message, err := client.Retreive(location)

	var processor mg.ActivityProcessor
	start := strings.Split(message.Subject, ":")[0]

	runtime.Debugf("[mailgun] received message: subject[\"%s\"] sender[%s]", message.Subject, message.From)

	switch start {
	case "image", "photo":
		processor = &mg.ImageProcessor{runtime.DB, runtime.Photos(), key}
	default:
		return runtime.LogicError("invalid-subject")
	}

	outputs := make(chan mg.ProcessedItem)
	go processor.Process(&message, outputs)

	for item := range outputs {
		if err := item.Error; err != nil {
			runtime.Debugf("failed processing item: %s", err.Error())
			continue
		}

		message := item.Message
		runtime.Debugf("[mailgun] successfully activity for object \"%s\", publishing...", message.Object.URL())
		runtime.Publish(message)
	}

	return nil
}
