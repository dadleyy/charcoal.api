package routes

import "fmt"
import "image"
import "strings"
import "net/http"

import _ "image/jpeg"
import _ "image/png"

import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/activity"

const MIN_PHOTO_LABEL_LENGTH = 6
const MIN_PHOTO_LABEL_MESSAGE = "BAD_PHOTO_LABEL"
const MAX_PHOTO_WIDTH = 2048
const MAX_PHOTO_HEIGHT = 2048

func CreatePhoto(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)

	header, err := runtime.FormFile("photo")

	// bad form file - error out
	if err != nil {
		return err
	}

	label := runtime.FormValue("label")

	// bad label - error out
	if len(label) < MIN_PHOTO_LABEL_LENGTH {
		return fmt.Errorf(MIN_PHOTO_LABEL_MESSAGE)
	}

	// make sure the mime type detected is an image
	mime, ok := header.Header["Content-Type"]

	if ok != true || len(mime) != 1 {
		return fmt.Errorf("MISSING_CONTENT_TYPE")
	}

	if isimg := strings.HasPrefix(mime[0], "image/"); isimg != true {
		return fmt.Errorf("BAD_MIME_TYPE")
	}

	// open the multipart file and defer it's closing
	source, err := header.Open()
	defer source.Close()

	if err != nil {
		return err
	}

	// attempt to decode the image and get it's dimensions
	image, _, err := image.DecodeConfig(source)

	if err != nil {
		runtime.Logger().Errorf("unable to decode image: %s", err.Error())
		return err
	}

	width, height := image.Width, image.Height

	if width == 0 || height == 0 || width > MAX_PHOTO_WIDTH || height > MAX_PHOTO_HEIGHT {
		runtime.Logger().Debugf("bad image sizes")
		return fmt.Errorf("BAD_IMAGE_SIZES")
	}

	// attempt to use the runtime to persist the file uploaded
	ormfile, err := runtime.PersistFile(source, mime[0])

	if err != nil {
		runtime.Logger().Error(err)
		return err
	}

	runtime.Logger().Debugf("UPLOADED \"%s\" (width: %d, height: %d)", ormfile.Key, width, height)

	// with the persisted file, create the new photo object that will be saved to the orm
	photo := models.Photo{
		Label: label,
		File: ormfile.ID,
		Width: width,
		Height: height,
	}

	if runtime.User.ID >= 1 {
		runtime.Logger().Debugf("associating user #%d with photo \"%s\"", runtime.User.ID, photo.Label)
		photo.Author.Scan(runtime.User.ID)
	}

	// attempt to create the photo record in the database
	if err := runtime.DB.Create(&photo).Error; err != nil {
		return err
	}

	// let the file record know that it is owned
	if err := runtime.DB.Model(&ormfile).Update("status", "OWNED").Error; err != nil {
		return err
	}

	// publish this event to the activity stream
	runtime.Publish(activity.Message{&runtime.User, &photo, "created"})

	// add our result to the response
	runtime.AddResult(&photo)

	return nil
}

func ViewPhoto(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)
	id, err := runtime.ParamInt("id")

	if err != nil {
		return fmt.Errorf("BAD_PHOTO_ID")
	}

	var photo models.Photo

	if err := runtime.DB.First(&photo, id).Error; err != nil {
		return fmt.Errorf("NOT_FOUND")
	}

	var file models.File

	if err := runtime.DB.First(&file, photo.File).Error; err != nil {
		return fmt.Errorf("NOT_FOUND")
	}

	url, err := runtime.FS.DownloadUrl(&file)

	if err != nil {
		return fmt.Errorf("BAD_DOWNLOAD_URL")
	}

	resp, err := http.Get(url)

	if err != nil {
		runtime.Logger().Debugf("unable to download file: %s", err.Error())
		return fmt.Errorf("BAD_DOWNLOAD_URL")
	}

	defer resp.Body.Close()

	runtime.Stream(200, "image/png", resp.Body)
	return nil
}

func FindPhotos(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)

	var results []models.Photo
	blueprint := runtime.Blueprint()

	total, err := blueprint.Apply(&results, runtime.DB)

	if err != nil {
		runtime.Logger().Debugf("bad photo lookup: %s", err.Error())
		return fmt.Errorf("BAD_QUERY")
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.AddMeta("total", total)
	return nil
}
