#!/bin/bash

# This script runs the API server

export HOST="0.0.0.0" BINDING="0.0.0.0" PORT=3000

cd ./api && go run main.go
