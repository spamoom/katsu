package s3

import (
	"bytes"
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

func PutFile(bucketName string, fileName string, fileContents []byte) error {
	ctx := context.Background()

	client := s3.NewFromConfig(katsuAws.GetConfig())

	reader := bytes.NewReader(fileContents)

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   reader,
	}

	_, err := client.PutObject(ctx, input)

	return err
}

func ListFiles(bucketName string) ([]string, error) {
	ctx := context.Background()

	client := s3.NewFromConfig(katsuAws.GetConfig())

	objects, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return nil, err
	}

	var files []string

	for _, object := range objects.Contents {
		files = append(files, *object.Key)
	}

	return files, nil
}
