package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/evanj/streambench/messages"
	"github.com/evanj/streambench/pubsubgrpc"
	"google.golang.org/protobuf/proto"
)

// BQ streaming API suggests a max of 500 rows / insert
const maxMessagesPerPull = 500

type bqSchema struct {
	GoroutineID        string    `bigquery:"goroutine_id"`
	Sequence           int64     `bigquery:"sequence"`
	Created            time.Time `bigquery:"created"`
	Published          time.Time `bigquery:"published"`
	SubscriberReceived time.Time `bigquery:"subscriber_received"`
}

func subscriberGoroutine(
	wg *sync.WaitGroup, subscriber pubsubpb.SubscriberClient, fullSubscriptionName string,
	inserter *bigquery.Inserter,
) {
	defer wg.Done()
	log.Printf("subscriber goroutine starting ...")

	ctx := context.Background()
	pullReq := &pubsubpb.PullRequest{
		Subscription:      fullSubscriptionName,
		MaxMessages:       maxMessagesPerPull,
		ReturnImmediately: false,
	}

	ackReq := &pubsubpb.AcknowledgeRequest{
		Subscription: fullSubscriptionName,
	}
	rows := []bqSchema{}
	msg := &messages.DuplicateTest{}
	for {
		start := time.Now()
		resp, err := subscriber.Pull(ctx, pullReq)
		if err != nil {
			panic(err)
		}
		pullTimestamp := time.Now()

		// convert the rows
		for _, receivedMsg := range resp.ReceivedMessages {
			ackReq.AckIds = append(ackReq.AckIds, receivedMsg.AckId)

			err = proto.Unmarshal(receivedMsg.Message.Data, msg)
			if err != nil {
				panic(err)
			}
			created := msg.Created.AsTime()
			published := receivedMsg.Message.PublishTime.AsTime()

			row := bqSchema{msg.GoroutineId, msg.Sequence, created, published, pullTimestamp}
			rows = append(rows, row)
		}

		if len(rows) > 0 {
			err = inserter.Put(ctx, rows)
			if err != nil {
				panic(err)
			}
			rows = rows[:0]

			_, err = subscriber.Acknowledge(ctx, ackReq)
			if err != nil {
				panic(err)
			}
			ackReq.AckIds = ackReq.AckIds[:0]
		}

		end := time.Now()
		duration := end.Sub(start)
		rate := float64(len(resp.ReceivedMessages)) / duration.Seconds()
		log.Printf("streamed %d messages in %s; %.1f msgs/sec; ",
			len(resp.ReceivedMessages), duration.String(), rate)
	}
}

func main() {
	projectID := flag.String("projectID", "", "Google Cloud Project to use (empty = default)")
	subscriptionID := flag.String("subscriptionID", "dupbench", "Pubsub subscription ID to pull from")
	datasetID := flag.String("datasetID", "dupbench", "BigQuery dataset ID to insert into")
	tableID := flag.String("tableID", "dupbench", "BigQuery table ID to insert into")
	goroutines := flag.Int("goroutines", 1, "Number of goroutines to use for publishing")
	flag.Parse()
	log.Printf("copying from pubsub projectID:%s subscriptionID:%s to BigQuery %s.%s using %d goroutines",
		*projectID, *subscriptionID, *datasetID, *tableID, *goroutines)

	fullSubscriptionName := fmt.Sprintf("projects/%s/subscriptions/%s", *projectID, *subscriptionID)

	ctx := context.Background()
	conn, err := pubsubgrpc.Dial(ctx)
	if err != nil {
		panic(err)
	}
	subscriber := pubsubpb.NewSubscriberClient(conn)

	bq, err := bigquery.NewClient(ctx, *projectID)
	if err != nil {
		panic(err)
	}
	inserter := bq.Dataset(*datasetID).Table(*tableID).Inserter()

	wg := &sync.WaitGroup{}
	for i := 0; i < *goroutines; i++ {
		wg.Add(1)
		go subscriberGoroutine(wg, subscriber, fullSubscriptionName, inserter)
	}
	// will never return
	wg.Wait()
}
