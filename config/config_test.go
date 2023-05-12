package config

import (
	"testing"

	"gopkg.in/yaml.v2"
)

var configYaml string = `srvconfig:
    debug: "256"
    adminpasswordfile: /config/passfile
    ldapport: "1389"
    srvtls:
        ldapsport: "1686"
        ldapstls:
            cafile: "/config/ca.crt"
            crtfile: "/config/cert.crt"
            crtkeyfile: "/config/cert.key"
database:
    - base: "dc=example,dc=org"
schemas:
    - path: "/config/schemas/guacConfigGroup.schema"`

func TestConfig(t *testing.T) {
	myConfig := Config{}
	err := yaml.Unmarshal([]byte(configYaml), &myConfig)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	got := myConfig.SrvConfig.LdapPort
	expected := "1389"

	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}
