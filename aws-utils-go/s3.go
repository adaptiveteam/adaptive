package aws_utils_go

import (
	"bytes"
	"fmt"
	"github.com/adaptiveteam/adaptive/core-utils-go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type S3Bucket struct {
	Name         *string    `json:"name"`
	CreationDate *time.Time `json:"creation_date"`
}

type S3Request struct {
	svc *s3.S3
	log *logger.Logger
}

func NewS3(region, endpoint, namespace string) *S3Request {
	session, config := sess(region, endpoint)
	return &S3Request{
		svc: s3.New(session, config),
		log: logger.WithNamespace(fmt.Sprintf("adaptive.s3.%s", namespace)),
	}
}

func (s *S3Request) errorLog(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case s3.ErrCodeBucketAlreadyExists:
			s.log.Error(s3.ErrCodeBucketAlreadyExists, aerr.Error())
		case s3.ErrCodeBucketAlreadyOwnedByYou:
			s.log.Error(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
		case s3.ErrCodeNoSuchBucket:
			s.log.Error(s3.ErrCodeNoSuchBucket, aerr.Error())
		case s3.ErrCodeNoSuchKey:
			s.log.Error(s3.ErrCodeNoSuchKey, aerr.Error())
		case s3.ErrCodeNoSuchUpload:
			s.log.Error(s3.ErrCodeNoSuchUpload, aerr.Error())
		case s3.ErrCodeObjectAlreadyInActiveTierError:
			s.log.Error(s3.ErrCodeObjectAlreadyInActiveTierError, aerr.Error())
		case s3.ErrCodeObjectNotInActiveTierError:
			s.log.Error(s3.ErrCodeObjectNotInActiveTierError, aerr.Error())
		default:
			s.log.Error(aerr.Error())
		}
	} else {
		s.log.Error(err.Error())
	}
}

func (s *S3Request) EnsureBucketExists(bucket string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}
	print(input, true)
	result, err2 := s.svc.CreateBucket(input)
	if err2 != nil {
		if aerr, ok := err2.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				// ErrCodeBucketAlreadyOwnedByYou means bucket has already been created by you,
				// So this error code should deal with success.
				s.log.Warnf("Bucket %s has already been created.\n", bucket)
				err2 = nil
			default:
				s.log.Warnf("Bucket %s creation error %v; CreateBucketOutput: %v\n", bucket, err2, result)
			}
		}
	} else {
		print(result, true)
		s.log.Infof("Bucket \"%s\" set up successfully.\n", bucket)	
	}
	return err2
}

func (s *S3Request) ListBuckets() (buckets []S3Bucket, err error) {
	result, err2 := s.svc.ListBuckets(&s3.ListBucketsInput{})
	err = err2
	if err != nil {
		s.errorLog(err)
	} else {
		for _, each := range result.Buckets {
			buckets = append(buckets, S3Bucket{Name: each.Name, CreationDate: each.CreationDate})
		}
	}
	return
}

func (s *S3Request) DeleteBucket(bucket string) error {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}
	_, err2 := s.svc.DeleteBucket(input)
	if err2 != nil {
		s.errorLog(err2)
	}
	return err2
}

func (s *S3Request) AddFile(filepath, bucket, key string) (err error) {
	// Open the file for use
	file, err2 := os.Open(filepath)
	err = err2
	if err2 == nil {
		defer file.Close()

		// Get file size and read the file content into a buffer
		fileInfo, _ := file.Stat()
		size := fileInfo.Size()
		buffer := make([]byte, size)
		file.Read(buffer)

		input := &s3.PutObjectInput{
			Bucket:             aws.String(bucket),
			Key:                aws.String(key),
			Body:               bytes.NewReader(buffer),
			ContentLength:      aws.Int64(size),
			ContentType:        aws.String(http.DetectContentType(buffer)),
			ContentDisposition: aws.String("attachment"),
		}
		_, err = s.svc.PutObject(input)
	}
	if err != nil {
		s.errorLog(err)
	}
	return 
}

func (s *S3Request) GetObject(bucket, key string) (body []byte, err error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	res, err2 := s.svc.GetObject(input)
	err = err2
	if err == nil {
		body, err = ioutil.ReadAll(res.Body)
	}
	if err != nil {
		s.errorLog(err)
	}
	return
}

func (s *S3Request) DeleteObject(bucket, key string) (*bool, error) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	res, err2 := s.svc.DeleteObject(input)
	if err2 != nil {
		s.errorLog(err2)
		return nil, err2
	}
	return res.DeleteMarker, nil
}

func (s *S3Request) ObjectExists(bucket, key string) bool {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err2 := s.svc.GetObject(input)
	return err2 == nil
	// if err2 != nil {
	// 	if aerr, ok := err2.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		case s3.ErrCodeNoSuchKey:
	// 			return false
	// 		default:
	// 			return false
	// 		}
	// 	} else {
	// 		return false
	// 	}
	// }
	// return true
}
