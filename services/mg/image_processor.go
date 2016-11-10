package mg

import "fmt"
import "sync"
import "image"
import "bytes"
import "strings"
import "net/http"
import "net/mail"
import "io/ioutil"
import _ "image/jpeg"
import _ "image/png"

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/activity"
import "github.com/sizethree/miritos.api/services"

type ImageProcessor struct {
	*db.Connection
	Photos services.PhotoSaver
	ApiKey string
}

type ProcessedItem struct {
	Error   error
	Message activity.Message
	Item    ContentItem
}

type semaphore chan struct{}
type outputs chan ProcessedItem

func (processor *ImageProcessor) Process(message *Message) ([]activity.Message, error) {
	activities := make([]activity.Message, 0)

	if message == nil {
		return activities, fmt.Errorf("INVALID_MESSAGE")
	}

	parts := strings.SplitN(message.Subject, ":", 2)

	if len(parts) < 2 || len(parts[1]) == 0 {
		return activities, fmt.Errorf("BAD_CAPTION")
	}

	caption := strings.TrimSpace(parts[1])

	email, err := mail.ParseAddress(message.From)

	if err != nil {
		return activities, err
	}

	var author models.User

	if err := processor.Where("email = ?", email.Address).Find(&author).Error; err != nil {
		return activities, err
	}

	images := message.Images()

	if len(images) != 1 {
		return activities, fmt.Errorf("BAD_IMAGES")
	}

	waitlist := make(semaphore, 10)
	var deferred sync.WaitGroup
	processed := make(outputs)

	for _, image := range images {
		deferred.Add(1)
		go processor.single(image, caption, processed, waitlist, &deferred)
	}

	finish := func() {
		deferred.Wait()
		close(processed)
	}

	go finish()

	for result := range processed {
		if result.Error != nil {
			continue
		}

		msg := result.Message
		msg.Actor = author

		activities = append(activities, msg)
	}

	return activities, nil
}

func (p *ImageProcessor) single(img ContentItem, caption string, out outputs, queue semaphore, def *sync.WaitGroup) {
	// push into our semaphore which blocks until there is an open "slot"
	queue <- struct{}{}

	finish := func() {
		def.Done()
		<-queue
	}

	// once we're done in here, let the wait group know
	defer finish()

	// start our http request to the item's url
	req, err := http.NewRequest("GET", img.Url, nil)

	// if we failed opening a request, be done
	if err != nil {
		out <- ProcessedItem{Error: fmt.Errorf("UNABLE_TO_LOAD"), Item: img}
		return
	}

	get := &http.Client{}

	// set our mailgun api key
	req.SetBasicAuth("api", p.ApiKey)

	// execute our get request
	response, err := get.Do(req)

	// once we're done with this function, clear out the body
	defer response.Body.Close()

	buffer, err := ioutil.ReadAll(response.Body)

	// if we failed opening a request, be done
	if err != nil {
		out <- ProcessedItem{Error: fmt.Errorf("UNABLE_TO_LOAD"), Item: img}
		return
	}

	reader := bytes.NewReader(buffer)

	// attempt to decode the image and get it's dimensions
	config, _, err := image.DecodeConfig(reader)

	// if we failed opening a request, be done
	if err != nil || config.Height == 0 || config.Width == 0 {
		out <- ProcessedItem{Error: fmt.Errorf("UNABLE_TO_LOAD"), Item: img}
		return
	}

	photo := models.Photo{Width: config.Width, Height: config.Height, Label: caption}

	if err := p.Photos.Persist(buffer, img.ContentType, &photo); err != nil {
		out <- ProcessedItem{Error: err, Item: img}
		return
	}

	message := activity.Message{&models.Client{}, &photo, "created"}
	out <- ProcessedItem{Error: nil, Item: img, Message: message}
}
