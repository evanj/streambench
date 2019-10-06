package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"time"

	"google.golang.org/api/bigquery/v2"
	"google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

const pubsubEndpoint = "pubsub.googleapis.com:443"
const pubsubScope = "https://www.googleapis.com/auth/pubsub"

func connectPublisherGRPC(ctx context.Context, creds credentials.PerRPCCredentials) (pubsub.PublisherClient, error) {
	grpcOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(creds),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithBlock(),
	}
	conn, err := grpc.DialContext(ctx, pubsubEndpoint, grpcOpts...)
	if err != nil {
		return nil, err
	}
	return pubsub.NewPublisherClient(conn), nil
}

func topicPath(projectID string, topicID string) string {
	return "projects/" + projectID + "/topics/" + topicID
}

func percentile(durations []time.Duration, percentile float64) time.Duration {
	index := int(float64(len(durations)) * percentile)
	return durations[index]
}

type inserter interface {
	insert(ctx context.Context, messages [][]byte) error
}

type pubsubInserter struct {
	publisher pubsub.PublisherClient
	topicPath string
}

func (p *pubsubInserter) insert(ctx context.Context, messages [][]byte) error {
	req := &pubsub.PublishRequest{Topic: p.topicPath}
	for _, message := range messages {
		req.Messages = append(req.Messages, &pubsub.PubsubMessage{Data: message})
	}

	md := metadata.Pairs("x-goog-request-params", "topic="+p.topicPath)
	outCtx := metadata.NewOutgoingContext(ctx, md)
	resp, err := p.publisher.Publish(outCtx, req)
	if err != nil {
		return err
	}
	if len(resp.MessageIds) != len(messages) {
		return fmt.Errorf("wrong number of message IDs: %d", len(resp.MessageIds))
	}
	return nil
}

type bqInserter struct {
	bq          *bigquery.Service
	projectID   string
	datasetID   string
	tableID     string
	rowCounter  int
	useInsertID bool
}

func (b *bqInserter) insert(ctx context.Context, messages [][]byte) error {
	tableInsertReq := &bigquery.TableDataInsertAllRequest{}
	for _, message := range messages {
		b.rowCounter += 1
		rowData := map[string]bigquery.JsonValue{
			"record_id": bigquery.JsonValue(b.rowCounter),
			"data":      bigquery.JsonValue(message),
		}

		row := &bigquery.TableDataInsertAllRequestRows{
			Json: rowData,
		}
		if b.useInsertID {
			row.InsertId = strconv.Itoa(b.rowCounter)
		}
		tableInsertReq.Rows = append(tableInsertReq.Rows, row)
	}

	resp, err := b.bq.Tabledata.InsertAll(b.projectID, b.datasetID, b.tableID, tableInsertReq).Context(ctx).Do()
	if err != nil {
		return err
	}
	if len(resp.InsertErrors) != 0 {
		return fmt.Errorf("%d insert errors", len(resp.InsertErrors))
	}
	return nil
}

func main() {
	projectID := flag.String("projectID", "", "Google Cloud Project to use (empty = default)")
	topicID := flag.String("topicID", "", "Pubsub topic ID to publish to")
	datasetID := flag.String("datasetID", "", "BigQuery dataset")
	tableID := flag.String("tableID", "", "BigQuery table")
	flag.Parse()
	fmt.Printf("publishing to projectID:%s topicID:%s\n", *projectID, *topicID)
	fmt.Printf("BigQuery datasetID:%s tableID:%s\n", *datasetID, *tableID)

	ctx := context.Background()
	credentials, err := oauth.NewApplicationDefault(ctx, pubsubScope)
	if err != nil {
		panic(err)
	}

	// using the raw BigQuery API instead of the friendly Go API: Allows us to build the request directly
	bq, err := bigquery.NewService(ctx)
	if err != nil {
		panic(err)
	}
	bqInsert := &bqInserter{bq, *projectID, *datasetID, *tableID, 0, true}
	bqInsertNoInsertID := &bqInserter{bq, *projectID, *datasetID, *tableID, 0, false}

	// create a raw gRPC publisher connection: the Go pubsub client defers work to background threads
	publisher, err := connectPublisherGRPC(ctx, credentials)
	if err != nil {
		panic(err)
	}
	pubsubInsert := &pubsubInserter{publisher, topicPath(*projectID, *topicID)}

	// warm up the pub sub topic: it needs some traffic to get fast;
	// the p90 publishing time is about 12ms, so this should happen fairly quickly
	const warmUpMinDuration = 12 * time.Millisecond
	const warmUpConsecutive = 10
	totalWarmUps := 0
	fmt.Println("warming up the topic/publisher ...")
	for fastCount := 0; fastCount < warmUpConsecutive; {
		totalWarmUps += 1
		start := time.Now()
		err := pubsubInsert.insert(ctx, [][]byte{[]byte("hello")})
		end := time.Now()
		if err != nil {
			panic(err)
		}
		duration := end.Sub(start)
		if duration > warmUpMinDuration {
			fastCount = 0
		} else {
			fastCount += 1
		}
	}
	fmt.Printf("sent %d pubsub warmup requests\n", totalWarmUps)

	for i := 0; i < 10; i++ {
		start := time.Now()
		err := pubsubInsert.insert(ctx, [][]byte{[]byte("hello")})
		end := time.Now()
		if err != nil {
			panic(err)
		}
		duration := end.Sub(start)
		fmt.Printf("bigquery warmup: %d = %s\n", i, duration.String())
	}

	const requestCount = 1000
	for _, inserter := range []inserter{bqInsert, bqInsertNoInsertID, pubsubInsert} {
		fmt.Printf("\n# Publishing messages using %T\n", inserter)

		for _, byteLength := range []int{10, 100, 1000, 10000} {
			fakeRequestData := make([]byte, byteLength)
			_, err = rand.Read(fakeRequestData)
			if err != nil {
				panic(err)
			}

			fmt.Println()
			for _, messageCount := range []int{100, 10, 5, 2, 1} {
				messages := make([][]byte, messageCount)
				for i := range messages {
					messages[i] = fakeRequestData
				}

				fmt.Printf("publishing batches of %d messages with %d bytes (%d kiB/publish)\n",
					len(messages), len(messages[0]), (len(messages)*len(messages[0]))/1024)
				durations := make([]time.Duration, requestCount)
				for i := range durations {
					start := time.Now()
					err := inserter.insert(ctx, messages)
					end := time.Now()
					if err != nil {
						panic(err)
					}

					durations[i] = end.Sub(start)
				}

				sort.Slice(durations, func(i int, j int) bool { return durations[i] < durations[j] })
				sum := time.Duration(0)
				for _, v := range durations {
					sum += v
				}
				fmt.Printf("%d requests; mean:%s; min:%s p25:%s p50:%s p75:%s p90:%s p95:%s p99:%s max:%s\n",
					len(durations), sum/time.Duration(len(durations)), durations[0],
					percentile(durations, 0.25), percentile(durations, 0.5), percentile(durations, 0.75),
					percentile(durations, 0.9), percentile(durations, 0.95), percentile(durations, 0.99),
					durations[len(durations)-1])
			}
		}
	}

}
