package interfaces

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Tags {
//   [name: string]: Input<string>;
// }

// RootFlags describes a struct that holds flags that can be set on root level of the command
type RootFlags struct {
	DebugLogging       bool
	TraceLogging       bool
	TimestampedLogging bool
	Version            bool
	Config    				 *App
}

type EnabledPlugins struct {
  Kubernetes bool
  Docker bool
	Aws bool
	Azure bool
	Hetzner bool
	// BareMetal bool
	// VmWare bool
}

type App struct {
  metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec  AppArgs   `json:"spec,omitempty"`
  Status AppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppStatus defines the observed state of Cluster
type AppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// type Networking struct {
// }

type AppArgs struct {
	ProjectName string `mapstructure:"projectName" yaml:"projectName" json:"projectName,omitempty"`
	StackName string `mapstructure:"stackName" yaml:"stackName" json:"stackName,omitempty"`
  // Cidr CidrInterface `mapstructure:"cidr" yaml:"cidr" json:"cidr,omitempty"`
  Providers VpcOutputInterface `mapstructure:"providers" yaml:"providers" json:"providers,omitempty"`
  SecurityGroup SecurityGroupArgs `mapstructure:"sg" yaml:"sg" json:"sg,omitempty"`
  Plugins EnabledPlugins
}

type VpcAwsOutputInterface struct {
  // enabled: boolean;
  VpcId string `mapstructure:"vpcId" yaml:"vpcId" json:"vpcId,omitempty"`
  SubnetsIds []string
  SecureGroupId string
  Cidr string
  Region string `mapstructure:"region" yaml:"region" json:"region,omitempty"`
}

type VpcAzureOutputInterface struct {
  // enabled: boolean;
  VpcId string `mapstructure:"vpcId" yaml:"vpcId" json:"vpcId,omitempty"`
  ResourceGroupId string
  Cidr string
  Region string `mapstructure:"region" yaml:"region" json:"region,omitempty"`
}

type VpcHetznerOutputInterface struct {
  // enabled: boolean;
  VpcId string `mapstructure:"vpcId" yaml:"vpcId" json:"vpcId,omitempty"`
  Cidr string
  Region string `mapstructure:"region" yaml:"region" json:"region,omitempty"`
}

type VpcGcpOutputInterface struct {
  // enabled: boolean;
  VpcId string `mapstructure:"vpcId" yaml:"vpcId" json:"vpcId,omitempty"`
  Cidr string
  Region string `mapstructure:"region" yaml:"region" json:"region,omitempty"`
}

type VpcBareMetalInterface struct {
  // enabled: boolean;
  VpcId string `mapstructure:"vpcId" yaml:"vpcId" json:"vpcId,omitempty"`
  Cidr string
  Region string `mapstructure:"region" yaml:"region" json:"region,omitempty"`
}

type VpcOutputInterface struct {
  Aws VpcAwsOutputInterface `mapstructure:"aws" yaml:"aws" json:"aws,omitempty"`
  Hetzner VpcHetznerOutputInterface `mapstructure:"hetzner" yaml:"hetzner" json:"hetzner,omitempty"`
  Azure VpcAzureOutputInterface `mapstructure:"azure" yaml:"azure" json:"azure,omitempty"`
  BareMetal VpcBareMetalInterface `mapstructure:"vm" yaml:"vm" json:"vm,omitempty"`
  // Gcp
}

type CidrInterface struct {
  Aws string
  Hetzner string
  Azure string
  Gcp string
  VmWare string
  BareMetal string
}

type PortSecurityGroupArgs struct {
  Name string
  Port int64
  Allows []string
}

type SecurityGroupArgs struct {
  Name string
  Http bool
  Https bool
  SshAllows []string `mapstructure:"ssh" yaml:"ssh" json:"ssh,omitempty"`
  PostgreSqlAllows []string `mapstructure:"postgresql" yaml:"postgresql" json:"postgresql,omitempty"`
  MySqlAllows []string `mapstructure:"mysql" yaml:"mysql" json:"mysql,omitempty"`
  MongoAllows []string `mapstructure:"mongo" yaml:"mongo" json:"mongo,omitempty"`
  MsSqlAllows []string `mapstructure:"mssql" yaml:"mssql" json:"mssql,omitempty"`
	KubeApi []string `mapstructure:"kubeapi" yaml:"kubeapi" json:"kubeapi,omitempty"`
	Vault []string `mapstructure:"vault" yaml:"vault" json:"vault,omitempty"`
  Custom []PortSecurityGroupArgs `mapstructure:"custom" yaml:"custom" json:"custom,omitempty"`
}
