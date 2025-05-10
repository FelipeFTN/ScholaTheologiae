#!/bin/bash

# This script runs the web application

export HOST="0.0.0.0" BINDING="0.0.0.0" PORT=3000

cd ./app && bin/rails server -b $BINDING -p $PORT
