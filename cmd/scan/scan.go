package scan

import (
	"fmt"
	"os"
	"time"

	janitorconfig "github.com/lendi-au/helm-janitor/internal/config"
	internalhelm "github.com/lendi-au/helm-janitor/internal/helm"
	"github.com/lendi-au/helm-janitor/pkg/utils"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// main struct for this file.

func RunV2(sr InputRun) {

	b := sr.Makeekscfg()
	mycfg := sr.Config()
	cfg := sr.Makeawscfg()
	cluster := sr.Getekscluster(cfg, b)
	if !mycfg.DebugFlag {
		log.Info("should clean " + cluster.CAFile)
		defer os.Remove(cluster.CAFile)
	} else {
		log.Info("DEBUG flag set - won't actually delete anything.")
	}
	defer os.Remove(cluster.CAFile)

	actionConfig := new(action.Configuration)

	newList := internalhelm.New()
	rel := sr.Getreleases(cluster, actionConfig, newList)
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
			del := internalhelm.NewDelete()
			sr.Deleterelease(cluster, actionConfig, release, del)
		}
	}

	// Finally throw the last error!.
	if errorCount > 0 {
		log.Errorf("Encountered %d errors while cleaning up helm releases - investigation required.", errorCount)
	}
}

// Annotations
type Annotations struct {
	Expiry string `json:"janitor/expires,omitempty"`
	Ttl    string `json:"janitor/ttl,omitempty"`
}

// CheckReleaseExpired will return true if the release should be deleted.
// Safely returns false for any errors that occur.
func CheckReleaseExpired(r release.Release) (bool, error) {
	log.Debugf("Processing release: %s in namespace: %s", r.Name, r.Namespace)
	ttlKey := janitorconfig.TTLKey
	expiryKey := janitorconfig.ExpiryKey
	now := time.Now()
	deployedTime := r.Info.LastDeployed
	annotations := r.Config[janitorconfig.AnnotationKey]
	var output Annotations
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &output,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	decoder.Decode(annotations)
	// Process the `janitor/ttl` key
	if output.Ttl != "" {
		log.Debugf("found %s: %s", ttlKey, output.Ttl)
		timeLeft, err := utils.ParseTime(output.Ttl)
		if err != nil {
			log.Errorf("%s key value: %s not valid - using default 7 days instead", ttlKey, output.Ttl)
			timeLeft.Seconds = janitorconfig.DefaultTTL
		}
		var ttl int64 = int64(timeLeft.Seconds)
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
		//
		if now.Unix()-deployedTime.Unix()-ttl > 0 {
			log.Debugf("release has expired - last deployed at: %v", deployedTime)
			return true, nil
		}
		return false, nil
		// work off janitor/expires key instead.
	} else if output.Expiry != "" {
		log.Debugf("found %s: %s", expiryKey, output.Expiry)
		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, output.Expiry)
		if err != nil {
			log.Error(err)
		}
		if now.Unix()-t.Unix() > 0 {
			log.Debugf("release has expired - expired at: %v", t.Local())
			return true, nil
		}
	} else {
		return false, fmt.Errorf("no %s or %s found", ttlKey, expiryKey)
		// silently skip only - don't panic
	}
	// TODO: remove last catch all...
	log.Debugf("Nothing found here... %v", output)
	return false, nil
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
