#! /usr/bin/sh

# WeatherDB waits at most 30 secs for requests to finish processing, hence 32 as timeout, 2 secs extra to be safe.
#
PID=$1
TIMEOUT=32

if [ -z "$PID" ]; then
    echo "Need to give pid as command line argument"
    exit
fi

kill -2 $PID
timeout $TIMEOUT tail --pid $PID -f /dev/null
