#!/bin/bash

set -ex

#create a registry to store the actual image rootfs of LXR
mkdir -p /home/LXR/LXR-registry/$IMAGE/rootfs

LXR_IMAGE_REG=/home/LXR/LXR-registry/$IMAGE


# request an token for the specific repo to fetch oci images
#-s ignores the progress bar and error msg ,just give the json token response and that passed as input to the jq(JSON parser) that parse that JSON response and return only the token field (-r means return raw data (without any quotes)
TOKEN=$(curl -s \
"https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/$IMAGE:pull" \
| jq -r .token)


#Add token in http header and request with it to get manifest list
#manifest list is an list of references (digests) to image manifests for different platforms & architecture 
#the manifest JSON response is stored to the manifest.json
curl -s 
-H "Authorization: Bearer $TOKEN" \
-H "Accept: application/vnd.docker.distribution.manifest.v2+json, application/vnd.docker.distribution.manifest.list.v2+json" \
https://registry-1.docker.io/v2/library/$IMAGE/manifests/latest > manifest.json


#extract only the manifest json that has arm64 and linux in manifest list from the manifest.json using jq 
#extract the digest (unique ID for manifest) from the retrieved manifest json 
MANIFEST_DIGEST=$(jq -r '.manifests[] | select(.platform.architecture == "arm64" and .platform.os == "linux") | .digest' manifest.json)


#request manifest for the specific platform & architecure with its digest
#JSON response contains layers list that is known as blob (Binary Large Object) ,each layers has binary of the actual rootfs
#iterate layers list and extract only the digest for each layers and that has been stored to the variable LAYERS
LAYERS=$(curl -s \
-H "Authorization: Bearer $TOKEN" \
-H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
https://registry-1.docker.io/v2/library/$IMAGE/manifests/$MANIFEST_DIGEST\
| jq -r '.layers[].digest')


#for loop to iterate the digest list
COUNT=1
for DIGEST in $LAYERS;do
    LAYER_NAME="$LXR_IMAGE_REG/layer$COUNT.tar.gz"

    #request blob for each digest in a list and store it as a layerx.tar.gz
    curl -L -s \
    -H "Authorization: Bearer $TOKEN" \
    "https://registry-1.docker.io/v2/library/$IMAGE/blobs/$DIGEST" \
    -o "$LAYER_NAME"

    #untar the layer and store it to rootfs
    tar -xf $LAYER_NAME -C $LXR_IMAGE_REG/rootfs

    rm $LAYER_NAME
   (( COUNT++))
done

