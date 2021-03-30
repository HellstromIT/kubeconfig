package kcaws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

type EKSClusters struct {
	Client      *eks.Client
	Params      *eks.ListClustersInput
	Config      *aws.Config
	clusterList []string
	Clusters    []*eks.DescribeClusterOutput
}

func (e *EKSClusters) EKSListClusters() {
	paginator := eks.NewListClustersPaginator(e.Client, e.Params)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("error paginating")
		}

		e.clusterList = append(e.clusterList, output.Clusters...)
	}
}

func (e *EKSClusters) GetClusterInfo() {
	for _, c := range e.clusterList {
		input := &eks.DescribeClusterInput{
			Name: &c,
		}
		cluster, err := e.Client.DescribeCluster(context.TODO(), input)
		if err != nil {
			log.Fatalf("Error describing cluster")
		}

		e.Clusters = append(e.Clusters, cluster)
	}
}

func AWSConfig(p string, r string) *EKSClusters {
	clusters := EKSClusters{}
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(p),
		config.WithRegion(r),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	clusters.Config = &cfg

	clusters.Client = eks.NewFromConfig(cfg)

	clusters.Params = &eks.ListClustersInput{}

	return &clusters
}
