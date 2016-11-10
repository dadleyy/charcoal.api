package routes

import "fmt"
import "bytes"
import "image"
import "strings"
import "io/ioutil"
import _ "image/jpeg"
import _ "image/png"

import "github.com/albrow/forms"

import "github.com/sizethree/miritos.api/net"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/activity"

func FindInstagramPosts(runtime *net.RequestRuntime) error {
	var accounts []models.InstagramPhoto
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&accounts, runtime.Database())

	if err != nil {
		return err
	}

	for _, account := range accounts {
		runtime.AddResult(account.Public())
	}

	runtime.SetMeta("toal", total)

	return nil
}

func CreateInstagramPost(runtime *net.RequestRuntime) error {
	data, err := forms.Parse(runtime.Request)

	// bad form file - error out
	if err != nil {
		runtime.Debugf("unable to parse instagram post body data")
		return runtime.AddError(err)
	}

	validator := data.Validator()

	validator.Require("id")
	validator.Require("caption")
	validator.Require("owner")
	validator.MinLength("id", 8)
	validator.MinLength("owner", 8)

	if validator.HasErrors() {
		runtime.Debugf("bad form: %s", strings.Join(validator.Messages(), " | "))
		return runtime.AddError(fmt.Errorf("BAD_DATA"))
	}

	gramid := data.Get("id")
	file := data.GetFile("photo")
	dupes := 0

	ig := models.InstagramPhoto{InstagramID: gramid}

	runtime.Debugf("checking for duplicate instagram id: %s", gramid)

	check := runtime.Database().Model(&ig).Where("instagram_id = ?", gramid)

	if err := check.Count(&dupes).Error; err != nil || dupes > 0 {
		runtime.Debugf("duplicate instagram record: %s", gramid)
		return runtime.AddError(fmt.Errorf("DUPLICATE_RECORD"))
	}

	if file == nil {
		runtime.Debugf("no \"photo\" value found in post form")
		return runtime.AddError(fmt.Errorf("BAD_FILE"))
	}

	mime := file.Header.Get("Content-Type")

	if len(mime) == 0 {
		return runtime.AddError(fmt.Errorf("BAD_FILE_CONTENT_TYPE"))
	}

	content, err := file.Open()
	defer content.Close()

	if err != nil {
		runtime.Debugf("unable to open photo file: %s", err.Error())
		return runtime.AddError(fmt.Errorf("BAD_FILE"))
	}

	buffer, err := ioutil.ReadAll(content)

	// bad form file - error out
	if err != nil || len(buffer) == 0 {
		return runtime.AddError(err)
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

	photos := runtime.Photos()

	// with the persisted file, create the new photo object that will be saved to the orm
	photo := models.Photo{Width: width, Height: height}

	// attempt to use the runtime to persist the file uploaded

	if err := photos.Persist(buffer, mime, &photo); err != nil {
		runtime.Debugf("failed persisting: %s", err.Error())
		return runtime.AddError(err)
	}

	ig = models.InstagramPhoto{
		Photo:       photo.ID,
		Owner:       data.Get("owner"),
		Caption:     data.Get("caption"),
		InstagramID: data.Get("id"),
	}

	ig.Client.Scan(runtime.Client.ID)

	// attempt to create the photo record in the database
	if err := runtime.Database().Create(&ig).Error; err != nil {
		runtime.Debugf("failed saving instagram: %s", err.Error())
		runtime.Database().Unscoped().Delete(&photo)
		return runtime.AddError(err)
	}

	runtime.Debugf("uploaded \"%s\" (width: %d, height: %d, size)", photo.ID, width, height)
	runtime.AddResult(ig.Public())
	runtime.Publish(activity.Message{&runtime.Client, &ig, "created"})
	return nil
}
