#!/bin/bash
PORTS=(8000 8001 8002 8003 8004 8005 8006 8007 8008 8009)

for PORT in "${PORTS[@]}"
do
  PORT=$PORT ./gocommerce & # Start each server in the background
done

