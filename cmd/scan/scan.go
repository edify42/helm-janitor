package scan

import (
	"fmt"
	"os"
	"time"

	janitorconfig "github.com/edify42/helm-janitor/internal/config"
	"github.com/edify42/helm-janitor/pkg/utils"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// main struct for this file.

func RunV2(sr InputRun) {
	mycfg := sr.Config()
	cfg := sr.Makeawscfg()
	cluster := sr.Getekscluster(cfg)
	if !mycfg.DebugFlag {
		log.Info("should clean " + cluster.CAFile)
		defer os.Remove(cluster.CAFile)
	} else {
		log.Info("no clean")
	}
	defer os.Remove(cluster.CAFile)
	releaseNamespace := ""

	// direct copy...
	actionConfig := new(action.Configuration)
	settings := cli.New()
	settings.KubeAPIServer = cluster.Endpoint
	settings.KubeToken = cluster.Token
	settings.KubeCaFile = cluster.CAFile
	if err := actionConfig.Init(settings.RESTClientGetter(), releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v)
	}); err != nil {
		panic(err)
	}

	iCli := action.NewList(actionConfig)
	iCli.Selector = mycfg.JanitorLabel
	rel, err := iCli.Run()
	if err != nil {
		panic(err)
	}
	releaseList := NameList(rel)
	log.Debugf("Got a list of releases: %v", releaseList)

	// loop through releases, track errors
	errorCount := 0

	for _, release := range rel {
		log := log.WithFields(log.Fields{
			"namespace": release.Namespace,
			"release":   release.Name,
		})
		expired, err := CheckReleaseExpired(*release)
		if err != nil {
			errorCount++
			log.Error(err)
		}
		if expired {
			log.Infof("deleting release %s in namespace %s", release.Name, release.Namespace)
			if err := actionConfig.Init(settings.RESTClientGetter(), release.Namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			}); err != nil {
				log.Fatal(err)
			}
			cli := action.NewUninstall(actionConfig)
			cli.DryRun = true
			// rel, err := cli.Run(release.Name)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// log.Infof("deleted: %s", rel.Info)
		}
	}

	// Finally throw the last error!.
	if errorCount > 0 {
		log.Errorf("Encountered %d errors while cleaning up helm releases - investigation required.", errorCount)
	}
}

// CheckReleaseExpired will return true if the release should be deleted.
// Safely returns false for any errors that occur.
func CheckReleaseExpired(r release.Release) (bool, error) {
	log.Debugf("Processing release: %s in namespace: ", r.Name, r.Namespace)
	ttlKey := janitorconfig.TTLKey
	expiryKey := janitorconfig.ExpiryKey
	now := time.Now()
	deployedTime := r.Info.LastDeployed
	if val, ok := r.Labels[ttlKey]; ok {
		log.Debugf("found %s: %s", ttlKey, val)
		timeLeft, err := utils.ParseTime(val)
		if err != nil {
			log.Errorf("%s key value: %s not valid - using default 7 days instead", ttlKey, val)
			timeLeft.Seconds = janitorconfig.DefaultTTL
		}
		var expirySeconds int64 = int64(timeLeft.Seconds)
		// work off the modifiedAt key - not required.
		// if modifiedTime, ok := r.Labels["modifiedAt"]; ok {
		// 	log.Debugf("release: %s was modifiedAt %s", r.Name, modifiedTime)
		// 	if n, err := strconv.ParseInt(modifiedTime, 10, 64); err == nil {
		// 		if currentTime-n-expirySeconds > 0 {
		// 			return true, nil
		// 		} else {
		// 			return false, nil
		// 		}
		// 	} else {
		// 		return false, fmt.Errorf("modifiedTime cannot me made to int64")
		// 	}
		if now.Unix()-deployedTime.Unix()-expirySeconds > 0 {
			log.Debugf("release has expired - last deployed at: %v", deployedTime)
			return true, nil
		}
		return false, nil
		// work off helm-janitor/expiry key instead.
	} else if val, ok := r.Labels[expiryKey]; ok {
		log.Debugf("found %s: %s", expiryKey, val)
	} else {
		return false, fmt.Errorf("no %s or %s found", ttlKey, expiryKey)
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

func cleanup(cafile string) {
	log.Infof("Cleaning CAFile: %s", cafile)
	os.Remove(cafile)
}
