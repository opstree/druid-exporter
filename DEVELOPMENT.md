## Development

## Prerequisites

#### Access to a druid cluster

First of all, you will need access to a druid cluster. The easiest way is to start is docker-compose.

- [Docker](https://docs.docker.com/engine/install/) - Layer on which druid cluster will run
- [Docker Compose](https://docs.docker.com/compose/install/) - to create druid cluster over docker containers

#### Tools to build and test druid exporter

- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/dl/)
- [Make](https://www.gnu.org/software/make/manual/make.html)

### Build Locally

To achieve this, execute this command:-

```shell
make build-code
```

### Build Docker Image

Druid exporter gets packaged as a container image for running on Kubernetes cluster. These instructions will guide you to build image.

```shell
make build-image
```

## Testing

To test the druid exporter, first we need a druid cluster. For creating the cluster we will use the official docker compose file of druid:-

https://github.com/apache/druid/blob/master/distribution/docker/docker-compose.yml

Before creating the druid cluster we have to make few changes in [environment](https://github.com/apache/druid/blob/master/distribution/docker/environment) file of druid docker compose.

https://github.com/apache/druid/blob/master/distribution/docker/environment

```properties
druid_emitter_http_recipientBaseUrl=http://<druid_exporter_url>:<druid_exporter_port>/druid
druid_emitter=http
```

After making these changes you can up the cluster by executing this command:-

```
docker-compose up -d
```

### Run Tests

```shell
make test
```
