#! /bin/sh

cd docker
docker-compose stop

rm -f /root/taep-data/kapacitor/agent.sock
