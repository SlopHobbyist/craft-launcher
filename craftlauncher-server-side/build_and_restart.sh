#!/bin/bash

# Run the build and copy script from its directory
echo "Running build_and_copy.sh..."
(cd /Users/farfur/Downloads/clientside-mod && ./build_and_copy.sh)

# Check if the previous command succeeded
if [ $? -eq 0 ]; then
    echo "Build and copy completed successfully."
    echo "Running restart.sh..."
    (cd /Users/farfur/Downloads/macserverside/local && ./restart.sh)
else
    echo "Build and copy failed. Not running restart.sh."
    exit 1
fi
