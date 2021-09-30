package automation

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

// EnsurePlugins убедитесь, что плагины запускаются один раз перед загрузкой сервера
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
