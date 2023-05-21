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
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aescanero/openldap-node/apiserver"
	"github.com/aescanero/openldap-node/config"
	"github.com/aescanero/openldap-node/ldaputils"
)

func Start(myConfig config.Config) {
	var wg sync.WaitGroup
	pid := make(chan string)
	stateError := make(chan error)
	//stateError <- nil
	//pid <- ""

	err := ldaputils.Prepare(myConfig)
	if err != nil {
		log.Fatalf("Can't prepare ldap: %s", err)
	}

	wg.Add(1)

	go func() {
		ldaputils.Start(myConfig)
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
