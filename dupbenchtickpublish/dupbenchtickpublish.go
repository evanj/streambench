package main

import (
	"flag"
	"log"
	"time"

	"github.com/evanj/streambench/dupbench"
)

func main() {
	publishInterval := flag.Duration("interval", 2*time.Minute, "interval to publish a new batch")
	client, err := dupbench.ParseClient()
	if err != nil {
		panic(err)
	}

	now := time.Now().UTC()
	nextWakeUp := now.Truncate(*publishInterval).Add(*publishInterval)
	log.Printf("now:%s starting at %s ...",
		now.Format(time.RFC3339), nextWakeUp.Format(time.RFC3339))
	for {
		time.Sleep(nextWakeUp.Sub(time.Now()))
		nextWakeUp = nextWakeUp.Add(*publishInterval)

		err = client.PublishBatch()
		if err != nil {
			panic(err)
		}
	}

	// technically probably does not need to be closed since we are exiting
	err = client.Close()
	if err != nil {
		panic(err)
	}
}
