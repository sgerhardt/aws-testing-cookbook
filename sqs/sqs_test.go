package sqs

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestReceiveMessageViaTestifyAndMockery(t *testing.T) {
	cases := []struct {
		name                string
		queueURL            string
		maxNumberOfMessages int32
		waitTimeSeconds     int32
		expect              []types.Message
	}{
		{
			name:                "GetMessages gets messages from the queue",
			queueURL:            "demo-queue-url",
			maxNumberOfMessages: 1,
			waitTimeSeconds:     10,
			expect:              []types.Message{{Body: aws.String("this is the body foo bar baz")}},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			expectedInput := &sqs.ReceiveMessageInput{
				QueueUrl:            &tt.queueURL,
				MaxNumberOfMessages: tt.maxNumberOfMessages,
				WaitTimeSeconds:     tt.waitTimeSeconds,
			}

			mockReturn := &sqs.ReceiveMessageOutput{
				Messages:       []types.Message{{Body: aws.String("this is the body foo bar baz")}},
				ResultMetadata: middleware.Metadata{},
			}

			ctx := context.TODO()
			mockSqsClient := NewMockReceiveMessageAPI(t)
			mockSqsClient.On("ReceiveMessage", ctx, expectedInput, mock.Anything).Return(mockReturn, nil)

			sqsActor := &Actions{api: mockSqsClient}

			messages, err := sqsActor.GetMessages(ctx, tt.queueURL, tt.maxNumberOfMessages, tt.waitTimeSeconds)
			require.NoError(t, err)
			assert.Equal(t, tt.expect, messages)
		})
	}
}
