# systemd_exporter

Proof of concept systemd monitor daemon.  It is modeled as a cross between https://github.com/prometheus/node_exporter and https://github.com/google/cadvisor.

### Run Locally
```shell
make docker
SHA=<sha>
docker run -p 8080:8080 --volume=/run/systemd:/run/systemd:ro --volume=/sys/fs/cgroup:/sys/fs/cgroup:ro $SHA
curl localhost:8080/metrics
```

### Deploy to kubernetes
```shell
make docker
# Tag and push your image
# Edit deploy/kubernetes/daemonset.yaml to specify your pushed image
kubectl apply -f deploy/kubernetes/daemonset.yaml
kubectl get po
# pick a systemd-exporter pod to query
POD=systemd-exporter-12345
kubectl get --raw /api/v1/namespaces/default/pods/$POD/proxy/metrics
```