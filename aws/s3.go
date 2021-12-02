package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client interface {
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func CreateS3Bucket(client S3Client, bucketName *string, accountId *string) *string {

	if *bucketName == "" {
		*bucketName = fmt.Sprintf("sagemaker-edgemanager-%s", *accountId)
	}

	_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintUsWest2,
		},
	})
	if err != nil {
		var bne *types.BucketAlreadyOwnedByYou
		var be *types.BucketAlreadyExists
		if errors.As(err, &bne) || errors.As(err, &be) {
			return bucketName
		}
		log.Fatal("Error: ", err)
	}

	return bucketName
}

func DownloadFileFromS3ToPath(client S3Client, bucketName *string, key *string, filePath *string) *string {
	downloader := manager.NewDownloader(client)
	os.MkdirAll(filepath.Dir(*filePath), os.ModePerm)
	fd, err := os.Create(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{
		Bucket: bucketName,
		Key:    key,
	})

	if err != nil {
		log.Fatal(err)
	}
	return filePath
}

func DownloadFileFromS3(client S3Client, bucketName *string, key *string) *string {
	filePath := fmt.Sprintf("/tmp/%s", *key)
	return DownloadFileFromS3ToPath(client, bucketName, key, &filePath)
}

func ListBucket(client S3Client, bucketName *string, prefix *string) *s3.ListObjectsOutput {

	listObjectsInput := &s3.ListObjectsInput{
		Bucket: bucketName,
		Prefix: prefix,
	}
	output, err := client.ListObjects(context.TODO(), listObjectsInput)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return output
}
