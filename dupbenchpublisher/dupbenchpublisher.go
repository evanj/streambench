package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	"cloud.google.com/go/pubsub"
	"github.com/evanj/streambench/messages"
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

func publisherGoroutine(wg *sync.WaitGroup, topic *pubsub.Topic, idString string, numMessages int) {
	defer wg.Done()
	log.Printf("goroutine %s starting; %d numMessages", idString, numMessages)

	ctx := context.Background()

	msg := &messages.DuplicateTest{
		GoroutineId: idString,
		Created:     ptypes.TimestampNow(),
	}

	results := []*pubsub.PublishResult{}
	for i := 0; i < numMessages; i++ {
		msg.Sequence = int64(i)
		setTimestampNow(msg.Created)

		// must create new pubsub.Message/bytes since Publish takes a reference
		msgBytes, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}

		psMsg := &pubsub.Message{Data: msgBytes}
		results = append(results, topic.Publish(ctx, psMsg))
		if len(results) >= publishWaitAfterMessages {
			waitForPublish(ctx, results)
			results = results[:0]
		}
	}
	waitForPublish(ctx, results)
	log.Printf("goroutine %s exiting", idString)
}

func main() {
	projectID := flag.String("projectID", "", "Google Cloud Project to use (empty = default)")
	topicID := flag.String("topicID", "dupbench", "Pubsub topic ID to publish to")
	goroutines := flag.Int("goroutines", 5, "Number of goroutines to use for publishing")
	numMessages := flag.Int("numMessages", 1000, "numMessages each goroutine will publish")
	flag.Parse()
	log.Printf("publishing to projectID:%s topicID:%s; %d numMessages * %d goroutines = %d total numMessages",
		*projectID, *topicID, *numMessages, *goroutines, *numMessages**goroutines)

	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		panic(err)
	}
	topic := pubsubClient.Topic(*topicID)

	start := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < *goroutines; i++ {
		idBytes := make([]byte, goroutineIDLength)
		_, err = rand.Read(idBytes)
		if err != nil {
			panic(err)
		}
		idString := hex.EncodeToString(idBytes)

		wg.Add(1)
		go publisherGoroutine(wg, topic, idString, *numMessages)
	}
	wg.Wait()

	// technically probably does not need to be closed since we are exiting
	topic.Stop()
	err = pubsubClient.Close()
	if err != nil {
		panic(err)
	}
	end := time.Now()

	duration := end.Sub(start)
	rate := float64(*goroutines**numMessages) / duration.Seconds()
	log.Printf("published %d total messages in %s ; %.1f msgs/sec",
		*goroutines**numMessages, duration.String(), rate)
}
