package service

import (
	_ "embed"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
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
	Auth                  []GrypeRegistryAuthConfig `yaml:"auth"`
}

type GrypeRegistryAuthConfig struct {
	Authority string `yaml:"authority"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type GrypeLogConfig struct {
	File       string `yaml:"file"`
	Level      string `yaml:"level"`
	Structured bool   `yaml:"structured"`
}

func SaveGrypeConfigFile(conf *core.Config) error {
	var grypeConfig GrypeConfig
	err := yaml.Unmarshal([]byte(GrypeConfigTemplate), &grypeConfig)
	if err != nil {
		return err
	}

	auths, err := db.ListAllRegistryAuthSettings(conf)
	if err != nil {
		return err
	}

	for _, auth := range auths {
		grypeConfig.Registry.Auth = append(grypeConfig.Registry.Auth, GrypeRegistryAuthConfig{
			Authority: auth.Authority,
			Username:  auth.Username,
			Password:  auth.Password,
		})
	}

	//Configuration search paths:
	//  .grype.yaml
	//  .grype/config.yaml
	//  ~/.grype.yaml
	//  <XDG_CONFIG_HOME>/grype/config.yaml
	bytes, err := yaml.Marshal(grypeConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(".grype.yaml", bytes, 0777)
	if err != nil {
		return err
	}
	return nil
}
