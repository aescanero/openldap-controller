#!/bin/sh

podman run -v ./standalone:/config --name node --entrypoint '["./controller","start","--config_file=/config/config.yaml"]' --rm -it localhost/openldap-node:0.1.1
