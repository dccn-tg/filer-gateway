# official Docker (CLI) image
FROM docker:latest

# launch parameters
ENV LAUNCH_IMAGE            hello-world
ENV LAUNCH_PULL             false
ENV LAUNCH_CONTAINER_NAME=
ENV LAUNCH_PRIVILEGED       false
ENV LAUNCH_INTERACTIVE      false
ENV LAUNCH_TTY              false
ENV LAUNCH_HOST_NETWORK     false
ENV LAUNCH_ENVIRONMENT=
ENV LAUNCH_VOLUMES=
ENV LAUNCH_EXTRA_ARGS=
ENV LAUNCH_IMAGE_REPO       hub.docker.com
ENV LAUNCH_IMAGE_REPO_USER=
ENV LAUNCH_IMAGE_REPO_PASS=
ENV LAUNCH_IMAGE_COMMAND=
ENV LAUNCH_IMAGE_ARGS=

# add entrypoint.sh launcher script
ADD entrypoint.sh   /

# run the image
ENTRYPOINT /entrypoint.sh
