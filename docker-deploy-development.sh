#!/bin/bash

# export variables defined in env.sh
set -a && source env.sh && set +a
docker stack up -c docker-compose.yml -c docker-compose.swarm.yml --prune --with-registry-auth --resolve-image always filer-gateway