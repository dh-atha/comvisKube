# Cloud Native Go App

Simple Go web app for Kubernetes practice. It exposes a dashboard at `/`, a persistent state view at `/state`, a crash trigger at `/panic`, a CPU load endpoint at `/work`, and a pod identity endpoint at `/pod`.

## Prerequisites

- Docker or another container runtime
- Kubernetes cluster such as Minikube, kind, or a local cluster
- `kubectl`

## Build the image

```bash
docker build -t cloud-go:latest .
```

If you use a local Kubernetes cluster, make sure the image is available to the cluster.

- Minikube:

```bash
minikube image load cloud-go:latest
```

- kind:

```bash
kind load docker-image cloud-go:latest
```

## Deploy to Kubernetes

Apply the manifest:

```bash
kubectl apply -f k8s-deployment.yaml
```

If you want to re-apply all manifests in the folder:

```bash
kubectl apply -R -f .
```

## Verify the deployment

Check the rollout and pod status:

```bash
kubectl get pods
kubectl get hpa
kubectl get svc
```

## Access the app

The service is exposed on port `7001` and forwards to container port `8080`.

- If you are using `LoadBalancer`, get the external address with:

```bash
kubectl get svc cloud-go-service
```

- If you are using Minikube, you can also port-forward:

```bash
kubectl port-forward service/cloud-go-service 7001:7001
```

Then open:

- `http://localhost:7001/`
- `http://localhost:7001/state`

## Test behavior

- Click `Run long load` to create CPU pressure for the HPA.
- Click `Trigger panic` to force a restart.
- Click `Show pod` to see which pod handled the request.
- Open `/state` to inspect persisted metadata stored on the PVC.

## Notes

- The deployment starts with `replicas: 2` and the HPA range is `minReplicas: 2`, `maxReplicas: 4` in `k8s-deployment.yaml`.
- The app writes metadata to `/data/metadata/events.jsonl`.
- `imagePullPolicy: Always` means the node will try to pull the image on each pod start.
