# systemd_exporter

```shell
make docker
docker run -p 8080:8080 --volume=/run/systemd:/run/systemd:ro <sha>
curl localhost:8080
```