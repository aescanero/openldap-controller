package ldaputils

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"log"
	"testing"
	"text/template"

	"github.com/aescanero/openldap-node/config"
	"github.com/aescanero/openldap-node/utils"
	"gopkg.in/yaml.v2"
)

//go:embed templates/slapd.conf.test.tmpl
var slapdConfTemplateTest string

var configYaml string = `srvconfig:
    debug: "256"
    adminpassword: random
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

func TestBase(t *testing.T) {
	myConfig := config.Config{}
	err := yaml.Unmarshal([]byte(configYaml), &myConfig)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	got := "cn=admin," + myConfig.Database[0].Base

	expected := "cn=admin,dc=example,dc=org"

	if got != expected {
		t.Errorf("got %s, expected %s", got, expected)
	}
}

func TestGetPassword(t *testing.T) {
	myConfig := config.Config{}
	err := yaml.Unmarshal([]byte(configYaml), &myConfig)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	got, _ := myConfig.SrvConfig.GetAdminPassword()

	expected := "random"

	if got != expected {
		t.Errorf("got %s, expected %s", got, expected)
	}
}

func TestTemplateMatch(t *testing.T) {
	myConfig := config.Config{}
	err := yaml.Unmarshal([]byte(configYaml), &myConfig)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	myConfig.SrvConfig.AdminPasswordSHA = "test"

	slapdConf, err := template.New("slapdConf").Parse(slapdConfTemplateTest) //template.ParseFS(conf, "templates/slapd.conf.tmpl")
	if err != nil {
		log.Fatal("Error loading templates:" + err.Error())
		panic(err)
	}

	var tpl bytes.Buffer
	slapdConf.Execute(&tpl, myConfig)
	fmt.Printf("slapdConf: %s\n", tpl.String())

	got := tpl.String()
	expected := "rootpw test"

	if got != expected {
		t.Errorf("got %s, expected %s", got, expected)
	}
}

func TestPasswordMatch(t *testing.T) {
	myConfig := config.Config{}
	err := yaml.Unmarshal([]byte(configYaml), &myConfig)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	encode := utils.Encode{}
	passwordTest, _ := myConfig.SrvConfig.GetAdminPassword()
	fmt.Printf("passwordTest: %s\n", passwordTest)
	passwordSSHA := encode.MakeSSHAEncode([]byte(passwordTest))
	passwordSSHAB64 := fmt.Sprintf("{SSHA}%s", base64.StdEncoding.EncodeToString(passwordSSHA))
	fmt.Printf("passwordSSHA: %s\n", string(passwordSSHA))
	fmt.Printf("passwordSSHAB64: %s\n", passwordSSHAB64)
	myConfig.SrvConfig.AdminPasswordSHA = passwordSSHAB64

	slapdConf, err := template.New("slapdConf").Parse(slapdConfTemplateTest) //template.ParseFS(conf, "templates/slapd.conf.tmpl")
	if err != nil {
		log.Fatal("Error loading templates:" + err.Error())
		panic(err)
	}

	var tpl bytes.Buffer
	slapdConf.Execute(&tpl, myConfig)
	fmt.Printf("slapdConf: %s\n", tpl.String())

	got := encode.Matches([]byte(passwordSSHAB64), []byte(passwordTest), true)
	expected := true

	if got != expected {
		t.Errorf("got %t, expected %t", got, expected)
	}
}
