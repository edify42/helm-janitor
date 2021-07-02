package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	client "github.com/edify42/helm-janitor/internal/eks"
	"github.com/edify42/helm-janitor/pkg/utils"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {

	// logging setup
	log.SetFormatter(&log.JSONFormatter{})

	// TODO: configuration block which we should define better
	name := "development-lendi-eks-cluster"
	region := "ap-southeast-2"
	janitorLabel := "helm-janitor=true"
	tmpFileLocation := "/tmp"
	tmpFilePrefix := "k8s-ca-*"
	debuggingFlag := false

	// k8s connectivity
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		// handle error :(
		log.Panic("config management issue...")
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
	log.Debugf("There are %d nodes associated with cluster %s", len(nodes.Items), name)

	// first generate a temp file to write the cluster CA to
	file, err := ioutil.TempFile(tmpFileLocation, tmpFilePrefix)
	if err != nil {
		log.Fatal(err)
	}
	if !debuggingFlag {
		defer os.Remove(file.Name())
	}
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
	iCli.Selector = janitorLabel
	rel, err := iCli.Run()
	if err != nil {
		panic(err)
	}
	releaseList := NameList(rel)
	log.Debugf("Got a list of releases: %v", releaseList)

	// loop through releases, track errors
	errorCount := 0

	for _, release := range rel {
		expired, err := CheckReleaseExpired(*release)
		if err != nil {
			errorCount++
			log.Error(err)
		}
		if expired {
			log.Infof("deleting release %s in namespace %s", release.Name, release.Namespace)
			cli := action.NewUninstall(actionConfig)
		}
	}

	// Finally throw the panic.
	if errorCount > 0 {
		log.Fatalf("Encountered %d errors while cleaning up helm releases - investigation required.", errorCount)
	}
}

// CheckReleaseExpired will return true if the release should be deleted.
// Safely returns false for any errors that occur.
func CheckReleaseExpired(r release.Release) (bool, error) {
	log.Debugf("Processing release: %s in namespace: ", r.Name, r.Namespace)
	ttlKey := "helm-janitor/ttl"
	expiryKey := "helm-janitor/expiry"
	now := time.Now()
	currentTime := now.Unix()
	if val, ok := r.Labels[ttlKey]; ok {
		log.Debugf("found %s: %s", ttlKey, val)
		timeLeft, err := utils.ParseTime(val)
		if err != nil {
			log.Errorf("%s key value: %s not valid - using default 7 days instead", ttlKey, val)
			timeLeft.Seconds = 7 * 24 * 60 * 60 // 7 days in seconds.
		}
		var expirySeconds int64 = int64(timeLeft.Seconds) // TODO: think about making this env variable? some kind of autocleanup after X days.
		// work off the modifiedAt key
		if modifiedTime, ok := r.Labels["modifiedAt"]; ok {
			log.Debugf("release: %s was modifiedAt %s", r.Name, modifiedTime)
			if n, err := strconv.ParseInt(modifiedTime, 10, 64); err == nil {
				if currentTime-n-expirySeconds > 0 {
					return true, nil
				} else {
					return false, nil
				}
			} else {
				return false, fmt.Errorf("modifiedTime cannot me made to int64")
			}
		} else {
			return false, fmt.Errorf("no modifiedAt label to work off on release: %s", r.Name)
		}
		// work off
	} else if val, ok := r.Labels[expiryKey]; ok {
		log.Debugf("found %s: %s", expiryKey, val)
	} else {
		return false, fmt.Errorf("no %s or %s found on release", ttlKey, expiryKey)
		// silently skip only - don't panic
	}
	// TODO: remove last catch all...
	return true, nil
}

// NameList loops through the releases and returns a []string of the
// Release[*].Name values
func NameList(r []*release.Release) []string {
	var list []string
	for _, user := range r {
		list = append(list, user.Name)
	}
	return list
}
