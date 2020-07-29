#!/bin/sh
docker run -v $(git rev-parse --show-toplevel)/examples/elasticsearch/input:/config -v $(git rev-parse --show-toplevel)/examples/elasticsearch/output:/output protoc
