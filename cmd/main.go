package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/grengojbo/pulumi-modules/automation"
	"github.com/grengojbo/pulumi-modules/cmd/infra"
	"github.com/grengojbo/pulumi-modules/config"
	"github.com/grengojbo/pulumi-modules/interfaces"
	cliutil "github.com/grengojbo/pulumi-modules/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version is the string that contains version
var Version string
var stack string
var configFile string
var Flags = interfaces.RootFlags{}

var rootCmd = &cobra.Command{
	Use:   "mgr",
	Short: "mgr management Cloud Infrastructure command line",
	Long:  `Pupiter is an interpreter for Pulumi.`,
	Args:  cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Infof("Inside rootCmd PersistentPreRun with args: %v\n", args)
		log.Infof("stack: %s", stack)

  	config.InitConfig(configFile)
		cfg, err := config.FromViperConfig(stack, config.CfgViper)
		if err != nil {
			log.Fatalln(err)
		}
		log.Debugf("========== Simple Config ==========\n%+v\n==========================\n", cfg)
		// log.Debugf("kind: %s", cfg.Kind)
		Flags.Config = &cfg

		automation.EnsurePlugins(&Flags.Config.Spec.Plugins)
		// automation.EnsurePlugins(&cfg.Spec.Plugins)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start...")
		if Flags.Version {
			printVersion()
		} else {
			if err := cmd.Usage(); err != nil {
				log.Fatalln(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if len(os.Args) > 1 {
		parts := os.Args[1:]
		// Check if it's a built-in command, else try to execute it as a plugin
		if _, _, err := rootCmd.Find(parts); err != nil {
			pluginFound, err := cliutil.HandlePlugin(context.Background(), parts)
			if err != nil {
				log.Errorf("Failed to execute plugin '%+v'", parts)
				log.Fatalln(err)
			} else if pluginFound {
				os.Exit(0)
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVar(&Flags.DebugLogging, "verbose", false, "Enable verbose output (debug logging)")
	rootCmd.PersistentFlags().BoolVar(&Flags.TraceLogging, "trace", false, "Enable super verbose output (trace logging)")
	rootCmd.PersistentFlags().BoolVar(&Flags.TimestampedLogging, "timestamps", false, "Enable Log timestamps")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Show run command and skip execute")
	_ = viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))

	// add local flags
	rootCmd.Flags().BoolVar(&Flags.Version, "version", false, "Show mgr version")

	rootCmd.PersistentFlags().StringVar(&stack, "stack", "dev", "Current Pulumi stack")
	
	/***************
	 * Config File *
	 ***************/

	 rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path of a config file to use")
	 if err := cobra.MarkFlagFilename(rootCmd.PersistentFlags(), "config", "yaml", "yml"); err != nil {
		 log.Fatalln("Failed to mark flag 'config' as filename flag")
	 }
	// // add subcommands
	rootCmd.AddCommand(NewCmdCompletion())
	rootCmd.AddCommand(infra.NewCmdInfra(&Flags))
	// // rootCmd.AddCommand(kubeconfig.NewCmdKubeconfig())
	// // rootCmd.AddCommand(node.NewCmdNode())
	// // rootCmd.AddCommand(image.NewCmdImage())
	// // rootCmd.AddCommand(cfg.NewCmdConfig())
	// // rootCmd.AddCommand(registry.NewCmdRegistry())

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show mgr version",
		Long:  "Show mgr version",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	})

	// Init
	cobra.OnInitialize(initLogging)
	// cobra.OnInitialize(initLogging, initRuntime)
}

// initLogging initializes the logger
func initLogging() {
	if Flags.TraceLogging {
		log.SetLevel(log.TraceLevel)
	} else if Flags.DebugLogging {
		log.SetLevel(log.DebugLevel)
	} else {
		switch logLevel := strings.ToUpper(os.Getenv("LOG_LEVEL")); logLevel {
		case "TRACE":
			log.SetLevel(log.TraceLevel)
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
		case "WARN":
			log.SetLevel(log.WarnLevel)
		case "ERROR":
			log.SetLevel(log.ErrorLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}
	}
	log.SetOutput(ioutil.Discard)
	log.AddHook(&writer.Hook{
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
			log.TraceLevel,
		},
	})

	formatter := &log.TextFormatter{
		ForceColors: true,
	}

	if Flags.TimestampedLogging || os.Getenv("LOG_TIMESTAMPS") != "" {
		formatter.FullTimestamp = true
	}

	log.SetFormatter(formatter)

}

// GetVersion returns the version for cli, it gets it from "git describe --tags" or returns "dev" when doing simple go build
func GetVersion() string {
	if len(Version) == 0 {
		return "v1-dev"
	}
	return Version
}

func printVersion() {
	fmt.Printf("mgr version %s\n", GetVersion())
}

func generateFishCompletion(writer io.Writer) error {
	return rootCmd.GenFishCompletion(writer, true)
}

// Completion
var completionFunctions = map[string]func(io.Writer) error{
	"bash": rootCmd.GenBashCompletion,
	"zsh": func(writer io.Writer) error {
		if err := rootCmd.GenZshCompletion(writer); err != nil {
			return err
		}

		fmt.Fprintf(writer, "\n# source completion file\ncompdef _k3d k3d\n")

		return nil
	},
	"psh":        rootCmd.GenPowerShellCompletion,
	"powershell": rootCmd.GenPowerShellCompletion,
	"fish":       generateFishCompletion,
}

// NewCmdCompletion creates a new completion command
func NewCmdCompletion() *cobra.Command {
	// create new cobra command
	cmd := &cobra.Command{
		Use:   "completion SHELL",
		Short: "Generate completion scripts for [bash, zsh, fish, powershell | psh]",
		Long:  `Generate completion scripts for [bash, zsh, fish, powershell | psh]`,
		Args:  cobra.ExactArgs(1), // TODO: NewCmdCompletion: add support for 0 args = auto detection
		Run: func(cmd *cobra.Command, args []string) {
			if completionFunc, ok := completionFunctions[args[0]]; ok {
				if err := completionFunc(os.Stdout); err != nil {
					log.Fatalf("Failed to generate completion script for shell '%s'", args[0])
				}
				return
			}
			log.Fatalf("Shell '%s' not supported for completion", args[0])
		},
	}
	return cmd
}

func main() {
	Execute()
	// if err := Execute(); err != nil {
	// 	log.Fatalln(err)
	// }
}