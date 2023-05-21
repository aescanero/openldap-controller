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
	"bytes"
	"io/ioutil"
	"strconv"
	"syscall"
)

func Stop() {

	pid, err := ioutil.ReadFile("/var/lib/ldap/slapd.pid")
	if err != nil {
		panic(err)
	}
	ipid, err := strconv.ParseInt(string(bytes.TrimRight(pid, "\n")), 10, 64)
	if err != nil {
		panic(err)
	}
	syscall.Kill(int(ipid), syscall.SIGTERM)
}
