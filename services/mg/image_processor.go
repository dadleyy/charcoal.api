package mg

import "fmt"
import "sync"
import "image"
import "bytes"
import "strings"
import "net/http"
import "net/mail"
import "io/ioutil"
import _ "image/gif"
import _ "image/jpeg"
import _ "image/png"

import "github.com/dadleyy/charcoal.api/db"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/services"

const ErrUnkownUser = "unkown user"

type ImageProcessor struct {
	*db.Connection
	Photos services.PhotoSaver
	ApiKey string
}

type semaphore chan struct{}

type processedPhoto struct {
	ProcessedItem
	Photo models.Photo
}

func (processor *ImageProcessor) Process(message *Message, outputs chan ProcessedItem) {
	// close once we're done
	defer close(outputs)

	if message == nil {
		outputs <- ProcessedItem{Error: fmt.Errorf("INVALID_MESSAGE")}
		return
	}

	parts := strings.SplitN(message.Subject, ":", 2)

	if len(parts) < 2 || len(parts[1]) == 0 {
		outputs <- ProcessedItem{Error: fmt.Errorf("BAD_SUBJECT_LINE")}
		return
	}

	caption := strings.TrimSpace(parts[1])

	email, err := mail.ParseAddress(message.From)

	if err != nil {
		outputs <- ProcessedItem{Error: err}
		return
	}

	var author models.User

	if err := processor.Where("email = ?", email.Address).Find(&author).Error; err != nil {
		outputs <- ProcessedItem{Error: fmt.Errorf(ErrUnkownUser)}
		return
	}

	images := message.Images()

	if len(images) == 0 {
		return
	}

	waitlist := make(semaphore, 10)
	var deferred sync.WaitGroup
	processed := make(chan processedPhoto)

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
			outputs <- result.ProcessedItem
			continue
		}

		// set the photo's author
		result.Photo.Author.Scan(author.ID)
		processor.Save(&result.Photo)
		result.Message.Actor = &author

		outputs <- result.ProcessedItem
	}
}

func (p *ImageProcessor) single(img ContentItem, caption string, out chan processedPhoto, queue semaphore, def *sync.WaitGroup) {
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
		out <- processedPhoto{ProcessedItem: ProcessedItem{Error: err, Item: img}}
		return
	}

	get := &http.Client{}

	// set our mailgun api key
	req.SetBasicAuth("api", p.ApiKey)

	// execute our get request
	response, err := get.Do(req)

	// if we failed opening a request, be done
	if err != nil {
		out <- processedPhoto{ProcessedItem: ProcessedItem{Error: err, Item: img}}
		return
	}

	// once we're done with this function, clear out the body
	defer response.Body.Close()

	buffer, err := ioutil.ReadAll(response.Body)

	// if we failed opening a request, be done
	if err != nil {
		out <- processedPhoto{ProcessedItem: ProcessedItem{Error: err, Item: img}}
		return
	}

	reader := bytes.NewReader(buffer)

	// attempt to decode the image and get it's dimensions
	config, _, err := image.DecodeConfig(reader)

	// if we failed opening a request, be done
	if err != nil || config.Height == 0 || config.Width == 0 {
		out <- processedPhoto{ProcessedItem: ProcessedItem{Error: fmt.Errorf("BAD_IMAGE"), Item: img}}
		return
	}

	photo := models.Photo{Width: config.Width, Height: config.Height, Label: caption}

	if err := p.Photos.Persist(buffer, img.ContentType, &photo); err != nil {
		out <- processedPhoto{ProcessedItem: ProcessedItem{Error: err, Item: img}}
		return
	}

	message := activity.Message{&models.Client{}, &photo, "created"}
	out <- processedPhoto{Photo: photo, ProcessedItem: ProcessedItem{Error: nil, Item: img, Message: message}}
}
