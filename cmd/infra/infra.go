package infra

import (
	"github.com/grengojbo/pulumi-modules/interfaces"
	pkgAws "github.com/grengojbo/pulumi-modules/pkg/aws"
	"github.com/grengojbo/pulumi-modules/pkg/util"

	// "github.com/grengojbo/pulumi-modules/automation"
	// "github.com/grengojbo/pulumi-modules/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdInfra(rootFlags *interfaces.RootFlags) *cobra.Command {

	// create new cobra command
	cmd := &cobra.Command{
		Use:   "infra",
		Short: "Manage Infrastructure(s)",
		Long:  `Manage Infrastructure(s)`,
		Run: func(cmd *cobra.Command, args []string) {
			// Если невыполняется то раскоментировать
			// if err := cmd.Help(); err != nil {
			// 	log.Errorln("Couldn't get help text")
			// 	log.Fatalln(err)
			// }
			log.Infoln("RUN sync infrastructure...")
			// args := []string{"sample"}
  		// config.InitConfig(args)
	
			log.Infof("ProjectName: %s", rootFlags.Config.Spec.ProjectName)
			
			sgList := util.ListSecurityGroup(rootFlags.Config.Spec.SecurityGroup)

			if len(rootFlags.Config.Spec.Providers.Aws.Cidr) > 0 {
				// aws, err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, rootFlags.Config.Spec.Providers.Aws, &sgList)
				// aws, err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, &sgList)
				err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, &rootFlags.Config.Spec.Providers.Aws, &sgList)
				if err != nil {
					log.Errorln(err)
				}
				log.Debugf("========== aws ==========\n%+v\n==========================\n", rootFlags.Config.Spec.Providers.Aws)
			}

			if len(rootFlags.Config.Spec.Providers.Azure.Cidr) > 0 {
				log.Warningf("TODO Azure cidr: %s", rootFlags.Config.Spec.Providers.Azure.Cidr)
			}
			if len(rootFlags.Config.Spec.Providers.Hetzner.Cidr) > 0 {
				log.Warningf("TODO Hetzner cidr: %s", rootFlags.Config.Spec.Providers.Hetzner.Cidr)
			}

			// if rootFlags.Config.Spec.SecurityGroup {
				// log.Warningf("TODO SecurityGroup sg: %v", rootFlags.Config.Spec.SecurityGroup)
			// }

			// cfg, err := config.FromViperConfig(config.CfgViper)
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			// log.Debugf("========== Simple Config ==========\n%+v\n==========================\n", cfg)
			// log.Debugf("kind: %s", cfg.Kind)

			// automation.EnsurePlugins(&cfg.Spec.Plugins)
		},
	}

	// add subcommands
	// cmd.AddCommand(NewCmdClusterCreate())
	// cmd.AddCommand(NewCmdClusterStart())
	// cmd.AddCommand(NewCmdClusterStop())
	// cmd.AddCommand(NewCmdClusterDelete())
	// cmd.AddCommand(NewCmdClusterList())

	// add flags

	// done
	return cmd
}