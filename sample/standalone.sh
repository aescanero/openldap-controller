#!/bin/sh

podman run -v ./standalone:/config --name micropki --entrypoint '["./micropki","cert","new","--cafile","/config/ca.crt","--cakeyfile","/config/ca.key","--certfile","/config/cert.crt","--certkeyfile","/config/cert.key"]' --rm -it docker.io/aescanero/micropki:0.1.2-linux-amd64

podman run -v ./standalone:/config --name node --entrypoint '["./controller","start","--config_file=/config/config.yaml"]' -p 9090:9090,1389:1389,1686:1686 --rm -it docker.io/aescanero/openldap-node:0.1.3-linux-amd66

sleep 10

ldapsearch -H ldap://127.0.0.1:1389 -x -w password -D "cn=admin,dc=example,dc=org" -b "dc=example,dc=org"
