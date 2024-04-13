package setup

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Entry struct {
	Path   string
	Config *SetupResult
}

func (r *Entry) Save() error {
	buf, err := yaml.Marshal(r.Config)
	if err != nil {
		return err
	}

	err = os.WriteFile(r.Path, buf, 0755)
	return err
}

func (r *Entry) Load() (error, *SetupResult) {

	buf, err := os.ReadFile(r.Path)
	if err != nil {
		return err, nil
	}

	c := &SetupResult{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return err, nil
	}

	r.Config = c
	return nil, c
}
