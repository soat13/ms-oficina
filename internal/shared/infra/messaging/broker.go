package messaging

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sharedMsg "github.com/soat13/fase-1-oficina/internal/shared/messaging"
)

type (
	Broker struct {
		client    *sqs.Client
		baseURL   string
		consumers []*consumer
		cancel    context.CancelFunc
		wg        sync.WaitGroup
	}

	consumer struct {
		client   *sqs.Client
		queueURL string
		topic    string
		handler  sharedMsg.Handler
	}
)

func NewNoopBroker() *Broker {
	return &Broker{}
}

func NewBroker(ctx context.Context) (*Broker, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
			o.BaseEndpoint = &endpoint
		}
	})

	baseURL := os.Getenv("SQS_BASE_URL")
	if baseURL != "" && !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return &Broker{client: client, baseURL: baseURL}, nil
}

func (b *Broker) Publish(ctx context.Context, topic string, payload []byte) error {
	if b.client == nil {
		return nil
	}

	queueURL := b.queueURL(topic)
	_, err := b.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(payload)),
	})
	if err != nil {
		log.Error().Err(err).Str("topic", topic).Str("queue", queueURL).Msg("failed to publish message to SQS")
		return fmt.Errorf("sqs publish to %s: %w", topic, err)
	}

	log.Debug().Str("topic", topic).Str("queue", queueURL).Msg("message published to SQS")
	return nil
}

func (b *Broker) Subscribe(topic string, handler sharedMsg.Handler) {
	b.consumers = append(b.consumers, &consumer{
		client:   b.client,
		queueURL: b.queueURL(topic),
		topic:    topic,
		handler:  handler,
	})
}

func (b *Broker) Start(ctx context.Context) {
	ctx, b.cancel = context.WithCancel(ctx)
	for _, c := range b.consumers {
		b.wg.Add(1)
		go func(c *consumer) {
			defer b.wg.Done()
			c.poll(ctx)
		}(c)
	}
	log.Info().Int("count", len(b.consumers)).Msg("all SQS consumers started")
}

func (b *Broker) Shutdown() {
	if b.cancel != nil {
		b.cancel()
	}
	b.wg.Wait()
	log.Info().Msg("all SQS consumers stopped")
}

func (b *Broker) queueURL(topic string) string {
	return b.baseURL + strings.ToLower(strings.ReplaceAll(topic, ".", "-"))
}

func (c *consumer) poll(ctx context.Context) {
	logger := log.With().Str("queue", c.queueURL).Str("topic", c.topic).Logger()
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("SQS consumer stopped")
			return
		default:
		}

		messages, err := c.receiveMessages(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Error().Err(err).Msg("failed to receive messages from SQS")
			continue
		}

		for _, msg := range messages {
			c.processMessage(ctx, logger, msg)
		}
	}
}

func (c *consumer) receiveMessages(ctx context.Context) ([]types.Message, error) {
	result, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		return nil, err
	}
	return result.Messages, nil
}

func (c *consumer) processMessage(ctx context.Context, logger zerolog.Logger, msg types.Message) {
	if err := c.handler(ctx, c.topic, []byte(aws.ToString(msg.Body))); err != nil {
		logger.Error().Err(err).Str("message_id", aws.ToString(msg.MessageId)).Msg("failed to handle message")
		return
	}

	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		logger.Error().Err(err).Str("message_id", aws.ToString(msg.MessageId)).Msg("failed to delete message from SQS")
	}
}
