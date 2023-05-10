#!/bin/bash

function stop_servers() {
  echo 'Stopping the services'
  if [ -n "$UI_PID" ]; then
    kill -INT "$UI_PID"
    UI_PID=
  fi

  if [ -n "$API_PID" ]; then
    kill -INT "$API_PID"
    API_PID=
  fi
  wait
}

API_PID=
UI_PID=

trap stop_servers INT TERM EXIT

go run ./cmd/generator/main.go -devel &
API_PID=$!

cd ui && npm install && npm start &
UI_PID=$!

wait
