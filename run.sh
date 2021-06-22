#!/bin/bash

# Build the docker image and run it
# HLS playback will be available on port 8081
# The RTMP server will be available on port 1935
docker build -t livestream-rtmp . && docker run --rm -p 8081:8081 -p 1935:1935 -it livestream-rtmp
