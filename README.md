# Streaming Data: Google Cloud Pub/Sub and BigQuery

This is a quick-and-dirty benchmark comparing Google Cloud Pub/Sub publish and BigQuery streaming insert APIs. The summary is:

* Pubsub publish is fast: p50: 6m; p95: 13ms
* BigQuery streaming inserts are pretty slow: p50: 75ms; p95: 110ms
* Setting insertID or not does not change BigQuery performance substantially.

*ARCHIVED 2024-01-07*: I haven't used this in years. I did one last set of updates, and it passed the CI tests, but I am marking this archived. It might be useful to me or others in the future ... maybe, but that becomes less likely as time goes on.


## Create necessary resources

1. Create a topic: `gcloud pubsub topics create streamtest`
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
gcloud pubsub subscriptions delete streamtest-subscription
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

However, as of 2019-11-23 this is still listed as "Beta" in the Quotas page, so it is possible I never tested the V2 inserts: https://cloud.google.com/bigquery/quotas#streaming_inserts

"you must complete the BigQuery Streaming V2 beta enrollment form in order to use [the V2 streaming insert quotas]"


# Duplicate Message Testing

Pub/Sub sometimes duplicates messages. To try and estimate how often, I created a quick test:

publisher -> pub/sub -> subscriber -> BigQuery table

This allows us to query the table after the fact to find messages that may have been duplicated. This
lets us experiment with what might cause duplicates

## Useful Queries

```sql
-- Find batches with duplicates and show "progress on an "in flight" batch
SELECT goroutine_id, COUNT(*) AS messages, COUNT(DISTINCT sequence) AS unique_sequence_nums,
  MAX(sequence) AS max_sequenc, MIN(sequence) AS min_sequence
FROM `dupbench.dupbench`
GROUP BY goroutine_id

SELECT goroutine_id, COUNT(*) as count, MAX(sequence) as max_seq, MIN(sequence) as min_seq,
  MIN(created) as min_created, MAX(created) as max_created,
  MIN(published) AS min_published, MAX(published) AS max_published,
  MIN(subscriber_received) as min_received, MAX(subscriber_received) as max_received
FROM `dupbench.dupbench`
GROUP BY goroutine_id

-- Find duplicated sequence numbers
SELECT goroutine_id, sequence, COUNT(*) as count
FROM `dupbench.dupbench`
GROUP BY goroutine_id, sequence
HAVING count > 1
ORDER BY goroutine_id, sequence
```

## Random observations

* When deleting and recreating a BigQuery table with streaming inserts, it is unclear exactly where your rows are going to go. I'e seen cases where the next batch of inserts went missing. It seems possible they went to the deleted table?

* Pubsub makes no guarantee of order, and when publishing batches in order, they arrive at the subscribers out of order (e.g. we receive some of the very low and very high sequence numbers, while still missing 50% of the numbers in a batch). However, some other experiments have shown that when there is a large backlog, messages do arrive "mostly in order". It might be nice to quantify this.

* The most common cause of duplication seems to be "dropped acks." For example, we pull a batch of messages, ack them, then all those messages are later redelivered. This means while the duplication rate is pretty low, it is bursty, which can cause issues.

* Sometimes some messages get "stuck" for the ack duration. E.g. put in 5000 messages, we get 4500 messages out, and the last 500 show up ~10 minutes later. It is as if something pulled them, then they have to retry. I saw this once when we had a "race" between the timeout on the pull request: our puller thinks we timed out, but pubsub thinks we got a batch.
