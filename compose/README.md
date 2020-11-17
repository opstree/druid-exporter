# docker-compose test setup

This directory contains a docker-compose file for a full dockerized druid cluster taken from the 
official [repository of druid](https://github.com/apache/druid/blob/master/distribution/docker/docker-compose.yml), 
in addition to that, it has druid-exporter as well as Prometheus server and Grafana with druid-dashboard installed. 

### Quick start
- Clone the repo and cd to compose directory `cd compose`
- Run docker-compose command to start the services : `docker-compose up -d`
- After a couple seconds (time for services to start.), browse to Grafana web UI `http://localhost:3000` and authenticate
with the user `admin` and password `SomePassword`. 
- You will see the druid-dashboard, you can try to ingest some data into druid, and you'll the dashboard metrics get populated.

Enjoy !