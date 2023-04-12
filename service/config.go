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

/*

ports:
  ldap: 1389
  ldaps: 1686
  ca.pem: ...
  crt.pem: ...
  crt.key: ...
  ca.pem.file: ...
  crt.pem.file: ...
  crt.key.file: ...
databases:
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
/*
olcSynrepl:
syncrepl
 ...
 provider=ldaps://ldap.example.com
 bindmethod=simple
 binddn="cn=goodguy,dc=example,dc=com"
 credentials=dirtysecret
 starttls=critical
 schemachecking=on
 scope=sub
 searchbase="dc=example,dc=com"
 tls_cacert=/path/to/file
 tls_cert=/path/to/file.ext
 tls_key=/path/to/file.ext
 tls_protocol_min=1.2
 tls_reqcert=demand
 type=refreshAndPersist

*/
package service

import (
	"errors"
	"os"
)

type tls struct {
	Ca         string
	Crt        string
	CrtKey     string
	CaFile     string
	CrtFile    string
	CrtKeyFile string
}

type serverTls struct {
	LdapsPort string
	LdapsTls  tls
}

type replicaTls struct {
	ReplicaPasswordFile string
	LdapsTls            tls
	ReplicaUrl          string
}

type ServerConfig struct {
	AdminPassword       string
	AdminPasswordFile   string
	ReplicaPassword     string
	ReplicaPasswordFile string
	LdapPort            string
	Srvtls              serverTls
	Debug               string
}

type DatabaseConfig struct {
	Base       string
	Replicatls replicaTls
}

type Config struct {
	SrvConfig ServerConfig
	Database  []DatabaseConfig
}

func (scIn *ServerConfig) ImportNotNull(sc *ServerConfig) {
	if sc.AdminPassword != "" {
		scIn.AdminPassword = sc.AdminPassword
	}
	if sc.AdminPasswordFile != "" {
		scIn.AdminPasswordFile = sc.AdminPasswordFile
	}
	if sc.Debug != "" {
		scIn.Debug = sc.Debug
	}
	if sc.LdapPort != "" {
		scIn.LdapPort = sc.LdapPort
	}
	if sc.ReplicaPassword != "" {
		scIn.ReplicaPassword = sc.ReplicaPassword
	}
	if sc.ReplicaPasswordFile != "" {
		scIn.ReplicaPasswordFile = sc.ReplicaPasswordFile
	}
	scIn.Srvtls.ImportNotNull(&sc.Srvtls)
}

func (stIn *serverTls) ImportNotNull(st *serverTls) {
	if st.LdapsPort != "" {
		stIn.LdapsPort = st.LdapsPort
	}
	stIn.LdapsTls.ImportNotNull(&st.LdapsTls)
}

func (tIn *tls) ImportNotNull(t *tls) {
	if t.Ca != "" {
		tIn.Ca = t.Ca
	}
	if t.CaFile != "" {
		tIn.CaFile = t.CaFile
	}
	if t.Crt != "" {
		tIn.Crt = t.Crt
	}
	if t.CrtFile != "" {
		tIn.CrtFile = t.CrtFile
	}
	if t.CrtKey != "" {
		tIn.CrtKey = t.CrtKey
	}
	if t.CrtKeyFile != "" {
		tIn.CrtKeyFile = t.CrtKeyFile
	}
}

func (dbIn *DatabaseConfig) ImportNotNull(db *DatabaseConfig) {
	if db.Base != "" {
		dbIn.Base = db.Base
	}
	dbIn.Replicatls.ImportNotNull(&db.Replicatls)
}

func (rtIn *replicaTls) ImportNotNull(rt *replicaTls) {
	if rt.ReplicaPasswordFile != "" {
		rtIn.ReplicaPasswordFile = rt.ReplicaPasswordFile
	}
	rtIn.LdapsTls.ImportNotNull(&rt.LdapsTls)
}

func (scIn *ServerConfig) GetAdminPassword() (string, error) {
	if scIn.AdminPassword != "" {
		return scIn.AdminPassword, nil
	}
	if scIn.AdminPasswordFile != "" {
		adminPass, err := os.ReadFile(scIn.AdminPasswordFile)
		if err != nil {
			return "", err
		}
		return string(adminPass), nil
	}
	return "", errors.New("admin password is required")
}

func (scIn *ServerConfig) GetReplicaPassword() (string, error) {
	if scIn.ReplicaPassword != "" {
		return scIn.ReplicaPassword, nil
	}
	if scIn.ReplicaPasswordFile != "" {
		ReplicaPass, err := os.ReadFile(scIn.ReplicaPasswordFile)
		if err != nil {
			return "", err
		}
		return string(ReplicaPass), nil
	}
	return "", errors.New("replica password is required")
}
