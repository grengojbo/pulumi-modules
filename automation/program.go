package automation

import (
	// "encoding/base64"
	// "fmt"
	// "path/filepath"
	// "strings"
	// "time"

	// "github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ecr"
	// "github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	// appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	// corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	// metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"

	pkgAws "github.com/grengojbo/pulumi-modules/pkg/aws"

	"github.com/grengojbo/pulumi-modules/interfaces"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	log "github.com/sirupsen/logrus"
)

func ApplyVPC(appArgs *interfaces.AppArgs, sg *[]interfaces.PortSecurityGroupArgs) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {

		if len(appArgs.Providers.Aws.Cidr) > 0 {
			// // aws, err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, rootFlags.Config.Spec.Providers.Aws, &sgList)
			// // aws, err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, &sgList)
			// // err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, &rootFlags.Config.Spec.Providers.Aws, &sgList)
			// err := pkgAws.CreateVPC(&rootFlags.Config.Spec, &sgList)
			// if err != nil {
			// 	log.Errorln(err)
			// }

			_, err := pkgAws.NewVPC(ctx, appArgs.ProjectName, &appArgs.Providers.Aws, sg)
			if err != nil {
				return err
			}

			log.Debugf("========== aws ==========\n%+v\n==========================\n", appArgs.Providers.Aws)
		}

		if len(appArgs.Providers.Azure.Cidr) > 0 {
			log.Warningf("TODO Azure cidr: %s", appArgs.Providers.Azure.Cidr)
		}
		if len(appArgs.Providers.Hetzner.Cidr) > 0 {
			log.Warningf("TODO Hetzner cidr: %s", appArgs.Providers.Hetzner.Cidr)
		}

		// if appArgs.SecurityGroup {
		// 	log.Warningf("TODO SecurityGroup sg: %v", appArgs.SecurityGroup)
		// }

		// _, err := NewPloyDeployment(ctx, name, &PloyDeploymentArgs{
		// 	Directory: directory,
		// 	Nlb:       nlb,
		// })
		// if err != nil {
		// 	return err
		// }

		return nil
	}

}
