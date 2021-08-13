package main

import (

	// "encoding/json"

	// "log"
	// "net/http"

	// "github.com/gorilla/mux"
	// "github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"

	log "github.com/sirupsen/logrus"

	// "github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	// "github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	// "github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	// "github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/grengojbo/pulumi-modules/automation"
	"github.com/grengojbo/pulumi-modules/config"
)

func main() {
	log.SetLevel(log.DebugLevel)

	log.Infoln("Start main...")
	
	// app := interfaces.AppArgs{
	// 	ProjectName: "myops",
	// 	StackName: "dev",
	// }

	// app.Plugins.Aws = true
	// app.Plugins.Azure = true


	stackName :="dev"
	cnfFile := "sample"
  config.InitConfig(cnfFile)
	
	cfg, err := config.FromViperConfig(stackName,	config.CfgViper)
	if err != nil {
		log.Fatalln(err)
	}
	log.Debugf("========== Simple Config ==========\n%+v\n==========================\n", cfg)
	log.Debugf("kind: %s", cfg.Kind)

	automation.EnsurePlugins(&cfg.Spec.Plugins)
	
	log.Infoln("Finish :)")
}