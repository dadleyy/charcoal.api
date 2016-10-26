package routes

import "fmt"
import "image"
import "bytes"
import "io/ioutil"

import _ "image/jpeg"
import _ "image/png"

import "github.com/albrow/forms"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/activity"

const MIN_PHOTO_LABEL_LENGTH = 2
const MAX_PHOTO_WIDTH = 2048
const MAX_PHOTO_HEIGHT = 2048

const defaultMemory = 32 << 20

func CreatePhoto(runtime *net.RequestRuntime) error {
	data, err := forms.Parse(runtime.Request)

	// bad form file - error out
	if err != nil {
		return runtime.AddError(err)
	}

	validator := data.Validator()
	validator.Require("label")

	if validator.HasErrors() {
		return runtime.AddError(fmt.Errorf("MISSING_LABEL"))
	}

	file := data.GetFile("photo")

	if file == nil {
		return runtime.AddError(fmt.Errorf("BAD_FILE"))
	}

	content, err := file.Open()
	mime := file.Header.Get("Content-Type")

	if err != nil || len(mime) == 0 {
		return runtime.AddError(fmt.Errorf("BAD_FILE"))
	}

	buffer, err := ioutil.ReadAll(content)

	// bad form file - error out
	if err != nil || len(buffer) == 0 {
		return runtime.AddError(err)
	}

	label := data.Get("label")

	// bad label - error out
	if len(label) < MIN_PHOTO_LABEL_LENGTH {
		return runtime.AddError(fmt.Errorf("MISSING_LABEL"))
	}

	reader := bytes.NewReader(buffer)

	// attempt to decode the image and get it's dimensions
	image, _, err := image.DecodeConfig(reader)

	if err != nil {
		runtime.Errorf("unable to decode image: %s", err.Error())
		return runtime.AddError(err)
	}

	width, height := image.Width, image.Height

	if width == 0 || height == 0 || width > MAX_PHOTO_WIDTH || height > MAX_PHOTO_HEIGHT {
		runtime.Debugf("bad image sizes")
		return runtime.AddError(fmt.Errorf("BAD_IMAGE_SIZES"))
	}

	// attempt to use the runtime to persist the file uploaded
	ormfile, err := runtime.PersistFile(buffer, mime)

	if err != nil {
		runtime.Debugf("failed persisting: %s", err.Error())
		return runtime.AddError(err)
	}

	runtime.Debugf("uploaded \"%s\" (width: %d, height: %d, size)", ormfile.Key, width, height)

	// with the persisted file, create the new photo object that will be saved to the orm
	photo := models.Photo{
		Label:  label,
		File:   ormfile.ID,
		Width:  width,
		Height: height,
	}

	if runtime.User.ID >= 1 {
		runtime.Debugf("associating user #%d with photo \"%s\"", runtime.User.ID, photo.Label)
		photo.Author.Scan(runtime.User.ID)
	}

	// attempt to create the photo record in the database
	if err := runtime.Database().Create(&photo).Error; err != nil {
		return runtime.AddError(err)
	}

	// let the file record know that it is owned
	if err := runtime.Database().Model(&ormfile).Update("status", "OWNED").Error; err != nil {
		return runtime.AddError(err)
	}

	runtime.AddResult(photo.Public())

	if runtime.User.ID >= 1 {
		// publish this event to the activity stream
		runtime.Publish(activity.Message{&runtime.User, &photo, "created"})
		return nil
	}

	runtime.Debugf("publishing photo upload w/ client not user")
	runtime.Publish(activity.Message{&runtime.Client, &photo, "created"})

	return nil
}

func ViewPhoto(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_PHOTO_ID"))
	}

	var photo models.Photo

	if err := runtime.Database().First(&photo, id).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	var file models.File

	if err := runtime.Database().First(&file, photo.File).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	url, err := runtime.DownloadUrl(&file)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_DOWNLOAD_URL"))
	}

	runtime.Proxy(url)
	return nil
}

func FindPhotos(runtime *net.RequestRuntime) error {
	var results []models.Photo
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results, runtime.Database())

	if err != nil {
		runtime.Debugf("bad photo lookup: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.AddResult(item.Public())
	}

	runtime.SetMeta("total", total)
	return nil
}
