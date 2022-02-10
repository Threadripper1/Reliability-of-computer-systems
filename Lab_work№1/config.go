package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppCfg struct {
	Selection    []float64 `yaml:"selection"`
	IntervalSize int       `yaml:"interval_size"`
	Gamma        float64   `yaml:"gamma"`
	Hours        []float64 `yaml:"hours"`
}

func (c *AppCfg) String() string {
	return fmt.Sprintf("{selection_size: %d, interval_size: %d, y: %.3f, timelist: %v}", len(c.Selection), c.IntervalSize, c.Gamma, c.Hours)
}

func NewAppConfig(fileName string) (AppCfg, error) {
	var appConfig AppCfg
	cfgFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return appConfig, err
	}
	err = yaml.Unmarshal(cfgFile, &appConfig)
	if err != nil {
		return appConfig, errors.New("cannot decode file")
	}

	return appConfig, nil
}
