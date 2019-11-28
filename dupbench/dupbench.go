package dupbench

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/evanj/streambench/messages"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// 16 = 128 bits which should make collisions "impossible"
const goroutineIDLength = 16

// pubsub has a maximum of 1000 messages per publish request. Hopefully this is enough messages
// to "fill" the pipeline and keep the publishing busy
const publishWaitAfterMessages = 2000

func setTimestampNow(ts *timestamp.Timestamp) {
	// stupidly over-optimized to avoid allocations (probably unnecessary)
	t := time.Now()
	ts.Seconds = t.Unix()
	ts.Nanos = int32(t.Nanosecond())
}

func waitForPublish(ctx context.Context, results []*pubsub.PublishResult) {
	for _, r := range results {
		_, err := r.Get(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Client) publisherGoroutine(wg *sync.WaitGroup, idString string) {
	defer wg.Done()

	ctx := context.Background()
	msg := &messages.DuplicateTest{
		GoroutineId: idString,
		Created:     ptypes.TimestampNow(),
	}

	results := []*pubsub.PublishResult{}
	for i := 0; i < c.numMessages; i++ {
		msg.Sequence = int64(i)
		setTimestampNow(msg.Created)

		// must create new pubsub.Message/bytes since Publish takes a reference
		msgBytes, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}

		psMsg := &pubsub.Message{Data: msgBytes}
		results = append(results, c.topic.Publish(ctx, psMsg))
		if len(results) >= publishWaitAfterMessages {
			waitForPublish(ctx, results)
			results = results[:0]
		}
	}
	waitForPublish(ctx, results)
	log.Printf("goroutine %s exiting", idString)
}

// Client publishes batches of messages to Pub/Sub.
type Client struct {
	goroutines   int
	numMessages  int
	pubsubClient *pubsub.Client
	topic        *pubsub.Topic
}

// PublishBatch publishes one batch of messages.
func (c *Client) PublishBatch() error {
	start := time.Now()
	log.Printf("starting batch ...")
	wg := &sync.WaitGroup{}
	for i := 0; i < c.goroutines; i++ {
		idBytes := make([]byte, goroutineIDLength)
		_, err := rand.Read(idBytes)
		if err != nil {
			panic(err)
		}
		idString := hex.EncodeToString(idBytes)

		wg.Add(1)
		go c.publisherGoroutine(wg, idString)
	}
	wg.Wait()
	end := time.Now()

	duration := end.Sub(start)
	totalMessages := c.goroutines * c.numMessages
	rate := float64(totalMessages) / duration.Seconds()
	log.Printf("published %d total messages in %s ; %.1f msgs/sec",
		totalMessages, duration.String(), rate)
	return nil
}

// ParseClient parses command line arguments and returns a configured client.
func ParseClient() (*Client, error) {
	projectID := flag.String("projectID", "", "Google Cloud Project ID (empty = default)")
	topicID := flag.String("topicID", "dupbench", "Pub/Sub topic ID to publish to")
	goroutines := flag.Int("goroutines", 1, "number of goroutines to use to publish at a time")
	numMessages := flag.Int("numMessages", 1000, "numMessages each goroutine will publish in a batch")
	flag.Parse()

	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		return nil, err
	}
	topic := pubsubClient.Topic(*topicID)

	return &Client{*goroutines, *numMessages, pubsubClient, topic}, nil
}

// Close flushes messages from the topic and stops the pub/sub client.
func (c *Client) Close() error {
	c.topic.Stop()
	return c.pubsubClient.Close()
}
