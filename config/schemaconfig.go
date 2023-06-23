package config

type SchemaConfig struct {
	Path  string   `yaml:"path" json:"path"`
	Index []string `yaml:"index" json:"index"`
}
