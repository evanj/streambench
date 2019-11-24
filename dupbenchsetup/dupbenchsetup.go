package main

import (
	"bufio"
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

func gcloudLines(args ...string) ([]string, error) {
	cmd := exec.Command("gcloud", args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer output.Close()
	log.Println(strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	lines := []string{}
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	err = output.Close()
	if err != nil {
		return nil, err
	}
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func mustBQ(args ...string) {
	mustExec(exec.Command("bq", args...))
}

type cliArgs struct {
	projectID      string
	topicID        string
	subscriptionID string
	datasetID      string
	tableID        string
}

func setUp(args cliArgs) error {
	log.Println("creating pubsub topic/subscription")
	mustGcloud("--project="+args.projectID, "pubsub", "topics", "create", args.topicID)
	mustGcloud("--project="+args.projectID, "pubsub", "subscriptions", "create", "--topic="+args.topicID, args.subscriptionID)

	log.Println("creating bigquery table")
	mustBQ("--project="+args.projectID, "mk", args.datasetID)
	mustBQ("--project="+args.projectID, "mk", "--table", args.datasetID+"."+args.tableID,
		"goroutine_id:STRING,sequence:INTEGER,created:TIMESTAMP,subscriber_received:TIMESTAMP")

	log.Println("building and publishing container images")
	containerURL := fmt.Sprintf("gcr.io/%s/dupbenchpublisher", args.projectID)
	mustGcloud("--project="+args.projectID, "builds", "submit", ".", "--tag="+containerURL)

	fmt.Printf("gcloud --project=%s compute instances create-with-container publisher-vm \\\n"+
		"    --container-image=%s --zone=us-east1-c \\\n"+
		"    --container-restart-policy=NEVER --machine-type=n1-highcpu-2 --preemptible \\\n"+
		"    --container-arg=--projectID=%s --container-arg=--numMessages=10000000 \\\n"+
		"    --container-arg=--goroutines=8\n",
		args.projectID, containerURL, args.projectID,
	)

	return nil
}

func tearDown(args cliArgs) error {
	log.Println("deleting pubsub topic/subscription")
	mustGcloud("--project="+args.projectID, "pubsub", "subscriptions", "delete", args.subscriptionID)
	mustGcloud("--project="+args.projectID, "pubsub", "topics", "delete", args.topicID)

	log.Println("deleting bigquery dataset")
	mustBQ("--project="+args.projectID, "rm", "-r", "-f", args.datasetID)

	log.Println("deleting container images")
	containerURL := fmt.Sprintf("gcr.io/%s/dupbenchpublisher", args.projectID)
	shaDigests, err := gcloudLines("--project="+args.projectID, "container", "images", "list-tags",
		containerURL, "--format=value[no-transforms](digest)")
	if err != nil {
		return err
	}

	images := []string{}
	for _, digest := range shaDigests {
		images = append(images, containerURL+"@"+digest)
	}
	deleteArgs := append([]string{"--project=" + args.projectID, "container", "images", "delete"}, images...)
	mustGcloud(deleteArgs...)

	return nil
}

func main() {
	args := cliArgs{}
	flag.StringVar(&args.projectID, "projectID", "", "Google Cloud project ID")
	flag.StringVar(&args.topicID, "topicID", "dupbench", "Pub/Sub Topic ID")
	flag.StringVar(&args.subscriptionID, "subscriptionID", "dupbench", "Pub/Sub Subscription ID")
	flag.StringVar(&args.datasetID, "datasetID", "dupbench", "BigQuery dataset ID")
	flag.StringVar(&args.tableID, "tableID", "dupbench", "BigQuery table ID")
	tearDownArg := flag.Bool("tearDown", false, "if specified: delete resources instead of creating them")
	flag.Parse()

	var err error
	if *tearDownArg {
		err = tearDown(args)
	} else {
		err = setUp(args)
	}
	if err != nil {
		panic(err)
	}
}
