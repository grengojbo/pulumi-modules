package aws

import (
	"github.com/grengojbo/pulumi-modules/interfaces"
	// "github.com/grengojbo/pulumi-modules/automation"
	// "github.com/grengojbo/pulumi-modules/config"
	log "github.com/sirupsen/logrus"
	// "github.com/spf13/cobra"
)

func CreateVPC(projectName string, stackName string, providerArgs *interfaces.VpcAwsOutputInterface, sg *[]interfaces.PortSecurityGroupArgs) (err error) {
// func CreateVPC(projectName string, stackName string, providerArgs interfaces.VpcAwsOutputInterface, sg *[]interfaces.PortSecurityGroupArgs) (vpc *interfaces.VpcAwsOutputInterface, err error) {
	// log.Warningf("TODO AWS cidr: %s", providerArgs.Cidr)
	log.Infoln("Start CreateVPC...")
	
	// providerArgs.VpcId = "2134215423534534"
	return nil
}