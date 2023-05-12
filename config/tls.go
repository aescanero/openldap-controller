package config

type tls struct {
	Ca         string
	Crt        string
	CrtKey     string
	CaFile     string
	CrtFile    string
	CrtKeyFile string
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
