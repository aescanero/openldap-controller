package config

import (
	"fmt"
	"testing"
)

var configYaml string = `
srvconfig:
  debug: "256"
  admin_password_file: "/config/passfile"
  ldap_port: "1389"
  srvtls:
    ldaps_port: "1686"
    ldaps_tls:
      ca_file: "/config/ca.crt"
      crt_file: "/config/cert.crt"
      crt_key_file: "/config/cert.key"
  database:
    - base: "dc=example,dc=org"
  schemas:
    - path: "/config/schemas/guacConfigGroup.schema"
    - path: "/config/schemas/rfc2307bis.schema"
    - path: "/config/schemas/sudo.schema"
  modules:
    - name: unique
    - name: pw-sha2
    - name: ppolicy
`

func TestConfig(t *testing.T) {
	myConfig := NewConfig()
	_, err := myConfig.fromYaml("", configYaml)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	fmt.Printf("myConfig: %s\n", myConfig)

	got := myConfig.SrvConfig.LdapPort
	expected := "1389"

	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}
