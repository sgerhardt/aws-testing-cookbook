package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options),
	) (*ec2.DescribeInstancesOutput, error)
}

func ListInstances(ctx context.Context, api DescribeInstancesAPI) ([]ec2types.Instance, error) {
	var (
		token     *string
		instances []ec2types.Instance
	)

	for {
		out, err := api.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
			NextToken: token,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range out.Reservations {
			instances = append(instances, r.Instances...)
		}
		if out.NextToken == nil {
			return instances, nil
		}
		token = out.NextToken
	}
}
