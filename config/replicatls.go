package config

type ReplicaTls struct {
	ReplicaPasswordFile string `yaml:"replicaPasswordFile" json:"replicaPasswordFile"`
	LdapsTls            tls    `yaml:"ldapsTls" json:"ldapsTls"`
	ReplicaUrl          string `yaml:"replicaUrl" json:"replicaUrl"`
}

func (rtIn *ReplicaTls) ImportNotNull(rt *ReplicaTls) {
	if rt.ReplicaPasswordFile != "" {
		rtIn.ReplicaPasswordFile = rt.ReplicaPasswordFile
	}
	rtIn.LdapsTls.ImportNotNull(&rt.LdapsTls)
}
