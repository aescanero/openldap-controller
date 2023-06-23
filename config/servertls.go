package config

type serverTls struct {
	LdapsPort string `yaml:"ldaps_port" json:"ldaps_port"`
	LdapsTls  tls    `yaml:"ldaps_tls" json:"ldaps_tls"`
}

func (stIn *serverTls) ImportNotNull(st *serverTls) {
	if st.LdapsPort != "" {
		stIn.LdapsPort = st.LdapsPort
	}
	stIn.LdapsTls.ImportNotNull(&st.LdapsTls)
}
