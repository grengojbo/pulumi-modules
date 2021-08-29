package util

import (
	"fmt"
	"os"

	// "github.com/grengojbo/k3ctl/pkg/types"
	// homedir "github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
	// "github.com/spf13/afero"
)

func GerConfigFileName(configFile string) (configFilePath string) {
	messageError := "Is NOT cluster config file:"
	if configFile == "sample" {
		return "./samples/infra-example.yaml"
	}
	// file, err := afero.ReadFile(v.fs, filename)
	// if err != nil {
	// 	return err
	// }
	// if _, err := afero.Exists(configFile); err != nil {
	// 	log.Errorf("")
	// }
	if _, err := os.Stat(configFile); err != nil {
		messageError = fmt.Sprintf("%s %s", messageError, configFile)
		configFileCurrentDir := fmt.Sprintf("./%s.yaml", configFile)
		if _, err := os.Stat(configFileCurrentDir); err != nil {
			messageError = fmt.Sprintf("%s, %s", messageError, configFileCurrentDir)
		// 	configFileHomeDir := fmt.Sprintf("~/%s/cluster.yaml", configFile)
		// 	if _, err := os.Stat(configFileHomeDir); err != nil {
		// 		messageError = fmt.Sprintf("%s, %s", messageError, configFileHomeDir)
		// 		configFileDefaultFile := fmt.Sprintf("~/%s/%s.yaml", types.DefaultConfigDirName, configFile)
		// 		if _, err := os.Stat(configFileDefaultFile); err != nil {
		// 			messageError = fmt.Sprintf("%s, %s", messageError, configFileDefaultFile)
		// 			configFileDefaultDir := fmt.Sprintf("~/%s/%s/cluster.yaml", types.DefaultConfigDirName, configFile)
		// 			if _, err := os.Stat(configFileDefaultDir); err != nil {
		// 				messageError = fmt.Sprintf("%s, %s", messageError, configFileDefaultDir)
		// 			} else {
		// 				return configFileDefaultDir
		// 			}
		// 		} else {
		// 			return configFileDefaultFile
		// 		}
		// 	} else {
		// 		return configFileHomeDir
		// 	}
		} else {
			return configFileCurrentDir
		}
		// log.Fatalf("%+v", err)
		log.Fatalln(messageError)
	}
	return configFile
}