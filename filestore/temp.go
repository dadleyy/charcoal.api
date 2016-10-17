package filestore

import "github.com/sizethree/miritos.api/models"

type TempStore struct {}


func (store TempStore) DownloadUrl(target *models.File) (string, error) {
	return "", nil
}

func (store TempStore) Upload(target File, mime string) (models.File, error) {
	var result models.File
	return result, nil
}
