# Go Spew

This is a dumb container that randomly generates and emits structured log messages on a regular interval for monitoring a log pipeline. This is also useful for watermarking in kafka topics.

Watermark is a moving threshold in event-time that trails behind the maximum event-time seen by the query in the processed data. The trailing gap defines how long we will wait for late data to arrive. By knowing the point at which no more data will arrive for a given group, we can limit the total amount of state that we need to maintain for a query. For example, suppose the configured maximum lateness is 10 minutes. That means the events that are up to 10 minutes late will be allowed to aggregate. And if the maximum observed event time is 12:33, then all the future events with event-time older than 12:23 will be considered as "too late" and dropped.

Because the messages are written to disk with a timestamp, and there's another timestamp from when the message is consumed, while allowing for clock skew, this will give good insight into the latency of the log pipeline.

Similarly, if messages fail to show up in a reasonable window, this can indicate a failure of the log pipeline.

## How to run it in kubernetes

```
kubectl apply -f k8s/go-spew.yaml
```

You can adjust the `INTERVAL` environment variable (default to 30 seconds) for how frequently.

## How to run it locally

```
docker-compose up
```

## How to stop it locally

```bash
image_name=stevenacoffman/go_spew
docker stop "$(docker ps --format '{{.ID}}' --filter "ancestor=${image_name}")"
```
