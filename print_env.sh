#!/bin/bash

echo "# version"
echo "DOCKER_IMAGE_TAG=$DOCKER_IMAGE_TAG"
echo 
echo "# docker registry endpoint"
echo "DOCKER_REGISTRY=$DOCKER_REGISTRY"
echo 
echo "# docker registry username"
echo "DOCKER_REGISTRY_USER=$DOCKER_REGISTRY_USER"
echo 
echo "# docker registry password"
echo "DOCKER_REGISTRY_PASSWORD=$DOCKER_REGISTRY_PASSWORD"
echo 
echo "# volume for home directory"
echo "HOME_VOL=$HOME_VOL"
echo
echo "# volume for project directory"
echo "PROJECT_VOL=$PROJECT_VOL"
echo
echo "# volume for project_freenas directory"
echo "PROJECT_FREENAS_VOL=$PROJECT_FREENAS_VOL"
echo
echo "# volume for project_cephfs directory"
echo "PROJECT_CEPHFS_VOL=$PROJECT_CEPHFS_VOL"
echo
echo "# configuration file for api-server"
echo "CFG_API_SERVER=$CFG_API_SERVER"
echo
echo "# configuration file for worker"
echo "CFG_WORKER=$CFG_WORKER"
echo
echo '# service port for external client'
echo "FILER_GATEWAY_EXTERNAL_PORT=$FILER_GATEWAY_EXTERNAL_PORT"
