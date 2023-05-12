/*Copyright [2023] [Alejandro Escanero Blanco <aescanero@disasterproject.com>]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.*/

package service

import (
	_ "embed"
	"encoding/base64"
	"errors"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aescanero/openldap-node/apiserver"
	"github.com/aescanero/openldap-node/config"
	"github.com/aescanero/openldap-node/utils"
	"github.com/go-ldap/ldap"
)

//go:embed templates/slapd.conf.tmpl
var slapdConfTemplate string

//go:embed templates/base.ldif.tmpl
var baseLdifTemplate string

func Start(myConfig config.Config) {
	var wg sync.WaitGroup
	pid := make(chan string)
	stateError := make(chan error)
	//stateError <- nil
	//pid <- ""

	createConfiguration(myConfig)

	wg.Add(1)

	go func() {
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
		println("Starting Openldap: " + app + " " + strings.Join(args[:], " "))
		cmd := exec.Command(app, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Print("Can't execute " + app + " " + strings.Join(args[:], " ") + " cause: " + err.Error())
			panic(err)
		}
		stateError <- errors.New("openldap ended")
	}()

	//Raising LDAP Service
	go func() {
		for <-pid == "" {
			time.Sleep(100 * time.Millisecond)
			source, err := os.Open("/var/lib/ldap/slapd.pid")
			if err != nil {
				stateError <- err
			}
			BUFFERSIZE := 4096
			buf := make([]byte, BUFFERSIZE)
			_, err = source.Read(buf)
			if err != nil && err != io.EOF {
				stateError <- err
			}
			pid <- string(buf)
		}
	}()

	//Raise dashboard
	go func() {
		apiserver.Server(myConfig)
	}()

	//Monitor when Done
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if <-stateError != nil {
				log.Print(<-stateError)
				wg.Done()
			}
		}
	}()

	wg.Wait()

	log.Print("Openldap Terminated")
}

