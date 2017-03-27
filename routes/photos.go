package routes

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

func CreatePhoto(runtime *net.RequestRuntime) *net.ResponseBucket {
	data, err := forms.Parse(runtime.Request)

	// bad form file - error out
	if err != nil {
		return runtime.LogicError("bad-request")
	}

	validator := data.Validator()
	validator.Require("label")

	if validator.HasErrors() {
		return runtime.LogicError("bad-label")
	}

	file := data.GetFile("photo")

	if file == nil {
		return runtime.LogicError("missing-photo-file")
	}

	mime := file.Header.Get("Content-Type")

	if len(mime) == 0 {
		return runtime.LogicError("invalid-file")
	}

	content, err := file.Open()
	defer content.Close()

	if err != nil {
		runtime.Warnf("[create photo] unable to open uploaded multipart file")
		return runtime.LogicError("invalid-file")
	}

	buffer, err := ioutil.ReadAll(content)

	// bad form file - error out
	if err != nil || len(buffer) == 0 {
		return runtime.LogicError("bad-file")
	}

	label := data.Get("label")

	// bad label - error out
	if len(label) < MIN_PHOTO_LABEL_LENGTH {
		return runtime.LogicError("invalid-label")
	}

	reader := bytes.NewReader(buffer)

	// attempt to decode the image and get it's dimensions
	image, _, err := image.DecodeConfig(reader)

	if err != nil {
		runtime.Errorf("unable to decode image: %s", err.Error())
		return runtime.LogicError("invalid-image-type")
	}

	width, height := image.Width, image.Height

	if width == 0 || height == 0 || width > MAX_PHOTO_WIDTH || height > MAX_PHOTO_HEIGHT {
		runtime.Debugf("bad image sizes")
		return runtime.LogicError("invalid-image-size")
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
		runtime.Errorf("[create photo] failed persisting: %s", err.Error())
		return runtime.ServerError()
	}

	if runtime.User.ID >= 1 {
		runtime.Debugf("associating user #%d with photo \"%s\"", runtime.User.ID, photo.Label)
		photo.Author.Scan(runtime.User.ID)
		runtime.Save(&photo)
	}

	// publish this event to the activity stream
	if runtime.User.ID >= 1 {
		runtime.Publish(activity.Message{&runtime.User, &photo, activity.VerbCreated})
		return runtime.SendResults(1, photo)
	}

	runtime.Publish(activity.Message{&runtime.Client, &photo, "created"})
	return runtime.SendResults(1, photo)
}

func ViewPhoto(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("invalid-id")
	}

	var photo models.Photo

	if err := runtime.First(&photo, id).Error; err != nil {
		runtime.Errorf("[view photo] unable to find photo: %s", err.Error())
		return runtime.LogicError("invalid-id")
	}

	var file models.File

	if err := runtime.First(&file, photo.File).Error; err != nil {
		runtime.Errorf("[view photo] unable to find photo file: %s", err.Error())
		return runtime.ServerError()
	}

	url, err := runtime.DownloadUrl(&file)

	if err != nil {
		runtime.Errorf("[view photo] unable to find download url: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.Proxy(url)
}

func DestroyPhoto(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("not-found")
	}

	var photo models.Photo

	if err := runtime.First(&photo, id).Error; err != nil {
		runtime.Errorf("[destroy photo] error deleting photo: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	uman := services.UserManager{runtime.DB, runtime.Logger}
	admin := uman.IsAdmin(&runtime.User)

	// if we arent an admin, and the photo has an author, make sure its the current user
	if admin != true && photo.Author.Valid && photo.Author.Int64 != int64(runtime.User.ID) {
		runtime.Warnf("user %d attempted to delete photo %d w/o permission", runtime.User.ID, photo.Author.Int64)
		return runtime.LogicError("not-found")
	}

	photos := runtime.Photos()

	if err := photos.Destroy(&photo); err != nil {
		runtime.Errorf("[delete photo] unable to delete photo: %s", err.Error())
		return runtime.ServerError()
	}

	runtime.Publish(activity.Message{&runtime.Client, &photo, activity.VerbDeleted})
	return nil
}

func FindPhotos(runtime *net.RequestRuntime) *net.ResponseBucket {
	var results []models.Photo
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results)

	if err != nil {
		runtime.Errorf("[find photos] bad photo lookup: %s", err.Error())
		return runtime.LogicError("bad-request")
	}

	return runtime.SendResults(total, results)
}
