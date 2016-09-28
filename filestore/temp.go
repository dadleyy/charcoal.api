package filestore

import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

type TempStore struct {}

func (store TempStore) Upload(target context.File) (models.File, error) {
	var result models.File
	return result, nil
}
