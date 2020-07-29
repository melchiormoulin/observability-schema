#!/bin/sh
docker build $(git rev-parse --show-toplevel) -t protoc
