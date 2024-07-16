#!/bin/bash

# Default values for parameters

PORT=8087

# Default values for the motivating scenario event logs
LOG="testing_logs/motivating/hospital.xes"
#LOG="testing_logs/motivating/specialized.xes"
#LOG="testing_logs/motivating/pharma.xes"

# Default values for the the merge key storing the case id of the event log
MERGEKEY="concept:name"
# Default values for the Secure Miner's measurement
MEASUREMENT=$(ego uniqueid app)
#If set to false, the provisioner will skip the remote attestation
SKIPATTESTATION=true

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -port) PORT="$2"; shift ;;
        -log) LOG="$2"; shift ;;
        -mergekey) MERGEKEY="$2"; shift ;;
        -measurement) MEASUREMENT="$2"; shift ;;
        -skipattestation) SKIPATTESTATION="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# Build the project
CGO_CFLAGS=-I/opt/ego/include CGO_LDFLAGS=-L/opt/ego/lib go build -o logprovision provisioner/log-provision/log_provision.go

# Run the application with parameterized arguments
./logprovision -port "$PORT" -log "$LOG" -mergekey "$MERGEKEY" -measurement "$MEASUREMENT" -skipattestation "$SKIPATTESTATION"
