## Development

## Prerequisites

#### Access to a druid cluster

First, you'll need connections to a cluster of druids. Docker-compose is the easiest way to get started.

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

Druid exporter for running on Kubernetes cluster is packaged as a container file. These instructions will help you in the image making process.

```shell
make build-image
```

## Testing

First we need a cluster of druid to check the druid-exporter. We will use the official docker-compose druid file to build the cluster:-

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
