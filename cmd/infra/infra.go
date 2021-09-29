package infra

import (
	"os"

	"github.com/grengojbo/pulumi-modules/interfaces"
	pkgAws "github.com/grengojbo/pulumi-modules/pkg/aws"
	"github.com/grengojbo/pulumi-modules/pkg/util"

	// "github.com/grengojbo/pulumi-modules/automation"
	// "github.com/grengojbo/pulumi-modules/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var destroy string

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

			if len(destroy) > 0 {
				// https://github.com/pulumi/automation-api-examples/blob/fb4ffddd8dd5eade5a7d1454e23123ae432fa3ac/go/multi_stack_orchestration/main.go#L119
				log.Warningln("TODO: add destroy")
				// // for destroying, we must remove stack in reverse order
				// // this means retrieving any dependend outputs first
				// fmt.Println("getting bucketID for object stack")

				// // wire up our destroy to stream progress to stdout
				// stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)

				// outs, err := websiteStack.Outputs(ctx)
				// if err != nil {
				// 	fmt.Printf("failed to get website outputs: %v\n", err)
				// 	os.Exit(1)
				// }
					os.Exit(0)
			
			}


			if len(rootFlags.Config.Spec.Providers.Aws.Cidr) > 0 {
				// aws, err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, rootFlags.Config.Spec.Providers.Aws, &sgList)
				// aws, err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, &sgList)
				// err := pkgAws.CreateVPC(rootFlags.Config.Spec.ProjectName, rootFlags.Config.Spec.StackName, &rootFlags.Config.Spec.Providers.Aws, &sgList)
				err := pkgAws.CreateVPC(&rootFlags.Config.Spec, &sgList)
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
	cmd.Flags().StringVarP(&destroy, "destroy", "d", "", "Destroy stack")

	// done
	return cmd
}