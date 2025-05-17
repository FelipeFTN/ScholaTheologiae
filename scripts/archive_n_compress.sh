#!/bin/bash

# This script archives and compresses the pdf files in the current directory.
# It creates a tar.gz archive with the name "archive.tar.gz" and includes all the necessary pdf files

# Usage: ./archive_n_compress.sh
# Make sure to run this script in the directory containing the pdf files.
# The script will create a tar.gz archive named "archive.tar.gz" in the current directory.
# This file will be commited to the git repository.

# Check if there are any pdf files in the current directory
if ls ../books/*_Vol_*.pdf 1> /dev/null 2>&1; then
  # Create a tar.gz archive with the name "archive.tar.gz"
  # The archive will include all the pdf files that match the pattern *_Vol_*.pdf
  # The archive will be compressed using gzip with the highest compression level
  export GZIP=-9
  tar -cvzf ../archived/books.tar.gz ../books/*_Vol_*.pdf
  echo "Archive created successfully: books.tar.gz"
else
  echo "No pdf files found matching the pattern *_Vol_*.pdf"
  exit 1
fi

# Compress the book folder 
if ls ../books/summa_theologiae/* 1> /dev/null 2>&1; then
  # Create a tar.gz archive with the name "summa_theologiae.tar.gz"
  # The archive will include all the md files that was generated using the scripts
  # The archive will be compressed using gzip with the highest compression level
  export GZIP=-9
  tar -cvzf ../books/summa_theologiae.tar.gz data/summa_theologiae/ data/summa_theologiae.db
  echo "Archive created successfully: summa_theologiae.tar.gz"
else
  echo "No pdf files found matching the pattern summa_theologiae/*"
  exit 1
fi
