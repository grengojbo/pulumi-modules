package main

import (
	"context"
	// "encoding/json"

	// "log"
	// "net/http"

	// "github.com/gorilla/mux"
	// "github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	log "github.com/sirupsen/logrus"
	// "github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	// "github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	// "github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	// "github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	// "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// https://github.com/pulumi/pulumi-random/releases
var pluginVersionRandom = "v4.2.0"
// https://github.com/pulumi/pulumi-aws/releases
var pluginVersionAws = "v4.15.0"
// https://github.com/pulumi/pulumi-azure/releases
var pluginVersionAzure = "v4.13.0"

type EnabledPlugins struct {
	Aws bool
	Azure bool
	// Hetzner bool
	// BareMetal bool
	// VmWare bool
}

// убедитесь, что плагины запускаются один раз перед загрузкой сервера
// убеждаемся, что установлены правильные плагины Pulumi
func EnsurePlugins(plugins *EnabledPlugins) {
	ctx := context.Background()
	w, err := auto.NewLocalWorkspace(ctx)
	if err != nil {
		log.Fatalf("Failed to setup and run http server: %v\n", err)
		// os.Exit(1)
	}
	
	err = w.InstallPlugin(ctx, "random", pluginVersionRandom)
	if err != nil {
		log.Fatalf("Failed to install program plugins: %v\n", err)
	}
	
	// err = w.InstallPlugin(ctx, "azure", "v4.0.0")
	// if err != nil {
	// 	log.Fatalf("Failed to install program plugins: %v\n", err)
	// }

	// err = s.SetConfig(ctx, "azure:location", auto.ConfigValue{Value: "westus"})
	// if err != nil {
	// 	log.Fatalf("Failed to set config: %v\n", err)
	// }
	
	if plugins.Aws {		
		err = w.InstallPlugin(ctx, "aws", pluginVersionAws)
		if err != nil {
			log.Fatalf("Failed to install program plugins: %v\n", err)
			// log.Fatalln(err)
			// fmt.Printf("Failed to install program plugins: %v\n", err)
			// os.Exit(1)
		}
	}
}


func main() {
	log.Infoln("Start main...")
	plugins := EnabledPlugins{
		Aws: true,
	} 
	
	EnsurePlugins(&plugins)

	log.Infoln("Finish :)")
}