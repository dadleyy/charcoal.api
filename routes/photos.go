package routes

import "fmt"
import "errors"
import "github.com/labstack/echo"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

const MIN_PHOTO_LABEL_LENGTH int = 4
const MIN_PHOTO_LABEL_MESSAGE string = "must provide a \"label\" at least %d characters long"

func CreatePhoto(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Miritos)

	file, err := runtime.FormFile("photo")
	label := runtime.FormValue("label")

	if len(label) < MIN_PHOTO_LABEL_LENGTH {
		message := fmt.Sprintf(MIN_PHOTO_LABEL_MESSAGE, MIN_PHOTO_LABEL_LENGTH)
		return runtime.ErrorOut(errors.New(message))
	}

	if err != nil {
		return runtime.ErrorOut(err)
	}

	source, err := file.Open()

	if err != nil {
		return runtime.ErrorOut(err)
	}

	ormfile, err := runtime.PersistFile(source)

	if err != nil {
		return runtime.ErrorOut(err)
	}

	source.Close()
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
