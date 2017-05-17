#!/bin/bash
#
# Simple benchmark test suite
#
# You must have installed vegeta:
# go get github.com/tsenart/vegeta
#

# Default port to listen
port=8088

# Start the server
./bin/imaginary -p $port & > /dev/null
pid=$!

suite() {
  echo "$1 --------------------------------------"
  echo "GET http://10.60.6.12:8068/$2" | vegeta attack \
    -duration=100s \
    -rate=50 \ | vegeta report
  sleep 1
}

# Run suites
suite "Avatar" "profile_avatar/19384"

# Kill the server
kill -9 $pid
