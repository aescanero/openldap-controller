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
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/aescanero/openldap-controller/utils"
	"github.com/go-ldap/ldap"
)

//go:embed templates/slapd.conf.tmpl
var slapdConfTemplate string

//go:embed templates/base.ldif.tmpl
var baseLdifTemplate string

func Start(base, adminPassword, port, debug string) {
	var wg sync.WaitGroup
	pid := make(chan string)
	stateError := make(chan error)
	//stateError <- nil
	//pid <- ""

	createConfiguration(base, adminPassword, port, debug)

	wg.Add(1)

	go func() {
		out, _ := exec.Command("/usr/sbin/slapd", "-d", "256", "-F", "/etc/ldap/slapd.d", "-h", "ldap://0.0.0.0:1389").Output()
		fmt.Printf("RES: %s\n", out)
		stateError <- errors.New("Openldap Ended")
	}()

	go func() {
		for <-pid == "" {
			time.Sleep(100)
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

	go func() {
		for {
			time.Sleep(100)
			if <-stateError != nil {
				log.Print(<-stateError)
				wg.Done()
			}
		}
	}()

	wg.Wait()

	log.Print("Openldap Terminated")
}

func createConfiguration(base, adminPassword, port, debug string) {

	//var conf embed.FS

	encode := utils.Encode{}
	adminPasswordSHA := encode.MakeSSHAEncode([]byte(adminPassword))

	config := map[string]string{
		"ldapRoot":                         base,
		"ldapEncryptedConfigAdminPassword": "{SSHA}" + base64.StdEncoding.EncodeToString(adminPasswordSHA),
	}
	fmt.Println("base64 password")
	fmt.Println(base64.StdEncoding.EncodeToString(adminPasswordSHA))

	slapdConf, err := template.New("slapdConf").Parse(slapdConfTemplate) //template.ParseFS(conf, "templates/slapd.conf.tmpl")
	if err != nil {
		log.Fatal("Error loading templates:" + err.Error())
	}

	f, err := os.Create("/tmp/slapd.conf")
	if err != nil {
		log.Print("Can't create ", "/tmp/slapd.conf")
	}

	err = slapdConf.Execute(io.Writer(f), config)
	if err != nil {
		log.Print("Can't execute ", "templates/slapd.conf.tmpl")
	}

	err = utils.CreateDirs([]string{"/etc/ldap", "/etc/ldap/slapd.d", "/var/lib/ldap/0", "/etc/ldap/schema"})
	if err != nil {
		log.Println(err)
	}

	err = utils.CopyFiles(
		[]string{"/etc/openldap/schema/core.schema",
			"/etc/openldap/schema/cosine.schema",
			"/etc/openldap/schema/misc.schema",
			"/etc/openldap/schema/inetorgperson.schema",
			"/etc/openldap/schema/nis.schema"},
		"/etc/ldap/schema",
	)
	if err != nil {
		log.Println(err)
	}

	out, _ := exec.Command("/usr/sbin/slaptest", "-f", "/tmp/slapd.conf", "-F", "/etc/ldap/slapd.d").Output()
	fmt.Printf("RES %s\n", out)

	out, err = exec.Command("ls", "-l", "/etc/ldap").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Directory %s\n", out)

	f, err = os.Create("/tmp/base.ldif")
	if err != nil {
		log.Print("Can't create ", "/tmp/base.ldif")
	}

	parsedDN, err := ldap.ParseDN(base)
	if err != nil || len(parsedDN.RDNs) == 0 {
		fmt.Println(err)
	}
	switch parsedDN.RDNs[0].Attributes[0].Type {
	case "o":
		baseLdifTemplate = `dn: ` + base + `
objectClass: organization
o: ` + parsedDN.RDNs[0].Attributes[0].Value + "\n\n" + baseLdifTemplate

	case "dc":
		baseLdifTemplate = `dn: ` + base + `
objectClass: dcObject
dc: ` + parsedDN.RDNs[0].Attributes[0].Value + "\n\n" + baseLdifTemplate

	}
	fmt.Println("Ldif: " + baseLdifTemplate)

	baseLdap, err := template.New("baseLdap").Parse(baseLdifTemplate) //template.ParseFS(conf, "templates/base.ldif.tmpl")
	if err != nil {
		log.Fatal("Error loading templates:" + err.Error())
	}

	config = map[string]string{
		"ldapRoot": base,
	}

	err = baseLdap.Execute(io.Writer(f), config)
	if err != nil {
		log.Print("Can't execute ", "/tmp/base.ldif")
	}

	out, _ = exec.Command("/usr/sbin/slapadd", "-F", "/etc/ldap/slapd.d", "-l", "/tmp/base.ldif").Output()
	fmt.Printf("RES: %s\n", out)

	log.Print("Configuring Openldap")
}
