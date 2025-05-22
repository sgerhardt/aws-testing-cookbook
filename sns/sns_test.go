package sns

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnsActions_PublishViaTestifyAndMockery(t *testing.T) {
	cases := []struct {
		name     string
		topicARN string
		message  string

		expectedSequenceNumber *string
		expectedMessageID      *string
	}{
		{
			name:                   "Publish message",
			topicARN:               "apples",
			message:                "this is the message",
			expectedSequenceNumber: aws.String("1"),
			expectedMessageID:      aws.String("mock-message-id"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			expectedInput := &sns.PublishInput{
				Message:  &tt.message,
				TopicArn: &tt.topicARN,
			}

			mockReturn := &sns.PublishOutput{
				MessageId:      tt.expectedMessageID,
				SequenceNumber: tt.expectedSequenceNumber,
				ResultMetadata: middleware.Metadata{},
			}

			ctx := context.TODO()
			mockPublishAPI := NewMockPublishAPI(t)
			mockPublishAPI.On("Publish", ctx, expectedInput).Return(mockReturn, nil)

			snsActor := &SnsActions{Publisher: mockPublishAPI}

			output, err := snsActor.Publish(ctx, tt.topicARN, tt.message, "", "", "", "")
			require.NoError(t, err)
			assert.Equal(t, tt.expectedSequenceNumber, output.SequenceNumber)
			assert.Equal(t, tt.expectedMessageID, output.MessageId)
		})
	}
}
