package delete

// bit of a mash between delete and scane

import (
	"os"

	"github.com/edify42/helm-janitor/cmd/scan"
	internalhelm "github.com/edify42/helm-janitor/internal/helm"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// main struct for this file.

func RunDeleteSet(sr scan.InputRun) {

	b := sr.Makeekscfg()
	mycfg := sr.Config()
	cfg := sr.Makeawscfg()
	cluster := sr.Getekscluster(cfg, b)
	if !mycfg.DebugFlag {
		log.Info("should remove the CA file " + cluster.CAFile)
		defer os.Remove(cluster.CAFile)
	} else {
		log.Info("DEBUG flag set - won't remove the cluster CA file")
	}

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
		log.Infof("deleting release %s in namespace %s", release.Name, release.Namespace)
		del := internalhelm.NewDelete()
		sr.Deleterelease(cluster, actionConfig, release, del)
	}

	// Finally throw the last error!.
	if errorCount > 0 {
		log.Errorf("Encountered %d errors while cleaning up helm releases - investigation required.", errorCount)
	}
}
