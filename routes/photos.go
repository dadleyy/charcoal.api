package routes

import "github.com/labstack/echo"

import "github.com/sizethree/miritos.api/context"
import "github.com/sizethree/miritos.api/services"

func CreatePhoto(ectx echo.Context) error {
	runtime, _ := ectx.(*context.Miritos)

	file, err := runtime.FormFile("photo")
	label := runtime.FormValue("label")

	if err != nil {
		runtime.Error(err)
		return nil
	}

	source, err := file.Open()

	if err != nil {
		runtime.Error(err)
		return nil
	}

	result, err := services.UploadFile(source)

	if err != nil {
		runtime.Error(err)
		return nil
	}

	defer source.Close()

	runtime.Logger().Infof("creating photo \"%s\"", label)
	runtime.Result(&result)

	return nil
}

func UpdatePhoto(ectx echo.Context) error {
	return nil
}

func FindPhotos(ectx echo.Context) error {
	return nil
}
