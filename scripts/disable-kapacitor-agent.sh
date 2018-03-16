#! /bin/sh

export TASKS_DIR=/root/taep-analytics/docker/kapacitor/tasks

docker run --rm --net=host -it kapacitor:taep kapacitor disable agent
echo --------- disabled agent task ---------

curl -X DELETE 'http://127.0.0.1:8100/divert'
echo "\n--------- divert table resetted ---------\n"
