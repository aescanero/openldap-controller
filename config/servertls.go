package config

type serverTls struct {
	LdapsPort string
	LdapsTls  tls
}

func (stIn *serverTls) ImportNotNull(st *serverTls) {
	if st.LdapsPort != "" {
		stIn.LdapsPort = st.LdapsPort
	}
	stIn.LdapsTls.ImportNotNull(&st.LdapsTls)
}
