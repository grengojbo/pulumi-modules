package aws

import (
	AWS "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

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

// GetRegions список регионов
// https://github.com/pulumi/examples/blob/master/aws-go-console-slack-notification/main.go
func GetRegions() ([]string, error) {
	svc := ec2.New(session.New())
	result, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{
		Filters: []*ec2.Filter{
			{
				Name:   AWS.String("opt-in-status"),
				Values: []*string{AWS.String("opt-in-not-required"), AWS.String("opted-in")},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	regions := make([]string, len(result.Regions))
	for i, region := range result.Regions {
		regions[i] = *region.RegionName
	}

	return regions, nil
}