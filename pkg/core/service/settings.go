package service

import (
	_ "embed"
	"encoding/base64"
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

type ConfigRegistry struct {
	Hub  string               `yaml:"hub,omitempty"`
	Auth []ConfigRegistryAuth `yaml:"auth,omitempty"`
}

type ConfigRegistryAuth struct {
	Authority string `yaml:"authority,omitempty"`
	Username  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
}

func SaveGrypeConfigFile(registry string) error {
	var grypeConfig GrypeConfig
	err := yaml.Unmarshal([]byte(GrypeConfigTemplate), &grypeConfig)
	if err != nil {
		return err
	}

	if registry != "" {
		decode, err := base64.StdEncoding.DecodeString(registry)
		if err != nil {
			return err
		}

		var configRegistry ConfigRegistry
		err = yaml.Unmarshal(decode, &configRegistry)
		if err != nil {
			return err
		}

		for _, auth := range configRegistry.Auth {
			grypeConfig.Registry.Auth = append(grypeConfig.Registry.Auth, GrypeRegistryAuthConfig{
				Authority: auth.Authority,
				Username:  auth.Username,
				Password:  auth.Password,
			})
		}
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
