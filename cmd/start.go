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

	"github.com/aescanero/openldap-node/service"
	"github.com/aescanero/openldap-node/utils"
	"github.com/spf13/cobra"
)

var (
	baseEnv             string = utils.GetEnv("LDAP_BASE", "dc=example")
	adminPasswordEnv    string = utils.GetEnv("LDAP_ADMIN_PASSWORD", utils.Random(10))
	replicaPasswordEnv  string = utils.GetEnv("LDAP_REPLICA_PASSWORD", utils.Random(10))
	portEnv             string = utils.GetEnv("LDAP_PORT", "1389")
	debugEnv            string = utils.GetEnv("LDAP_DEBUG", "256")
	base                string
	adminPassword       string
	port                string
	debug               string
	config              string
	configFile          string
	replicaPassword     string
	adminPasswordFile   string
	replicaPasswordFile string
	ldapPort            string
	ldapsPort           string
	ca                  string
	crt                 string
	crtKey              string
	caFile              string
	crtFile             string
	crtKeyFile          string
	replicaUrl          string
	replicaCa           string
	replicaCrt          string
	replicaCrtKey       string
	replicaCaFile       string
	replicaCrtFile      string
	replicaCrtKeyFile   string
)

func init() {
	startCmd.Flags().StringVarP(&base, "base", "", "dc=example", "LDAP base RDN")
	startCmd.Flags().StringVarP(&config, "config", "", "", "Yaml config")
	startCmd.Flags().StringVarP(&configFile, "config_file", "", "", "Yaml config file")
	startCmd.Flags().StringVarP(&adminPassword, "adminpassword", "", "", "LDAP admin Password (for cn=admin, base DN)")
	startCmd.Flags().StringVarP(&replicaPassword, "replicapassword", "", "", "LDAP replica Password (for cn=replica, base DN)")
	startCmd.Flags().StringVarP(&adminPasswordFile, "adminpassword_file", "", "", "File with LDAP admin Password (for cn=admin, base DN)")
	startCmd.Flags().StringVarP(&replicaPasswordFile, "replicapassword_file", "", "", "File with LDAP replica Password (for cn=replica, base DN)")
	startCmd.Flags().StringVarP(&debug, "debug", "d", "", "Openldap debug (default 256 = Show all queries)")

	startCmd.Flags().StringVarP(&ldapPort, "ldap_port", "", "", "LDAP port")
	startCmd.Flags().StringVarP(&ldapsPort, "ldaps_port", "", "", "LDAPS port")
	startCmd.Flags().StringVarP(&ca, "ca", "", "", "CA certificate")
	startCmd.Flags().StringVarP(&crt, "crt", "", "", "CERT certificate")
	startCmd.Flags().StringVarP(&crtKey, "crt_key", "", "", "CERT Private Key")
	startCmd.Flags().StringVarP(&caFile, "ca_file", "", "", "File with CA certificate")
	startCmd.Flags().StringVarP(&crtFile, "crt_file", "", "", "File with CERT certificate")
	startCmd.Flags().StringVarP(&crtKeyFile, "crt_key_file", "", "", "File with CERT Private Key")

	startCmd.Flags().StringVarP(&base, "base", "", "dc=example", "LDAP base RDN")
	startCmd.Flags().StringVarP(&replicaUrl, "replica_url", "", "ldaps://ldap.example.com", "LDAP base RDN")
	startCmd.Flags().StringVarP(&replicaCa, "replica_ca", "", "", "CA certificate for Replica")
	startCmd.Flags().StringVarP(&replicaCrt, "replica_crt", "", "", "CERT certificate for Replica")
	startCmd.Flags().StringVarP(&replicaCrtKey, "replica_crt_key", "", "", "CERT Private Key for Replica")
	startCmd.Flags().StringVarP(&replicaCaFile, "replica_ca_file", "", "", "File with CA certificate for Replica")
	startCmd.Flags().StringVarP(&replicaCrtFile, "replica_crt_file", "", "", "File with CERT certificate for Replica")
	startCmd.Flags().StringVarP(&replicaCrtKeyFile, "replica_crt_key_file", "", "", "File with CERT Private Key for Replica")

	/*   databases:
	  - base: ...
		replicas:
		- url:
		  ca.pem: ...
		  crt.pem: ...
		  crt.key: ...
		  ca.pem.file: ...
		  crt.pem.file: ...
		  crt.key.file: ...
		  attrs: "*,+"
	*/
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
