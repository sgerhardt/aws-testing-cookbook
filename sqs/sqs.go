package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type ReceiveMessageAPI interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

type Actions struct {
	api ReceiveMessageAPI
}

// GetMessages uses the ReceiveMessage action to get messages from an Amazon SQS queue.
func (actor Actions) GetMessages(ctx context.Context, queueUrl string, maxMessages int32, waitTime int32) ([]types.Message, error) {
	var messages []types.Message
	result, err := actor.api.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTime,
	})
	if err != nil {
		log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", queueUrl, err)
	} else {
		messages = result.Messages
	}
	return messages, err
}