func createConfiguration(myConfig config.Config) error {

	//var conf embed.FS
	base := myConfig.Database[0].Base
	adminPassword, err := myConfig.SrvConfig.GetAdminPassword()
	if err != nil {
		log.Fatal("Error loading config:" + err.Error())
		panic(err)
	}

	encode := utils.Encode{}
	adminPasswordSHA := encode.MakeSSHAEncode([]byte(adminPassword))
	myConfig.SrvConfig.AdminPasswordSHA = "{SSHA}" + base64.StdEncoding.EncodeToString(adminPasswordSHA)

	if myConfig.Database[0].Replicatls[0].ReplicaUrl == "" {
		myConfig.Database[0].Replicatls = nil
	}

	slapdConf, err := template.New("slapdConf").Parse(slapdConfTemplate) //template.ParseFS(conf, "templates/slapd.conf.tmpl")
	if err != nil {
		log.Fatal("Error loading templates:" + err.Error())
		panic(err)
	}

	err = utils.CreateDirs([]string{"/etc/ldap", "/etc/ldap/slapd.d", "/var/lib/ldap/0", "/etc/ldap/schema", "/etc/ldap/certs"})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	schemaFiles := []string{"/etc/openldap/schema/core.schema",
		"/etc/openldap/schema/cosine.schema",
		"/etc/openldap/schema/misc.schema",
		"/etc/openldap/schema/inetorgperson.schema",
		"/etc/openldap/schema/nis.schema"}

	for _, schema := range myConfig.Schemas {
		schemaFiles = append(schemaFiles, schema.Path)
	}

	err = utils.CopyFiles(
		schemaFiles,
		"/etc/ldap/schema",
	)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if myConfig.SrvConfig.Srvtls.LdapsTls.CaFile != "" && myConfig.SrvConfig.Srvtls.LdapsTls.CrtFile != "" && myConfig.SrvConfig.Srvtls.LdapsTls.CrtKeyFile != "" {
		tlsfiles := []string{myConfig.SrvConfig.Srvtls.LdapsTls.CaFile, myConfig.SrvConfig.Srvtls.LdapsTls.CrtFile, myConfig.SrvConfig.Srvtls.LdapsTls.CrtKeyFile}
		err = utils.CopyFiles(
			tlsfiles,
			"/etc/ldap/certs",
		)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		caFilename := filepath.Base(myConfig.SrvConfig.Srvtls.LdapsTls.CaFile)
		crtFilename := filepath.Base(myConfig.SrvConfig.Srvtls.LdapsTls.CrtFile)
		crtKeyFilename := filepath.Base(myConfig.SrvConfig.Srvtls.LdapsTls.CrtKeyFile)
		myConfig.SrvConfig.Srvtls.LdapsTls.CaFile = "/etc/ldap/certs/" + caFilename
		myConfig.SrvConfig.Srvtls.LdapsTls.CrtFile = "/etc/ldap/certs/" + crtFilename
		myConfig.SrvConfig.Srvtls.LdapsTls.CrtKeyFile = "/etc/ldap/certs/" + crtKeyFilename
	} else if myConfig.SrvConfig.Srvtls.LdapsTls.CaFile != "" || myConfig.SrvConfig.Srvtls.LdapsTls.CrtFile != "" || myConfig.SrvConfig.Srvtls.LdapsTls.CrtKeyFile != "" {
		panic(errors.New("the cafile, crtfile, crtkeyfile must be set to obtain TLS support"))
	}

	f, err := os.Create("/tmp/slapd.conf")
	if err != nil {
		log.Print("Can't create ", "/tmp/slapd.conf")
		panic(err)
	}

	err = slapdConf.Execute(io.Writer(f), myConfig)
	if err != nil {
		log.Print("Can't execute ", "templates/slapd.conf.tmpl")
		panic(err)
	}

	log.Print("Populating slapd conf")
	///usr/sbin/slaptest -f /tmp/slapd.conf -F /etc/ldap/slapd.d
	out, _ := exec.Command("/usr/sbin/slaptest", "-f", "/tmp/slapd.conf", "-F", "/etc/ldap/slapd.d").Output()
	log.Printf("RES %s\n", out)

	log.Print("Initializing database")
	f, err = os.Create("/tmp/base.ldif")
	if err != nil {
		log.Print("Can't create ", "/tmp/base.ldif")
		panic(err)
	}

	baseLdifTemplate = `dn: ou=templates,` + base + `
objectClass: organizationalUnit
ou: templates` + "\n\n" + baseLdifTemplate

	parsedDN, err := ldap.ParseDN(base)
	if err != nil || len(parsedDN.RDNs) == 0 {
		log.Println(err)
		panic(err)
	}
	switch parsedDN.RDNs[0].Attributes[0].Type {
	case "o":
		baseLdifTemplate = `dn: ` + base + `
objectClass: organization
o: ` + parsedDN.RDNs[0].Attributes[0].Value + "\n\n" + baseLdifTemplate

	case "dc":
		baseLdifTemplate = `dn: ` + base + `
objectClass: dcObject
objectClass: organization
dc: ` + parsedDN.RDNs[0].Attributes[0].Value + `
o: ` + parsedDN.RDNs[0].Attributes[0].Value + "\n\n" + baseLdifTemplate
	}

	baseLdap, err := template.New("baseLdap").Parse(baseLdifTemplate) //template.ParseFS(conf, "templates/base.ldif.tmpl")
	if err != nil {
		log.Fatal("Error loading templates:" + err.Error())
		panic(err)
	}

	config := map[string]string{
		"ldapRoot": base,
	}

	err = baseLdap.Execute(io.Writer(f), config)
	if err != nil {
		log.Print("Can't execute ", "/tmp/base.ldif")
		panic(err)
	}

	_, err = exec.Command("/usr/sbin/slapadd", "-F", "/etc/ldap/slapd.d", "-l", "/tmp/base.ldif").Output()
	if err != nil {
		log.Fatal("Can't execute /usr/sbin/slapadd -F /etc/ldap/slapd.d -l /tmp/base.ldif")
		panic(err)
	}

	log.Print("Configuring Openldap")
	return nil
}
