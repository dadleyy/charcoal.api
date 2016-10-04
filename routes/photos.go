package routes

import "fmt"
import "strings"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

const MIN_PHOTO_LABEL_LENGTH int = 4
const MIN_PHOTO_LABEL_MESSAGE string = "must provide a \"label\" at least %d characters long"

func CreatePhoto(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Miritos)

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

	runtime.Logger().Infof("successfully saved file \"%s\", updating photo orm", ormfile.Key)

	photo := models.Photo{
		Label: label,
		File: ormfile.ID,
	}

	if err := runtime.DB.Create(&photo).Error; err != nil {
		return runtime.ErrorOut(err)
	}

	runtime.Result(&photo)

	return nil
}

func UpdatePhoto(ectx echo.Context) error {
	return nil
}

func FindPhotos(ectx echo.Context) error {
	return nil
}
