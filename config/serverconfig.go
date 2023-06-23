package config

import (
	"errors"
	"os"
	"strings"
)

type ServerConfig struct {
	AdminPassword       string    // deprecated
	AdminPasswordSHA    string    // deprecated
	AdminPasswordFile   string    `yaml:"admin_password_file" json:"admin_password_file"`
	ReplicaPassword     string    // deprecated
	ReplicaPasswordSHA  string    // deprecated
	ReplicaPasswordFile string    `yaml:"replica_password_file" json:"replica_password_file"`
	LdapPort            string    `yaml:"ldap_port" json:"ldap_port"`
	Srvtls              serverTls `yaml:"srvtls" json:"srvtls"`
	Debug               string    `yaml:"debug" json:"debug"`
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

func (scIn *ServerConfig) GetAdminPassword() (string, error) {
	if scIn.AdminPassword != "" {
		return scIn.AdminPassword, nil
	}
	if scIn.AdminPasswordFile != "" {
		adminPass, err := os.ReadFile(scIn.AdminPasswordFile)
		if err != nil {
			return "", err
		}
		return strings.TrimSuffix(string(adminPass), "\n"), nil
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
