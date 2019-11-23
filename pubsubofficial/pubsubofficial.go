package main

import (
	"context"
	"flag"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

func main() {
	projectID := flag.String("projectID", "", "Google Cloud Project to use (empty = default)")
	topicID := flag.String("topicID", "", "Pubsub topic ID to publish to")
	flag.Parse()
	fmt.Printf("publishing to projectID:%s topicID:%s\n", *projectID, *topicID)

	ctx := context.Background()
	publisher, err := pubsub.NewClient(ctx, *projectID)
	it := publisher.Topics(ctx)
	for {
		topic, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			panic(err)
		}
		fmt.Printf("topic: %s\n", topic.ID())
	}

	result := publisher.Topic(*topicID).Publish(ctx, &pubsub.Message{Data: []byte("hello")})
	msgID, err := result.Get(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("message ID: %s\n", msgID)
}
