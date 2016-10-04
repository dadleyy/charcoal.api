package context

import "io"
import "github.com/sizethree/miritos.api/models"

type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type FileSaver interface {
	Upload(File, string) (models.File, error)
}
