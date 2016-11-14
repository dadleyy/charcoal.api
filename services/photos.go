package services

import "github.com/sizethree/miritos.api/db"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/filestore"

type PhotoSaver struct {
	*db.Connection
	filestore.FileSaver
}

func (saver *PhotoSaver) Persist(buffer []byte, mime string, out *models.Photo) error {
	ofile, err := saver.Upload(buffer, mime)

	if err != nil {
		return err
	}

	if err := saver.Create(&ofile).Error; err != nil {
		return err
	}

	// assign the newly created file to this photo
	out.File = ofile.ID

	// attempt to create the photo record in the database
	if err := saver.Create(out).Error; err != nil {
		return err
	}

	// let the file record know that it is owned
	if err := saver.Model(&ofile).Update("status", "OWNED").Error; err != nil {
		return err
	}

	return nil
}

func (saver *PhotoSaver) Destroy(photo *models.Photo) error {
	cursor := saver.Model(&models.File{}).Where("id = ?", photo.File)

	if err := cursor.Update("status", "ABANDONDED").Error; err != nil {
		return err
	}

	return saver.Unscoped().Delete(photo).Error
}
