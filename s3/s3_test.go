package s3

import (
	"awsInterfaces/s3/mocks"
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"strconv"
	"testing"
)

type mockGetObjectAPI func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)

func (m mockGetObjectAPI) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetObjectFromS3(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) GetObjectAPI
		bucket string
		key    string
		expect []byte
	}{
		{
			client: func(t *testing.T) GetObjectAPI {
				return mockGetObjectAPI(func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
					t.Helper()
					if params.Bucket == nil {
						t.Fatal("expect bucket to not be nil")
					}
					if e, a := "amzn-s3-demo-bucket", *params.Bucket; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}
					if params.Key == nil {
						t.Fatal("expect key to not be nil")
					}
					if e, a := "barKey", *params.Key; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					return &s3.GetObjectOutput{
						Body: io.NopCloser(bytes.NewReader([]byte("this is the body foo bar baz"))),
					}, nil
				})
			},
			bucket: "amzn-s3-demo-bucket",
			key:    "barKey",
			expect: []byte("this is the body foo bar baz"),
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := GetObjectFromS3(ctx, tt.client(t), tt.bucket, tt.key)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; bytes.Compare(e, a) != 0 {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestGetObjectFromS3ViaTestifyAndMockery(t *testing.T) {
	cases := []struct {
		name   string
		bucket string
		key    string
		expect []byte
	}{
		{
			name:   "GetObjectFromS3 returns expected content",
			bucket: "amzn-s3-demo-bucket",
			key:    "barKey",
			expect: []byte("this is the body foo bar baz"),
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()

			mockS3Client := mocks.NewGetObjectAPI(t)

			expectedInput := &s3.GetObjectInput{
				Bucket: &tt.bucket,
				Key:    &tt.key,
			}

			mockReturn := &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader([]byte("this is the body foo bar baz"))),
			}

			mockS3Client.On("GetObject", ctx, expectedInput, mock.Anything).Return(mockReturn, nil)

			content, err := GetObjectFromS3(ctx, mockS3Client, tt.bucket, tt.key)
			require.NoError(t, err)
			assert.Equal(t, tt.expect, content)
		})
	}
}
