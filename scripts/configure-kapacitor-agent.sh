#! /bin/sh

export TASKS_DIR=/root/taep-analytics/docker/kapacitor/tasks

docker run --rm --net=host -it -v $TASKS_DIR:/tasks kapacitor:taep kapacitor define agent -type batch -tick /tasks/agent.tick -dbrp telegraf.autogen
echo --------- uploaded task ---------

echo --------- show task ---------
docker run --rm --net=host -it kapacitor:taep kapacitor show agent

echo --------- list tasks ---------
docker run --rm --net=host -it kapacitor:taep kapacitor list tasks