#!/bin/bash

function cnt_containers() {
    docker network inspect -f '{{range $k, $v := .Containers}}{{println $k}}{{end}}' $1 | grep -v -e '^$' | wc -l
}

source env.sh

docker stack rm ${STACK_NAME}

while true; do
    nc=$( cnt_containers ${STACK_NAME}_default )
    echo "containers attached to default network: $nc"
    [ $nc -gt 0 ] && sleep 1 || break
done

echo "shutting down network ${STACK_NAME}_default ..."
docker network rm ${STACK_NAME}_default
