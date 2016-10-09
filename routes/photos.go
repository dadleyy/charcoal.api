package routes

import "fmt"
import "strings"
import "net/http"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/activity"

const MIN_PHOTO_LABEL_LENGTH int = 4
const MIN_PHOTO_LABEL_MESSAGE string = "must provide a \"label\" at least %d characters long"

func CreatePhoto(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)

	header, err := runtime.FormFile("photo")

	// bad form file - error out
	if err != nil {
		return runtime.ErrorOut(err)
	}

	label := runtime.FormValue("label")

	// bad label - error out
	if len(label) < MIN_PHOTO_LABEL_LENGTH {
		return runtime.ErrorOut(fmt.Errorf(MIN_PHOTO_LABEL_MESSAGE, MIN_PHOTO_LABEL_LENGTH))
	}

	// make sure the mime type detected is an image
	mime, ok := header.Header["Content-Type"]

	if ok != true || len(mime) != 1 {
		return runtime.ErrorOut(fmt.Errorf("unable to look up file type from multipart header"))
	}

	if isimg := strings.HasPrefix(mime[0], "image/"); isimg != true {
		return runtime.ErrorOut(fmt.Errorf("bad mime type"))
	}

	source, err := header.Open()
	defer source.Close()

	if err != nil {
		return runtime.ErrorOut(err)
	}

	ormfile, err := runtime.PersistFile(source, mime[0])

	if err != nil {
		runtime.Logger().Error(err)
		return runtime.ErrorOut(err)
	}

	runtime.Logger().Debugf("successfully saved file \"%s\", updating photo orm", ormfile.Key)

	photo := models.Photo{
		Label: label,
		File: ormfile.ID,
	}

	if runtime.User.ID >= 1 {
		runtime.Logger().Debugf("associating user #%d with photo \"%s\"", runtime.User.ID, photo.Label)
		photo.Author.Scan(runtime.User.ID)
	}

	if err := runtime.DB.Create(&photo).Error; err != nil {
		return runtime.ErrorOut(err)
	}

	if err := runtime.DB.Model(&ormfile).Update("status", "OWNED").Error; err != nil {
		return runtime.ErrorOut(err)
	}

	runtime.Publish(activity.Message{&runtime.User, &photo, "created"})
	runtime.AddResult(&photo)

	return nil
}

func ViewPhoto(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Runtime)
	id, err := runtime.ParamInt("id")

	if err != nil {
		return runtime.ErrorOut(fmt.Errorf("BAD_PHOTO_ID"))
	}

	var photo models.Photo

	if err := runtime.DB.First(&photo, id).Error; err != nil {
		return runtime.ErrorOut(fmt.Errorf("NOT_FOUND"))
	}

	var file models.File

	if err := runtime.DB.First(&file, photo.File).Error; err != nil {
		return runtime.ErrorOut(fmt.Errorf("NOT_FOUND"))
	}

	url, err := runtime.FS.DownloadUrl(&file)

	if err != nil {
		return runtime.ErrorOut(fmt.Errorf("BAD_DOWNLOAD_URL"))
	}

	resp, err := http.Get(url)

	if err != nil {
		runtime.Logger().Debugf("unable to download file: %s", err.Error())
		return runtime.ErrorOut(fmt.Errorf("BAD_DOWNLOAD_URL"))
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
		return runtime.ErrorOut(fmt.Errorf("BAD_QUERY"))
	}

	for _, item := range results {
		runtime.AddResult(item)
	}

	runtime.AddMeta("total", total)
	return nil
}
