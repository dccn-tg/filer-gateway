#!/bin/bash

docker-compose -p filer-gateway down

#function cnt_containers() {
#    docker network inspect -f '{{range $k, $v := .Containers}}{{println $k}}{{end}}' $1 | grep -v -e '^$' | wc -l
#}
#
#source env.sh
#
#docker stack rm filer-gateway
#
#while true; do
#    nc=$( cnt_containers filer-gateway_default )
#    echo "containers attached to default network: $nc"
#    [ $nc -gt 0 ] && sleep 1 || break
#done
#
#echo "shutting down network filer-gateway_default ..."
#docker network rm filer-gateway_default
