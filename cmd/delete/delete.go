package delete

import (
	"os"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
)

// RunV2 is the main exported method to delete a release
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

	actionConfig := new(action.Configuration)
	sr.Deleterelease(actionConfig, mycfg.)
}

// ValidateArgs should check the argument (release)
func ValidateArgs(a []string) {
	return
}