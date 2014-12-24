package wally

import "gopkg.in/yaml.v2"

type Config struct {
	Database Db     `yaml:"database"`
	Tables   Tables `yaml:"tables"`
}

type Db struct {
	Host string `yaml:"host"`
	Name string `yaml:"name`
}

type Tables struct {
	DocumentTable string `yaml:"document_table"`
	IndexTable    string `yaml:"index_table"`
}

func LoadConfig(file []byte) (*Config, error) {
	c := new(Config)

	if err := yaml.Unmarshal(file, &c); err != nil {
		return nil, err
	}

	return c, nil
}
