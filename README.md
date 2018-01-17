## Get kubernetes cert and token
```
KOPS_NAME=k8s-cluster.test.cirrostratus.org
bless_scp ubuntu@api.internal.$KOPS_NAME:/srv/kubernetes/ca.crt $HOME/Downloads/ca.crt
TOKEN=$(kubectl describe secret "$(kubectl get secrets | grep default-token | cut -f1 -d ' ')" | grep -E '^token' | cut -f2 -d':' | tr -d '\t' | tr -d ' ')
printf "$TOKEN" > ~/Downloads/token
```

## How to run it locally

```
docker-compose up
```

## How to stop it locally

```bash
image_name=docker-registry.acorn.cirrostratus.org/playground/go-eureka
docker stop "$(docker ps --format '{{.ID}}' --filter "ancestor=${image_name}")"
```
