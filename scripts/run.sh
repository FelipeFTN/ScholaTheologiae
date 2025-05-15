#!bin/bash

### Here's the recommended way to run the script
### from the command line, with timestamp logging
### and a log file to capture output and errors

time python ./summa_processor.py 2>&1 | tee $(date +%s).log
