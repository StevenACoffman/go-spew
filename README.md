# Go Spew

This is a dumb container that randomly generates and emits structured log messages on a regular interval for testing a log pipeline. This is useful for watermarking in kafka topics

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
