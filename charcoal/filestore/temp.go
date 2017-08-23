package filestore

import "fmt"
import "io/ioutil"
import "github.com/pborman/uuid"
import "github.com/dadleyy/charcoal.api/charcoal/models"

// TempStore is an implementation of the Saver interface using the local filesystem.
type TempStore struct {
	Root string
}

// DownloadURL returns a url in the api that can be used to download the file.
func (store TempStore) DownloadURL(target *models.File) (string, error) {
	return "", nil
}

// Upload saves the given byes to the file system.
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
