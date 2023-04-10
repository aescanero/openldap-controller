FROM docker.io/golang:alpine3.17 AS builder

LABEL org.opencontainers.image.authors="Alejandro Escanero Blanco <alejandro.escanero@accenture.com>"

USER 0

RUN apk --no-cache add ca-certificates && mkdir /data

WORKDIR /data/
COPY . .
#COPY go.sum .
#COPY app.go .

RUN go build -a -installsuffix cgo -o controller .

FROM docker.io/debian:stable-20230227-slim

LABEL org.opencontainers.image.authors="Alejandro Escanero Blanco <alejandro.escanero@accenture.com>"

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y \
        slapd ldap-utils gettext-base procps ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    mv /etc/ldap /etc/openldap && \
    rm -f /var/lib/ldap/*

COPY --from=builder /data/controller /.

VOLUME [ "/etc/ldap" ]
VOLUME [ "/var/lib/ldap" ]

RUN chgrp -R 0 /var/lib/ldap && chmod -R g=u /var/lib/ldap && chmod u+x /var/lib/ldap && \
    chgrp -R 0 /etc/ldap && chmod -R g=u /etc/ldap && chmod u+x /etc/ldap
  
#USER 1001

WORKDIR /

EXPOSE 1389 1636

ENTRYPOINT ["/controller"]
CMD ["start"]
