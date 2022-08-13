#!/bin/bash

SERVICE_NAME="chat-service"

getpid() {
  ps -ef | grep -w ${SERVICE_NAME} | grep -v grep | awk '{print $2}'
}

start() {
  pid=$(getpid)
  [ ! -z "$pid" ] && {
    echo -e "$SERVICE_NAME has already started"
    return
  }
  nohup ./$SERVICE_NAME >/dev/null 2>&1 &
}

stop() {
  pid=$(getpid)
  kill $pid 1>/dev/null
}

usage() {
  echo "Usage:$(SERVICE_NAME $0) <start|stop|restart>"
  exit 1
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    stop
    start
    ;;
  *)
    usage
esac

