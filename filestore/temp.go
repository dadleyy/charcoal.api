package filestore

import "fmt"
import "io/ioutil"
import "github.com/pborman/uuid"
import "github.com/dadleyy/charcoal.api/models"

type TempStore struct {
	Root string
}

func (store TempStore) DownloadUrl(target *models.File) (string, error) {
	return "", nil
}

func (store TempStore) Upload(buffer []byte, mime string) (models.File, error) {
	var result models.File
	photoid := uuid.NewRandom()
	path := fmt.Sprintf("%s/%s", store.Root, photoid.String())

	if err := ioutil.WriteFile(path, buffer, 0644); err != nil {
		return result, err
	}

	result = models.File{
		Key:    photoid.String(),
		Mime:   "poop",
		Status: "TEMPORARY",
	}

	return result, nil
}
