package s3

import (
	"context"
	"fmt"
	"io"
)
import "github.com/aws/aws-sdk-go-v2/service/s3"

type GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func GetObjectFromS3(ctx context.Context, api GetObjectAPI, bucket, key string) ([]byte, error) {
	object, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("s3.GetObjectFromS3: error closing body:%+v", err)
		}
	}(object.Body)

	return io.ReadAll(object.Body)
}
