#!/bin/sh
case "$1" in
  "api")
    echo "Starting API service..."
    exec ./api
    ;;
  "job")
    echo "Starting Job service..."
    exec ./job
    ;;
  *)
    echo "Usage: $0 {api|job}"
    echo "Starting API service by default..."
    exec ./api
    ;;
esac