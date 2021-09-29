package aws

import (
	"context"
	"os"

	AWS "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/grengojbo/pulumi-modules/automation"
	"github.com/grengojbo/pulumi-modules/interfaces"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	// "github.com/grengojbo/pulumi-modules/config"
	log "github.com/sirupsen/logrus"
	// "github.com/spf13/cobra"
)

// CreateVPC - Сщздаем сеть
// пример https://github.com/pulumi/automation-api-examples/blob/fb4ffddd8dd5eade5a7d1454e23123ae432fa3ac/go/multi_stack_orchestration/main.go#L119
func CreateVPC(appArgs *interfaces.AppArgs, sg *[]interfaces.PortSecurityGroupArgs) (err error) {
// func CreateVPC(projectName string, stackName string, providerArgs *interfaces.VpcAwsOutputInterface, sg *[]interfaces.PortSecurityGroupArgs) (err error) {
// func CreateVPC(projectName string, stackName string, providerArgs interfaces.VpcAwsOutputInterface, sg *[]interfaces.PortSecurityGroupArgs) (vpc *interfaces.VpcAwsOutputInterface, err error) {
	// log.Warningf("TODO AWS cidr: %s", providerArgs.Cidr)
	log.Infoln("Start AWS CreateVPC...")
	
	ctx := context.Background()

	log.Debugln("preparing aws vpc stack")
	// websiteStack := createOrSelectWebsiteStack(ctx, stackName)
	vpc := automation.CreateOrSelectStack(ctx, appArgs, vpcFunc)
	log.Debugln("aws vpc stack ready to deploy")
	// providerArgs.VpcId = "2134215423534534"

	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// run the update to deploy our s3 website
	vpcRes, err := vpc.Up(ctx, stdoutStreamer)
	if err != nil {
		log.Fatalf("Failed to update stack: %v\n\n", err)
		// os.Exit(1)
	}
	// get the bucketID output that object stack depends on
	// bucketID, ok := resRes.Outputs["bucketID"].Value.(string)
	// if !ok {
	// 	fmt.Println("failed to get bucketID output")
	// 	os.Exit(1)
	// }

	log.Infof("AWS VPC result: %v\n",vpcRes)
	log.Infoln("aws vpc stack update succeeded!")

	return nil
}

func vpcFunc(ctx *pulumi.Context) error {
	// // export the website URL
	// ctx.Export("websiteUrl", siteBucket.WebsiteEndpoint)
	// // also export the bucketID for Object stack to refer to
	// ctx.Export("bucketID", siteBucket.ID())
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