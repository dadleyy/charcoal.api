package filestore

import "os"
import "fmt"
import "time"
import "bytes"
import "errors"
import "strings"
import "github.com/pborman/uuid"
import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/service/s3"
import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/aws/credentials"

import "github.com/dadleyy/charcoal.api/charcoal/defs"
import "github.com/dadleyy/charcoal.api/charcoal/models"

// S3FileStore implements the Saver interface with an aws/s3 backend.
type S3FileStore struct {
	AccessID    string
	AccessKey   string
	AccessToken string
}

// DownloadURL returns a generate s3 url w/ dwonload permissions.
func (store S3FileStore) DownloadURL(target *models.File) (string, error) {
	var creds *credentials.Credentials

	switch {
	case len(store.AccessID) >= 1:
		creds = credentials.NewStaticCredentials(store.AccessID, store.AccessKey, store.AccessToken)
	default:
		creds = credentials.NewEnvCredentials()
	}

	_, err := creds.Get()

	if err != nil {
		return "", err
	}

	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET")
	storepath := os.Getenv("AWS_S3_STORAGE_PATH")

	if len(region) == 0 {
		region = "us-east-1"
	}

	config := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	client := s3.New(session.New(), config)
	path := fmt.Sprintf("%s/%s", storepath, target.Key)

	req, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})

	urlstr, err := req.Presign(time.Minute * 600)

	if err != nil {
		return "", err
	}

	return urlstr, nil
}

// Upload accepts a slice of bytes and a file type and returns a saved file model.
func (store S3FileStore) Upload(buffer []byte, mime string) (models.File, error) {
	photoid := uuid.NewRandom()
	var result models.File
	size := int64(len(buffer))

	reader := bytes.NewReader(buffer)

	if strings.TrimSpace(photoid.String()) == "" {
		return result, errors.New(defs.ErrBadImageUuid)
	}

	var creds *credentials.Credentials

	switch {
	case len(store.AccessID) >= 1:
		creds = credentials.NewStaticCredentials(store.AccessID, store.AccessKey, store.AccessToken)
	default:
		creds = credentials.NewEnvCredentials()
	}

	_, err := creds.Get()

	if err != nil {
		return result, err
	}

	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET")
	storepath := os.Getenv("AWS_S3_STORAGE_PATH")

	if len(region) == 0 {
		region = "us-east-1"
	}

	config := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	client := s3.New(session.New(), config)
	path := fmt.Sprintf("%s/%s", storepath, photoid)

	request := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(path),
		Body:          reader,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(mime),
	}

	resp, err := client.PutObject(request)

	if err != nil {
		return result, err
	}

	if len(resp.String()) < 1 {
		return result, errors.New(defs.ErrBadS3Response)
	}

	result = models.File{
		Key:    photoid.String(),
		Mime:   mime,
		Status: "TEMPORARY",
	}

	return result, nil
}
