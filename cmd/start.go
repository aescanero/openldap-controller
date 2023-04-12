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
	"log"
	"os"

	"github.com/aescanero/openldap-node/service"
	"github.com/aescanero/openldap-node/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	database  service.DatabaseConfig
	srvConfig service.ServerConfig

	baseEnv            string = utils.GetEnv("LDAP_BASE", "dc=example")
	adminPasswordEnv   string = utils.GetEnv("LDAP_ADMIN_PASSWORD", utils.Random(10))
	replicaPasswordEnv string = utils.GetEnv("LDAP_REPLICA_PASSWORD", utils.Random(10))
	portEnv            string = utils.GetEnv("LDAP_PORT", "1389")
	debugEnv           string = utils.GetEnv("LDAP_DEBUG", "256")
	config             string
	configFile         string
)

func init() {
	startCmd.Flags().StringVarP(&database.Base, "base", "", "dc=example", "LDAP base RDN")
	startCmd.Flags().StringVarP(&config, "config", "", "", "Yaml config")
	startCmd.Flags().StringVarP(&configFile, "config_file", "", "", "Yaml config file")
	startCmd.Flags().StringVarP(&srvConfig.AdminPassword, "admin_password", "", "", "LDAP admin Password (for cn=admin, base DN) Has priority over file")
	startCmd.Flags().StringVarP(&srvConfig.ReplicaPassword, "replica_password", "", "", "LDAP replica Password (for cn=replica, base DN) Has priority over file")
	startCmd.Flags().StringVarP(&srvConfig.AdminPasswordFile, "admin_password_file", "", "", "File with LDAP admin Password (for cn=admin, base DN)")
	startCmd.Flags().StringVarP(&srvConfig.ReplicaPasswordFile, "replica_password_file", "", "", "File with LDAP replica Password (for cn=replica, base DN)")
	startCmd.Flags().StringVarP(&srvConfig.Debug, "debug", "d", "", "Openldap debug (default 256 = Show all queries)")

	startCmd.Flags().StringVarP(&srvConfig.LdapPort, "ldap_port", "", "", "LDAP port")
	startCmd.Flags().StringVarP(&srvConfig.Srvtls.LdapsPort, "ldaps_port", "", "", "LDAPS port")
	startCmd.Flags().StringVarP(&srvConfig.Srvtls.LdapsTls.Ca, "ca", "", "", "CA certificate. Has priority over file")
	startCmd.Flags().StringVarP(&srvConfig.Srvtls.LdapsTls.Crt, "crt", "", "", "CERT certificate. Has priority over file")
	startCmd.Flags().StringVarP(&srvConfig.Srvtls.LdapsTls.CrtKey, "crt_key", "", "", "CERT Private Key. Has priority over file")
	startCmd.Flags().StringVarP(&srvConfig.Srvtls.LdapsTls.CaFile, "ca_file", "", "", "File with CA certificate")
	startCmd.Flags().StringVarP(&srvConfig.Srvtls.LdapsTls.CrtFile, "crt_file", "", "", "File with CERT certificate")
	startCmd.Flags().StringVarP(&srvConfig.Srvtls.LdapsTls.CrtKeyFile, "crt_key_file", "", "", "File with CERT Private Key")

	startCmd.Flags().StringVarP(&database.Replicatls.ReplicaUrl, "replica_url", "", "ldaps://ldap.example.com", "LDAP base RDN")
	startCmd.Flags().StringVarP(&database.Replicatls.LdapsTls.Ca, "replica_ca", "", "", "CA certificate for Replica. Has priority over file")
	startCmd.Flags().StringVarP(&database.Replicatls.LdapsTls.Crt, "replica_crt", "", "", "CERT certificate for Replica. Has priority over file")
	startCmd.Flags().StringVarP(&database.Replicatls.LdapsTls.CrtKey, "replica_crt_key", "", "", "CERT Private Key for Replica. Has priority over file")
	startCmd.Flags().StringVarP(&database.Replicatls.LdapsTls.CaFile, "replica_ca_file", "", "", "File with CA certificate for Replica")
	startCmd.Flags().StringVarP(&database.Replicatls.LdapsTls.CrtFile, "replica_crt_file", "", "", "File with CERT certificate for Replica")
	startCmd.Flags().StringVarP(&database.Replicatls.LdapsTls.CrtKeyFile, "replica_crt_key_file", "", "", "File with CERT Private Key for Replica")

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

		myConfig := service.Config{}
		if configFile != "" {
			data, err := os.ReadFile(configFile)
			if err != nil {
				log.Fatalf("error: %v", err)
				panic(err.Error())
			}
			err = yaml.Unmarshal(data, myConfig)
			if err != nil {
				log.Fatalf("error: %v", err)
				panic(err.Error())
			}
		} else {
			err := yaml.Unmarshal([]byte(config), myConfig)
			if err != nil {
				log.Fatalf("error: %v", err)
				panic(err.Error())
			}
		}

		myConfig.SrvConfig.ImportNotNull(&srvConfig)
		myConfig.Database[0].ImportNotNull(&database)

		fmt.Println("Base: " + myConfig.Database[0].Base)
		fmt.Println("Port: " + myConfig.SrvConfig.LdapPort)
		service.Start(myConfig)
	},
}
