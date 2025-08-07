#!/bin/bash
# Wrapper script to run tests with Go 1.24 compatibility

set -e

# Compile the test binary
go test -c

# Run the test binary directly to avoid flag conflicts
./dawg.test "$@"

# Clean up
rm -f dawg.test
