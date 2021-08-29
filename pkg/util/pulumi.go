package util

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ToPulumiStringArray Conver string array to pulumi input string array
func ToPulumiStringArray(a []string) pulumi.StringArrayInput {
	var res []pulumi.StringInput
	for _, s := range a {
		res = append(res, pulumi.String(s))
	}
	return pulumi.StringArray(res)
}

// GetEnv searches for the requested key in the pulumi context and provides 
// either the value of the key or the fallback.
// https://www.retgits.com/2020/01/how-to-create-a-vpc-in-aws-using-pulumi-and-golang/
func GetEnv(ctx *pulumi.Context, key string, fallback string) string {
	if value, ok := ctx.GetConfig(key); ok {
		return value
	}
	return fallback
}

// // Prepare the tags that are used for each individual resource so they can be found
// 		// using the Resource Groups service in the AWS Console
// 		tags := make(map[string]interface{})
// 		tags["version"] = getEnv(ctx, "tags:version", "unknown")
// 		tags["author"] = getEnv(ctx, "tags:author", "unknown")
// 		tags["team"] = getEnv(ctx, "tags:team", "unknown")
// 		tags["feature"] = getEnv(ctx, "tags:feature", "unknown")
// 		tags["region"] = getEnv(ctx, "aws:region", "unknown")

// 		// Create a VPC for the EKS cluster
// 		cidrBlock := getEnv(ctx, "vpc:cidr-block", "unknown")