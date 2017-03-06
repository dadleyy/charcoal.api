package routes

import "fmt"
import "image"
import "bytes"
import "io/ioutil"

import _ "image/jpeg"
import _ "image/png"

import "github.com/albrow/forms"

import "github.com/dadleyy/charcoal.api/net"
import "github.com/dadleyy/charcoal.api/models"
import "github.com/dadleyy/charcoal.api/activity"
import "github.com/dadleyy/charcoal.api/services"

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

	mime := file.Header.Get("Content-Type")

	if len(mime) == 0 {
		return runtime.AddError(fmt.Errorf("BAD_FILE_CONTENT_TYPE"))
	}

	content, err := file.Open()
	defer content.Close()

	if err != nil {
		runtime.Debugf("unable to open uploaded multipart file")
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
	photos := runtime.Photos()

	// with the persisted file, create the new photo object that will be saved to the orm
	photo := models.Photo{
		Label:  label,
		Width:  width,
		Height: height,
	}

	if err := photos.Persist(buffer, mime, &photo); err != nil {
		runtime.Debugf("failed persisting: %s", err.Error())
		return runtime.AddError(err)
	}

	if runtime.User.ID >= 1 {
		runtime.Debugf("associating user #%d with photo \"%s\"", runtime.User.ID, photo.Label)
		photo.Author.Scan(runtime.User.ID)
		runtime.Save(&photo)
	}

	runtime.AddResult(photo.Public())

	if runtime.User.ID >= 1 {
		// publish this event to the activity stream
		runtime.Publish(activity.Message{&runtime.User, &photo, activity.VerbCreated})
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

	if err := runtime.First(&photo, id).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	var file models.File

	if err := runtime.First(&file, photo.File).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	url, err := runtime.DownloadUrl(&file)

	if err != nil {
		return runtime.AddError(fmt.Errorf("BAD_DOWNLOAD_URL"))
	}

	runtime.Proxy(url)
	return nil
}

func DestroyPhoto(runtime *net.RequestRuntime) error {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.AddError(fmt.Errorf("BAD_PHOTO_ID"))
	}

	var photo models.Photo

	if err := runtime.First(&photo, id).Error; err != nil {
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	uman := services.UserManager{runtime.DB, runtime.Logger}
	admin := uman.IsAdmin(&runtime.User)

	// if we arent an admin, and the photo has an author, make sure its the current user
	if admin != true && photo.Author.Valid && photo.Author.Int64 != int64(runtime.User.ID) {
		runtime.Debugf("user %d attempted to delete photo %d w/o permission", runtime.User.ID, photo.Author.Int64)
		return runtime.AddError(fmt.Errorf("NOT_FOUND"))
	}

	photos := runtime.Photos()

	if err := photos.Destroy(&photo); err != nil {
		return runtime.AddError(err)
	}

	runtime.Publish(activity.Message{&runtime.Client, &photo, activity.VerbDeleted})
	return nil
}

func FindPhotos(runtime *net.RequestRuntime) error {
	var results []models.Photo
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results)

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
