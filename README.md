# Streaming Data: Google Cloud Pub/Sub and BigQuery

This is a quick-and-dirty benchmark comparing Google Cloud Pub/Sub publish and BigQuery streaming insert APIs. The summary is:

* Pubsub publish is fast: p50: 6m; p95: 13ms
* BigQuery streaming inserts are pretty slow: p50: 75ms; p95: 110ms
* Setting insertID or not does not change BigQuery performance substantially.


## Create necessary resources

1. Create a topic: `gcloud pubsub topic create streamtest`
2. Create a subscription for it: `gcloud pubsub subscriptions create --topic=streamtest streamtest-subscription`
3. Create a BigQuery dataset and table:
```
bq mk streamtest
bq mk --table --time_partitioning_type=DAY streamtest.streamtest 'record_id:INTEGER,data:BYTES'
```


## Run it

go run streambench.go --projectID=project --topicID=streamtest --datasetID=streamtest --tableID=streamtest


## Clean up resources

```
gcloud pubsub topics delete streamtest
bq rm streamtest.streamtest
bq rm streamtest
```


## Pubsub Performance notes

* There is a warmup effect. The first say 1000 API calls for a given topic/subscription are pretty slow, but fairly quickly converge to a "stable" state.
* There is something weird when publishing a batch with a single message. It seems to be slower than publishing a batch with 2 messages. This is even true when I reversed the order of message batch sizes!


## BigQuery Streaming

The BigQuery team released a new API version in August 2019. Release Notes: https://cloud.google.com/bigquery/docs/release-notes#August_19_2019

"If you stream data into BigQuery without populating the insertId field, you get the following higher quotas in the US multi-region location. These higher quotas are currently in beta."

It does not look seem like skipping the insertID field makes BigQuery publishing any lower latency.
