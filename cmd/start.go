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

package cmd

import (
	"fmt"

	"github.com/aescanero/openldap-controller/service"
	"github.com/aescanero/openldap-controller/utils"
	"github.com/spf13/cobra"
)

var (
	baseEnv          string = utils.GetEnv("LDAP_BASE", "dc=example")
	adminPasswordEnv string = utils.GetEnv("LDAP_ADMIN_PASSWORD", utils.Random(10))
	portEnv          string = utils.GetEnv("LDAP_PORT", "1389")
	debugEnv         string = utils.GetEnv("LDAP_DEBUG", "256")
	base             string
	adminPassword    string
	port             string
	debug            string
)

func init() {
	fmt.Println("Port1: " + port)
	startCmd.Flags().StringVarP(&base, "base", "b", "", "LDAP base RDN")
	startCmd.Flags().StringVarP(&adminPassword, "adminpassword", "a", "", "LDAP admin Password (for cn=admin, base DN)")
	startCmd.Flags().StringVarP(&port, "port", "p", "", "LDAP port")
	startCmd.Flags().StringVarP(&debug, "debug", "d", "", "Openldap debug (default 256 = Show all queries)")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Openldap Node",
	Long:  `Start Openldap Node`,
	Run: func(cmd *cobra.Command, args []string) {
		if adminPassword == "" {
			adminPassword = adminPasswordEnv
		}
		if base == "" {
			base = baseEnv
		}
		if port == "" {
			port = portEnv
		}
		if debug == "" {
			debug = debugEnv
		}
		fmt.Println("Base: " + base)
		fmt.Println("Port2: " + port)
		service.Start(base, adminPassword, port, debug)
	},
}
