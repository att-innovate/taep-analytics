#! /bin/sh

export TASKS_DIR=/root/taep-analytics/docker/kapacitor/tasks

docker run --rm --net=host -it kapacitor:taep kapacitor enable agent
echo --------- enabled agent task ---------

echo --------- show task ---------
docker run --rm --net=host -it kapacitor:taep kapacitor show agent
