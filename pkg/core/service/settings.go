package service

import (
	_ "embed"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//go:embed templates/grype-config.yaml
var GrypeConfigTemplate string

type GrypeConfig struct {
	CheckForAppUpdate bool                `yaml:"check-for-app-update"`
	FailOnSeverity    string              `yaml:"fail-on-severity"`
	Output            string              `yaml:"output"`
	Scope             string              `yaml:"scope"`
	Quiet             bool                `yaml:"quiet"`
	Db                GrypeDbConfig       `yaml:"db"`
	Registry          GrypeRegistryConfig `yaml:"registry"`
	Log               GrypeLogConfig      `yaml:"log"`
}

type GrypeDbConfig struct {
	AutoUpdate bool   `yaml:"auto-update"`
	CacheDir   string `yaml:"cache-dir"`
	UpdateUrl  string `yaml:"update-url"`
}

type GrypeRegistryConfig struct {
	InsecureSkipTlsVerify bool                      `yaml:"insecure-skip-tls-verify"`
	InsecureUseHttp       bool                      `yaml:"insecure-use-http"`
	auth                  []GrypeRegistryAuthConfig `yaml:"auth"`
}

type GrypeRegistryAuthConfig struct {
	authority string `yaml:"authority"`
	username  string `yaml:"username"`
	password  string `yaml:"password"`
}

type GrypeLogConfig struct {
	file       string `yaml:"file"`
	level      string `yaml:"level"`
	structured bool   `yaml:"structured"`
}

func SaveGrypeConfigFile() error {
	var conf GrypeConfig
	err := yaml.Unmarshal([]byte(GrypeConfigTemplate), &conf)
	if err != nil {
		return err
	}

	//Configuration search paths:
	//  .grype.yaml
	//  .grype/config.yaml
	//  ~/.grype.yaml
	//  <XDG_CONFIG_HOME>/grype/config.yaml
	bytes, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(".grype.yaml", bytes, 0777)
	if err != nil {
		return err
	}
	return nil
}
