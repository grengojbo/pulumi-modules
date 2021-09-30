package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/grengojbo/pulumi-modules/automation"
	"github.com/grengojbo/pulumi-modules/interfaces"
	"github.com/grengojbo/pulumi-modules/pkg/util"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"

	// "github.com/grengojbo/pulumi-modules/automation"
	// "github.com/grengojbo/pulumi-modules/config"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var destroy string

func NewCmdInfra(rootFlags *interfaces.RootFlags) *cobra.Command {

	// create new cobra command
	command := &cobra.Command{
		Use:   "infra",
		Short: "Manage Infrastructure(s)",
		Long:  `Manage Infrastructure(s)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Если невыполняется то раскоментировать
			// if err := cmd.Help(); err != nil {
			// 	log.Errorln("Couldn't get help text")
			// 	log.Fatalln(err)
			// }
			log.Infoln("RUN sync infrastructure...")
			// args := []string{"sample"}
			// config.InitConfig(args)
			dryrun := viper.GetBool("dry-run")
			stackName := rootFlags.Config.Spec.StackName
			projectName := rootFlags.Config.Spec.ProjectName

			ctx := context.Background()

			log.Infof("ProjectName: %s", projectName)

			sgList := util.ListSecurityGroup(rootFlags.Config.Spec.SecurityGroup)

			if len(destroy) > 0 {
				// https://github.com/pulumi/automation-api-examples/blob/fb4ffddd8dd5eade5a7d1454e23123ae432fa3ac/go/multi_stack_orchestration/main.go#L119

				if destroy != stackName {
					log.Fatalf("IS NOT stack: %s\n", destroy)
				}

				label := fmt.Sprintf("This will delete the stack %s. Are you sure you wish to continue?", stackName)

				prompt := promptui.Prompt{
					Label:     label,
					IsConfirm: true,
				}

				result, err := prompt.Run()

				if err != nil {
					fmt.Printf("User cancelled, not deleting %v\n", err)
					os.Exit(0)
				}
				log.Debug("User confirmed, continuing: %s", result)
				log.Infof("Deleting stack: %s", stackName)

				err = automation.DestroyStack(ctx, &rootFlags.Config.Spec, rootFlags.DebugLogging)
				if err != nil {
					log.Fatalln(err)
				}
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

			pulumiStack := automation.CreateOrSelectStack(ctx, &rootFlags.Config.Spec)
			// Set up the workspace and install all the required plugins the user needs
			workspace := pulumiStack.Workspace()

			// err := pulumi.EnsurePlugins(workspace)
			// if err != nil {
			// 	return err
			// }

			// Now, we set the pulumi program that is going to run
			// workspace.SetProgram(pulumi.Deploy(name, directory, nlb))
			workspace.SetProgram(automation.ApplyVPC(&rootFlags.Config.Spec, &sgList))

			if dryrun {
				_, err := pulumiStack.Preview(ctx, optpreview.Message("------- Running dryrun -------"))
				if err != nil {
					log.Fatalf("error creating stack: %v", err)
				}
			} else {
				// Wire up our update to stream progress to stdout
				// We give the user the option to actually view the Pulumi output
				var streamer optup.Option
				if rootFlags.DebugLogging {
					streamer = optup.ProgressStreams(os.Stdout)
				} else {
					upChannel := make(chan events.EngineEvent)
					go collectEvents(upChannel)

					streamer = optup.EventStreams(upChannel)
				}
				// log.Infof("Creating ploy application: %s", name)
				_, err := pulumiStack.Up(ctx, streamer)

				if err != nil {
					return err
				}
			}

			// cfg, err := config.FromViperConfig(config.CfgViper)
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			// log.Debugf("========== Simple Config ==========\n%+v\n==========================\n", cfg)
			// log.Debugf("kind: %s", cfg.Kind)

			// automation.EnsurePlugins(&cfg.Spec.Plugins)

			return nil
		},
	}

	// add subcommands
	// command.AddCommand(NewCmdClusterCreate())
	// command.AddCommand(NewCmdClusterStart())
	// command.AddCommand(NewCmdClusterStop())
	// command.AddCommand(NewCmdClusterDelete())
	// command.AddCommand(NewCmdClusterList())

	// add flags
	command.Flags().StringVarP(&destroy, "destroy", "d", "", "Destroy stack")

	// done
	return command
}

func collectEvents(eventChannel <-chan events.EngineEvent) {

	for {

		var event events.EngineEvent
		var ok bool

		createLogger := log.WithFields(log.Fields{"event": "CREATING"})
		completeLogger := log.WithFields(log.Fields{"event": "COMPLETE"})

		event, ok = <-eventChannel
		if !ok {
			return
		}

		if event.ResourcePreEvent != nil {
			createLogger.WithFields(log.Fields{"resource": event.ResourcePreEvent.Metadata.Type}).Infof("event.ResourcePreEvent.Metadata.Type=%v\n", event.ResourcePreEvent.Metadata.Type)
			// switch event.ResourcePreEvent.Metadata.Type {
			// case "aws:ecr/repository:Repository":
			// 	createLogger.WithFields(log.Fields{"resource": event.ResourcePreEvent.Metadata.Type}).Info("Creating ECR repository")
			// case "kubernetes:core/v1:Namespace":
			// 	createLogger.WithFields(log.Fields{"resource": event.ResourcePreEvent.Metadata.Type}).Info("Creating Kubernetes Namespace")
			// case "kubernetes:core/v1:Service":
			// 	createLogger.WithFields(log.Fields{"resource": event.ResourcePreEvent.Metadata.Type}).Info("Creating Kubernetes Service")
			// case "kubernetes:apps/v1:Deployment":
			// 	createLogger.WithFields(log.Fields{"resource": event.ResourcePreEvent.Metadata.Type}).Info("Creating Kubernetes Deployment")
			// case "docker:image:Image":
			// 	createLogger.WithFields(log.Fields{"resource": event.ResourcePreEvent.Metadata.Type}).Info("Creating Docker Image")
			// }
		}

		if event.ResOutputsEvent != nil {
			completeLogger.WithFields(log.Fields{"name": event.ResOutputsEvent.Metadata.New.Outputs["repositoryUrl"], "resource": event.ResOutputsEvent.Metadata.Type}).Infof("event.ResOutputsEvent.Metadata.Type=%v\n", event.ResOutputsEvent.Metadata.Type)
			// switch event.ResOutputsEvent.Metadata.Type {
			// case "aws:ecr/repository:Repository":
			// 	completeLogger.WithFields(log.Fields{"name": event.ResOutputsEvent.Metadata.New.Outputs["repositoryUrl"], "resource": event.ResOutputsEvent.Metadata.Type}).Info("Created ECR repository")
			// case "kubernetes:core/v1:Namespace":
			// 	completeLogger.WithFields(log.Fields{"resource": event.ResOutputsEvent.Metadata.Type}).Info("Created Kubernetes Namespace")
			// case "kubernetes:core/v1:Service":
			// 	completeLogger.WithFields(log.Fields{"resource": event.ResOutputsEvent.Metadata.Type}).Info("Created Kubernetes Service")
			// case "kubernetes:apps/v1:Deployment":
			// 	completeLogger.WithFields(log.Fields{"resource": event.ResOutputsEvent.Metadata.Type}).Info("Created Kubernetes Deployment")
			// case "docker:image:Image":
			// 	completeLogger.WithFields(log.Fields{"name": event.ResOutputsEvent.Metadata.New.Outputs["baseImageName"], "resource": event.ResOutputsEvent.Metadata.Type}).Info("Created Docker Image")
			// }

		}
	}
}
