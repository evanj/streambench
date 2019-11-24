package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func mustExec(cmd *exec.Cmd) {
	log.Println(strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func mustGcloud(args ...string) {
	mustExec(exec.Command("gcloud", args...))
}

func mustBQ(args ...string) {
	mustExec(exec.Command("bq", args...))
}

func main() {
	projectID := flag.String("projectID", "", "Google Cloud project ID")
	// topicID := flag.String("topicID", "dupbench", "Pub/Sub Topic ID")
	// subscriptionID := flag.String("subscriptionID", "dupbench", "Pub/Sub Subscription ID")
	// datasetID := flag.String("datasetID", "dupbench", "BigQuery dataset ID")
	// tableID := flag.String("tableID", "dupbench", "BigQuery table ID")
	flag.Parse()

	// log.Println("creating pubsub topic/subscription")
	// mustGcloud("--project="+*projectID, "pubsub", "topics", "create", *topicID)
	// mustGcloud("--project="+*projectID, "pubsub", "subscriptions", "create", "--topic="+*topicID, *subscriptionID)

	// log.Println("creating bigquery table")
	// mustBQ("--project="+*projectID, "mk", *datasetID)
	// mustBQ("--project="+*projectID, "mk", "--table", *datasetID+"."+*tableID,
	// 	"goroutine_id:STRING,sequence:INTEGER,created:TIMESTAMP,subscriber_received:TIMESTAMP")

	log.Println("building and publishing container images")
	containerURL := fmt.Sprintf("gcr.io/%s/dupbenchpublisher", *projectID)
	mustGcloud("--project="+*projectID, "builds", "submit", ".", "--tag="+containerURL)

	fmt.Printf("gcloud --project=%s compute instances create-with-container publisher-vm \\\n"+
		"    --container-image=%s --zone=us-east1-c \\\n"+
		"    --container-restart-policy=NEVER --machine-type=n1-highcpu-2 --preemptible \\\n"+
		"    --container-arg=--projectID=%s --container-arg=--numMessages=10000000 \\\n"+
		"    --container-arg=--goroutine=8\n",
		*projectID, containerURL, *projectID,
	)
}
