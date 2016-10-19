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

import "github.com/sizethree/miritos.api/models"

const ERR_BAD_IMAGE_TYPE = "BAD_IMAGE_TYPE"
const ERR_BAD_IMAGE_UUID = "BAD_UUID_GENERATED"
const ERR_BAD_S3_RESPONSE = "BAD_S3_RESPONSE"

type S3FileStore struct {
	AccessID    string
	AccessKey   string
	AccessToken string
}

func (store S3FileStore) DownloadUrl(target *models.File) (string, error) {
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

func (store S3FileStore) Upload(buffer []byte, mime string) (models.File, error) {
	photoid := uuid.NewRandom()
	var result models.File
	size := int64(len(buffer))

	reader := bytes.NewReader(buffer)

	if strings.TrimSpace(photoid.String()) == "" {
		return result, errors.New(ERR_BAD_IMAGE_UUID)
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
		return result, errors.New(ERR_BAD_S3_RESPONSE)
	}

	result = models.File{
		Key:    photoid.String(),
		Mime:   mime,
		Status: "TEMPORARY",
	}

	return result, nil
}
