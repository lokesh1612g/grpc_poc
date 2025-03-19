#!/bin/bash

# Check if client ID is passed as a parameter
if [ -z "$1" ]; then
  echo "Usage: $0 <client_id>"
  exit 1
fi

client_id=$1

# Print the client ID
echo "executed script for Client ID: $client_id"