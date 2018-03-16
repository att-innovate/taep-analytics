#! /bin/sh
echo --------- grafana ---------
docker logs docker_grafana_1

echo --------- influxdb ---------
docker logs docker_influxdb_1

echo --------- kapacitor ---------
docker logs docker_kapacitor_1

echo --------- telegraf ---------
docker logs docker_telegraf_1

echo --------- agent ---------
docker logs docker_agent_1
