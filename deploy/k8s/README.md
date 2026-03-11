# Kubernetes Deployment

## Prerequisites

- Kubernetes 1.28+
- Ingress controller (nginx)
- cert-manager (for TLS)

## Deploy

```bash
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secrets.yaml
kubectl apply -f phoenix-api.yaml
kubectl apply -f phoenix-frontend.yaml
kubectl apply -f ingress.yaml
```

## Verify

```bash
kubectl get pods -n phoenix
kubectl get svc -n phoenix
kubectl get ingress -n phoenix
```

## Scaling

- API: `kubectl scale deployment phoenix-api --replicas=5 -n phoenix`
- Frontend: `kubectl scale deployment phoenix-frontend --replicas=3 -n phoenix`
