package config

import (
	"fmt"
	"strings"

	"github.com/grengojbo/pulumi-modules/interfaces"
	"github.com/grengojbo/pulumi-modules/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	// "gopkg.in/yaml.v2"
)

// Frankfurt
var DefaultAwsRegion = "eu-central-1"
// germany west central (Frankfurt)
var DefaultAzureRegion = "germanywestcentral"
var DefaultHetznerRegion = "nbg1"
// var DefaultHetznerRegion = "eu-central"

var configFile string
var CfgViper = viper.New()
// var ppViper = viper.New()
var DryRun bool

func InitConfig(cnfFile string) {
	DryRun = viper.GetBool("dry-run")
	// // Viper for pre-processed config options
	// ppViper.SetEnvPrefix("K3S")

	// // viper for the general config (file, env and non pre-processed flags)
	CfgViper.SetEnvPrefix("DEVOPS")
	CfgViper.AutomaticEnv()

	CfgViper.SetConfigType("yaml")

	configFile = util.GerConfigFileName(cnfFile)
	CfgViper.SetConfigFile(configFile)

	// try to read config into memory (viper map structure)
	if err := CfgViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file %s not found: %+v", configFile, err)
		}
		// config file found but some other error happened
		log.Fatalf("Failed to read config file %s: %+v", configFile, err)
	}

	log.Infof("Using config file %s", CfgViper.ConfigFileUsed())

	// if log.GetLevel() >= log.DebugLevel {
	// 	c, _ := yaml.Marshal(cfgViper.AllSettings())
	// 	log.Debugf("Configuration:\n%s", c)

	// 	c, _ = yaml.Marshal(ppViper.AllSettings())
	// 	log.Debugf("Additional CLI Configuration:\n%s", c)
	// }
	// log.Debugln("Config load succeeded!")
	log.Infof("Config load succeeded!")
}

func FromViperConfig(stack string, config *viper.Viper) (interfaces.App, error) {

	var cfg interfaces.App

	// determine config kind
	// if config.GetString("kind") != "" && strings.ToLower(config.GetString("kind")) != "cluster" {
	// 	return cfg, fmt.Errorf("Wrong `kind` '%s' != 'Cluster' in config file", config.GetString("kind"))
	// }
	if config.GetString("kind") != "" && strings.ToLower(config.GetString("kind")) != "infrastructure" {
		return cfg, fmt.Errorf("Wrong `kind` '%s' != 'Infrastructure' in config file", config.GetString("kind"))
	}


	if err := config.Unmarshal(&cfg); err != nil {
		log.Errorln("Failed to unmarshal File config")
		return cfg, err
	}

	cfg.Kind= config.GetString("kind")
	cfg.APIVersion = config.GetString("APIVersion")

	cfg.Spec.StackName = stack
	
	if len(cfg.Spec.Providers.Aws.Cidr) > 0 || len(cfg.Spec.Providers.Aws.VpcId) > 0 {
		cfg.Spec.Plugins.Aws = true
		if len(cfg.Spec.Providers.Aws.Region) == 0 {
			cfg.Spec.Providers.Aws.Region = DefaultAwsRegion
		}
	}
	if len(cfg.Spec.Providers.Azure.Cidr) > 0  || len(cfg.Spec.Providers.Azure.VpcId) >0 {
		cfg.Spec.Plugins.Azure = true
		if len(cfg.Spec.Providers.Azure.Region) == 0 {
			cfg.Spec.Providers.Azure.Region = DefaultAzureRegion
		}
	}
	if len(cfg.Spec.Providers.Hetzner.Cidr) > 0  || len(cfg.Spec.Providers.Hetzner.VpcId) > 0 {
		cfg.Spec.Plugins.Hetzner = true
		if len(cfg.Spec.Providers.Hetzner.Region) == 0 {
			cfg.Spec.Providers.Hetzner.Region = DefaultHetznerRegion
		}
	}

	return cfg, nil
}