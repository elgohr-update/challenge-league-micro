#!/bin/bash
set -euxo pipefail
docker build "$PWD" -t localhost:32000/micro:dkozlov
docker push localhost:32000/micro:dkozlov
