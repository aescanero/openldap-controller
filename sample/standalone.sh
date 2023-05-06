#!/bin/sh

podman run -v ./standalone:/config --name micropki --entrypoint '["./micropki","cert","new","--cafile","/config/ca.crt","--cakeyfile","/config/ca.key","--certfile","/config/cert.crt","--certkeyfile","/config/cert.key"]' --rm -it docker.io/aescanero/micropki:0.1.2-linux-amd64

podman run -v ./standalone:/config --name node --entrypoint '["./controller","start","--config_file=/config/config.yaml"]' --rm -it docker.io/aescanero/openldap-node:0.1.2-linux-amd64
