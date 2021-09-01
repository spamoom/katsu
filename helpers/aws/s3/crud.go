package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	katsuAws "github.com/netsells/katsu/helpers/aws"
)

func GetFile(bucketName string, path string) (*s3.GetObjectOutput, error) {
	ctx := context.Background()

	client := s3.NewFromConfig(katsuAws.GetConfig())

	return client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	})
}
