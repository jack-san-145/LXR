#!/bin/bash

set -ex

# request an token to fetch oci images
#-s ignores the progress bar and error msg ,just give the json token response and that passed as input to the jq(JSON parser) that parse that JSON response and return only the token field (-r means return raw data (without any quotes)
TOKEN=$(curl -s \
"https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/$IMAGE:pull" \
| jq -r .token)


