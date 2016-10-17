package filestore

import "github.com/sizethree/miritos.api/models"

type FileSaver interface {
	Upload([]byte, string) (models.File, error)
	DownloadUrl(*models.File) (string, error)
}
