package automation

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	// "encoding/json"
	// "log"
	// "net/http"

	// "github.com/gorilla/mux"
	// "github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	pWorkspace "github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	log "github.com/sirupsen/logrus"

	// "github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	// "github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	// "github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	// "github.com/pulumi/pulumi/sdk/v3/go/common/workspace"

	"github.com/grengojbo/pulumi-modules/interfaces"
)

// DestroyStack - удаление стека
func DestroyStack(ctx context.Context, appSetting *interfaces.AppArgs, verbose bool) error {
	projectName := appSetting.ProjectName
	stackName := appSetting.StackName

	project := pWorkspace.Project{
		Name:    tokens.PackageName(projectName),
		Runtime: pWorkspace.NewProjectRuntimeInfo("go", nil),
	}
	nilProgram := auto.Program(func(pCtx *pulumi.Context) error { return nil })
	nWorkspace, err := auto.NewLocalWorkspace(ctx, nilProgram, auto.Project(project))
	if err != nil {
		return fmt.Errorf("error creating local workspace: %v", err)
	}

	pulumiStack, err := auto.SelectStack(ctx, stackName, nWorkspace)

	if err != nil {
		return fmt.Errorf("error getting stack: %v", err)
	}

	if len(appSetting.Providers.Aws.Cidr) > 0 {
		log.Debugf("AWS Region: %s", appSetting.Providers.Aws.Region)
		err := pulumiStack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: appSetting.Providers.Aws.Region})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
		}
		// skip the metadata check
		err = pulumiStack.SetConfig(ctx, "aws:skipMetadataApiCheck", auto.ConfigValue{Value: "false"})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
			// return err
		}
	}
	if len(appSetting.Providers.Azure.Cidr) > 0 {
		log.Warningf("Azure Region: %s", appSetting.Providers.Azure.Region)
		err = pulumiStack.SetConfig(ctx, "azure:location", auto.ConfigValue{Value: "westus"})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
		}
	}
	if len(appSetting.Providers.Hetzner.Cidr) > 0 {
		log.Warningf("Hetzner Region: %s", appSetting.Providers.Hetzner.Region)
	}

	var streamer optdestroy.Option
	if verbose {
		streamer = optdestroy.ProgressStreams(os.Stdout)
	} else {
		streamer = optdestroy.ProgressStreams(ioutil.Discard)
	}
	_, err = pulumiStack.Destroy(ctx, streamer)

	if err != nil {
		return fmt.Errorf("error deleting stack resources: %v", err)
	}

	// destroy the stack so it's no longer listed
	// Then we delete the stack from earlier so we don't include it in our list
	nWorkspace.RemoveStack(ctx, projectName)

	return nil
}

// CreateOrSelectStack
func CreateOrSelectStack(ctx context.Context, appSetting *interfaces.AppArgs) auto.Stack {
	// ctx := context.Background()

	// projectName := "ploy"
	projectName := appSetting.ProjectName

	// Create a stack in our backend
	// We place all apps we deploy in the same project, so we can list them later
	// Each app is a stack, so we can do this multiple times
	// stackName := auto.FullyQualifiedStackName(org, projectName, name)
	stackName := appSetting.StackName

	pulumiStack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, nil)
	if err != nil {
		log.Fatalf("failed to create or select stack: %v", err)
	}

	log.Infof("Created/Selected stack %q\n", appSetting.StackName)

	// set stack configuration specifying the AWS region to deploy
	if len(appSetting.Providers.Aws.Cidr) > 0 {
		log.Debugf("AWS Region: %s", appSetting.Providers.Aws.Region)
		err := pulumiStack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: appSetting.Providers.Aws.Region})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
		}
		// skip the metadata check
		err = pulumiStack.SetConfig(ctx, "aws:skipMetadataApiCheck", auto.ConfigValue{Value: "false"})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
			// return err
		}
	}
	if len(appSetting.Providers.Azure.Cidr) > 0 {
		log.Warningf("Azure Region: %s", appSetting.Providers.Azure.Region)
		err = pulumiStack.SetConfig(ctx, "azure:location", auto.ConfigValue{Value: "westus"})
		if err != nil {
			log.Fatalf("Failed to set config: %v\n", err)
		}
	}
	if len(appSetting.Providers.Hetzner.Cidr) > 0 {
		log.Warningf("Hetzner Region: %s", appSetting.Providers.Hetzner.Region)
	}

	log.Infoln("Successfully set config")

	// log.Infoln("Starting refresh")
	// _, err = pulumiStack.Refresh(ctx)
	// if err != nil {
	// 	log.Fatalf("Failed to refresh stack: %v\n", err)
	// 	// os.Exit(1)
	// }
	// log.Infoln("Refresh succeeded!")

	return pulumiStack
}

// CreateOrSelectStackApp this function gets our stack ready for update/destroy by prepping the workspace, init/selecting the stack
// and doing a refresh to make sure state and cloud resources are in sync
// func CreateOrSelectStack(ctx context.Context, appSetting *interfaces.AppArgs, deployFunc pulumi.RunFunc) auto.Stack {
func CreateOrSelectStackApp(ctx context.Context, appSetting *interfaces.AppArgs, deployFunc pulumi.RunFunc) auto.Stack {
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
