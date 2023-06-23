package config

type tls struct {
	Ca         string `yaml:"ca" json:"ca,omitempty"`
	Crt        string `yaml:"crt" json:"crt,omitempty"`
	CrtKey     string `yaml:"crt_key" json:"crt_key,omitempty"`
	CaFile     string `yaml:"ca_file" json:"ca_file"`
	CrtFile    string `yaml:"crt_file" json:"crt_file"`
	CrtKeyFile string `yaml:"crt_key_file" json:"crt_key_file"`
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
