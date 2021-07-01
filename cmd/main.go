package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	client "github.com/edify42/helm-janitor/internal/eks"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {

	// k8s connectivity
	name := "development-lendi-eks-cluster"
	region := "ap-southeast-2"
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		// handle error :(
		fmt.Print("config management issue...")
	}

	eksClient := eks.NewFromConfig(cfg)

	result, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{Name: &name})
	if err != nil {
		log.Fatalf("Error calling DescribeCluster: %v", err)
	}
	clientset, err := client.New(result.Cluster)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error getting EKS nodes: %v", err)
	}
	log.Printf("There are %d nodes associated with cluster %s", len(nodes.Items), name)

	// first generate a temp file to write the cluster CA to
	file, err := ioutil.TempFile("/tmp", "k8s-ca-*")
	if err != nil {
		log.Fatal(err)
	}
	// defer os.Remove(file.Name())
	decoded, _ := base64.StdEncoding.DecodeString(*result.Cluster.CertificateAuthority.Data)
	file.Write([]byte(decoded))

	// configure the helm api to re-use the EKS auth above.
	releaseNamespace := ""
	newToken, err := client.GenToken(result.Cluster.Name)
	if err != nil {
		panic("Shieeet")
	}
	actionConfig := new(action.Configuration)
	settings := cli.New()
	settings.KubeAPIServer = *result.Cluster.Endpoint
	settings.KubeToken = newToken.Token
	settings.KubeCaFile = file.Name()
	if err := actionConfig.Init(settings.RESTClientGetter(), releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v)
	}); err != nil {
		panic(err)
	}

	iCli := action.NewList(actionConfig)
	iCli.Selector = "modifiedAt=1625049478"
	// iCli.ReleaseName = releaseName
	rel, err := iCli.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("Got a list of releases: ", rel[0].Name)

}
