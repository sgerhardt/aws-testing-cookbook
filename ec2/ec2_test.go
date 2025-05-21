package ec2

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEC2InstancesViaTestifyAndMockery(t *testing.T) {
	cases := []struct {
		name          string
		expectedInput *ec2.DescribeInstancesInput
		expectOutput  []types.Instance
	}{
		{
			name:          "ListInstances returns empty response when no instances are found",
			expectedInput: &ec2.DescribeInstancesInput{},
			expectOutput:  nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			mockReturn := &ec2.DescribeInstancesOutput{
				NextToken:      nil,
				Reservations:   nil,
				ResultMetadata: middleware.Metadata{},
			}

			ctx := context.TODO()
			mockEC2Client := NewMockDescribeInstancesAPI(t)
			mockEC2Client.On("DescribeInstances", ctx, tt.expectedInput).Return(mockReturn, nil)

			instances, err := ListInstances(ctx, mockEC2Client)
			require.NoError(t, err)
			assert.Equal(t, tt.expectOutput, instances)
		})
	}
}
