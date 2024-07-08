#!/bin/bash

# Default values for parameters
SEG_SIZE=2000
PORT=8094
TEST=true
SIMULATION_MODE=true

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -segsize) SEG_SIZE="$2"; shift ;;
        -port) PORT="$2"; shift ;;
        -test) TEST="$2"; shift ;;
        -simulation) SIMULATION_MODE="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# Build the project
ego-go build -buildvcs=false

# Sign the enclave.json
ego sign enclave.json

# Run the application with parameterized arguments and simulation mode
if [ "$SIMULATION_MODE" = true ]; then
    OE_SIMULATION=1  ego run ./app -segsize "$SEG_SIZE" -port "$PORT" -test "$TEST"
else
     ego run ./app -segsize "$SEG_SIZE" -port "$PORT" -test "$TEST"
fi