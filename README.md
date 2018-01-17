# Go Spew

This is a dumb container that emits structured log messages on a regular interval for testing a log pipeline.

## How to run it locally

```
docker-compose up
```

## How to stop it locally

```bash
image_name=stevenacoffman/go_spew
docker stop "$(docker ps --format '{{.ID}}' --filter "ancestor=${image_name}")"
```
