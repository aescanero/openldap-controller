package ldaputils

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aescanero/openldap-node/config"
)

func Start(myConfig config.Config) {
	portStr := ""
	if myConfig.SrvConfig.LdapPort != "" {
		portStr = "ldap://0.0.0.0:" + myConfig.SrvConfig.LdapPort + "/"
	}
	if myConfig.SrvConfig.Srvtls.LdapsPort != "" {
		if portStr != "" {
			portStr = portStr + " "
		}
		portStr = portStr + "ldaps://0.0.0.0:" + myConfig.SrvConfig.Srvtls.LdapsPort + "/"
	} else {
		portStr = portStr + ""
	}
	debug := myConfig.SrvConfig.Debug
	app := "/usr/sbin/slapd"
	args := []string{"-d", debug, "-F", "/etc/ldap/slapd.d", "-h", portStr}
	log.Println("Starting Openldap: " + app + " " + strings.Join(args[:], " "))
	cmd := exec.Command(app, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal("Can't execute " + app + " " + strings.Join(args[:], " ") + " cause: " + err.Error())
	}
}
