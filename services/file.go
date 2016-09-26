package services

import "os"
import "mime/multipart"

import "github.com/pborman/uuid"
import _ "github.com/aws/aws-sdk-go"
import "github.com/aws/aws-sdk-go/aws/credentials"

import "github.com/sizethree/miritos.api/models"

func UploadFile(Target multipart.File) (models.File, error) {
	key := uuid.NewRandom()

	access_key := os.Getenv("AWS_CREDENTIALS_KEY")
	access_id := os.Getenv("AWS_CREDENTIALS_ID")
	access_token := os.Getenv("AWS_CREDENTIALS_TOKEN")

	creds := credentials.NewStaticCredentials(access_id, access_key, access_token)

	if _, err := creds.Get(); err != nil {
		return models.File{}, err
	}

	return models.File{Key: key.String()}, nil
}
