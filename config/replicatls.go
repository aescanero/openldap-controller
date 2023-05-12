package config

type ReplicaTls struct {
	ReplicaPasswordFile string
	LdapsTls            tls
	ReplicaUrl          string
}

func (rtIn *ReplicaTls) ImportNotNull(rt *ReplicaTls) {
	if rt.ReplicaPasswordFile != "" {
		rtIn.ReplicaPasswordFile = rt.ReplicaPasswordFile
	}
	rtIn.LdapsTls.ImportNotNull(&rt.LdapsTls)
}
