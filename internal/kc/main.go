package kc

import (
	"fmt"
	"log"

	"github.com/HellstromIT/kubeconfig/cmd/kubeconfig/internal/kcaws"
	"github.com/alecthomas/kong"
)

type context struct {
	version    string
	configFile string
	config     *kubeconf
}

type addCmd struct {
}

type awsCmd struct {
	awsCmdSub `cmd`
}

type awsCmdSub struct {
	Region  string `short:"r" default:"eu-west-1" help:"Region to add EKS clusters from. Default: eu-west-1"`
	Profile string `short:"p" default:"default" optional help:"AWS Profile to use when connecting to AWS. Default: default "`
	ARN     string `short:"a" default:"" optional help:"Role ARN assumed to get access to cluster. Default: '' "`
}

type listCmd struct {
}

type versionCmd struct {
}

var cli struct {
	Add     addCmd     `cmd help:"Add new configuration."`
	Aws     awsCmd     `cmd help:"Add all clusters from AWS account."`
	List    listCmd    `cmd help:"List current configuration."`
	Version versionCmd `cmd help:"Print version."`
}

func (a *awsCmd) Run(ctx *context) error {

	clusters := kcaws.AWSConfig(a.Profile, a.Region)

	clusters.EKSListClusters()
	clusters.GetClusterInfo()

	for _, cluster := range clusters.Clusters {
		ctx.config.addCluster(*cluster.Cluster.Name, *cluster.Cluster.CertificateAuthority.Data, *cluster.Cluster.Endpoint)
		ctx.config.addUserEKS(*cluster.Cluster.Name, a.ARN, a.Profile, a.Region)
		ctx.config.addContext(
			*cluster.Cluster.Name,
			*cluster.Cluster.Name,
			a.Profile+"-"+*cluster.Cluster.Name,
		)
	}
	ctx.config.writeConf(ctx.configFile)
	return nil
}

func (v *listCmd) Run(ctx *context) error {
	printConf(ctx.config)
	return nil
}

func (v *versionCmd) Run(ctx *context) error {
	fmt.Println(ctx.version)
	return nil
}

func Cli(v string) {
	configfile := getConfigFile("/.kube/config")
	conf, err := readConf(configfile)
	if err != nil {
		log.Fatal(err)
	}

	ctx := kong.Parse(&cli)
	err = ctx.Run(&context{version: v, configFile: configfile, config: conf})
	ctx.FatalIfErrorf(err)
}
