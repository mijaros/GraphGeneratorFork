#!/bin/bash

function stop_dependent() {
  if [ -n "$PID" ]; then
    kill -TERM "$PID"
    wait "$PID"
    PID=
  fi
}

TEST_DIR="./tmp/test/run_$(date '+%F-%s')"
DB_DIR="${TEST_DIR}/db"


if [ ! -f ./generator ]; then
  echo "Generator not built run make build"
  exit 1
fi

if [ ! -d ui/dist ]; then
  echo "UI not build run make build"
  exit 1
fi

mkdir -p "$DB_DIR"

export OUT=$(mktemp)
export ERR=$(mktemp)


echo "Starting the application"
echo "$(pwd)"

./generator -dbRoot "$DB_DIR" -test >$OUT 2>$ERR &
PID=$!

trap  stop_dependent TERM INT EXIT
sleep 5
echo "Starting the tests"

bash -c 'cd test && poetry install && poetry run pytest --html=report.html --self-contained-html suite/'


if [ $? -ne 0 ]; then
  echo "Tests failed collecting results $(pwd)"
else
  stop_dependent
  rm -rf "${DB_DIR}"
fi

stop_dependent
mv "$OUT" "${TEST_DIR}/server_stdout.txt"
mv "$ERR" "${TEST_DIR}/server_stderr.txt"
mv test/report.html "${TEST_DIR}/"

echo "Ending tests"
