#!/bin/bash
set -euxo pipefail
docker build "$PWD" -f Dockerfile.dev -t localhost:32000/micro-dev:dkozlov
docker push localhost:32000/micro-dev:dkozlov
