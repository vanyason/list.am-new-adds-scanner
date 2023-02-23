#!/bin/bash -e

IMAGE_NAME="listam-parser"
CONTAINER_NAME="listam-parser"

help()
{
   # Display Help
   echo "Usage: ./run.sh [ build | start | stop | stats | att ]"
   echo
   echo "build      Build image."
   echo "start      Start container."
   echo "stats      See containers related info."
   echo "stop       Stop container."
   echo "att        Attach to container."
}

build() {
    if [ ! "$(docker images | grep "$IMAGE_NAME")" ]; then
        echo "Building image..."
        docker build . -t "$IMAGE_NAME"
    else
        echo "Image is built already"
    fi
}

start() {
    if [ ! "$(docker ps --all | grep "$CONTAINER_NAME")" ]; then
        echo "Starting container..."
        docker run -d -it --entrypoint /bin/bash -v"$(pwd)":/app --name $CONTAINER_NAME $IMAGE_NAME
    else
        echo "Container is already running"
    fi
}

stats() {
    docker ps --all | grep "$CONTAINER_NAME"
}

stop() {    
    if [ "$(stats)" ]; then 
        ID=$(docker ps --all | grep "$CONTAINER_NAME" | awk '{ print $1 }')
        docker stop $ID && docker rm $ID
    else
        echo "Nothing to stop..."
    fi
}

att() {
    if [ "$(stats)" ]; then 
        docker attach "$CONTAINER_NAME"
    else
        echo "Nowhere to attach"
    fi
}

##############
### Main #####
##############

# Doesn`t work wtf ?`
# cd "$(dirname "$(readlink -f "$0")")"
# cd "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if  [ "$1" == "build" ]; then
    build
elif [ "$1" == "start" ]; then
    start
elif [ "$1" == "stats" ]; then
    stats
elif [ "$1" == "stop" ]; then
    stop
elif [ "$1" == "att" ]; then
    att
else 
    help
fi
