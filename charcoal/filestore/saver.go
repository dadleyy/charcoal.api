package filestore

import "github.com/dadleyy/charcoal.api/charcoal/models"

// FileSaver defines an interface that allows for both uploading a file as well as the ability to download them.
type FileSaver interface {
	Upload([]byte, string) (models.File, error)
	DownloadURL(*models.File) (string, error)
}
