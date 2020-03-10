#!/bin/bash

echo "# version"
echo "DOCKER_IMAGE_TAG=$DOCKER_IMAGE_TAG"
echo 
echo "# docker registry endpoint"
echo "DOCKER_REGISTRY=$DOCKER_REGISTRY"
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
echo '# service port for external client'
echo "FILER_GATEWAY_EXTERNAL_PORT=$FILER_GATEWAY_EXTERNAL_PORT"