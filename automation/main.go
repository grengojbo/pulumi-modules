package automation

import (
	"context"
	// "encoding/json"
	// "log"
	// "net/http"

	// "github.com/gorilla/mux"
	// "github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	log "github.com/sirupsen/logrus"

	// "github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	// "github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	// "github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	// "github.com/pulumi/pulumi/sdk/v3/go/common/workspace"

	"github.com/grengojbo/pulumi-modules/interfaces"
)

// https://github.com/pulumi/pulumi-random/releases
var pluginVersionRandom = "v4.2.0"
// https://github.com/pulumi/pulumi-aws/releases
var pluginVersionAws = "v4.22.0"
// https://github.com/pulumi/pulumi-azure/releases
var pluginVersionAzure = "v4.20.1"
// https://github.com/pulumi/pulumi-hcloud/releases
var pluginVersionHetzner = "v1.5.0"
// https://github.com/pulumi/pulumi-kubernetes/releases
var pluginVersionKubernetes = "v3.7.2"
// https://github.com/pulumi/pulumi-docker/releases
var pluginVersionDocker = "v3.1.0"

// убедитесь, что плагины запускаются один раз перед загрузкой сервера
// убеждаемся, что установлены правильные плагины Pulumi
func EnsurePlugins(plugins *interfaces.EnabledPlugins) {
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
	log.Infof("Successfully installed random plugin")
	
	if plugins.Kubernetes {
		err = w.InstallPlugin(ctx, "kubernetes", pluginVersionKubernetes)
		if err != nil {
			log.Fatalf("error installing kubernetes plugin: %v", err)
		}
		log.Infof("Successfully installed kubernetes plugin")
	}
	
	if plugins.Docker {
		err = w.InstallPlugin(ctx, "docker", pluginVersionDocker)
		if err != nil {
			log.Fatalf("error installing docker plugin: %v", err)
		}
		log.Infof("Successfully installed docker plugin")
	}

	if plugins.Azure {		
		err = w.InstallPlugin(ctx, "azure", pluginVersionAzure)
		if err != nil {
			log.Fatalf("Failed to install program plugins: %v\n", err)
		}
		log.Infof("Successfully installed Azure plugin")
	}
	
	if plugins.Hetzner {		
		err = w.InstallPlugin(ctx, "hcloud", pluginVersionHetzner)
		if err != nil {
			log.Fatalf("Failed to install program plugins: %v\n", err)
		}
		log.Infof("Successfully installed Hetzner plugin")
	}
	
	if plugins.Aws {		
		err = w.InstallPlugin(ctx, "aws", pluginVersionAws)
		if err != nil {
			log.Fatalf("Failed to install program plugins: %v\n", err)
			// log.Fatalln(err)
			// fmt.Printf("Failed to install program plugins: %v\n", err)
			// os.Exit(1)
		}
		log.Infof("Successfully installed AWS plugin")
	}
}

// this function gets our stack ready for update/destroy by prepping the workspace, init/selecting the stack
// and doing a refresh to make sure state and cloud resources are in sync
// func CreateOrSelectStack(ctx context.Context, appSetting *interfaces.AppArgs, deployFunc pulumi.RunFunc) auto.Stack {
func CreateOrSelectStack(ctx context.Context, appSetting *interfaces.AppArgs, deployFunc pulumi.RunFunc) auto.Stack {
	// create or select a stack with an inline Pulumi program
	s, err := auto.UpsertStackInlineSource(ctx, appSetting.StackName, appSetting.ProjectName, deployFunc)
	if err != nil {
		log.Fatalf("Failed to create or select stack: %v\n", err)
		// os.Exit(1)
	}

	log.Infof("Created/Selected stack %q\n", appSetting.StackName)


	// set stack configuration specifying the AWS region to deploy
	if len(appSetting.Providers.Aws.Cidr) > 0 {
		log.Debugf("AWS Region: %s", appSetting.Providers.Aws.Region)
		err := s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: appSetting.Providers.Aws.Region})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
		}
		// skip the metadata check
		err = s.SetConfig(ctx, "aws:skipMetadataApiCheck", auto.ConfigValue{Value: "false"})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
			// return err
		}
	}
	if len(appSetting.Providers.Azure.Cidr) > 0 {
		log.Warningf("Azure Region: %s", appSetting.Providers.Azure.Region)
		err = s.SetConfig(ctx, "azure:location", auto.ConfigValue{Value: "westus"})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
		}
	}
	if len(appSetting.Providers.Hetzner.Cidr) > 0 {
		log.Warningf("Hetzner Region: %s", appSetting.Providers.Hetzner.Region)
	}
	
	log.Infoln("Successfully set config")

	log.Infoln("Starting refresh")

	_, err = s.Refresh(ctx)
	if err != nil {
		log.Fatalf("Failed to refresh stack: %v\n", err)
		// os.Exit(1)
	}

	log.Infoln("Refresh succeeded!")

	return s
}