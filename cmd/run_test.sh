#!/bin/bash
# RUN THIS FROM THE LOCAL DIRECTORY
# Exit immediately if a command exits with a non-zero status.
set -e

# define colors for pretty output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# set environment to LOCAL if not already set
ENVIRONMENT=${ENVIRONMENT:-LOCAL}

# Function to run tests and check for errors
run_tests() {
    local path=$1
    echo -e "${YELLOW}Running tests in $path...${NC}"
    if [ "$ENVIRONMENT" == "LOCAL" ]; then
        if go test -v "$path"; then
            echo -e "${GREEN}Tests passed in $path${NC}"
        else
            echo -e "${RED}Tests failed in $path${NC}"
            exit 1
        fi
    elif [ "$ENVIRONMENT" == "TEST" ]; then
        if go test -json "$path" > test_results.json; then
            echo -e "${GREEN}Tests passed in $path${NC}"
        else
            echo -e "${RED}Tests failed in $path${NC}"
            exit 1
        fi
    else
        echo -e "${RED}Unknown environment: $ENVIRONMENT${NC}"
        exit 1
    fi
}

# set prefix based on environment
if [ "$ENVIRONMENT" == "LOCAL" ]; then
    PREFIX="../"
elif [ "$ENVIRONMENT" == "TEST" ]; then
    PREFIX="./"
else
    PREFIX=""
fi

# domain tests
DOMAIN_TEST_PATH="${PREFIX}internal/domain"

# run all tests
echo -e "${YELLOW}Running all Go tests in environment: $ENVIRONMENT...${NC}"
run_tests "$DOMAIN_TEST_PATH"

echo -e "${GREEN}All tests passed successfully.${NC}"
