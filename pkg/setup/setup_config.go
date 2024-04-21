package setup

import (
	"os"

	"gopkg.in/yaml.v3"
)

type YAMLEntry struct {
	Path   string
	Config *SetupResult
}

// Save persiste the configuration to YAML file
func (r *YAMLEntry) Save() error {
	buf, err := yaml.Marshal(r.Config)
	if err != nil {
		return err
	}

	err = os.WriteFile(r.Path, buf, 0755)
	return err
}

// Get returns the in memory configuration
func (r *YAMLEntry) Get() *SetupResult {
	return r.Config
}

// Load retrive the configuration from YAML file
func (r *YAMLEntry) Load() (*SetupResult, error) {
	// If the file does exist we return a blan config
	_, err := os.Stat(r.Path)
	if err != nil {
		return &SetupResult{}, nil
	}

	buf, err := os.ReadFile(r.Path)
	if err != nil {
		return nil, err
	}

	c := &SetupResult{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}

	r.Config = c
	return c, nil
}
