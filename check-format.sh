#!/bin/bash

# Check for formatting issues
unformatted=$(goimports -l .)

if [ -n "$unformatted" ]; then
  echo "The following files are not properly formatted:"
  echo "$unformatted"
  echo "Please run 'goimports -w .' to format your code."
  exit 1
else
  echo "All files are properly formatted."
fi
