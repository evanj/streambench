package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

func gcloud(args ...string) error {
	log.Println("gcloud", strings.Join(args, " "))
	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	projectID := flag.String("projectID", "", "Google Cloud project ID")
	topicID := flag.String("topicID", "dupbench", "Pub/Sub Topic ID")
	subscriptionID := flag.String("subscriptionID", "dupbench", "Pub/Sub Subscription ID")
	flag.Parse()

	err := gcloud("--project="+*projectID, "pubsub", "topics", "create", *topicID)
	if err != nil {
		panic(err)
	}
	err = gcloud("--project="+*projectID, "pubsub", "subscriptions", "create",
		"--topic="+*topicID, *subscriptionID)
	if err != nil {
		panic(err)
	}
}
