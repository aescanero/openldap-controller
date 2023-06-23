package config

type ModuleConfig struct {
	Name  string   `yaml:"name" json:"name"`
	Index []string `yaml:"index" json:"index"`
}
