#!/bin/bash

if [ -f .env ]; then
    # Load Environment Variables
    export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
    echo $PORT && echo $CUSTOMURL
fi