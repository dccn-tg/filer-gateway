#!/bin/sh
# pull latest image version
if [ "$LAUNCH_PULL" = true ]; then

    # does a docker login first
    if [ -n "${LAUNCH_IMAGE_REPO_USER}" ] && [ -n "${LAUNCH_IMAGE_REPO_PASS}" ]; then
        echo "Logging in"
        echo "${LAUNCH_IMAGE_REPO_PASS}" | docker login -u "${LAUNCH_IMAGE_REPO_USER}" --password-stdin "${LAUNCH_IMAGE_REPO}"
    fi

    echo "Pulling $LAUNCH_IMAGE: docker pull $LAUNCH_IMAGE"
    docker pull $LAUNCH_IMAGE
fi

# build launch parameters
DOCKER_ARGS="run --rm"
[ -n "$LAUNCH_CONTAINER_NAME" ] && DOCKER_ARGS="$DOCKER_ARGS --name $LAUNCH_CONTAINER_NAME"
[ "$LAUNCH_PRIVILEGED" = true ] && DOCKER_ARGS="$DOCKER_ARGS --privileged"
[ "$LAUNCH_INTERACTIVE" = true ] && DOCKER_ARGS="$DOCKER_ARGS -i"
[ "$LAUNCH_TTY" = true ] && DOCKER_ARGS="$DOCKER_ARGS -t"
[ "$LAUNCH_HOST_NETWORK" = true ] && DOCKER_ARGS="$DOCKER_ARGS --net host"
[ "$LAUNCH_PRIVILEGED" = true ] && DOCKER_ARGS="$DOCKER_ARGS --privileged"
DOCKER_ARGS="$DOCKER_ARGS $LAUNCH_ENVIRONMENT $LAUNCH_VOLUMES $LAUNCH_EXTRA_ARGS $LAUNCH_IMAGE $LAUNCH_IMAGE_COMMAND $LAUNCH_IMAGE_ARGS"

echo "Running $LAUNCH_IMAGE: exec docker $DOCKER_ARGS"
exec docker $DOCKER_ARGS
